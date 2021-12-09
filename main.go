package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"patch_data/db"
	"patch_data/models"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	sd, ed, appID := "", "", ""
	md := 0
	flag.StringVar(&sd, "sd", time.Now().Add(time.Hour*(-24)).Format("2006-01-02 15:04:05"), "起始日期")
	flag.StringVar(&ed, "ed", time.Now().Add(time.Hour*24).Format("2006-01-02"), "结束日期")
	flag.StringVar(&appID, "appid", "", "店铺id")
	flag.IntVar(&md, "md", 0, "模式, 0 读库 | 1文件")

	flag.Parse()
	if sd > ed {
		panic("开始日期必须小于结束日期")
	}
	godotenv.Load(".env")
	db.Init()
	dtc := make(chan models.AllInfo, 200000)
	wg := sync.WaitGroup{}
	wgWorker := sync.WaitGroup{}
	d, _ := os.Getwd()
	fd, e := os.OpenFile(d+"/storage/compare_"+time.Now().Format("2006-01-02"), os.O_CREATE|os.O_RDWR, os.ModePerm)
	if e != nil {
		panic(e)
	}
	// 消费者
	for i := 0; i < 3000; i++ {
		fmt.Println("worker", i)
		wgWorker.Add(1)
		go func() {
			defer wgWorker.Done()
			for data := range dtc {
				// 读库补考式,训练营成绩
				if md == 2 {
					opg, tm := getScore(data.ShopId, data.UserId, data.ResourceId)
					iopg := math.Ceil(float64(opg))
					if opg > 0 {
						func(appid, uid, rid, pgs, tm string) {
							repair(appid, uid, rid, pgs, tm)
						}(data.ShopId, data.UserId, data.ResourceId, strconv.Itoa(int(iopg)), tm)
					}
				} else if md == 1 {
					// 读文件,补时长
					func(appid, uid, rid string, rl int) {
						repairVideo(appid, uid, rid, rl)
					}(data.ShopId, data.UserId, data.ResourceId, data.ResLen)
				} else {
					// 巡检
					ComparePressureTst(data, fd)
				}
			}
		}()
	}
	if md == 0 {
		db.Read(dtc, &wg, sd)
	} else if md == 1 {
		wg.Add(1)
		go ReadFromFile(dtc, &wg, 3)
		wg.Wait()
	} else {
		go db.ReadEmptyCampusScore(&wg, appID, dtc)
	}
	//消费
	wgWorker.Wait()
	fmt.Println("all work done!")
	return
}
func ComparePressureTst(info models.AllInfo, fd *os.File) {
	sql, ok := db.ReadRow(info.ShopId, info.UserId, info.ResourceId, info.AgentType, info.TableID)
	//  agentType ==1 并且 不为wxb店铺 并且不包含考试打卡表单才入文件
	needLog := func(info models.AllInfo) (ok bool) {
		if info.AgentType != 1 || info.ShopId == "appe0MEs6qX8480" || info.ResourceType == 16 || info.ResourceType == 27 || info.ResourceType == 26 {
			return
		}
		return true
	}
	if !ok && needLog(info) {
		pendix := fmt.Sprintf(" AND resource_type = %d;\n", info.ResourceType)
		_, err := fd.WriteString(sql[:len(sql)-1] + pendix)
		if err != nil {
			panic(err)
		}
	}

}
func ReadFromFile(dtc chan models.AllInfo, wg *sync.WaitGroup, rt int) {
	defer wg.Done()
	dir, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	fd, e := os.Open(dir + "/input.txt")
	if e != nil {
		panic(e)
	}
	res, e := ioutil.ReadAll(fd)
	if e != nil {
		panic(e)
	}
	arr := strings.Split(string(res), "\n")
	for _, v := range arr {
		dt := strings.Split(v, ",")
		it, err := strconv.Atoi(dt[3])
		it = db.ReadResLen(dt[2])
		if err != nil {
			panic(err)
		}
		dtc <- models.AllInfo{
			ShopId:       dt[0],
			UserId:       dt[1],
			ResourceId:   dt[2],
			ResLen:       it,
			ResourceType: rt,
		}
	}
	close(dtc)
}

func getScore(appID, userID, resourceID string) (float32, string) {
	return db.ReadExam(appID, userID, resourceID)
}

func repairVideo(appID, userID, resourceID string, rl int) { // 查时长
	// rl := db.ReadResLen(resourceID)
	fmt.Println("开始补", time.Now(), "应该结束时间:", time.Now().Add(time.Second*time.Duration(rl)), appID, userID, resourceID, "时长", rl)
	for ; rl > 0; rl -= 10 {
		// 修数据
		uri := "https://learnreport.xiaoeknow.com/v1/learnRecord/pushData"
		rsp, err := http.PostForm(uri, url.Values{
			"app_id":             []string{appID},
			"user_id":            []string{userID},
			"resource_id":        []string{resourceID},
			"resource_type":      []string{"3"},
			"channel_id":         []string{"patch"},
			"progress":           []string{"100"},
			"org_learn_progress": []string{""},
			"agent_type":         []string{"1"},
			"display_state":      []string{"1"},
			"is_first":           []string{"1"},
			"content_app_id":     []string{""},
			"is_play":            []string{"1"},
			"is_try":             []string{"0"},
		})
		if err != nil {
			panic(err)
		}
		r, e := ioutil.ReadAll(rsp.Body)
		if e != nil {
			panic(e)
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", appID, userID, resourceID, string(r))
		time.Sleep(time.Second * 10)
	}
}
func repair(appID, userID, resourceID, orgProgress, tm string) {
	uri := "https://learnreport.xiaoeknow.com/v1/learnRecord/pushData"
	rsp, err := http.PostForm(uri, url.Values{
		"app_id":             []string{appID},
		"user_id":            []string{userID},
		"resource_id":        []string{resourceID},
		"resource_type":      []string{"27"},
		"channel_id":         []string{"patch"},
		"progress":           []string{"100"},
		"org_learn_progress": []string{orgProgress},
		"stay_time":          []string{"0"},
		"agent_type":         []string{"1"},
		"display_state":      []string{"1"},
		"is_first":           []string{"1"},
		"content_app_id":     []string{""},
	})
	if err != nil {
		panic(err)
	}
	r, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		panic(e)
	}
	fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n", appID, userID, resourceID, orgProgress, tm, string(r))
}
