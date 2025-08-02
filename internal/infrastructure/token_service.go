package infrastructure

import (
	"time"
	"write_base/internal/domain"

	"github.com/dgrijalva/jwt-go"
)



type jwtService struct{
	jwtSecret []byte
}

func NewJWTService(secret []byte) domain.ITokenService {
	return &jwtService{jwtSecret: secret}
}
func (j *jwtService) GenerateAccessToken(user *domain.User) (string, error){
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}


func (j *jwtService) GenerateRefreshToken(user *domain.User) (string, error){
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(24 * time.Hour * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.jwtSecret)
}
func (j *jwtService) ValidateAccessToken(tokenString string) (*domain.AuthClaims, error){
	claim, err := j.parseToken(tokenString)
	if err != nil{
		return nil, err
	}

	authCalim := &domain.AuthClaims{
		UserID: claim["user_id"].(string),
		Role: claim["role"].(string),
	}

	return authCalim, nil
}

func (j *jwtService) ValidateRefreshToken(tokenString string) (*domain.AuthClaims, error){
	claim, err := j.parseToken(tokenString)
	if err != nil{
		return nil, err
	}

	authCalim := &domain.AuthClaims{
		UserID: claim["user_id"].(string),
		Role: claim["role"].(string),
	}

	return authCalim, nil
}

func (j *jwtService) parseToken(tokenString string)(jwt.MapClaims, error){
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnexpectedSigningMethod
		}
		return j.jwtSecret, nil
	})

	if err != nil{
		return nil, err
	}
	if !token.Valid {
		return nil, domain.ErrJWTExpired
	}
	claim , ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}
	return claim, nil

}