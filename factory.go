package factorybot

import (
	"fmt"
	"reflect"
)

type Factory struct {
	Struct                any
	AfterBuildCallbacks   []func(m any)
	BeforeCreateCallbacks []func(m any)
	AfterCreateCallbacks  []func(m any)
	Fields                map[string]func(f *Factory) any
	Traits                map[string]func(m any)
	traitsToProcess       []string
	persistCallback       func(m any)
}

// Build creates a struct instance according to defined factory behavior without persisting it.
func (f *Factory) Build() any {
	m := f.build()
	return m.Addr().Interface()
}

// Create creates a struct instance according to defined factory behavior and persists it
// using callback set by `Persist()`.
func (f *Factory) Create() any {
	m := f.create()
	return m.Addr().Interface()
}

// BuildList creates a list of struct instances according to defined factory behavior without persisting them.
func (f *Factory) BuildList(count int) any {
	return f.processList(count, "build")
}

// CreateList creates a list of struct instances according to defined factory behavior and persists them
func (f *Factory) CreateList(count int) any {
	return f.processList(count, "create")
}

// Field lets developers define a transient field that doesn't belong to the reference struct.
func (f *Factory) Field(name string) any {
	if callback, ok := f.Fields[name]; ok {
		return callback(f)
	}

	return nil
}

// Set hints the factory how to set the field value.
func (f *Factory) Set(name string, callback func(f *Factory) any) *Factory {
	d := *f

	d.Fields = make(map[string]func(f *Factory) any)
	for k, v := range f.Fields {
		d.Fields[k] = v
	}
	d.Fields[name] = callback

	return &d
}

// SetT allows developers to perform one-off overrides while building/creating the strcut instance.
func (f *Factory) SetT(name string, value any) *Factory {
	d := *f

	d.Fields = make(map[string]func(f *Factory) any)
	for k, v := range f.Fields {
		d.Fields[k] = v
	}

	d.Fields[name] = func(f *Factory) any { return value }
	return &d
}

// WithTrait lets developers build/create a struct instance with a specific trait.
func (f *Factory) WithTrait(name string) *Factory {
	d := *f
	d.traitsToProcess = append(d.traitsToProcess, name)
	return &d
}

// Trait lets developers define a trait for the factory.
func (f *Factory) Trait(name string, callback func(m any)) *Factory {
	d := *f

	d.Traits = make(map[string]func(m any))
	for k, v := range f.Traits {
		d.Traits[k] = v
	}

	d.Traits[name] = callback
	return &d
}

// Persist sets the callback function to be called in order to save/persiste the struct instance.
func (f *Factory) Persist(callback func(m any)) *Factory {
	d := *f
	d.persistCallback = callback
	return &d
}

// AfterBuild sets the callback function to be called after the struct instance is built/initialized.
func (f *Factory) AfterBuild(callback func(m any)) *Factory {
	d := *f
	d.AfterBuildCallbacks = append(d.AfterBuildCallbacks, callback)
	return &d
}

// AfterCreate sets the callback function to be called after the struct instance is created/persisted.
func (f *Factory) AfterCreate(callback func(m any)) *Factory {
	d := *f
	d.AfterCreateCallbacks = append(d.AfterCreateCallbacks, callback)
	return &d
}

func (f *Factory) build() reflect.Value {
	// copy referenced struct
	m := reflect.New(reflect.ValueOf(f.Struct).Elem().Type()).Elem()

	for field, v := range f.Fields {
		if m.FieldByName(field).IsValid() {
			m.FieldByName(field).Set(reflect.ValueOf(v(f)))
		}
	}

	f.processTraits(m)
	f.processAfterBuild(m)

	return m
}

func (f *Factory) create() reflect.Value {
	m := f.build()

	if f.persistCallback != nil {
		f.persistCallback(m.Addr().Interface())
		defer f.processAfterCreate(m)
	}

	return m
}

func (f *Factory) processList(count int, action string) any {
	t := reflect.TypeOf(f.Struct)
	list := reflect.New(reflect.SliceOf(t)).Elem()

	var m reflect.Value

	for index := 0; index < count; index++ {
		if action == "create" {
			m = f.create()
		} else {
			m = f.build()
		}

		list = reflect.Append(list, m.Addr())
	}

	return list.Interface()
}

func (f *Factory) processTraits(m reflect.Value) {
	for _, traitName := range f.traitsToProcess {
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

func (f *Factory) processCallbacks(callbacks []func(m any), m reflect.Value) {
	if len(callbacks) > 0 {
		for _, callback := range callbacks {
			callback(m.Addr().Interface())
		}
	}
}

// NewFactory creates new factory instance using the provided struct pointer as a reference.
//
//	f := NewFactory(&User{}).
//		Set("Name", func(f *Factory) interface{} {
//			return namesSeq.One()
//		}).
//		Trait("admin", func(m interface{}) {
//			x := m.(*User)
//			x.Admin = true
//	    })
//
//	user := f.Build().(*User)
func NewFactory(s any) *Factory {
	return &Factory{
		Struct:               s,
		Traits:               make(map[string]func(m any)),
		Fields:               make(map[string]func(f *Factory) any),
		AfterBuildCallbacks:  make([]func(m any), 0),
		AfterCreateCallbacks: make([]func(m any), 0),
	}
}
