package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDB *gorm.DB

func DB_Handler() {
	log.Println("------DB_Handler---------")
	//MysqlDB, err := gorm.Open("mysql", "root:password@tcp(127.0.0.1:3306)/solarepg?charset=utf8")
	password := os.Getenv("DB_PASSWORD")
	db, err := sql.Open("mysql", "root:"+password+"@tcp(127.0.0.1:3306)/solarepg")
	if err != nil {
		log.Println("connection to mysql failed:", err)
		return
	} else {
		log.Println("connect database success")
	}

	if err != nil {
		log.Fatal(err)
	}
	//rows, err := db.Query("SELECT name FROM stream_room WHERE age = ?", age)
	query(db)
	db.Close()
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

func query(db *sql.DB) {

	//var streamId uint64
	//var stid string
	var matchId uint64
	var roomId string
	var locationId string
	var state sql.NullInt64
	var score sql.NullInt64
	//var kps string
	var ip string

	log.Println("------query(db *sql.DB)---------")
	start := time.Now()

	rows, _ := db.Query("SELECT `matchId`,`roomId`,`locationId`,`state`,`score`,`ip` FROM stream_room GROUP by `matchId`")

	defer rows.Close()

	//log.Println(rows.ColumnTypes())
	for rows.Next() {

		if err := rows.Scan(&matchId, &roomId, &locationId, &state, &score, &ip); err != nil {
			log.Println("------error---------")
			log.Fatal(err)
		}
		fmt.Printf("matchId :%d ,roomId: %s ,locationId: %s , state: %v ,score: %v  , ip: %s\n", matchId, roomId, locationId, state, score, ip)
	}
	end := time.Now()
	log.Println("方式1 query total time:", end.Sub(start).Seconds())
}
