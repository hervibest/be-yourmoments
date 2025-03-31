package adapter

import (
	"be-yourmoments/user-svc/internal/entity"
	"be-yourmoments/user-svc/internal/helper/utils"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oklog/ulid/v2"
)

type JWTAdapter interface {
	GenerateAccessToken(userId string) (*entity.AccessToken, error)
	GenerateRefreshToken(userId string) (*entity.RefreshToken, error)
	VerifyAccessToken(token string) (*entity.AccessToken, error)
	VerifyRefreshToken(token string) (*entity.RefreshToken, error)
}

type jwtAdapter struct {
	accessSecretByte  []byte
	refreshSecretByte []byte
	accessExpireTime  time.Duration
	refreshExpireTime time.Duration
}

func NewJWTAdapter() JWTAdapter {
	accessSecret := utils.GetEnv("ACCESS_TOKEN_SECRET")
	refreshSecret := utils.GetEnv("REFRESH_TOKEN_SECRET")

	accessExpireStr := utils.GetEnv("ACCESS_TOKEN_EXP_MINUTE")
	refreshExpireStr := utils.GetEnv("REFRESH_TOKEN_EXP_DAY")

	accessExpireInt, _ := strconv.Atoi(accessExpireStr)
	refreshExpirInt, _ := strconv.Atoi(refreshExpireStr)

	return &jwtAdapter{
		accessSecretByte:  []byte(accessSecret),
		refreshSecretByte: []byte(refreshSecret),
		accessExpireTime:  time.Duration(accessExpireInt),
		refreshExpireTime: time.Duration(refreshExpirInt),
	}
}

func (c *jwtAdapter) GenerateAccessToken(userId string) (*entity.AccessToken, error) {
	expirationTime := time.Now().Add(time.Minute * c.accessExpireTime)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	stringToken, err := token.SignedString(c.accessSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.AccessToken{
		UserId:    userId,
		Token:     stringToken,
		ExpiresAt: expirationTime,
	}, nil
}

func (c *jwtAdapter) GenerateRefreshToken(userId string) (*entity.RefreshToken, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * c.refreshExpireTime)

	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(c.refreshSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.RefreshToken{
		UserId:    userId,
		Token:     stringToken,
		ExpiresAt: expirationTime,
	}, nil
}

func (c *jwtAdapter) VerifyAccessToken(token string) (*entity.AccessToken, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.accessSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	accessTokenDetail := &entity.AccessToken{}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if ok && tokenClaims.Valid {
		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			log.Println("user_id not a string")
			return nil, errors.New("invalid token claims")
		}

		authorized, ok := claims["authorized"].(bool)
		if !ok {
			log.Println("authorized is not a bool")
			return nil, errors.New("invalid token claims")
		}

		if !authorized {
			log.Println("unathorize")
			return nil, errors.New("invalid token claims")
		}

		_, err := ulid.Parse(userIdStr)
		if err != nil {
			log.Println("failed to parse ulid:", err)
			return nil, errors.New("invalid token claims")
		}

		accessTokenDetail.UserId = userIdStr
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			log.Println("exp is not a float")
			return nil, errors.New("invalid exp in token claims")
		}

		expiresAt := time.Unix(int64(expFloat), 0)
		accessTokenDetail.ExpiresAt = expiresAt
		accessTokenDetail.Token = token
	}

	return accessTokenDetail, nil

}

func (c *jwtAdapter) VerifyRefreshToken(token string) (*entity.RefreshToken, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.refreshSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	refreshTokenDetail := &entity.RefreshToken{}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if ok && tokenClaims.Valid {
		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			log.Println("user_id not a string")
			return nil, errors.New("invalid token claims")
		}

		_, err := ulid.Parse(userIdStr)
		if err != nil {
			log.Println("failed to parse ulid:", err)
			return nil, errors.New("invalid token claims")
		}

		refreshTokenDetail.UserId = userIdStr
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			log.Println("exp is not a float")
			return nil, errors.New("invalid exp in token claims")
		}

		expiresAt := time.Unix(int64(expFloat), 0)
		refreshTokenDetail.ExpiresAt = expiresAt
		refreshTokenDetail.Token = token
	}

	return refreshTokenDetail, nil
}
