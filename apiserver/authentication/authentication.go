package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"oliujunk/server/database"
	"time"
)

var (
	SignKey = "whxph"
)

func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get("token")
		if token == "" {
			context.JSON(http.StatusOK, gin.H{
				"status":  -1,
				"message": "请求未携带token，无权限访问",
			})
			context.Abort()
			return
		}

		j := &JWT{[]byte(SignKey)}
		claims, err := j.ParseToken(token)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"status":  -1,
				"message": "token校验失败: " + err.Error(),
			})
			context.Abort()
			return
		}
		context.Set("claims", claims)
	}
}

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
}

func GenerateToken(user database.User) string {
	j := &JWT{
		[]byte(SignKey),
	}
	claims := CustomClaims{
		user.ID,
		user.Username,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,     // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24, // 过期时间 一小时
			Issuer:    SignKey,                      // 签名的发行者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.SigningKey)

	if err != nil {
		return ""
	}

	return tokenStr
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
