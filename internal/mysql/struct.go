package mysql

import (
	"database/sql"
)

type StreamRoom struct {
	//	gorm.Model
	streamId      uint64        `gorm:"Column:streamId,AUTO_INCREMENT"`
	stid          string        `gorm:"Column:stid,NOT NULL`
	matchId       uint64        `gorm:"Column:matchId,NOT NULL"`
	roomId        string        `gorm:"Column:roomId,NOT NULL"` //房間號碼
	locationId    string        `gorm:"Column:locationId,NOT NULL"`
	category      int8          `gorm:"Column:category,NOT NULL"`
	state         int8          `gorm:"Column:state,NOT NULL"` //1:開 0:關
	score         int8          `gorm:"Column:score,`          //0=紅燈、1,2=黃燈、3,4,5=綠燈
	cost_conn     uint16        `gorm:"Column:cost_conn`
	cost_pull     sql.NullInt64 `gorm:"Column:cost_pull`
	length_pull   sql.NullInt64 `gorm:"Column:length_pull`
	kps           string        `gorm:"Column:kps,size:10"`
	stateCategory int8          `gorm:"Column:stateCategory,not null"` // 1:網宿 2:阿里 3:原始流 4:直播間smtv 5:直播間star 6:直播間hbzb
	created_at    sql.NullTime  `gorm:"Column:created_at`
	updated_at    sql.NullTime  `gorm:"Column:updated_at`
	ip            string        `gorm:"Column:ip,size:255"`
}
