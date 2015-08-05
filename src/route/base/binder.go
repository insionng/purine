package base

import "github.com/tango-contrib/binding"

// base binder struct
type BaseBinder struct {
	binding.Binder
}

// bind data and return first error
func (b *BaseBinder) Bind(f interface{}) error {
	errors := b.Binder.Bind(f)
	if errors.Len() > 0 {
		return errors[0]
	}
	return nil
}
