package jwt

import (
    "time"

    "github.com/golang-jwt/jwt"
    "go-tools/tiktok-backend/internal/utils/md5"
)

var jwtSecret []byte

type Claims struct {
    Username string `json:"username"`
    Password string `json:"password"`
    jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username, password string) (string, error) {
    nowTime := time.Now()
    expireTime := nowTime.Add(3 * time.Hour)

    claims := Claims{
        md5.MD5(username),
        md5.MD5(password),
        jwt.StandardClaims{
            ExpiresAt: expireTime.Unix(),
            Issuer:    "gin-blog",
        },
    }

    tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err := tokenClaims.SignedString(jwtSecret)

    return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
    tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if tokenClaims != nil {
        if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
            return claims, nil
        }
    }

    return nil, err
}
