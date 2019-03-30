package config

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/go-xorm/xorm"
)

func DbInit() (*xorm.Engine, error) {
	dbcfg, err := ini.Load("./conf/config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	dbcfgSec := dbcfg.Section("db")
	x, err := xorm.NewEngine(dbcfgSec.Key("driver_name").String(), dbcfgSec.Key("data_source_name").String())
	if err != nil {
		log.Fatal("CONNETION DB FAIL ： ", err)
	}
	idle, _ := dbcfgSec.Key("max_idle_conns").Int()
	open, _ := dbcfgSec.Key("max_open_conns").Int()
	showsql, _ := dbcfgSec.Key("show_sql").Bool()
	x.SetMaxIdleConns(idle)
	x.SetMaxOpenConns(open)
	x.ShowSQL(showsql)
	return x, nil
}
