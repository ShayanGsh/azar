package perms

import "github.com/ShayanGsh/azar/core/errors"

type Policy interface{
	IsAllowed(action string) bool
	AddPolicy(name string, description string, action string, resource string) error
	RemovePolicy(name string) error
	GetPolicy(name string) (PolicyData, error)
}

type PolicyData struct{
	Name string
	Description string
	Action string
	Resource string
}

type PolicyList struct{
	Policies map[string]PolicyData
}

func (pl *PolicyList) IsAllowed(action string, resource string) bool{
	for _, policy := range pl.Policies {
		if policy.Action == action && policy.Resource == resource {
			return true
		}
	}
	return false
}

func (pl *PolicyList) AddPolicy(name string, description string, action string, resource string) error{
	if _, ok := pl.Policies[name]; ok {
		return errors.ErrPolicyNameExists
	}

	for _, policy := range pl.Policies {
		if policy.Action == action && policy.Resource == resource {
			return errors.ErrPolicyExists
		}
	}

	pl.Policies[name] = PolicyData{
		Name: name,
		Description: description,
		Action: action,
		Resource: resource,
	}
	return nil
}

func (pl *PolicyList) RemovePolicy(name string) error{
	if _, ok := pl.Policies[name]; !ok {
		return errors.ErrPolicyNotFound
	}

	delete(pl.Policies, name)
	return nil
}

func (pl *PolicyList) GetPolicy(name string) (PolicyData, error){
	if _, ok := pl.Policies[name]; !ok {
		return PolicyData{}, errors.ErrPolicyNotFound
	}

	return pl.Policies[name], nil
}