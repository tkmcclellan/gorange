package gorange

import (
	"math"
	"testing"
)

// PARSING:
// Parses empty list
func TestParseEmptyRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{}
	collection, err := ParseRangeCollection([]string{}, ":")

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// Parses non-empty list
func TestParseNonEmptyRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 4, End: math.Inf(1)}}
	collection, err := ParseRangeCollection([]string{"1:2", "4:"}, ":")

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// TODO: VALUES:
// Gets infinite range values
// Gets infinite start range values
// Gets infinite end range values
// Gets finite range values
// Gets singleton range values
// Each value in empty range
// Each value in finite range
// Each value in start infinite range
// Each value in end infinite range
// Each value in infinite range
// Map value in empty range
// Map value in finite range
// Map value in start infinite range
// Map value in end infinite range
// Map value in infinite range
// Gets infinite range values in range
// Gets infinite start range values in range
// Gets infinite end range values in range
// Gets finite range values in range
// Gets singleton range values in range
// Each value in empty range in range
// Each value in finite range in range
// Each value in start infinite range in range
// Each value in end infinite range in range
// Each value in infinite range in range
// Map value in empty range in range
// Map value in finite range in range
// Map value in start infinite range in range
// Map value in end infinite range in range
// Map value in infinite range in range
// Each value in singleton range
// Map value in singleton range
// Each value in singleton range in range
// Map value in singleton range in range

// MERGING:
// Merges empty list
func TestMergeEmptyRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{}
	unmerged, err := ParseRangeCollection([]string{}, ":")
	collection := unmerged.Merge()

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// Merges non-overlapping list
func TestMergeNonOverlappingRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 4, End: math.Inf(1)}}
	unmerged, err := ParseRangeCollection([]string{"1:2", "4:"}, ":")
	collection := unmerged.Merge()

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// Merges overlapping at start
func TestMergeStartOverlappingRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{Range{Start: 1, End: math.Inf(1)}}
	unmerged, err := ParseRangeCollection([]string{"1:2", "1:"}, ":")
	collection := unmerged.Merge()

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// Merges overlapping in middle
func TestMergeMiddleOverlappingRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{Range{Start: 1, End: math.Inf(1)}}
	unmerged, err := ParseRangeCollection([]string{"1:5", "3:"}, ":")
	collection := unmerged.Merge()

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// Merges with infinite beginning
func TestMergeInfiniteBeginningRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{Range{Start: math.Inf(-1), End: 5}}
	unmerged, err := ParseRangeCollection([]string{"1:5", ":3"}, ":")
	collection := unmerged.Merge()

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// Merges with infinite end
func TestMergeInfiniteEndingRangeCollection(t *testing.T) {
	expectedCollection := RangeCollection{Range{Start: 1, End: 6}, Range{Start: 7, End: math.Inf(1)}}
	unmerged, err := ParseRangeCollection([]string{"1:5", "7:", "4:6"}, ":")
	collection := unmerged.Merge()

	if err != nil || !collection.Equal(expectedCollection) {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedCollection, collection)
	}
}

// COMPARING:
// Compares two equal RangeCollections
func TestCompareEqualRangeCollections(t *testing.T) {
	firstCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 3, End: 4}}
	secondCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 3, End: 4}}

	if !firstCollection.Equal(secondCollection) {
		t.Errorf("Failed! %v does not compare to %v", firstCollection, secondCollection)
	}
}

// Compares two unequal RangeCollections of equal length
func TestCompareUnequalRangeCollections(t *testing.T) {
	firstCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 3, End: 4}}
	secondCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 5, End: 7}}

	if firstCollection.Equal(secondCollection) {
		t.Errorf("Failed! %v does not compare to %v", firstCollection, secondCollection)
	}
}

// Compares two RangeCollections of different lengths
func TestCompareUnequalLengthRangeCollections(t *testing.T) {
	firstCollection := RangeCollection{Range{Start: 1, End: 2}}
	secondCollection := RangeCollection{Range{Start: 1, End: 2}, Range{Start: 5, End: 7}}

	if firstCollection.Equal(secondCollection) {
		t.Errorf("Failed! %v does not compare to %v", firstCollection, secondCollection)
	}
}
