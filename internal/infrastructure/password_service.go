package infrastructure

import (
	"regexp"
	"write_base/internal/domain"

	"golang.org/x/crypto/bcrypt"
)




type passwordService struct{}

func NewPasswordService() domain.IPasswordService{
	return &passwordService{}
}
func (p *passwordService) HashPassword(password string) (string, error){

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err

}

func (p *passwordService) VerifyPassword(hashedPassword, password string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
	
	return err == nil
}
func (ps * passwordService) IsPasswordStrong(password string) bool {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[\W_]`).MatchString(password)
	)
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}