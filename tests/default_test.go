package tests

import (
	"fmt"
	"testing"

	"github.com/songhq211949/beego-api/utils"
)

func TestBee(t *testing.T) {
	utils.TestJwt()
}
func TestBee2(t *testing.T) {
	utils.TestJwt()
}
func TestCheck(t *testing.T){
	sid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxMjMifQ.eHuIb_iOD7OM4vK9TXS27GoRdZNMPelKBVpvtenOjao"
	result := utils.CheckToken("123",sid)
	fmt.Println(result)
}

