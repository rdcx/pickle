package pickle

import "fmt"

type Validator struct {
	Errors []string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: []string{},
	}
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *Validator) PrintErrors() {
	for _, err := range v.Errors {
		fmt.Println(err)
	}
}

func (v *Validator) AddError(err string) {
	v.Errors = append(v.Errors, err)
}
