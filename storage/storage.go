package storage

import (
	"blockpost/genprotos/authorization"
)

type StorageI interface {
	AddUser(id string, entity *authorization.CreateUserRequest) error
	GetUserByID(id string) (*authorization.User, error)
	GetUserList(offset, limit int, search string) (resp *authorization.GetUserListResponse, err error)
	UpdateUser(entity *authorization.UpdateUserRequest) error
	DeleteUser(id string) error
	GetUserByUsername(username string) (*authorization.User, error)
}
