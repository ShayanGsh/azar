package perms

import (
	"github.com/ShayanGsh/azar/core/errors"
	"github.com/ShayanGsh/azar/core/utils"
)

type Role interface {
	IsAllowed(action string, resource string) bool
	GetRoleName() string
	AddPermissions(policy Policy)
	RemovePermissions(policy Policy)
	RemoveAllPermissions()
	GetPermissions() Policies
	RemovePermissionByName(name string)
}

type Roleable interface {
	GetRole() Role
	SetRole(role Role)
	ChangeRole(role Role)
}

type RoleManager interface {
	GetRole(roleName string) Role
	GetRoles() []Role
	AssignRole(roleName string, roleable Roleable)
	UnassignRole(roleable Roleable)
	SetDefaultRole(role Role)
	SetDefaultRoleByName(roleName string)
	AddRole(role Role)
	AddBulkRoles(roles []Role)
	RemoveRole(roleName string)
	Clear()
}

type RoleData struct {
	Name        string
	Permissions Policies
}

func (rd *RoleData) IsAllowed(action string, resource string) bool {
	return rd.Permissions.IsAllowed(action, resource)
}

func (rd *RoleData) GetRoleName() string {
	return rd.Name
}

func (rd *RoleData) AddPermissions(policy Policy) {
	rd.Permissions.AddPolicyByObject(policy)
}

func (rd *RoleData) RemovePermissions(policy Policy) {
	rd.Permissions.RemovePolicy(policy.GetPolicyName())
}

func (rd *RoleData) RemoveAllPermissions() {
	for policy := range rd.Permissions.IterPolicies() {
		rd.Permissions.RemovePolicy(policy.Name)
	}
}

func (rd *RoleData) GetPermissions() Policies {
	return rd.Permissions
}

func (rd *RoleData) RemovePermissionByName(name string) {
	rd.Permissions.RemovePolicy(name)
}

type RoleMap struct {
	Roles       map[string]Role
	DefaultRole Role
}

func (rl *RoleMap) GetRole(roleName string) Role {
	if rl.RoleExists(roleName) {
		return rl.Roles[roleName]
	}
	return nil
}

func (rl *RoleMap) GetRoles() []Role {
	roles := make([]Role, 0, len(rl.Roles))
	for _, roleData := range rl.Roles {
		roles = append(roles, roleData)
	}
	return roles
}

func (rl *RoleMap) AssignRole(roleName string, roleable Roleable) {
	roleable.SetRole(rl.GetRole(roleName))
}

func (rl *RoleMap) UnassignRole(roleable Roleable) {
	roleable.SetRole(rl.DefaultRole)
}

func (rl *RoleMap) SetDefaultRole(role Role) {
	rl.DefaultRole = role
}

func (rl *RoleMap) SetDefaultRoleByName(roleName string) error {
	if role := rl.GetRole(roleName); role != nil {
		rl.SetDefaultRole(role)
		return nil
	} else {
		return utils.Error(errors.ErrRoleNotFound, roleName)
	}
}

func (rl *RoleMap) RoleExists(roleName string) bool {
	if _, ok := rl.Roles[roleName]; ok {
		return true
	}
	return false
}

func (rl *RoleMap) AddRole(role Role) {
	if !rl.RoleExists(role.GetRoleName()) {
		rl.Roles[role.GetRoleName()] = role
	}
}

func (rl *RoleMap) AddBulkRoles(roles []Role) {
	for _, role := range roles {
		rl.AddRole(role)
	}
}

func (rl *RoleMap) RemoveRole(roleName string) {
	if rl.RoleExists(roleName) {
		delete(rl.Roles, roleName)
	}
}

func (rl *RoleMap) Clear() {
	rl.Roles = make(map[string]Role)
}

func NewRoleMap() RoleMap {
	return RoleMap{
		Roles: make(map[string]Role),
	}
}

func InitRoleMap(roles []Role) RoleMap {
	roleMap := NewRoleMap()
	for _, role := range roles {
		roleMap.AddRole(role)
	}
	return roleMap
}
