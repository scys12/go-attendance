package service

import (
	"github.com/scys12/go-attendance/internal/models"
	"github.com/scys12/go-attendance/internal/repositories"
	"github.com/scys12/go-attendance/internal/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.IUserRepo
}

func NewUserService(repo repositories.IUserRepo) IUserService {
	return &UserService{
		userRepo: repo,
	}
}

func (u *UserService) RegisterUser(req models.RegisterRequest) (*models.User, error) {
	isEmailExist, err := u.userRepo.IsEmailExist(req.Email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "RegisterUser",
			"error":    err.Error(),
		}).Errorln("[UserService] Failed to get existing email")

		return nil, err
	}

	if isEmailExist {
		return nil, err
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "RegisterUser",
			"error":    err.Error(),
		}).Errorln("[UserService] Failed to generate from password bcrypt")

		return nil, err
	}

	user := models.User{
		Email:    req.Email,
		Password: string(pwd),
		FullName: req.FullName,
	}

	id, err := u.userRepo.InsertUser(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "RegisterUser",
			"error":    err.Error(),
		}).Errorln("[UserService] Failed to insert database")

		return nil, err
	}

	user.ID = id

	return &user, nil

}

func (u *UserService) LoginUser(req models.LoginRequest) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "LoginUser",
			"error":    err.Error(),
		}).Errorln("[UserService] Failed to get user email")

		return nil, types.ErrEmailNotFound
	}

	pwdCheck := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if pwdCheck != nil {
		logrus.WithFields(logrus.Fields{
			"function": "LoginUser",
			"error":    pwdCheck.Error(),
		}).Errorln("[UserService] Failed to get user email")

		return nil, types.ErrEmailPasswordWrong
	}

	return user, nil
}
