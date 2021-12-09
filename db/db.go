package db

import (
	"database/sql"
	"fmt"
	"log"
	"patch_data/models"
	"sync"
)

func Init() {
	initMySQL()
	initBuz()
	initCore()
	// 压测数据库
	initPT()
}
func ReadRow(appID, uid, resID string, agentType int, tableID string) (res string, ok bool) {
	SQL := `SELECT id,app_id, user_id,resource_id,resource_type,learn_progress,max_learn_progress,org_learn_progress,is_finish,finished_at,stay_time,spend_time,last_learn_time,state,created_at,updated_at,content_app_id,display_state,product_id  FROM db_ex_learn_record.t_learn_record_%s WHERE app_id= '%s' AND user_id = '%s' AND resource_id = '%s' AND agent_type = %d;`
	sqlStr := fmt.Sprintf(SQL, tableID, appID, uid, resID, agentType)
	row := DBptconn.QueryRow(sqlStr)
	rcd := models.AllInfo{}
	err := row.Scan(&rcd.Id, &rcd.ShopId, &rcd.UserId, &rcd.ResourceId, &rcd.ResourceType, &rcd.OrgLearnProgress, &rcd.MaxLearnProgress, &rcd.OrgLearnProgress, &rcd.IsFinish, &rcd.Fa, &rcd.St, &rcd.Spt, &rcd.Llt, &rcd.State, &rcd.Ca, &rcd.Ua, &rcd.Cappid, &rcd.Dstate, &rcd.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			res = sqlStr
			return
		}
		panic(err)
	}
	ok = true
	return
}
func Read(dc chan<- models.AllInfo, wg *sync.WaitGroup, startDate string) {
	// 直接扫就好了.
	SQL := `SELECT id,app_id, user_id,resource_id,agent_type,resource_type, learn_progress,max_learn_progress ,org_learn_progress ,is_finish,finished_at,stay_time,spend_time,last_learn_time,state,created_at,updated_at,content_app_id,display_state,product_id  FROM db_ex_learn_record.t_learn_record_%02d WHERE created_at > '%s'`
	fmt.Println(startDate, "---")
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sqlStr := fmt.Sprintf(SQL, i, startDate)
			rows, err := DBconn.Query(sqlStr)
			if err != nil {
				log.Fatal(err)
			}

			for rows.Next() {
				rcd := models.AllInfo{}
				err = rows.Scan(&rcd.Id, &rcd.ShopId, &rcd.UserId, &rcd.ResourceId, &rcd.AgentType, &rcd.ResourceType, &rcd.OrgLearnProgress, &rcd.MaxLearnProgress, &rcd.OrgLearnProgress, &rcd.IsFinish, &rcd.Fa, &rcd.St, &rcd.Spt, &rcd.Llt, &rcd.State, &rcd.Ca, &rcd.Ua, &rcd.Cappid, &rcd.Dstate, &rcd.ProductID)
				rcd.TableID = fmt.Sprintf("%02d", i)
				if err != nil {
					log.Fatal(err)
				}
				dc <- rcd
			}
		}(i)
	}
	wg.Wait()
	close(dc)
}

func ReadResLen(ID string) int {
	SQL := `SELECT id,video_length FROM db_ex_business.t_video WHERE id = ?`
	row := DBCoreconn.QueryRow(SQL, ID)
	id := ""
	vlen := 0
	err := row.Scan(&id, &vlen)
	if err != nil && sql.ErrNoRows != err {
		panic(err)
	}
	return vlen
}
func ReadExam(appID, userID, resourceID string) (float32, string) {
	SQL := `select score,updated_at from db_ex_examination.t_participate_exam_user where app_id = ? and exam_id = ? and user_id = ? order by score desc limit 1`

	row := DBbuzconn.QueryRow(SQL, appID, resourceID, userID)

	score := float32(0)
	tm := ""
	err := row.Scan(&score, &tm)
	if err != nil && sql.ErrNoRows != err {
		panic(err)
	}
	return score, tm
}

func ReadEmptyCampusScore(wg *sync.WaitGroup, appID string, dc chan models.AllInfo) {
	defer wg.Wait()
	wg.Add(1)
	pageSize := 1000
	SQLtmpl := `select app_id,user_id,resource_id,term_id FROM db_ex_camp.t_task_user WHERE app_id = ? AND is_task_finish = 1 AND resource_type = 27 AND exam_score = 0 LIMIT %d,%d;`
	// 每次1000条 直到扫描完成
	go func() {
		defer func() {
			close(dc)
			wg.Done()
		}()

		for i := 0; ; i++ {
			SQL := fmt.Sprintf(SQLtmpl, i*pageSize, pageSize)
			rows, err := DBbuzconn.Query(SQL, appID)
			if err != nil {
				log.Fatal(err)
			}
			total := 0
			for ; rows.Next(); total++ {
				rcd := models.AllInfo{}
				err = rows.Scan(&rcd.ShopId, &rcd.UserId, &rcd.ResourceId, &rcd.ProductID)
				if err != nil {
					log.Fatal(err)
				}
				dc <- rcd
			}
			if total == 0 {
				break
			}
		}

	}()
}
