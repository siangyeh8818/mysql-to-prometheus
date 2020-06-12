package lib

import (
	"time"
)

type Data []Datum

type Datum struct {
	MatchID         uint64
	SpecialID       string
	LocationID      string
	Location        string
	Score           int
	StateCategory   string
	UpdatedAt       time.Time
	DBtable         string
	StreamAPIStatus int
	SportName       string
}

type BaseConfig struct {
	MYSQL_ADDRESS  string
	MYSQL_PASSWORD string
}

func (instance *Datum) SetDBtable(newtable string) {
	instance.DBtable = newtable
}

func (instance *Datum) SetMatchID(id uint64) {
	instance.MatchID = id
}

func (instance *Datum) SetSpecialID(id string) {
	instance.SpecialID = id
}

func (instance *Datum) SetStreamAPIStatus(status int) {
	instance.StreamAPIStatus = status
}

func (instance *Datum) SetLocationID(id string) {
	instance.LocationID = id
}

func (instance *Datum) SetLocation(location string) {
	instance.Location = location
}

func (instance *Datum) SetSportName(name string) {
	instance.SportName = name
}

func (instance *Datum) SetScore(score int) {
	instance.Score = score
}

func (instance *Datum) SetStateCategory(category string) {
	instance.StateCategory = category
}

func MergeSlice(s1 Data, s2 Data) []Datum {
	slice := make([]Datum, len(s1)+len(s2))

	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
