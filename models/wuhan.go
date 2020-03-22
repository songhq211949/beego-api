package models

import (
	//"github.com/astaxie/beego/orm" //初始化数据库用的

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //使用的mysql数据驱动
)

//武汉VO
type Wuhan struct {
	Date        string `orm:"pk" json:"date"` //以空格区分打多个tag
	SureAdd     string `json:"sureAdd"`
	SureSum     string `json:"sureSum"`
	DieAdd      string `json:"dieAdd"`
	DieSum      string `json:"dieSum"`
	RecoveryAdd string `json:"recoveryAdd"`
	RecoverySum string `json:"recoverySum"`
}
//初始化数据
func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//set default database
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(songhq.club:3306)/him?charset=utf8")

	//register model
	orm.RegisterModel(new(Wuhan))

	//create table
	//orm.RunSyncdb("default", false, true)
}
