package usecase

import (
	"context"
	"strings"
	"write_base/internal/domain"
)



type userUsercase struct{
	userRepo domain.IUserRepository
	passwordService domain.IPasswordService
	tokenService domain.ITokenService
}

// func NewUserUsecase(repo domain.IUserRepository,pass domain.IPasswordService, tk domain.ITokenService ) domain.IUserUsecase{
// 	return &userUsercase{userRepo: repo, passwordService: pass, tokenService: tk}
// }


func (uu *userUsercase) Register(ctx context.Context, user *domain.User) error{
if !uu.passwordService.IsPasswordStrong(user.Password){
	return domain.ErrWeakPassword

}
hashed, err := uu.passwordService.HashPassword(user.Password)
if err != nil{
	return err
}
user.Password = hashed
return uu.userRepo.CreateUser(ctx, user)
}

func(uu *userUsercase) Login(ctx context.Context, user *domain.User) (*domain.LoginResult, error){
    if strings.Contains(user.Email, "@") {
        existingUser, err := uu.userRepo.GetByEmail(ctx, user.Email)
		if err != nil{
			return 
		}
    } else {
        user = repo.FindByUsername(input.Username)
    }

	existingUser, err := uu.userRepo.Login(ctx, user)
	if err != nil {
		return nil, err
	}

	if !u.passwordService.VerifyPassword(existingUser.Password, user.Password){
		return nil, domain.ErrInvalidCredentials 
	}
	token, err:= u.jwtservice.GenerateToken(user)
	if err != nil{
		return nil, err
	}
	return token, nil
}