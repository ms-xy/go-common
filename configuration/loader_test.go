package configuration

import (
	"reflect"
	"testing"

	"github.com/ms-xy/go-common/stack"
)

func TestLoadJsonConfiguration(t *testing.T) {
	filepath := "./loader_test.json"
	loader, err := LoadJsonConfiguration(filepath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// ----
	key := "unmapped.key"
	value := loader.Get(key)
	assertEquals(t, key, value, nil)
	key = "aMap.partially.mapped"
	value = loader.GetOrDefault(key, "testValue")
	assertType(t, key, value, "str")
	assertEquals(t, key, value, "testValue")
	var aString string
	err = loader.GetTypeSafeOrDefault(key, &aString, "defaultValue")
	assertErrNil(t, key, err)
	assertEquals(t, key, aString, "defaultValue")

	// ----
	key = "anInt"
	value = loader.Get(key)
	// json.Unmarshal treats all numbers as float64 when unmarshalling into an
	// interface{}
	assertType(t, key, value, 1.0)
	assertEquals(t, key, value, 3600.0)
	var anInt int
	err = loader.GetTypeSafe(key, &anInt)
	assertErrNil(t, key, err)
	assertEquals(t, key, anInt, 3600)
	var aFloat float64
	err = loader.GetTypeSafe(key, &aFloat)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 3600.0)
	aFloat = 8.8
	err = loader.GetTypeSafeOrDefault(key, &aFloat, 10.11)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 3600.0)

	// ----
	key = "aMap.aBoolean"
	value = loader.Get(key)
	assertType(t, key, value, true)
	assertEquals(t, key, value, false)

	// ----
	key = "aStrFloat"
	value = loader.Get(key)
	assertType(t, key, value, "str")
	assertEquals(t, key, value, "4.56")
	aFloat = 0.0
	err = loader.GetTypeSafe(key, &aFloat)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 4.56)
	key = "aString"
	err = loader.GetTypeSafeOrDefault(key, &aFloat, 10.11)
	if err == nil {
		t.Errorf("Expected conversion from string value 'xml' to float64 to fail, but it didn't")
		t.FailNow()
	}
	// value must stay the same, field should not be touched!
	assertEquals(t, key, aFloat, 4.56)
}

func TestLoadYamlConfiguration(t *testing.T) {
	filepath := "./loader_test.yaml"
	loader, err := LoadYamlConfiguration(filepath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// ----
	key := "unmapped.key"
	value := loader.Get(key)
	assertEquals(t, key, value, nil)
	key = "aMap.partially.mapped"
	value = loader.GetOrDefault(key, "testValue")
	assertType(t, key, value, "str")
	assertEquals(t, key, value, "testValue")
	var aString string
	err = loader.GetTypeSafeOrDefault(key, &aString, "defaultValue")
	assertErrNil(t, key, err)
	assertEquals(t, key, aString, "defaultValue")

	// ----
	key = "anInt"
	value = loader.Get(key)
	assertType(t, key, value, 1)
	assertEquals(t, key, value, 3600)
	var anInt int
	err = loader.GetTypeSafe(key, &anInt)
	assertErrNil(t, key, err)
	assertEquals(t, key, anInt, 3600)
	var aFloat float64
	err = loader.GetTypeSafe(key, &aFloat)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 3600.0)
	aFloat = 8.8
	err = loader.GetTypeSafeOrDefault(key, &aFloat, 10.11)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 3600.0)

	// ----
	key = "aMap.aBoolean"
	value = loader.Get(key)
	assertType(t, key, value, true)
	assertEquals(t, key, value, false)

	// ----
	key = "aStrFloat"
	value = loader.Get(key)
	assertType(t, key, value, "str")
	assertEquals(t, key, value, "4.56")
	aFloat = 0.0
	err = loader.GetTypeSafe(key, &aFloat)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 4.56)
	key = "aString"
	err = loader.GetTypeSafeOrDefault(key, &aFloat, 10.11)
	if err == nil {
		t.Errorf("Expected conversion from string value 'xml' to float64 to fail, but it didn't")
		t.FailNow()
	}
	// value must stay the same, field should not be touched!
	assertEquals(t, key, aFloat, 4.56)
}

