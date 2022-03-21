package gorange

import (
	"math"
	"sort"
)

// RangeCollection represents a collection of Ranges
type RangeCollection []Range

// NewRangeCollection creates a new RangeCollection
func NewRangeCollection(ranges []Range) RangeCollection {
	collection := RangeCollection{}

	for _, grange := range ranges {
		collection = append(collection, grange)
	}

	return collection
}

// Len returns the length of a RangeCollection
func (collection RangeCollection) Len() int {
	return len(collection)
}

// Less tests if element i in a RangeCollection is less than element j
func (collection RangeCollection) Less(i, j int) bool {
	if collection[i].Start == collection[j].Start {
		return collection[i].End < collection[j].End
	} else {
		return collection[i].Start < collection[j].Start
	}
}

// Swap swaps elements i and j in a RangeCollection
func (collection RangeCollection) Swap(i, j int) {
	collection[i], collection[j] = collection[j], collection[i]
}

// IsMerged tests if a RangeCollection has been merged
func (collection RangeCollection) IsMerged() bool {
	if !sort.IsSorted(collection) {
		return false
	}

	for i := 0; i < len(collection)-1; i++ {
		if collection[i].Overlap(collection[i+1]) {
			return false
		}
	}

	return true
}

// Values returns all values represented by the Ranges in a RangeCollection
func (collection RangeCollection) Values() []float64 {
	if !collection.IsMerged() {
		collection = collection.Merge()
	}

	values := []float64{}

	for _, grange := range collection {
		grange.EachValue(func(value float64) float64 {
			values = append(values, value)
			return value
		})
	}

	return values
}

// ValuesInRange returns all values contained within this RangeCollection that are also
// contained in the supplied Range
func (collection RangeCollection) ValuesInRange(r Range) []float64 {
	if !collection.IsMerged() {
		collection = collection.Merge()
	}

	values := []float64{}

	for _, grange := range collection {
		grange.EachValueInRange(r, func(value float64) float64 {
			values = append(values, value)

			return value
		})
	}

	return values
}

// Merge merges the Ranges in this RangeCollection so that all Ranges are
// in order and non-overlapping
func (collection RangeCollection) Merge() RangeCollection {
	if len(collection) == 0 {
		return collection
	}

	if collection.IsMerged() {
		return collection
	}

	sort.Sort(collection)

	newCollection := RangeCollection{}

	currentRange := collection[0]
	if currentRange.End == math.Inf(1) {
		return append(newCollection, currentRange)
	}

	for i := 1; i < len(collection); i++ {
		if currentRange.Overlap(collection[i]) {
			currentRange, _ = currentRange.Merge(collection[i])

			if currentRange.End == math.Inf(1) {
				newCollection = append(newCollection, currentRange)
				return newCollection
			}
		} else {
			newCollection = append(newCollection, currentRange)
			currentRange = collection[i]
		}
	}

	newCollection = append(newCollection, currentRange)

	return newCollection
}

// Equal tests if two RangeCollections contain the same Ranges
func (collection RangeCollection) Equal(other RangeCollection) bool {
	if len(collection) != len(other) {
		return false
	}

	for i := 0; i < len(collection); i++ {
		if collection[i] != other[i] {
			return false
		}
	}

	return true
}

// ParseRangeCollection parses a list of Ranges in string form. If any range is not in the
// correct format, this function will return an error
func ParseRangeCollection(collection []string, delimiter string) (RangeCollection, error) {
	rcollection := RangeCollection{}

	for _, srange := range collection {
		grange, err := ParseRange(srange, delimiter)
		if err != nil {
			return rcollection, err
		}
		rcollection = append(rcollection, grange)
	}

	return rcollection, nil
}
