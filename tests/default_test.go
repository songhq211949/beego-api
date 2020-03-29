package tests

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/songhq211949/beego-api/models"
	"github.com/songhq211949/beego-api/utils"
)

func TestBee(t *testing.T) {
	utils.TestJwt()
}
func TestBee2(t *testing.T) {
	utils.TestJwt()
}
func TestCheck(t *testing.T) {
	sid := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiIxMjMifQ.eHuIb_iOD7OM4vK9TXS27GoRdZNMPelKBVpvtenOjao"
	result := utils.CheckToken("123", sid)
	fmt.Println(result)
}
func TestParseResult(t *testing.T) {
	vo := new(models.QqOpenIdResVO)
	result := `callback({"client_id":"YOUR_APPID","openid":"YOUR_OPENID"})`
	index1 := strings.Index(result, "(")
	index2 := strings.Index(result, ")")
	jsonStr := result[index1+1 : index2]
	json.Unmarshal([]byte(jsonStr), vo)
	fmt.Println(vo)

}
