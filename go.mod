module github.com/songhq211949/beego-api

go 1.13

require (
	github.com/astaxie/beego v1.12.1
	github.com/beego/samples v0.0.0-20180420090448-1c696ee4905b // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.5
	github.com/gorilla/websocket v1.4.2
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/smartwalle/alipay v1.0.2 // indirect
	google.golang.org/appengine v1.6.5 // indirect
)

//replace golang.org/x/text => github.com/golang/text v0.3.0 //github.com/golang/text v0.3.0 解决这个包下不下来
