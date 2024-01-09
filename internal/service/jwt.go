package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zetsux/gin-gorm-template-clean/common/standard"
)

type JWTService interface {
	GenerateToken(id string, role string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetIDByToken(token string) (string, error)
	GetRoleByToken(token string) (string, error)
}

type jwtCustomClaim struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    standard.ENUM_ROLE_ADMIN,
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "jwt_secret_key"
	}
	return secretKey
}

func (j *jwtService) GenerateToken(id string, role string) string {
	claims := &jwtCustomClaim{
		id,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return t
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (any, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) GetIDByToken(token string) (string, error) {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}
	claims := t_Token.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["id"])
	return id, nil
}

func (j *jwtService) GetRoleByToken(token string) (string, error) {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}
	claims := t_Token.Claims.(jwt.MapClaims)
	role := fmt.Sprintf("%v", claims["role"])
	return role, nil
}
