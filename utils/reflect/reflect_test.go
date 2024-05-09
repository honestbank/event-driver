package reflect_test

import (
	"testing"

	"github.com/lukecold/event-driver/utils/reflect"
	"github.com/stretchr/testify/assert"
)

type testInterface interface{}

type testStruct struct{}

func makeTestStruct() testInterface {
	return testStruct{}
}

func TestGetType(t *testing.T) {
	t.Run("primitive", func(t *testing.T) {
		integer := 1
		integerType := reflect.GetType(integer)
		assert.Equal(t, "int", integerType)

		stringType := reflect.GetType("test string")
		assert.Equal(t, "string", stringType)
	})

	t.Run("primitive pointer", func(t *testing.T) {
		str := "test string"
		stringType := reflect.GetType(&str)
		assert.Equal(t, "*string", stringType)
	})

	t.Run("struct from interface", func(t *testing.T) {
		objectType := reflect.GetType(makeTestStruct())
		assert.Equal(t, "testStruct", objectType)
	})

	t.Run("struct pointer", func(t *testing.T) {
		objectType := reflect.GetType(&testStruct{})
		assert.Equal(t, "*testStruct", objectType)
	})

	t.Run("struct pointer of pointer", func(t *testing.T) {
		objectPtr := &testStruct{}
		objectType := reflect.GetType(&objectPtr)
		assert.Equal(t, "**testStruct", objectType)
	})
}
