package gorange

import (
	"math"
	"testing"
)

// PARSING:
// Parses basic range
func TestParseFullRange(t *testing.T) {
	expectedRange := Range{Start: 3, End: 4}
	grange, err := ParseRange("3:4", ":")

	if err != nil || grange != expectedRange {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedRange, grange)
	}
}

// Parses open start range
func TestParseOpenStartRange(t *testing.T) {
	expectedRange := Range{Start: math.Inf(-1), End: 4}
	grange, err := ParseRange(":4", ":")

	if err != nil || grange != expectedRange {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedRange, grange)
	}
}

// Parses open end range
func TestParseOpenEndRange(t *testing.T) {
	expectedRange := Range{Start: 3, End: math.Inf(1)}
	grange, err := ParseRange("3:", ":")

	if err != nil || grange != expectedRange {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedRange, grange)
	}
}

// Parses infinite range
func TestParseInfiniteRange(t *testing.T) {
	expectedRange := Range{Start: math.Inf(-1), End: math.Inf(1)}
	grange, err := ParseRange(":", ":")

	if err != nil || grange != expectedRange {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedRange, grange)
	}
}

// Parses singular range
func TestParseSingularRange(t *testing.T) {
	expectedRange := Range{Start: 3, End: 3}
	grange, err := ParseRange("3", ":")

	if err != nil || grange != expectedRange {
		t.Errorf("Failed! Expected: %v, Got: %v", expectedRange, grange)
	}
}

// Fails to parse invalid range (fuzzing??)
func TestParseInvalidRange(t *testing.T) {
	grange, err := ParseRange("3:1:1", ":")

	if err == nil {
		t.Errorf("Failed! Expected failure with: %v", grange)
	}
}

// TODO: Parses range with strange separators (fuzzing??)
// Fails to parse start before end range
func TestDoNotParseInvalidRange(t *testing.T) {
	grange, err := NewRange(10, 1)

	if err == nil {
		t.Errorf("Failed! Did not return error when creating invalid range: %v", grange)
	}
}

// MERGING:
// Does not merge non-overlapping ranges
func TestDoNotMergeNonOverlappingRanges(t *testing.T) {
	rangeOne, _ := NewRange(1, 4)
	rangeTwo, _ := NewRange(5, 7)

	_, err := rangeOne.Merge(rangeTwo)

	if err == nil {
		t.Errorf("Failure! Range %v does not overlap range %v", rangeOne, rangeTwo)
	}
}

// Merges start overlapping ranges
func TestMergeStartOverlappingRanges(t *testing.T) {
	rangeOne, _ := NewRange(3, 5)
	rangeTwo, _ := NewRange(2, 4)
	expectedRange, _ := NewRange(2, 5)

	testRange, err := rangeOne.Merge(rangeTwo)

	if err != nil || testRange != expectedRange {
		t.Errorf("Failure! Range %v failed to merge with start overlapping range %v", rangeOne, rangeTwo)
	}
}

// Merges end overlapping ranges
func TestMergeEndOverlappingRanges(t *testing.T) {
	rangeOne, _ := NewRange(2, 4)
	rangeTwo, _ := NewRange(3, 5)
	expectedRange, _ := NewRange(2, 5)

	testRange, err := rangeOne.Merge(rangeTwo)

	if err != nil || testRange != expectedRange {
		t.Errorf("Failure! Range %v failed to merge with end overlapping range %v", rangeOne, rangeTwo)
	}
}

// Merges enveloping range
func TestMergeEnvelopingRanges(t *testing.T) {
	rangeOne, _ := NewRange(3, 5)
	rangeTwo, _ := NewRange(2, 6)
	expectedRange, _ := NewRange(2, 6)

	testRange, err := rangeOne.Merge(rangeTwo)

	if err != nil || testRange != expectedRange {
		t.Errorf("Failure! Range %v failed to merge with enveloping range %v", rangeOne, rangeTwo)
	}
}

// VALUES:
func rangeValueTest(t *testing.T, grange Range, expectedValues []float64) {
	values := grange.Values()

	if len(expectedValues) != len(values) {
		t.Errorf("Failure! Range %v values %v did not equal values %v length", grange, values, expectedValues)
	}

	for i := range values {
		if values[i] != expectedValues[i] {
			t.Errorf("Failure! Range %v values %v did not equal values %v", grange, values, expectedValues)
		}
	}
}

// Gets infinite range values
func TestInfiniteRangeValues(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), math.Inf(1))
	expectedValues := []float64{math.Inf(-1), math.Inf(1)}
	rangeValueTest(t, grange, expectedValues)
}

// Gets infinite start range values
func TestInfiniteStartRangeValues(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), 3)
	expectedValues := []float64{math.Inf(-1), 3}
	rangeValueTest(t, grange, expectedValues)
}

