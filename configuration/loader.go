package configuration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/ms-xy/go-common/log"
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
	// Print information about all contained config values
	DumpConfig()
}

// -------------------------------------------------------------------------- //

type CombinedLoader struct {
	loaders    []ConfigurationLoader
	loaderInfo string
}

var _ ConfigurationLoader = (*CombinedLoader)(nil)

// NewCombinedLoader supplies the means of loading a variadic set of settings
//
// Example:
//
//		NewCombinedLoader().MustLoadJSON(filepath).LoadEnv().CanLoadYaml(filepath)
func NewCombinedLoader() *CombinedLoader {
	cl := new(CombinedLoader)
	cl.loaders = make([]ConfigurationLoader, 0)
	cl.loaderInfo = ""
	return cl
}

func (cl *CombinedLoader) DumpConfig() {
	for _, l := range cl.loaders {
		l.DumpConfig()
	}
}

func (cl *CombinedLoader) updateLoaderInfo(loader ConfigurationLoader) {
	if cl.loaderInfo != "" {
		cl.loaderInfo += ", "
	}
	cl.loaderInfo = fmt.Sprintf("%s%s", cl.loaderInfo, loader)
}

// LoadEnv attempts to load env variables as configuration for a MapConfigLoader
//
// If the given variables create ambiguous config paths, an error will be
// returned instead (e.g. my.path=10 my.path.subpath=12 would conflict for the
// my.path part)
func (cl *CombinedLoader) LoadEnv() error {
	if mcl, err := LoadEnvConfiguration(); err != nil {
		return err
	} else {
		cl.loaders = append(cl.loaders, mcl)
		cl.updateLoaderInfo(mcl)
		return nil
	}
}

// MustLoadEnv attempts to load the present environment variables as settings
// and panics if an error occurs, else the loader is returned
func (cl *CombinedLoader) MustLoadEnv() *CombinedLoader {
	if err := cl.LoadEnv(); err != nil {
		panic(err)
	}
	return cl
}

// CanLoadEnv attempts to load the present environment variables as settings,
// logs any occuring error and returns the loader
func (cl *CombinedLoader) CanLoadEnv() *CombinedLoader {
	if err := cl.LoadEnv(); err != nil {
		log.Warn("error loading env settings", err)
	}
	return cl
}

// LoadJSON attempts to load the given filepath as a JSON file and create a
// MapConfigLoader based on it
func (cl *CombinedLoader) LoadJSON(filepath string) error {
	if mcl, err := LoadJsonConfiguration(filepath); err != nil {
		return err
	} else {
		cl.loaders = append(cl.loaders, mcl)
		cl.updateLoaderInfo(mcl)
		return nil
	}
}

// MustLoadJSON uses LoadJSON under the hood, but panics if an error is
// returned, otherwise it returns the loader again for call chaining
func (cl *CombinedLoader) MustLoadJSON(filepath string) *CombinedLoader {
	if err := cl.LoadJSON(filepath); err != nil {
		panic(err)
	}
	return cl
}

// CanLoadJSON uses LoadJSON to attempt to load the settings
//
// If this fails, it will log the error as a warning and return the loader
// for further chaining
func (cl *CombinedLoader) CanLoadJSON(filepath string) *CombinedLoader {
	if err := cl.LoadJSON(filepath); err != nil {
		log.WithField("filepath", filepath).Debug("no such settings file", err)
	}
	return cl
}

// LoadYaml attempts to load the given filepath as a YAML config file and
// returns any errors encountered
func (cl *CombinedLoader) LoadYaml(filepath string) error {
	if mcl, err := LoadYamlConfiguration(filepath); err != nil {
		return err
	} else {
		cl.loaders = append(cl.loaders, mcl)
		cl.updateLoaderInfo(mcl)
		return nil
	}
}

// MustLoadYaml attempts to load the given filepath and panics if it fails
func (cl *CombinedLoader) MustLoadYaml(filepath string) *CombinedLoader {
	if err := cl.LoadYaml(filepath); err != nil {
		panic(err)
	}
	return cl
}

