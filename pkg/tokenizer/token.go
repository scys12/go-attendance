package tokenizer

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/scys12/go-attendance/config"
	"github.com/scys12/go-attendance/internal/types"
)

var cfg config.JWTConfig

func init() {
	cfg, _ = config.InitJWTConfig()
}

func CreateToken(userid int64, email string) (string, error) {

	atClaims := jwt.MapClaims{
		"user_id":    userid,
		"email":      email,
		"exp":        time.Now().Add(time.Minute * time.Duration(cfg.JWTExpires)).Unix(),
		"authorized": true,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func DecodeToken(token string) (map[string]interface{}, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, types.ErrTokenUnexpectedMethod
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok && !tokenParsed.Valid {
		return nil, types.ErrTokenFailedClaim
	}
	tokenDetails := make(map[string]interface{})
	tokenDetails["user_id"] = claims["user_id"]
	tokenDetails["email"] = claims["email"]
	return tokenDetails, nil
}
