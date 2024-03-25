package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/geedotrar/mygram/internal/model"
	"github.com/geedotrar/mygram/internal/repository"
	"github.com/geedotrar/mygram/pkg/helper"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersByID(ctx context.Context, id uint64) (model.User, error)
	DeleteUsersById(ctx context.Context, id uint64) (model.User, error)
	EditUser(ctx context.Context, id uint64, user model.User) (model.User, error)

	SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.UserView, error)
	GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error)
	CheckCredentials(ctx context.Context, email string, password string) (model.User, error)
}

type userServiceImpl struct {
	repo repository.UserQuery
}

func NewUserService(repo repository.UserQuery) UserService {
	return &userServiceImpl{repo: repo}
}

func (u *userServiceImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	users, err := u.repo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, err

}
func (u *userServiceImpl) GetUsersByID(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	return user, err
}

func (u *userServiceImpl) SignUp(ctx context.Context, userSignUp model.UserSignUp) (model.UserView, error) {
	dob, err := time.Parse("2006-01-02", userSignUp.Dob)
	if err != nil {
		return model.UserView{}, errors.New("invalid date of birth format")
	}

	// count age
	today := time.Now()
	age := today.Year() - dob.Year()
	if today.Month() < dob.Month() || (today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}

	// check age < 8
	if age < 8 {
		return model.UserView{}, errors.New("age must be at least 8 years old")
	}
	user := model.User{
		Username: userSignUp.Username,
		Email:    userSignUp.Email,
		Dob:      dob,
	}
	// encryption password
	// hashing
	pass, err := helper.GenerateHash(userSignUp.Password)
	if err != nil {
		return model.UserView{}, err
	}
	user.Password = pass

	getUserByEmail, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return model.UserView{}, err
	}
	if getUserByEmail.Email == user.Email {
		return model.UserView{}, errors.New("email already exist")
	}

	// store to db
	createdUser, err := u.repo.SignUp(ctx, user)
	if err != nil {
		return model.UserView{}, err
	}
	printUser := model.UserView{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Dob:      createdUser.Dob,
	}

	return printUser, err
}

func (u *userServiceImpl) GenerateUserAccessToken(ctx context.Context, user model.User) (token string, err error) {
	// generate claim
	now := time.Now()

	claim := model.StandardClaim{
		Jti: fmt.Sprintf("%v", time.Now().UnixNano()),
		Iss: "go-middleware",
		Aud: "golang-006",
		Sub: "access-token",
		Exp: uint64(now.Add(time.Hour).Unix()),
		Iat: uint64(now.Unix()),
		Nbf: uint64(now.Unix()),
	}

	userClaim := model.AccessClaim{
		StandardClaim: claim,
		UserID:        user.ID,
		Username:      user.Username,
		Dob:           user.Dob,
	}

	token, err = helper.GenerateToken(userClaim)
	return
}

func (u *userServiceImpl) CheckCredentials(ctx context.Context, email string, password string) (model.User, error) {
	// Retrieve user by email
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}

	// Check if user exists
	if user.ID == 0 {
		return model.User{}, errors.New("user not found")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, err
	}

	// Credentials are correct, return user
	return user, nil
}

func (u *userServiceImpl) EditUser(ctx context.Context, id uint64, user model.User) (model.User, error) {
	// Perform validation or additional checks here if necessary

	// Call repository to edit user
	updatedUser, err := u.repo.EditUser(ctx, id, user)
	if err != nil {
		return model.User{}, err
	}
	return updatedUser, nil
}

func (u *userServiceImpl) DeleteUsersById(ctx context.Context, id uint64) (model.User, error) {
	user, err := u.repo.GetUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	// if user doesn't exist, return
	if user.ID == 0 {
		return model.User{}, nil
	}

	// delete user by id
	err = u.repo.DeleteUsersByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, err
}
