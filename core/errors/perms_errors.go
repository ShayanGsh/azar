package errors

import (
	 "github.com/ShayanGsh/azar/core/utils"
)

var (
	ErrPolicyExists = utils.ErrorMsg("Policy already exists")
	ErrPolicyNameExists = utils.ErrorMsg("Policy name already exists")
	ErrPolicyNotFound = utils.ErrorMsg("Policy not found")
	ErrRoleNotFound = utils.ErrorMsg("Role not found")
)