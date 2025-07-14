package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUserRequest request params
type CreateUserRequest struct {
	Name string `json:"name" binding:""`
}

// UpdateUserByIDRequest request params
type UpdateUserByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name string `json:"name" binding:""`
}

// UserObjDetail detail
type UserObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

// CreateUserReply only for api docs
type CreateUserReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteUserByIDReply only for api docs
type DeleteUserByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateUserByIDReply only for api docs
type UpdateUserByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetUserByIDReply only for api docs
type GetUserByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		User UserObjDetail `json:"user"`
	} `json:"data"` // return data
}

// ListUsersRequest request params
type ListUsersRequest struct {
	query.Params
}

// ListUsersReply only for api docs
type ListUsersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users []UserObjDetail `json:"users"`
	} `json:"data"` // return data
}
