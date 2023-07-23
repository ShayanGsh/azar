package perms

type Role interface{
	IsAllowed(action string, resource string) bool
	GetRoleName() string
	AddPermissions(policy Policy)
	RemovePermissions(policy Policy)
	RemoveAllPermissions()
	GetPermissions() []Policy
	RemovePermissionByName(name string)
}

type Roleable interface{
	GetRole() Role
	SetRole(role Role)
	ChangeRole(role Role)
}

type RoleManager interface{
	GetRole(roleName string) Role
	GetRoles() []Role
	AssignRole(roleName string, roleable Roleable)
	UnassignRole(roleable Roleable)
	SetDefaultRole(role Role)
}

type RoleData struct{
	Name string
	Permissions Policy
}

func (rd *RoleData) IsAllowed(action string, resource string) bool{
	return rd.Permissions.IsAllowed(action, resource)
}

func (rd *RoleData) GetRoleName() string{
	return rd.Name
}

func (rd *RoleData) AddPermissions(policy Policy){
	rd.Permissions.AddPolicy(policy)
}

type RoleList struct{
	Roles map[string]RoleData
	DefaultRole Role
}

func (rl *RoleList) GetRole(roleName string) Role{
	if role, ok := rl.Roles[roleName]; ok {
		return role
	}
	return nil
}