package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"user-service/db/postgresql/managers"
	pb "user-service/services/genproto"
)

type UserService struct {
	UM managers.UserManager
	pb.UnimplementedUserServiceServer
}

const (
	emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	phoneRegexPattern = `^\+998(33|50|90|91|93|94|95|97|99)\d{7}$`
)

func NewUserService(db *sql.DB) *UserService {
	return &UserService{UM: *managers.NewUserManager(db)}
}

func (s *UserService) Create(ctx context.Context, user *pb.UserCreatedReq) (*pb.UserCreatedResp, error) {
	matchEmail, err := regexp.Match(emailRegexPattern, []byte(user.Email))
	if err != nil {
		return nil, err
	}
	if !matchEmail {
		return nil, errors.New("email is not valid")
	}
	matchPhone, err := regexp.Match(phoneRegexPattern, []byte(user.Phone))
	if err != nil {
		return nil, err
	}
	if !matchPhone {
		return nil, errors.New("phone is not valid")
	}
	id, err := s.UM.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if err := s.UM.CreateCart(id); err != nil {
		fmt.Println("Error creating cart:", err)
	}

	return &pb.UserCreatedResp{Id: id}, nil
}

func (s *UserService) GetById(context context.Context, id *pb.UserGetByIDReq) (*pb.UserGetByIDResp, error) {
	user, err := s.UM.GetUserByID(id.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAll(context context.Context, void *pb.Void) (*pb.UsersGetAllResp, error) {
	users, err := s.UM.GetAllUsers()

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) Update(context context.Context, user *pb.UserUpdatedReq) (*pb.Void, error) {
	if !s.isValidUserID(user.Id) {
		return nil, errors.New("user not found")
	}
	matchEmail, err := regexp.Match(emailRegexPattern, []byte(user.Email))
	if err != nil {
		return nil, err
	}
	if !matchEmail {
		return nil, errors.New("email is not valid")
	}
	matchPhone, err := regexp.Match(phoneRegexPattern, []byte(user.Phone))
	if err != nil {
		return nil, err
	}
	if !matchPhone {
		return nil, errors.New("phone is not valid")
	}
	err = s.UM.UpdateUser(user)
	return nil, err
}

func (s *UserService) Delete(context context.Context, id *pb.UserDeletedReq) (*pb.Void, error) {
	if !s.isValidUserID(id.GetId()) {
		return nil, errors.New("user not found")
	}
	err := s.UM.DeleteUser(id.GetId())
	return nil, err
}

func (s *UserService) isValidUserID(id string) bool {
	user, err := s.UM.GetUserByID(id)
	if err != nil {
		return false
	}
	return user != nil
}
