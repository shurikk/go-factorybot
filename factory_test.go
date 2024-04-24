package factorybot_test

import (
	. "factorybot"

	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStruct struct {
	AfterBuildSet  bool
	AfterCreateSet bool
	Persisted      bool
	TraitInvoked   bool
	FieldToSet     bool
	FieldToSet2    bool
	Name           string
}

func TestFactory(t *testing.T) {
	seq := NewSequence(func(n int) any {
		return fmt.Sprintf("name%d", n)
	})

	f0 := NewFactory(&mockStruct{}).
		AfterBuild(func(m any) {
			x := m.(*mockStruct)
			x.AfterBuildSet = true
		}).
		AfterCreate(func(m any) {
			x := m.(*mockStruct)
			x.AfterCreateSet = true
		}).
		Persist(func(m interface{}) {
			x := m.(*mockStruct)
			x.Persisted = true
		}).
		Set("FieldToSet", func(f *Factory) interface{} {
			return true
		}).
		Set("FieldToSet2", func(f *Factory) interface{} {
			return false
		}).
		Set("Name", func(f *Factory) interface{} {
			switch f.Field("CustomizeName") {
			case true:
				return "other-name"
			default:
				return seq.One()
			}
		}).
		Trait("some trait", func(m interface{}) {
			x := m.(*mockStruct)
			x.TraitInvoked = true
		})

	f1 := f0.
		Set("FieldToSet2", func(f *Factory) interface{} {
			return true
		}).
		Trait("new trait", func(m interface{}) {
			x := m.(*mockStruct)
			x.TraitInvoked = true
		})

	t.Run("Field", func(t *testing.T) {
		m1 := f0.Build().(*mockStruct)
		m2 := f0.SetT("FieldToSet", false).Build().(*mockStruct)
		require.True(t, m1.FieldToSet)
		require.False(t, m2.FieldToSet)
	})

	t.Run("Build", func(t *testing.T) {
		m, ok := f0.Build().(*mockStruct)
		require.True(t, ok)
		require.True(t, m.AfterBuildSet)

		m1, _ := f0.Build().(*mockStruct)
		require.NotEqual(t, m.Name, m1.Name)
	})

	t.Run("Build/WithTrait", func(t *testing.T) {
		m := f0.WithTrait("some trait").Build().(*mockStruct)
		require.True(t, m.TraitInvoked)
	})

	t.Run("Build/WithTrait/Panic", func(t *testing.T) {
		require.Panics(t, func() { f0.WithTrait("missing trait").Build() })
	})

	t.Run("Build/SetField", func(t *testing.T) {
		m := f0.Build().(*mockStruct)
		require.True(t, m.FieldToSet)
	})

	t.Run("Build/SetField/Override", func(t *testing.T) {
		m := f0.SetT("FieldToSet", false).Build().(*mockStruct)
		require.False(t, m.FieldToSet)
	})

	t.Run("Build/SetField/Transient", func(t *testing.T) {
		m := f0.SetT("CustomizeName", true).Build().(*mockStruct)
		require.Equal(t, "other-name", m.Name)
	})

	t.Run("Build/CompositeFactory", func(t *testing.T) {
		m0 := f0.Build().(*mockStruct)
		require.False(t, m0.FieldToSet2)

		m1 := f1.Build().(*mockStruct)
		require.True(t, m1.FieldToSet2)
	})

	t.Run("Create", func(t *testing.T) {
		m := f0.WithTrait("some trait").Create().(*mockStruct)

		require.True(t, m.AfterBuildSet)
		require.True(t, m.AfterCreateSet)
		require.True(t, m.TraitInvoked)
		require.True(t, m.Persisted)
	})

	t.Run("BuildList", func(t *testing.T) {
		list, ok := f0.BuildList(5).([]*mockStruct)
		require.True(t, ok)
		require.Len(t, list, 5)
	})

	t.Run("CreateList", func(t *testing.T) {
		list, ok := f0.CreateList(5).([]*mockStruct)
		require.True(t, ok)
		require.Len(t, list, 5)
	})
}
