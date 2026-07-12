package policy

import "fmt"

// Validator — policy doğrulama sözleşmesi.
type Validator interface {
	Validate(p Policy) error
}

// DefaultValidator — temel policy doğrulama.
type DefaultValidator struct{}

func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

func (v *DefaultValidator) Validate(p Policy) error {
	for i, rule := range p.Rules {
		if rule.CIDR == nil {
			return fmt.Errorf("kural %d: CIDR boş olamaz", i)
		}
		for _, port := range rule.Ports {
			if port == 0 {
				return fmt.Errorf("kural %d: geçersiz port 0", i)
			}
		}
	}
	return nil
}
