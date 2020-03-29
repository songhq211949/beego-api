package models

import (
	"github.com/astaxie/beego/orm"     //初始化数据库用的
	_ "github.com/go-sql-driver/mysql" //使用的mysql数据驱动
)

//初始化数据
func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//set default database 指定时区 loc=Asia%2FShanghai&parseTime=true
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(songhq.club:3306)/him?charset=utf8&loc=Asia%2FShanghai&parseTime=true")

	//register model
	orm.RegisterModel(new(Wuhan), new(GroupUser), new(Group), new(GroupMsg), new(User), new(UserProfile),
		new(UserFriend), new(UserFriendMsg), new(UserFriendAsk),new(UserQq))
	//设置打印日志
	//orm.Debug = true
	//orm.DebugLog = orm.NewLog(os.Stdout)

	//create table
	//orm.RunSyncdb("default", false, true)
}
