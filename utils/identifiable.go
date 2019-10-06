package utils

type Identifiable interface {
	Identifier() string
}

type IdentifiableIterable interface {
	Get(index int) Identifiable
	Len() int
}

type IdentityMap = map[string]Identifiable

func AsMap(objects IdentifiableIterable) *IdentityMap {
	var total = objects.Len()
	var interfaceSlice = make(IdentityMap, total)
	for i := 0; i < total; i++ {
		value := objects.Get(i)
		interfaceSlice[value.Identifier()] = value
	}
	return &interfaceSlice
}

func ToIdentifiable(objects IdentifiableIterable) *[]Identifiable {
	var total = objects.Len()
	var interfaceSlice = make([]Identifiable, total)
	for i := 0; i < total; i++ {
		interfaceSlice[i] = objects.Get(i)
	}
	return &interfaceSlice
}

func ContainsByIdentity(slice *[]Identifiable, value *string) bool {
	for _, v := range *slice {
		if *value == v.Identifier() {
			return true
		}
	}
	return false
}

func ContainsAll(slice *[]Identifiable, values *[]string) bool {
	for _, v := range *values {
		if !ContainsByIdentity(slice, &v) {
			return false
		}
	}
	return true
}

type ToIdentityConverter func(interface{}) Identifiable
