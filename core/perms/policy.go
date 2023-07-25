package perms

import "github.com/ShayanGsh/azar/core/errors"

type Policies interface{
	IsAllowed(action string, resource string) bool
	AddPolicy(name string, description string, action string, resource string) error
	AddPolicyByObject(policy Policy) error
	RemovePolicy(name string) error
	GetPolicy(name string) (PolicyData, error)
	IterPolicies() <- chan PolicyData
}

type Policy interface {
	IsAllowed(action string, resource string) bool
	GetPolicyName() string
	GetPolicyDescription() string
	GetPolicyAction() string
	GetPolicyResource() string
	SetPolicyName(name string)
	SetPolicyDescription(description string)
	SetPolicyAction(action string)
	SetPolicyResource(resource string)
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

func (pl *PolicyList) AddPolicyByObject(policy Policy) error{
	if _, ok := pl.Policies[policy.GetPolicyName()]; ok {
		return errors.ErrPolicyNameExists
	}

	for _, p := range pl.Policies {
		if p.Action == policy.GetPolicyAction() && p.Resource == policy.GetPolicyResource() {
			return errors.ErrPolicyExists
		}
	}

	pl.Policies[policy.GetPolicyName()] = PolicyData{
		Name: policy.GetPolicyName(),
		Description: policy.GetPolicyDescription(),
		Action: policy.GetPolicyAction(),
		Resource: policy.GetPolicyResource(),
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

func (pl *PolicyList) IterPolicies() <- chan PolicyData{
	ch := make(chan PolicyData)
	go func(){
		for _, policy := range pl.Policies {
			ch <- policy
		}
		close(ch)
	}();
	return ch
}

func (pd *PolicyData) IsAllowed(action string, resource string) bool{
	return pd.Action == action && pd.Resource == resource
}

func (pd *PolicyData) GetPolicyName() string{
	return pd.Name
}

func (pd *PolicyData) GetPolicyDescription() string{
	return pd.Description
}

func (pd *PolicyData) GetPolicyAction() string{
	return pd.Action
}

func (pd *PolicyData) GetPolicyResource() string{
	return pd.Resource
}

func (pd *PolicyData) SetPolicyName(name string){
	pd.Name = name
}

func (pd *PolicyData) SetPolicyDescription(description string){
	pd.Description = description
}

func (pd *PolicyData) SetPolicyAction(action string){
	pd.Action = action
}

func (pd *PolicyData) SetPolicyResource(resource string){
	pd.Resource = resource
}