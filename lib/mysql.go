package lib

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func NewEngine(config string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("mysql", config)
	if err != nil {
		return nil, err
	}
	if err = engine.Ping(); err != nil {
		return nil, err
	}
	return engine, nil
}