// Gets infinite end range values
func TestInfiniteEndRangeValues(t *testing.T) {
	grange, _ := NewRange(1, math.Inf(1))
	expectedValues := []float64{1, math.Inf(1)}
	rangeValueTest(t, grange, expectedValues)
}

// Gets finite range values
func TestFiniteRangeValues(t *testing.T) {
	grange, _ := NewRange(1, 3)
	expectedValues := []float64{1, 2, 3}
	rangeValueTest(t, grange, expectedValues)
}

// Gets singleton range values
func TestSingletonRangeValues(t *testing.T) {
	grange, _ := NewRange(1, 1)
	expectedValues := []float64{1}
	rangeValueTest(t, grange, expectedValues)
}

func rangeEachValueTest(t *testing.T, grange Range, expectedValues []float64) {
	values := []float64{}

	grange.EachValue(func(value float64) float64 {
		values = append(values, value*2)
		return value
	})

	if len(expectedValues) != len(values) {
		t.Errorf("Failure! Range %v values %v did not equal values %v", grange, values, expectedValues)
	}

	for i := range values {
		if values[i] != expectedValues[i] {
			t.Errorf("Failure! Range %v values %v did not equal values %v", grange, values, expectedValues)
		}
	}
}

// Each value in finite range
func TestFiniteRangeEachValues(t *testing.T) {
	grange, _ := NewRange(1, 3)
	expectedValues := []float64{2, 4, 6}
	rangeEachValueTest(t, grange, expectedValues)
}

// Each value in start infinite range
func TestStartInfiniteRangeEachValues(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), 3)
	expectedValues := []float64{math.Inf(-1), 6}
	rangeEachValueTest(t, grange, expectedValues)
}

// Each value in end infinite range
func TestEndInfiniteRangeEachValues(t *testing.T) {
	grange, _ := NewRange(3, math.Inf(1))
	expectedValues := []float64{6, math.Inf(1)}
	rangeEachValueTest(t, grange, expectedValues)
}

// Each value in infinite range
func TestInfiniteRangeEachValues(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), math.Inf(1))
	expectedValues := []float64{math.Inf(-1), math.Inf(1)}
	rangeEachValueTest(t, grange, expectedValues)
}

func rangeMapValueTest(t *testing.T, grange Range, expectedValues []float64) {
	values := grange.ValueMap(func(value float64) float64 {
		return value * 2
	})

	if len(expectedValues) != len(values) {
		t.Errorf("Failure! Range %v values %v did not equal values %v", grange, values, expectedValues)
	}

	for i := range values {
		if values[i] != expectedValues[i] {
			t.Errorf("Failure! Range %v values %v did not equal values %v", grange, values, expectedValues)
		}
	}
}

// Map value in finite range
func TestFiniteRangeMapValues(t *testing.T) {
	grange, _ := NewRange(1, 3)
	expectedValues := []float64{2, 4, 6}
	rangeMapValueTest(t, grange, expectedValues)
}

// Map value in start infinite range
func TestStartInfiniteRangeMapValues(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), 3)
	expectedValues := []float64{math.Inf(-1), 6}
	rangeMapValueTest(t, grange, expectedValues)
}

// Map value in end infinite range
func TestEndInfiniteRangeMapValues(t *testing.T) {
	grange, _ := NewRange(3, math.Inf(1))
	expectedValues := []float64{6, math.Inf(1)}
	rangeMapValueTest(t, grange, expectedValues)
}

// Map value in infinite range
func TestInfiniteRangeMapValues(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), math.Inf(1))
	expectedValues := []float64{math.Inf(-1), math.Inf(1)}
	rangeMapValueTest(t, grange, expectedValues)
}

func rangeInRangeTest(t *testing.T, grange Range, expectedValues [][]float64) {
	selectionRanges := []Range{
		{Start: math.Inf(-1), End: math.Inf(1)},
		{Start: math.Inf(-1), End: 3},
		{Start: 3, End: math.Inf(1)},
		{Start: 2, End: 5},
	}

	if len(selectionRanges) != len(expectedValues) {
		t.Errorf("Invalid test format")
	}

	for i := range selectionRanges {
		values := grange.ValuesInRange(selectionRanges[i])
		if len(values) != len(expectedValues[i]) {
			t.Errorf("Failed! Values %v from range %v from constraint range %v do not match expected values %v", values, grange, selectionRanges[i], expectedValues[i])
		}

		for j := range values {
			if values[j] != expectedValues[i][j] {
				t.Errorf("Failed! Values %v from range %v from constraint range %v do not match expected values %v", values, grange, selectionRanges[i], expectedValues[i])
			}
		}
	}
}

