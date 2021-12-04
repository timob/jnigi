package jnigi

// Types that can covert to Go values from object reference
type ToGoConverter interface {
	ConvertToGo(obj *ObjectRef) error
}


// ArrayRef just disables auto conversion of Java arrays to Go slices
type ArrayRef struct {
	*ObjectRef
}

// Just hold on to the reference to the array jobject
func (a *ArrayRef) ConvertToGo(obj *ObjectRef) error {
	a.ObjectRef = obj
	return nil
}