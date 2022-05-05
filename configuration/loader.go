package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type ConfigurationLoader interface {
	// Returns the value associated with key or nil
	Get(key string) interface{}
	// Returns the value associated with key or the default value
	GetOrDefault(key string, defaultValue interface{}) interface{}
	// Returns the value associated with key or panics
	Must(key string) interface{}
	// Writes the value found using key if - and only if - it either matches the
	// type of dest or if it is a string and can be unmarshelled to dest,
	// returns an error otherwise.
	// dest must be a pointer
	GetTypeSafe(key string, ptrDest interface{}) error
	// Same as GetTypeSafe, except it returns the default value if the key is
	// not present
	GetTypeSafeOrDefault(key string, ptrDest interface{}, defaultValue interface{}) error
}

// -------------------------------------------------------------------------- //

// MapConfigLoader is a config loader which will interprete keys with dots
// as key-chains within map structures.
// E.g. loader.get(my.fancy.key) is resolved to data[my][fancy][key]
type MapConfigLoader struct {
	data     map[string]interface{}
	filepath string
	sep      string
}

var _ ConfigurationLoader = (*MapConfigLoader)(nil)

func readFile(filepath string) ([]byte, error) {
	if file, err := os.Open(filepath); err != nil {
		return nil, err
	} else {
		return io.ReadAll(file)
	}
}

func NewMapConfigLoader(data map[string]interface{}, filepath string, keySeparator string) *MapConfigLoader {
	return &MapConfigLoader{
		data:     data,
		filepath: filepath,
		sep:      keySeparator,
	}
}

func LoadJsonConfiguration(filepath string) (*MapConfigLoader, error) {
	if buf, err := readFile(filepath); err != nil {
		return nil, err
	} else {
		data := make(map[string]interface{})
		if err := json.Unmarshal(buf, &data); err != nil {
			return nil, err
		}
		return NewMapConfigLoader(data, filepath, "."), nil
	}
}

func LoadYamlConfiguration(filepath string) (*MapConfigLoader, error) {
	if buf, err := readFile(filepath); err != nil {
		return nil, err
	} else {
		data := make(map[string]interface{})
		if err := yaml.Unmarshal(buf, &data); err != nil {
			return nil, err
		}
		return NewMapConfigLoader(data, filepath, "."), nil
	}
}

func LoadEnvConfiguration() (*MapConfigLoader, error) {
	data := make(map[string]interface{})
	for _, kv := range os.Environ() {
		keys, value, _ := strings.Cut(kv, "=")
		if err := loadKvRecursive(data, strings.Split(keys, "."), value, []string{}); err != nil {
			return nil, err
		}
	}
	return NewMapConfigLoader(data, "env", "."), nil
}
func loadKvRecursive(m map[string]interface{}, keys []string, value string, trail []string) error {
	if len(keys) > 1 {
		var _m map[string]interface{}
		if v, exists := m[keys[0]]; exists {
			var ok bool
			if _m, ok = v.(map[string]interface{}); !ok {
				return loadKvRecursiveError(trail, keys[0])
			}
		} else {
			_m = make(map[string]interface{})
			m[keys[0]] = _m
		}
		return loadKvRecursive(_m, keys[1:], value, append(trail, keys[0]))
	} else {
		if _, exists := m[keys[0]]; exists {
			return loadKvRecursiveError(trail, keys[0])
		}
		m[keys[0]] = value
		return nil
	}
}
func loadKvRecursiveError(trail []string, key string) error {
	_trail := append(trail, key)
	return errors.New(fmt.Sprintf("Conflicting entry in env variables for key %s",
		strings.Join(_trail, ".")))
}

func (jcl *MapConfigLoader) getTraverse(m map[string]interface{}, keys []string) (v interface{}, exists bool) {
	if v, ok := m[keys[0]]; ok {
		if len(keys) > 1 {
			if newMap, castOk := v.(map[string]interface{}); castOk {
				return jcl.getTraverse(newMap, keys[1:])
			} else {
				return nil, false
			}
		} else {
			return v, true
		}
	} else {
		return nil, false
	}
}

// Returns the value associated with key or nil
func (jcl *MapConfigLoader) Get(key string) interface{} {
	if value, exists := jcl.getTraverse(jcl.data, strings.Split(key, jcl.sep)); exists {
		return value
	} else {
		return nil
	}
}

// Returns the value associated with key or the default value
func (jcl *MapConfigLoader) GetOrDefault(key string, defaultValue interface{}) interface{} {
	if value, exists := jcl.getTraverse(jcl.data, strings.Split(key, jcl.sep)); exists {
		return value
	} else {
		return defaultValue
	}
}

// Returns the value associated with key or panics
func (jcl *MapConfigLoader) Must(key string) interface{} {
	if value, exists := jcl.getTraverse(jcl.data, strings.Split(key, jcl.sep)); exists {
		return value
	} else {
		panic("Key " + key + " not found by JsonConfigurationLoader in '" + jcl.filepath + "'")
	}
}

// Writes the value found using key if - and only if - it either matches the
// type of dest or if it is a string and can be unmarshelled to dest,
// returns an error otherwise.
func (jcl *MapConfigLoader) GetTypeSafe(key string, ptrDest interface{}) error {
	err, _ := jcl.getTypeSafeExists(key, ptrDest)
	return err
}
func (jcl *MapConfigLoader) getTypeSafeExists(key string, ptrDest interface{}) (error, bool) {
	var (
		retError  error = nil
		retExists bool  = false
		value     interface{}
	)
	if value, retExists = jcl.getTraverse(jcl.data, strings.Split(key, jcl.sep)); retExists {
		vVal := reflect.ValueOf(value)
		vRef := reflect.ValueOf(ptrDest).Elem()

		if vVal.Type().AssignableTo(vRef.Type()) {
			// if it's a perfect type match, simply copy
			vRef.Set(vVal)
		} else if vVal.CanConvert(vRef.Type()) {
			// if on the other hand it is convertible, convert
			vRef.Set(vVal.Convert(vRef.Type()))
		} else if str, castOk := value.(string); castOk {
			// try unmarshalling if it's a string
			if err := json.Unmarshal([]byte(str), ptrDest); err != nil {
				retError = NewError(fmt.Sprintf("Unable to unmarshal key(%s)='%v' to a field of type %s", key, str, vRef.Type().String()), err)
			}
		}
	} else {
		retError = errors.New("Key " + key + " not found by JsonConfigurationLoader in '" + jcl.filepath + "'")
		retExists = false
	}
	return retError, retExists
}

// Same as GetTypeSafe, except it returns the default value if the key is
// not present
func (jcl *MapConfigLoader) GetTypeSafeOrDefault(key string, ptrDest interface{}, defaultValue interface{}) error {
	if err, exists := jcl.getTypeSafeExists(key, ptrDest); !exists {
		reflect.ValueOf(ptrDest).Elem().Set(reflect.ValueOf(defaultValue))
		return nil
	} else {
		return err
	}
}
