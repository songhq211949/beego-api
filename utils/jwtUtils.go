package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const key = "beego-api-key"

//一个JWT(json web token)由三部分组成，头部、载荷与签名。
//头部就是说明采用的是什么算法等描述
//要加密的为m,加密中用claims，包装，即为载荷
//这个key要头部的算法，对荷载加密就是签名 ,这里最后都要转换成base64
//这里没有自定义头部，使用 HS256
func CreateToken(m map[string]interface{}) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	for index, val := range m {
		claims[index] = val
	}
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(key))
	return tokenString
}
func ParseToken(tokenString string) (interface{}, bool) {
	//这是一个匿名函数，整体是parse方法
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		return "", false
	}
}

// CheckToken 检查是否合法
func CheckToken(uid string, sid string) bool {
	claims, ok := ParseToken(sid)
	if ok {
		parseUID := claims.(jwt.MapClaims)["uid"]
		return parseUID == uid

	}
	return false

}

//TestJwt 测试jwt
func TestJwt() {
	//t := time.Now() //获取当前时间
	//key := "my key is secret"
	fmt.Println("加密的key是", key)
	//这个是过期时间
	//userInfo["exp"] = strconv.FormatInt(t.UTC().UnixNano(), 10)
	//userInfo["iat"] = "0" //在什么时候签发的
	userInfo := make(map[string]interface{})
	userInfo["uid"] = "2"
	tokenString := CreateToken(userInfo)
	claims, ok := ParseToken(tokenString)
	if ok {
		// oldT, _ := strconv.ParseInt(claims.(jwt.MapClaims)["exp"].(string), 10, 64)
		// ct := t.UTC().UnixNano()
		// c := ct - oldT
		// if c > expTime {
		// 	ok = false
		// 	tokenState = "Token 已过期"
		// } else {
		// 	tokenState = "Token 正常"
		// }
		//这里就判断过期时间了
		fmt.Println("解析成功", claims)
	} else {
		fmt.Println("解析失败")
	}
	fmt.Println("加密后的token 是", tokenString)
	fmt.Println("uid 是",claims.(jwt.MapClaims)["uid"])
}
