package errors

import (
	 "github.com/ShayanGsh/azar/core/utils"
)

var (
	ErrPolicyExists = utils.CustomError("Policy already exists")
	ErrPolicyNameExists = utils.CustomError("Policy name already exists")
	ErrPolicyNotFound = utils.CustomError("Policy not found")
)