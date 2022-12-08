package authorization

import (
	"blockpost/config"
	"blockpost/genprotos/authorization"
	"blockpost/storage"
	"blockpost/util"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// authorizationService is a struct that implements the server interface
type authorizationService struct {
	stg storage.StorageI
	cfg config.Config
	authorization.UnimplementedAuthServiceServer
}

//NewAuthService...
func NewAuthService(cfg config.Config, stg storage.StorageI) *authorizationService {
	return &authorizationService{
		cfg: cfg,
		stg: stg,
	}
}

// CreateUser ...
func (u *authorizationService) CreateUser(c context.Context, req *authorization.CreateUserRequest) (*authorization.User, error) {
	id := uuid.New()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword: %s", err.Error())
	}

	req.Password = hashedPassword

	err = u.stg.AddUser(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.AddUser: %s", err.Error())
	}
	user, err := u.stg.GetUserByID(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "u.stg.GetUserByID: %s", err.Error())
	}
	return user, nil
}

// GetUserByID ...
func (u *authorizationService) GetUserByID(c context.Context, req *authorization.GetUserByIDRequest) (*authorization.User, error) {

	res, err := u.stg.GetUserByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.GetUserByID: %s", err.Error())
	}
	return res, nil
}

// GetUserByUsername ...
func (u *authorizationService) GetUserByUsername(c context.Context, req *authorization.User) (*authorization.User, error) {

	res, err := u.stg.GetUserByUsername(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.GetUserByUsername: %s", err.Error())
	}
	return res, nil
}

// GetUserList ...
func (u *authorizationService) GetUserList(c context.Context, req *authorization.GetUserListRequest) (*authorization.GetUserListResponse, error) {
	res, err := u.stg.GetUserList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.GetUserList: %s", err.Error())
	}
	return res, nil
}

// UpdateUser ...
func (u *authorizationService) UpdateUser(c context.Context, req *authorization.UpdateUserRequest) (*authorization.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword: %s", err.Error())
	}

	req.Password = hashedPassword
	err = u.stg.UpdateUser(req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.UpdateUser: %s", err.Error())
	}
	res, e := u.stg.GetUserByID(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.GetUserByID: %s", e.Error())
	}
	return res, nil
}

// DeleteUser ...
func (u *authorizationService) DeleteUser(c context.Context, req *authorization.DeleteUserRequest) (*authorization.User, error) {
	res, e := u.stg.GetUserByID(req.Id)
	if e != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.GetUserByID: %s", e.Error())
	}
	err := u.stg.DeleteUser(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "u.stg.DeleteUser: %s", err.Error())
	}

	return res, nil
}
