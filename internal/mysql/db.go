package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/siangyeh8818/mysql-to-prometheus/internal/lib"
)

var MysqlDB *gorm.DB
var locationMapping map[string]string

func DB_Handler() lib.Data {

	var result lib.Data
	log.Println("------DB_Handler---------")
	//MysqlDB, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/solarepg?charset=utf8&parseTime=true")
	password := os.Getenv("MYSQL_PASSWORD")
	user := os.Getenv("MYSQL_USER")
	mysql_address := os.Getenv("MYSQL_ADDRESS")
	db_name := os.Getenv("MYSQL_DATABASE")
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+mysql_address+")/"+db_name+"?charset=utf8&parseTime=true")

	if err != nil {
		log.Println("connection to mysql failed:", err)
		//return

	} else {
		log.Println("connect database success")
	}

	if err != nil {
		log.Fatal(err)
	}

	locationMapping = querylocation(db)
	log.Println(locationMapping)
	//rows, err := db.Query("SELECT name FROM stream_room WHERE age = ?", age)

	room_data := make(chan lib.Data)
	armenia_data := make(chan lib.Data)
	nami_data := make(chan lib.Data)
	//WHERE `updated_at` BETWEEN '2020-05-18 09:20:00.158054 +0000 UTC' AND '2020-05-18 09:25:00.158054 +0000 UTC'
	now := time.Now()
	local1, err1 := time.LoadLocation("Asia/Taipei") //等同于"UTC"
	if err1 != nil {
		fmt.Println(err1)
	}
	log.Println(now.In(local1))
	//log.Println(now)

	d, _ := time.ParseDuration(os.Getenv("INTERNAL_TIME_TO_MYSQL"))
	now_5min := now.Add(d)
	log.Println(now_5min.In(local1))
	//log.Println(now_5min)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		//armenia_data <- query(&wg, db, "SELECT `matchId`,`armeniaId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_armenia_status", "stream_armenia_status")
		armenia_data <- query(&wg, db, "SELECT `matchId`,`armeniaId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_armenia_status WHERE `updated_at` BETWEEN '"+now_5min.In(local1).Format("2006-01-02 15:04:05")+"' AND '"+now.In(local1).Format("2006-01-02 15:04:05")+"'", "stream_armenia_status")
	}()
	go func() {
		//nami_data <- query(&wg, db, "SELECT `matchId`,`namiId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_nami_status", "stream_nami_status")
		nami_data <- query(&wg, db, "SELECT `matchId`,`namiId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_nami_status WHERE `updated_at` BETWEEN '"+now_5min.In(local1).Format("2006-01-02 15:04:05")+"' AND '"+now.In(local1).Format("2006-01-02 15:04:05")+"'", "stream_nami_status")
	}()
	go func() {
		//room_data <- query(&wg, db, "SELECT `matchId`,`roomId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_room_status", "stream_room_status")
		room_data <- query(&wg, db, "SELECT `matchId`,`roomId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_room_status WHERE `updated_at` BETWEEN '"+now_5min.In(local1).Format("2006-01-02 15:04:05")+"' AND '"+now.In(local1).Format("2006-01-02 15:04:05")+"'", "stream_room_status")
	}()

	wg.Wait()
	temp_room_data := <-room_data
	temp_nami_data := <-nami_data
	temp_armenia_data := <-armenia_data

	//esult = append(result,temp_room_data)
	//result = append(result,temp_nami_data)
	//result = append(result,temp_armenia_data)

	result = lib.MergeSlice(result, temp_room_data)
	result = lib.MergeSlice(result, temp_nami_data)
	result = lib.MergeSlice(result, temp_armenia_data)

	defer db.Close()
	return result
	/*
		//表名默认就是结构体名称的复数
		//禁用默认表名的复数形式，如果置为 true，则 `User` 的默认表名是 `user
		MysqlDB.SingularTable(true)
		var stream_room StreamRoom
		//log.Println("------Inspect table -1 ---------")
		//MysqlDB.HasTable("stream_room")
		log.Println("------Inspect table ---------")
		MysqlDB.HasTable(&stream_room)
		rows := MysqlDB.Table("stream_room").Select("matchId, roomId", "locationId", "state", "score", "ip")
		log.Println("------Find ---------")
		for rows.Next() {
			//log.Println("------Next () ---------")
			var matchId uint64
			var roomId string
			var locationId string
			var state int8
			var score int8
			var ip string
			rows.Scan(&matchId, &roomId, &locationId, &state, &score, &ip)
			//fmt.Printf("matchId :%d ,roomId: %s ,locationId: %s , state: %d ,score: %d  , ip: %s\n", matchId, roomId, locationId, state, score, ip)
		}

		//log.Println(MysqlDB.First(&stream_room))

		defer MysqlDB.Close()
	*/
}

func query(wg *sync.WaitGroup, db *sql.DB, sqlexpression string, table string) lib.Data {
	defer wg.Done()
	//var streamId uint64
	//var stid string
	var data2 lib.Data
	var matchId uint64
	var spacialID string
	var locationId string
	var statecategory string
	var score int
	var update_time sql.NullTime
	//var kps string
	//var ip string
	log.Println("------query(db *sql.DB)---------")

	log.Println("------sqlexpression---------")
	log.Println(sqlexpression)
	start := time.Now()
	rows, _ := db.Query(sqlexpression)
	//rows, _ := db.Query("SELECT `matchId`,`roomId`,`locationId`,`state`,`score`,`ip` FROM stream_room GROUP by `matchId`")

	defer rows.Close()

	//log.Println(rows.ColumnTypes())
	for rows.Next() {
		//SELECT `matchId`,`roomId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_room_status`
		var tempdatum lib.Datum
		if err := rows.Scan(&matchId, &spacialID, &locationId, &score, &statecategory, &update_time); err != nil {
			log.Println("------error---------")
			log.Fatal(err)
		}
		fmt.Printf("matchId :%d ,roomId: %s ,locationId: %s , state: %v ,score: %v , updated_at: %v \n", matchId, spacialID, locationId, statecategory, score, update_time)
		tempdatum.SetLocationID(locationId)
		tempdatum.SetLocation(locationMapping[locationId])
		tempdatum.SetMatchID(matchId)
		tempdatum.SetSpecialID(spacialID)
		tempdatum.SetStateCategory(statecategory)
		tempdatum.SetScore(score)
		tempdatum.SetDBtable(table)
		data2 = append(data2, tempdatum)
	}
	end := time.Now()
	log.Println("方式1 query total time:", end.Sub(start).Seconds())
	return data2
}

func querylocation(db *sql.DB) map[string]string {
	log.Println("------querylocation---------")
	var locationId string
	var memo string
	m := make(map[string]string)

	start := time.Now()
	rows, _ := db.Query("SELECT `locationId`,`memo` FROM stream_location")
	defer rows.Close()

	for rows.Next() {
		//SELECT `matchId`,`roomId`,`locationId`,`score`,`stateCategory`,`updated_at` FROM stream_room_status`
		if err := rows.Scan(&locationId, &memo); err != nil {
			log.Println("------error---------")
			log.Fatal(err)
		}
		fmt.Printf("locationId :%s ,memo: %s \n", locationId, memo)
		m[locationId] = memo
	}
	end := time.Now()
	log.Println(" query total time:", end.Sub(start).Seconds())
	return m
}
