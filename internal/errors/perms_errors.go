package errors



var (
	ErrPolicyExists = ErrorMsg("Policy already exists")
	ErrPolicyNameExists = ErrorMsg("Policy name already exists")
	ErrPolicyNotFound = ErrorMsg("Policy not found")
	ErrRoleNotFound = ErrorMsg("Role not found")
)