// Gets infinite range values in range
func TestInfiniteRangeValuesInRange(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), math.Inf(1))
	expectedValues := [][]float64{
		// -Inf, Inf
		{math.Inf(-1), math.Inf(1)},
		// -Inf, 3
		{math.Inf(-1), 3},
		// 3, Inf
		{3, math.Inf(1)},
		// 2, 5
		{2, 3, 4, 5},
	}

	rangeInRangeTest(t, grange, expectedValues)
}

// Gets infinite start range values in range
func TestInfiniteStartRangeValuesInRange(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), 4)
	expectedValues := [][]float64{
		// -Inf, Inf
		{math.Inf(-1), 4},
		// -Inf, 3
		{math.Inf(-1), 3},
		// 3, Inf
		{3, 4},
		// 2, 5
		{2, 3, 4},
	}

	rangeInRangeTest(t, grange, expectedValues)
}

// Gets infinite end range values in range
func TestInfiniteEndRangeValuesInRange(t *testing.T) {
	grange, _ := NewRange(4, math.Inf(1))
	expectedValues := [][]float64{
		// -Inf, Inf
		{4, math.Inf(1)},
		// -Inf, 3
		{},
		// 3, Inf
		{4, math.Inf(1)},
		// 2, 5
		{4, 5},
	}

	rangeInRangeTest(t, grange, expectedValues)
}

// Gets finite range values in range
func TestFiniteRangeValuesInRange(t *testing.T) {
	grange, _ := NewRange(3, 4)
	expectedValues := [][]float64{
		// -Inf, Inf
		{3, 4},
		// -Inf, 3
		{3},
		// 3, Inf
		{3, 4},
		// 2, 5
		{3, 4},
	}

	rangeInRangeTest(t, grange, expectedValues)
}

// Gets singleton range values in range
func TestSingletonRangeValuesInRange(t *testing.T) {
	grange, _ := NewRange(4, 4)
	expectedValues := [][]float64{
		// -Inf, Inf
		{4},
		// -Inf, 3
		{},
		// 3, Inf
		{4},
		// 2, 5
		{4},
	}

	rangeInRangeTest(t, grange, expectedValues)
}

// COMPARING:
// Compares two equal Ranges
func TestCompareEqualRanges(t *testing.T) {
	firstRange := Range{Start: 1, End: 2}
	secondRange := Range{Start: 1, End: 2}

	if firstRange != secondRange {
		t.Errorf("Failed! %v does not basic compare to %v", firstRange, secondRange)
	}

	if !firstRange.Equal(secondRange) {
		t.Errorf("Failed! %v does not function compare to %v", firstRange, secondRange)
	}
}

// Compares two unequal Ranges
func TestCompareUnequalRanges(t *testing.T) {
	firstRange := Range{Start: 1, End: 2}
	secondRange := Range{Start: 1, End: 3}

	if firstRange == secondRange {
		t.Errorf("Failed! %v does not basic compare to %v", firstRange, secondRange)
	}

	if firstRange.Equal(secondRange) {
		t.Errorf("Failed! %v does not function compare to %v", firstRange, secondRange)
	}
}

// Determines if overlapping ranges overlap
func TestOverlappingRangeOverlap(t *testing.T) {
	rangeOne, _ := NewRange(2, 5)
	rangeTwo, _ := NewRange(3, 6)

	if !rangeOne.Overlap(rangeTwo) {
		t.Errorf("Failed! Range %v should overlap range %v", rangeOne, rangeTwo)
	}
}

// Determines if non-overlapping ranges overlap
func TestNonOverlappingRangeOverlap(t *testing.T) {
	rangeOne, _ := NewRange(2, 5)
	rangeTwo, _ := NewRange(6, 10)

	if rangeOne.Overlap(rangeTwo) {
		t.Errorf("Failed! Range %v should not overlap range %v", rangeOne, rangeTwo)
	}
}

// Determines if infinite range is infinite
func TestInfiniteRangeIsInfinite(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), math.Inf(1))

	if !grange.Infinite() {
		t.Errorf("Failed! Range %v should be infinite", grange)
	}
}

// Determines if non-infinite range is infinite
func TestNonInfiniteRangeIsNotInfinite(t *testing.T) {
	grange, _ := NewRange(math.Inf(-1), 1)

	if grange.Infinite() {
		t.Errorf("Failed! Range %v should not be infinite", grange)
	}
}

// Determines if range contains inside value
func TestRangeContainsInsideValue(t *testing.T) {
	grange, _ := NewRange(3, 10)
	value := 5.0

	if !grange.Contains(value) {
		t.Errorf("Failed! Range %v does contain %f", grange, value)
	}
}

// Determines if range contains outside value
func TestRangeDoesNotContainOutsideValue(t *testing.T) {
	grange, _ := NewRange(3, 10)
	value := 12.0

	if grange.Contains(value) {
		t.Errorf("Failed! Range %v does not contain %f", grange, value)
	}
}
