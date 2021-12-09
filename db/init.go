package db

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"patch_data/utils"

	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

var DBconn *sql.DB
var DBbuzconn *sql.DB
var DBCoreconn *sql.DB
var DBptconn *sql.DB

func initMySQL() {
	dbName := utils.GetEnv("USER_LEARN_RECORD_DB_NAME", "")
	dbUser := utils.GetEnv("USER_LEARN_RECORD_DB_USER_NAME", "")
	dbPass := utils.GetEnv("USER_LEARN_RECORD_DB_PASSWORD", "")
	dbHost := utils.GetEnv("USER_LEARN_RECORD_DB_HOST", "")
	dbPort := utils.GetEnv("USER_LEARN_RECORD_DB_PORT", "")
	dbMaxConns := utils.Atoi(utils.GetEnv("DB_MAX_CONNS", "8"), 8)
	dbMaxIdleConns := utils.Atoi(utils.GetEnv("DB_MAX_IDLE_CONNS", "4"), 4)
	dbMaxLifeTimeSecond := utils.Atoi(utils.GetEnv("DB_MAX_LIFE_TIME_SECOND", "10"), 10)
	DBconn = newDB(dbName, dbUser, dbPass, dbHost, dbPort, dbMaxConns, dbMaxIdleConns, dbMaxLifeTimeSecond)
}
func initBuz() {
	dbUser := utils.GetEnv("BUZ_DB_USER_NAME", "")
	dbPass := utils.GetEnv("BUZ_DB_PASSWORD", "")
	dbHost := utils.GetEnv("BUZ_DB_HOST", "")
	dbPort := utils.GetEnv("BUZ_DB_PORT", "")
	dbMaxConns := utils.Atoi(utils.GetEnv("DB_MAX_CONNS", "8"), 8)
	dbMaxIdleConns := utils.Atoi(utils.GetEnv("DB_MAX_IDLE_CONNS", "4"), 4)
	dbMaxLifeTimeSecond := utils.Atoi(utils.GetEnv("DB_MAX_LIFE_TIME_SECOND", "10"), 10)
	DBbuzconn = newDB("", dbUser, dbPass, dbHost, dbPort, dbMaxConns, dbMaxIdleConns, dbMaxLifeTimeSecond)
}
func initCore() {
	dbUser := utils.GetEnv("CORE_DB_USER_NAME", "")
	dbPass := utils.GetEnv("CORE_DB_PASSWORD", "")
	dbHost := utils.GetEnv("CORE_DB_HOST", "")
	dbPort := utils.GetEnv("CORE_DB_PORT", "")
	dbMaxConns := utils.Atoi(utils.GetEnv("DB_MAX_CONNS", "8"), 8)
	dbMaxIdleConns := utils.Atoi(utils.GetEnv("DB_MAX_IDLE_CONNS", "4"), 4)
	dbMaxLifeTimeSecond := utils.Atoi(utils.GetEnv("DB_MAX_LIFE_TIME_SECOND", "10"), 10)
	DBCoreconn = newDB("", dbUser, dbPass, dbHost, dbPort, dbMaxConns, dbMaxIdleConns, dbMaxLifeTimeSecond)
}
func initPT() {
	dbUser := utils.GetEnv("USER_LEARN_RECORD_DB_USER_NAME_PT", "")
	dbPass := utils.GetEnv("USER_LEARN_RECORD_DB_PASSWORD_PT", "")
	dbHost := utils.GetEnv("USER_LEARN_RECORD_DB_HOST_PT", "")
	dbPort := utils.GetEnv("USER_LEARN_RECORD_DB_PORT_PT", "")
	dbMaxConns := utils.Atoi(utils.GetEnv("DB_MAX_CONNS", "8"), 8)
	dbMaxIdleConns := utils.Atoi(utils.GetEnv("DB_MAX_IDLE_CONNS", "4"), 4)
	dbMaxLifeTimeSecond := utils.Atoi(utils.GetEnv("DB_MAX_LIFE_TIME_SECOND", "10"), 10)
	DBptconn = newDB("", dbUser, dbPass, dbHost, dbPort, dbMaxConns, dbMaxIdleConns, dbMaxLifeTimeSecond)
}
func newDB(dbName string, user string, password string, host string, port string, maxOpenConns, maxIdleConns, maxLifeTimeSecond int) *sql.DB {
	iPort, _ := strconv.Atoi(port)
	db, err := manager.New(dbName, user, password, host).Set(
		manager.SetCharset("utf8"),
		manager.SetParseTime(true),
		manager.SetLoc("UTC"),
	).Port(iPort).Open(true)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Duration(maxLifeTimeSecond) * time.Second)

	return db
}
