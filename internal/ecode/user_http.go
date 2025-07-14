package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// user business-level http error codes.
// the userNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	userNO       = 23
	userName     = "user"
	userBaseCode = errcode.HCode(userNO)

	ErrCreateUser     = errcode.NewError(userBaseCode+1, "failed to create "+userName)
	ErrDeleteByIDUser = errcode.NewError(userBaseCode+2, "failed to delete "+userName)
	ErrUpdateByIDUser = errcode.NewError(userBaseCode+3, "failed to update "+userName)
	ErrGetByIDUser    = errcode.NewError(userBaseCode+4, "failed to get "+userName+" details")
	ErrListUser       = errcode.NewError(userBaseCode+5, "failed to list of "+userName)

	// error codes are globally unique, adding 1 to the previous error code
)
