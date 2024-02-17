package jnigi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testNum = 42
)

func PTestTypes(t *testing.T) {
	obj, err := env.NewObject("local/JnigiTestBase")
	if err != nil {
		t.Fatal(err)
	}

	{
		var val bool
		if err := obj.CallMethod(env, "boolToboolean", &val, bool(true)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, bool(true), val) {
			t.Fail()
		}
	}
	{
		var val byte
		if err := obj.CallMethod(env, "byteTobyte", &val, byte(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, byte(testNum), val) {
			t.Fail()
		}
	}
	{
		var val int16
		if err := obj.CallMethod(env, "int16Toshort", &val, int16(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, int16(testNum), val) {
			t.Fail()
		}
	}
	{
		var val uint16
		if err := obj.CallMethod(env, "uint16Tochar", &val, uint16(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, uint16(testNum), val) {
			t.Fail()
		}
	}
	{
		var val int
		if err := obj.CallMethod(env, "intToint", &val, int(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, int(testNum), val) {
			t.Fail()
		}
	}
	{
		var val int64
		if err := obj.CallMethod(env, "int64Tolong", &val, int64(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, int64(testNum), val) {
			t.Fail()
		}
	}
	{
		var val float32
		if err := obj.CallMethod(env, "float32Tofloat", &val, float32(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, float32(testNum), val) {
			t.Fail()
		}
	}
	{
		var val float64
		if err := obj.CallMethod(env, "float64Todouble", &val, float64(testNum)); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, float64(testNum), val) {
			t.Fail()
		}
	}

	{
		wanted := []bool{true, true, true}
		var got []bool
		if err := obj.CallMethod(env, "s_boolToboolean", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}

	{
		wanted := []byte{testNum, testNum, testNum}
		var got []byte
		if err := obj.CallMethod(env, "s_byteTobyte", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}
	{
		wanted := []int16{testNum, testNum, testNum}
		var got []int16
		if err := obj.CallMethod(env, "s_int16Toshort", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}
	{
		wanted := []uint16{testNum, testNum, testNum}
		var got []uint16
		if err := obj.CallMethod(env, "s_uint16Tochar", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}
	{
		wanted := []int{testNum, testNum, testNum}
		var got []int
		if err := obj.CallMethod(env, "s_intToint", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}
	{
		wanted := []int64{testNum, testNum, testNum}
		var got []int64
		if err := obj.CallMethod(env, "s_int64Tolong", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}
	{
		wanted := []float32{testNum, testNum, testNum}
		var got []float32
		if err := obj.CallMethod(env, "s_float32Tofloat", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}
	{
		wanted := []float64{testNum, testNum, testNum}
		var got []float64
		if err := obj.CallMethod(env, "s_float64Todouble", &got, wanted); err != nil {
			t.Fatal(err)
		}
		if !assert.Equal(t, wanted, got) {
			t.Fail()
		}
	}

}
