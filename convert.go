package jnigi

// Types that can covert to Go values from object reference
type ToGoConverter interface {
	// Method should delete reference if it is not needed anymore
	ConvertToGo(obj *ObjectRef) error
}

// ArrayRef just disables auto conversion of Java arrays to Go slices
type ArrayRef struct {
	*ObjectRef
	Type
}

func NewArrayRef(t Type) *ArrayRef {
	a := &ArrayRef{nil, t}
	return a
}

// Just hold on to the reference to the array jobject
func (a *ArrayRef) ConvertToGo(obj *ObjectRef) error {
	a.ObjectRef = obj
	return nil
}

func (a *ArrayRef) GetType() Type {
	return a.Type
}

// Types that can convert to Java Object
type ToJavaConverter interface {
	// Returned reference will be deleted
	ConvertToJava() (*ObjectRef, error)
}

// Type used to represent arg that has been converted using ConvertToJava
type convertedArg struct {
	*ObjectRef
}

func replaceConvertedArgs(args []interface{}) (err error) {
	for i, arg := range args {
		if v, ok := arg.(ToJavaConverter); ok {
			var convArg convertedArg
			convArg.ObjectRef, err = v.ConvertToJava()
			if err != nil {
				break
			}
			args[i] = &convArg
		}
	}
	return err
}
