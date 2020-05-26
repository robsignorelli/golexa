package golexa_test

import (
	"testing"

	"github.com/robsignorelli/golexa"
	"github.com/stretchr/testify/suite"
)

func TestSlotsSuite(t *testing.T) {
	suite.Run(t, new(SlotsSuite))
}

type SlotsSuite struct {
	suite.Suite
}

func (suite SlotsSuite) TestClone() {
	slots := golexa.Slots{}
	clone := slots.Clone()
	suite.Len(clone, 0,
		"Cloning an empty Slots should give another empty Slots")

	slots = golexa.Slots{
		"name": {Name: "name", Value: "Foo"},
	}
	clone = slots.Clone()
	suite.Len(clone, 1,
		"Cloning an Slots should give a new instance w/ the same length")
	suite.Equal("Foo", clone["name"].Value,
		"Cloned Slots should have the same values as the original")
	suite.Equal("", clone["askdlfjjalskdjf"].Value,
		"Cloned Slots should have the same values as the original")

	slots = golexa.NewSlots(
		golexa.NewSlot("name", "Foo"),
		golexa.NewSlot("age", "99"),
		golexa.NewSlot("blank", ""),
		golexa.NewResolvedSlot("beverage", "pop", "soda"),
	)
	clone = slots.Clone()
	suite.Len(clone, 4,
		"Cloning an Slots should give a new instance w/ the same length")
	suite.Equal("Foo", clone["name"].Value,
		"Cloned Slots should have the same values as the original")
	suite.Equal("99", clone["age"].Value,
		"Cloned Slots should have the same values as the original")
	suite.Equal("", clone["blank"].Value,
		"Cloned Slots should have the same values as the original")
	suite.Equal("", clone["askdlfjjalskdjf"].Value,
		"Cloned Slots should have the same values as the original")
	suite.Equal("soda", clone["beverage"].Value,
		"Cloned Slots should use the *resolved* value rather than the raw value")
}

func (suite SlotsSuite) TestResolve() {
	slots := golexa.Slots{}
	suite.Equal("", slots.Resolve(""),
		"Empty slots should always resolve to blank")
	suite.Equal("", slots.Resolve("foo"),
		"Empty slots should always resolve to blank")

	slots = golexa.NewSlots(
		golexa.NewSlot("name", "Bob Loblaw"),
		golexa.NewSlot("age", "99"),
		golexa.NewSlot("blank", ""),
		golexa.NewResolvedSlot("beverage", "pop", "soda"),
	)
	suite.Equal("", slots.Resolve(""),
		"Non-existent slots should resolve to empty")
	suite.Equal("", slots.Resolve("aslkdfj"),
		"Non-existent slots should resolve to empty")
	suite.Equal("Bob Loblaw", slots.Resolve("name"),
		"Valid slots w/o resolutions should resolve to the 'Value'")
	suite.Equal("99", slots.Resolve("age"),
		"Valid slots w/o resolutions should resolve to the 'Value'")
	suite.Equal("", slots.Resolve("blank"),
		"Valid slots w/o resolutions should resolve to the 'Value'")
	suite.Equal("soda", slots.Resolve("beverage"),
		"Valid slots w/ resolutions should delegate to the resolution authority")
}