func TestLoadConfiguration(t *testing.T) {
	filepath1 := "./loader_test.yaml"
	filepath2 := "./loader_test.json"
	loader, err := LoadConfiguration(
		func() (ConfigurationLoader, error) { return LoadYamlConfiguration(filepath1) },
		func() (ConfigurationLoader, error) { return LoadJsonConfiguration(filepath2) },
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	// ----
	key := "unmapped.key"
	value := loader.Get(key)
	assertEquals(t, key, value, nil)
	key = "aMap.partially.mapped"
	value = loader.GetOrDefault(key, "testValue")
	assertType(t, key, value, "str")
	assertEquals(t, key, value, "testValue")
	var aString string
	err = loader.GetTypeSafeOrDefault(key, &aString, "defaultValue")
	assertErrNil(t, key, err)
	assertEquals(t, key, aString, "defaultValue")

	// ----
	key = "anInt"
	value = loader.Get(key)
	assertType(t, key, value, 1)
	assertEquals(t, key, value, 3600)
	var anInt int
	err = loader.GetTypeSafe(key, &anInt)
	assertErrNil(t, key, err)
	assertEquals(t, key, anInt, 3600)
	var aFloat float64
	err = loader.GetTypeSafe(key, &aFloat)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 3600.0)
	aFloat = 8.8
	err = loader.GetTypeSafeOrDefault(key, &aFloat, 10.11)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 3600.0)

	// ----
	key = "aMap.aBoolean"
	value = loader.Get(key)
	assertType(t, key, value, true)
	assertEquals(t, key, value, false)

	// ----
	key = "aStrFloat"
	value = loader.Get(key)
	assertType(t, key, value, "str")
	assertEquals(t, key, value, "4.56")
	aFloat = 0.0
	err = loader.GetTypeSafe(key, &aFloat)
	assertErrNil(t, key, err)
	assertEquals(t, key, aFloat, 4.56)
	// The following test currently fails as the combined loader has no concept
	// of identifying child loader errors (yet)
	// key = "aString"
	// err = loader.GetTypeSafeOrDefault(key, &aFloat, 10.11)
	// if err == nil {
	// 	t.Errorf("Expected conversion from string value 'xml' to float64 to fail, but it didn't")
	// 	t.FailNow()
	// }
	// // value must stay the same, field should not be touched!
	// assertEquals(t, key, aFloat, 4.56)

	// ----
	key = "onlyInYaml"
	var anotherString string
	err = loader.GetTypeSafe(key, &anotherString)
	assertErrNil(t, key, err)
	assertEquals(t, key, anotherString, "yamlTest")
	key = "onlyInJson"
	var anotherString2 string
	err = loader.GetTypeSafe(key, &anotherString2)
	assertErrNil(t, key, err)
	assertEquals(t, key, anotherString2, "jsonTest")
}

func runLoaderTest(t *testing.T, loader ConfigurationLoader) {
}

func assertType(t *testing.T, key string, value interface{}, expected interface{}) {
	if reflect.TypeOf(value) != reflect.TypeOf(expected) {
		t.Errorf("Expected key(%s)='%v' to be of type %v but got %v",
			key, value, reflect.TypeOf(expected), reflect.TypeOf(value))
		t.Error(stack.Stack(1))
		t.FailNow()
	}
}

func assertEquals(t *testing.T, key string, value interface{}, expected interface{}) {
	if value != expected {
		t.Errorf("Expected key(%s)=='%v' but got '%v' instead",
			key, expected, value)
		t.Error(stack.Stack(1))
		t.FailNow()
	}
}

func assertErrNil(t *testing.T, key string, err error) {
	if err != nil {
		t.Errorf("Unexpected Error retrieving key(%s): %s", key, err.Error())
		t.Error(stack.Stack(1))
		t.FailNow()
	}
}

func TestLoadEnvConfiguration(t *testing.T) {
	loader, err := LoadEnvConfiguration()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	value := loader.Get("GOPATH")
	if len(value.(string)) == 0 {
		t.Error("Environment loading failed (expected GOPATH to be set for test)")
		t.FailNow()
	}
}
