package jnigi

import (
	"testing"
)

var obj *ObjectRef

func BenchmarkSimple(b *testing.B) {
	if obj == nil {
		nenv := jvm.AttachCurrentThread()
		var err error
		obj, err = nenv.NewObject("java/lang/Integer", 0)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var dummy int
		obj.CallMethod(env, "intValue", Int, &dummy)
	}
}
