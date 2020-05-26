package golexa

const resolutionSuccessCode = "ER_SUCCESS_MATCH"

// Slots represents a set of runtime values that the Alexa interaction model parsed out for you.
type Slots map[string]Slot

// Slot is one of placeholder for the important data used in an utterance/query. It can be a very
// open-ended bit of text as you'd have in an AMAZON.SearchQuery slot or one of the phrases w/
// synonyms you set up in a custom slot.
type Slot struct {
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Resolutions resolutions `json:"resolutions"`
}

// Clone creates a copy of all of the slots and their RESOLVED values. Typically you use this when you
// want to include a set of slots in your response w/o modifying the map in the request. Be aware that
// while it preserves the resolved value, you will lose the resolution authority data from the original.
func (s Slots) Clone() Slots {
	slots := Slots{}
	for slotName, slot := range s {
		slots[slotName] = Slot{Name: slot.Name, Value: slot.Resolve()}
	}
	return slots
}

// Resolve locates the specified slot entry and returns its resolved value.
func (s Slots) Resolve(slotName string) string {
	if slot, ok := s[slotName]; ok {
		return slot.Resolve()
	}
	return ""
}

// Resolve takes into account the synonyms and resolutions, returning the mapped value that
// the Alexa API thinks we want. If there was no resolution data, you'll simply get back the
// transcribed text from what the user actually said.
func (slot Slot) Resolve() string {
	// They spoke a synonym for your custom slot or they said something like "this month" and it
	// resolved to the ISO date string.
	if resolvedValue := slot.Resolutions.ResolutionPerAuthority.resolvedValue(); resolvedValue != "" {
		return resolvedValue
	}

	// There was no synonym/mapping for what the user spoke, so use their exact word(s)
	return slot.Value
}

type resolutions struct {
	ResolutionPerAuthority resolutionAuthorities `json:"resolutionsPerAuthority"`
}

type resolutionAuthorities []resolutionPerAuthority

func (resolutions resolutionAuthorities) resolvedValue() string {
	if len(resolutions) == 0 {
		return ""
	}

	resolution := resolutions[0]
	if resolution.Status.Code == resolutionSuccessCode && len(resolution.Values) > 0 {
		return resolution.Values[0].Value.Name
	}
	return ""
}

type resolutionPerAuthority struct {
	Authority string             `json:"authority"`
	Status    resolutionStatus   `json:"status"`
	Values    []resolutionValues `json:"values"`
}

type resolutionStatus struct {
	Code string `json:"code"`
}

type resolutionValues struct {
	Value resolutionValue `json:"value"`
}

type resolutionValue struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// NewResolvedSlot is mainly used for faking test data to create a slot that
// has an uttered value as well as a resolution value.
func NewResolvedSlot(name, value, resolvedValue string) Slot {
	return Slot{
		Name:  name,
		Value: value,
		Resolutions: resolutions{
			ResolutionPerAuthority: resolutionAuthorities{
				resolutionPerAuthority{
					Status: resolutionStatus{
						Code: resolutionSuccessCode,
					},
					Values: []resolutionValues{
						{Value: resolutionValue{Name: resolvedValue}},
					},
				},
			},
		},
	}
}

// NewSlot is used primarily for faking data in tests. It creates a simple slot/value with
// no alternate resolution info.
func NewSlot(name, value string) Slot {
	return Slot{
		Name:  name,
		Value: value,
	}
}

// NewSlots creates a Slots map containing entries for all of the individual slot values.
func NewSlots(values ...Slot) Slots {
	slots := Slots{}
	for _, value := range values {
		slots[value.Name] = value
	}
	return slots
}
