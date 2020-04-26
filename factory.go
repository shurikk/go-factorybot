package factorybot

import (
	"fmt"
	"reflect"
)

// Factory ...
type Factory struct {
	Struct                interface{}
	AfterBuildCallbacks   []func(m interface{})
	BeforeCreateCallbacks []func(m interface{})
	AfterCreateCallbacks  []func(m interface{})
	Fields                map[string]func(f *Factory) interface{}
	Traits                map[string]func(m interface{})
	TraitsToProcess       []string
	PersistCallback       func(m interface{})
}

// Build ...
func (f *Factory) Build() interface{} {
	m := f.build()
	f.processTraits(m)
	return m.Interface()
}

// Create ...
func (f *Factory) Create() interface{} {
	m := f.create()
	f.processTraits(m)
	return m.Interface()
}

// BuildList ...
func (f *Factory) BuildList(count int, params ...interface{}) interface{} {
	return f.processList(count, "build")
}

// CreateList ...
func (f *Factory) CreateList(count int, params ...interface{}) interface{} {
	return f.processList(count, "create")
}

// Field returns value of a field
func (f *Factory) Field(name string) interface{} {
	if callback, ok := f.Fields[name]; ok {
		return callback(f)
	}

	return nil
}

// Set ...
func (f *Factory) Set(name string, callback func(f *Factory) interface{}) *Factory {
	d := *f

	d.Fields = make(map[string]func(f *Factory) interface{})
	for k, v := range f.Fields {
		d.Fields[k] = v
	}
	d.Fields[name] = callback

	return &d
}

// SetT ...
func (f *Factory) SetT(name string, value interface{}) *Factory {
	d := *f

	d.Fields = make(map[string]func(f *Factory) interface{})
	for k, v := range f.Fields {
		d.Fields[k] = v
	}

	d.Fields[name] = func(f *Factory) interface{} { return value }
	return &d
}

// WithTrait ...
func (f *Factory) WithTrait(name string) *Factory {
	d := *f
	d.TraitsToProcess = append(d.TraitsToProcess, name)
	return &d
}

// Trait ...
func (f *Factory) Trait(name string, callback func(m interface{})) *Factory {
	d := *f

	d.Traits = make(map[string]func(m interface{}))
	for k, v := range f.Traits {
		d.Traits[k] = v
	}

	d.Traits[name] = callback
	return &d
}

// Persist ...
func (f *Factory) Persist(callback func(m interface{})) *Factory {
	d := *f
	d.PersistCallback = callback
	return &d
}

// AfterBuild  ...
func (f *Factory) AfterBuild(callback func(m interface{})) *Factory {
	d := *f
	d.AfterBuildCallbacks = append(d.AfterBuildCallbacks, callback)
	return &d
}

// AfterCreate ...
func (f *Factory) AfterCreate(callback func(m interface{})) *Factory {
	d := *f
	d.AfterCreateCallbacks = append(d.AfterCreateCallbacks, callback)
	return &d
}

func (f *Factory) build() reflect.Value {
	m := reflect.New(reflect.TypeOf(f.Struct).Elem()).Elem()

	for field, v := range f.Fields {
		if m.FieldByName(field).IsValid() {
			m.FieldByName(field).Set(reflect.ValueOf(v(f)))
		}
	}

	defer f.processAfterBuild(m)

	return m
}

func (f *Factory) create() reflect.Value {
	m := f.build()
	defer f.processAfterCreate(m)

	if f.PersistCallback != nil {
		f.PersistCallback(m.Addr().Interface())
	}

	return m
}

func (f *Factory) processList(count int, action string) interface{} {
	t := reflect.TypeOf(f.Struct).Elem()
	list := reflect.New(reflect.SliceOf(t)).Elem()

	var m reflect.Value

	for index := 0; index < count; index++ {
		if action == "create" {
			m = f.create()
		} else {
			m = f.build()
		}

		f.processTraits(m)
		list = reflect.Append(list, m)
	}

	return list.Interface()
}

func (f *Factory) processTraits(m reflect.Value) {
	for _, traitName := range f.TraitsToProcess {
		if callback, ok := f.Traits[traitName]; ok {
			callback(m.Addr().Interface())
		} else {
			panic(fmt.Errorf("trait '%s' is not defined", traitName))
		}
	}
}

func (f *Factory) processAfterBuild(m reflect.Value) {
	f.processCallbacks(f.AfterBuildCallbacks, m)
}

func (f *Factory) processAfterCreate(m reflect.Value) {
	f.processCallbacks(f.AfterCreateCallbacks, m)
}

func (f *Factory) processCallbacks(callbacks []func(m interface{}), m reflect.Value) {
	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			callback(m.Addr().Interface())
		}
	}
}

// NewFactory creates new factory
func NewFactory(s interface{}) *Factory {
	return &Factory{
		Struct:               s,
		Traits:               make(map[string]func(m interface{})),
		Fields:               make(map[string]func(f *Factory) interface{}),
		AfterBuildCallbacks:  make([]func(m interface{}), 0),
		AfterCreateCallbacks: make([]func(m interface{}), 0),
	}
}
