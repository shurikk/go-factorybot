# Factory Bot

[![Build Status](https://travis-ci.org/shurikk/go-factorybot.svg)](https://travis-ci.org/shurikk/go-factorybot)

Basic factory bot to generate complext stucts in tests

## Sequence

```go
IDSeq := factorybot.NewSequence()

NameSeq := factorybot.NewSequence(func(n int) interface{} {
	return fmt.Sprintf("name%d", n)
})

FancyNameSeq := factorybot.NewSequence(func(n int) interface{} {
	return fmt.Sprintf("fancy-name%d", n)
})
```

## Factory

```go
type Thing struct {
	ID   int    `json:"id"`
	Name string `json:"email"`
}

ThingsFactory := factorybot.NewFactory(&Thing{}).
	Persist(func(m interface{}) {
		// your favorite persistence layer to store m.(*Thing)
		thing := m.(*Thing)
		thing.ID = IDSeq.N()
	}).
	Set("Name", func(f *factorybot.Factory) interface{} {
		return NameSeq.One()
	}).
	AfterBuild(func(m interface{}) {
		// do something with m.(*Thing)
	}).
	AfterCreate(func(m interface{}) {
		// do something with m.(*Thing)
	}).
	Trait("something special", func(m interface{}) {
		// do something with m.(*Thing)
	}).
	Trait("something else", func(m interface{}) {
		// do something with m.(*Thing)
	})
```

## Composite Factory

```go
FancyThingsFactory := ThingsFactory.
	Set("Name", func(f *factorybot.Factory) interface{} {
		return FancyNameSeq.One()
	})
```

## Example

```go
// unique thing with all fields set
oneThing := ThingsFactory.Build()

// with a preset field and a trait
oneThing = ThingsFactory.
	Set("Name", "some name").
	WithTrait("something else").
	Build()

// persist: same as build but persisted
oneThing = ThingsFactory.Create()

// one thing with a trait
oneThing = ThingsFactory.WithTrait("something special").Build()

// one thing with multiple traits
oneThing = ThingsFactory.
	WithTrait("something special").
	WithTrait("something else").
	Build()

// persist
oneThing = ThingsFactory.
	WithTrait("something special").
	Create()

// build multiple things with a trait
NameSeq.Rewind()
things := ThingsFactory.
	WithTrait("something else").
	BuildList(2)

// create multiple things
NameSeq.Rewind()
things = ThingsFactory.CreateList(2)
```
