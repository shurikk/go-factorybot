package factorybot_test

import (
	. "factorybot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockStruct struct {
	AfterBuildSet  bool
	AfterCreateSet bool
	Persisted      bool
	TraitInvoked   bool
	FieldToSet     bool
	FieldToSet2    bool
	Name           string
}

var _ = Describe("Factory", func() {
	mockFactory := NewFactory(&MockStruct{}).
		AfterBuild(func(m interface{}) {
			x := m.(*MockStruct)
			x.AfterBuildSet = true
		}).
		AfterCreate(func(m interface{}) {
			x := m.(*MockStruct)
			x.AfterCreateSet = true
		}).
		Persist(func(m interface{}) {
			x := m.(*MockStruct)
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
				return "name"
			}
		}).
		Trait("some trait", func(m interface{}) {
			x := m.(*MockStruct)
			x.TraitInvoked = true
		})

	mockFactory1 := mockFactory.
		Set("FieldToSet2", func(f *Factory) interface{} {
			return true
		}).
		Trait("new trait", func(m interface{}) {
			x := m.(*MockStruct)
			x.TraitInvoked = true
		})

	Describe("#Field", func() {
		It("returns value of a field", func() {
			f := mockFactory.SetT("SomeField", true)
			Expect(f.Field("SomeField")).To(BeTrue())
		})
	})

	Describe("#Build", func() {
		It("invokes AfterBuild callback", func() {
			m := mockFactory.Build()
			_ = mockFactory1.Build()
			Expect(m.(MockStruct).AfterBuildSet).To(BeTrue())
		})

		It("invokes Trait callback", func() {
			m := mockFactory.WithTrait("some trait").Build()
			Expect(m.(MockStruct).TraitInvoked).To(BeTrue())
		})

		It("panices when invoked Trait is not defined", func() {
			Expect(func() { mockFactory.WithTrait("missing trait").Build() }).To(Panic())
		})

		It("sets fields as defined in the factory", func() {
			m := mockFactory.Build()
			Expect(m.(MockStruct).FieldToSet).To(BeTrue())
		})

		It("accepts alternative value", func() {
			m := mockFactory.SetT("FieldToSet", false).Build()
			Expect(m.(MockStruct).FieldToSet).To(BeFalse())
		})

		It("works with transient fields", func() {
			m := mockFactory.SetT("CustomizeName", true).Build()
			Expect(m.(MockStruct).Name).To(Equal("other-name"))
		})

		It("works in composite factories", func() {
			m := mockFactory.Build()
			Expect(m.(MockStruct).FieldToSet2).To(BeFalse())

			m1 := mockFactory1.Build()
			Expect(m1.(MockStruct).FieldToSet2).To(BeTrue())
		})
	})

	Describe("#Create", func() {
		It("invokes AfterCreate callback", func() {
			m := mockFactory.WithTrait("some trait").Create()
			Expect(m.(MockStruct).AfterBuildSet).To(BeTrue())
			Expect(m.(MockStruct).AfterCreateSet).To(BeTrue())
			Expect(m.(MockStruct).TraitInvoked).To(BeTrue())
			Expect(m.(MockStruct).Persisted).To(BeTrue())
		})
	})

	Describe("#BuildList", func() {
		It("generates multiple structs", func() {
			ms := mockFactory.BuildList(5)
			Expect(len(ms.([]MockStruct))).To(Equal(5))
		})
	})

	Describe("#CreateList", func() {
		It("generates multiple structs", func() {
			ms := mockFactory.CreateList(5)
			Expect(len(ms.([]MockStruct))).To(Equal(5))
		})
	})
})
