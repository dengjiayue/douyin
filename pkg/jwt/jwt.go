package my_jwt

import (
	"time"

	allerrors "github.com/808-not-found/tik_duck/pkg/allerrors"
	jwt "github.com/dgrijalva/jwt-go"
)

// 设置token过期时间:24小时(时间戳)
const TokenExpireDuration = time.Hour * 24

// 设置jwt加密密钥
var MySecret = []byte("dengjiayue0804")

// 自定义声明结构体并内嵌jwt.StandardClaims
type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenToken 生成JWT.
func GenToken(user_id int64) (string, error) {
	c := MyClaims{
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "dengjiayue",
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT.
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i any, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, allerrors.ErrJWTParseTokenRun()
}