// CanLoadYaml attempts to load the given filepath
//
// If loading fails it will log the error as a warning and return the loader
func (cl *CombinedLoader) CanLoadYaml(filepath string) *CombinedLoader {
	if err := cl.LoadYaml(filepath); err != nil {
		log.WithField("filepath", filepath).Debug("no such settings file", err)
	}
	return cl
}

// Returns the value associated with key or nil
func (cl *CombinedLoader) Get(key string) interface{} {
	for _, loader := range cl.loaders {
		if v := loader.Get(key); v != nil {
			return v
		}
	}
	return nil
}

// Returns the value associated with key or the default value
func (cl *CombinedLoader) GetOrDefault(key string, defaultValue interface{}) interface{} {
	if v := cl.Get(key); v != nil {
		return v
	}
	return defaultValue
}

// Returns the value associated with key or panics
func (cl *CombinedLoader) Must(key string) interface{} {
	if v := cl.Get(key); v != nil {
		return v
	}
	panic("Key " + key + " not found in " + cl.loaderInfo)
}

// Writes the value found using key if - and only if - it either matches the
// type of dest or if it is a string and can be unmarshelled to dest,
// returns an error otherwise.
// dest must be a pointer
func (cl *CombinedLoader) GetTypeSafe(key string, ptrDest interface{}) error {
	_err := make([]error, 0)
	for _, loader := range cl.loaders {
		if err := loader.GetTypeSafe(key, ptrDest); err == nil {
			return nil
		} else {
			_err = append(_err, err)
		}
	}
	_errors := make([]string, len(_err))
	for i, err := range _err {
		_errors[i] = err.Error()
	}
	return errors.New("Unable to resolve key " + key + ":\n" + strings.Join(_errors, "\n"))
}

// Same as GetTypeSafe, except it returns the default value if the key is
// not present
func (cl *CombinedLoader) GetTypeSafeOrDefault(key string, ptrDest interface{}, defaultValue interface{}) error {
	if err := cl.GetTypeSafe(key, ptrDest); err == nil {
		return nil
	}
	reflect.ValueOf(ptrDest).Elem().Set(reflect.ValueOf(defaultValue))
	return nil
}

// -------------------------------------------------------------------------- //

// MapConfigLoader is a config loader which will interprete keys with dots
// as key-chains within map structures.
// E.g. loader.get(my.fancy.key) is resolved to data[my][fancy][key]
type MapConfigLoader struct {
	data       map[string]interface{}
	filepath   string
	sep        string
	configType string
}

var _ ConfigurationLoader = (*MapConfigLoader)(nil)

func (mcl *MapConfigLoader) DumpConfig() {
	log.WithFields(mcl.data).Infof("configuration values obtained from '%s':", mcl.filepath)
}

func readFile(filepath string) ([]byte, error) {
	if file, err := os.Open(filepath); err != nil {
		return nil, err
	} else {
		return io.ReadAll(file)
	}
}

func NewMapConfigLoader(data map[string]interface{}, filepath, configType, keySeparator string) *MapConfigLoader {
	return &MapConfigLoader{
		data:       data,
		filepath:   filepath,
		sep:        keySeparator,
		configType: configType,
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
		return NewMapConfigLoader(data, filepath, "json", "."), nil
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
		return NewMapConfigLoader(data, filepath, "yaml", "."), nil
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
	return NewMapConfigLoader(data, "environment variables", "env", "."), nil
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
	return fmt.Errorf("conflicting entry in env variables for key %s",
		strings.Join(_trail, "."))
}

func (jcl *MapConfigLoader) String() string {
	str := "MapConfigLoader"
	if jcl.configType != "" {
		str += "<" + jcl.configType + ">"
	}
	if jcl.filepath != "" {
		str += "[" + jcl.filepath + "]"
	}
	return str
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
	if v := jcl.Get(key); v != nil {
		return v
	}
	return defaultValue
}

// Returns the value associated with key or panics
func (jcl *MapConfigLoader) Must(key string) interface{} {
	if v := jcl.Get(key); v != nil {
		return v
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
