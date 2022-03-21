package gorange

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Range is a struct for representing infinite, open-ended, and finite ranges of values
type Range struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// Equal tests if two ranges are identical
func (r Range) Equal(other Range) bool {
	return r.Start == other.Start && r.End == other.End
}

// Overlap tests if the values of one range overlap the values of another
func (r Range) Overlap(other Range) bool {
	return r.End >= other.Start || other.Start <= r.End
}

// Infinite tests if a range is infinite in both directions
func (r Range) Infinite() bool {
	return r.Start == math.Inf(-1) && r.End == math.Inf(1)
}

// Contains tests if a range contains a given value
func (r Range) Contains(float float64) bool {
	return float >= r.Start && float <= r.End
}

// Merge merges one range with another.
// It will return an error if the ranges do not overlap
func (r Range) Merge(other Range) (Range, error) {
	if !r.Overlap(other) {
		return r, errors.New(fmt.Sprintf("Range %v does not overlap range %v", r, other))
	} else {
		newRange, err := NewRange(math.Min(r.Start, other.Start), math.Max(r.End, other.End))
		if err != nil {
			return r, err
		} else {
			return newRange, nil
		}
	}
}

func (r Range) values(fn func(float64) float64) []float64 {
	valueList := []float64{}

	if r.Start == math.Inf(-1) {
		valueList = append(valueList, fn(math.Inf(-1)))
		valueList = append(valueList, fn(r.End))
		return valueList
	} else if r.End == math.Inf(1) {
		valueList = append(valueList, fn(r.Start))
		valueList = append(valueList, fn(math.Inf(1)))
		return valueList
	} else {
		for i := r.Start; i <= r.End; i++ {
			valueList = append(valueList, fn(i))
		}
	}

	return valueList
}

// Values returns the values in a range. If one end of the
// range is open-ended, this function will return a list of the
// range's start and end.
func (r Range) Values() []float64 {
	return r.values(func(v float64) float64 { return v })
}

// EachValue accepts a function to be applied to each value of a range.
func (r Range) EachValue(fn func(float64) float64) {
	r.values(fn)
}

// ValueMap accepts a function to be applied to each value of a range the modified
// values will be returned in place of the original values.
func (r Range) ValueMap(fn func(float64) float64) []float64 {
	return r.values(fn)
}

func (r Range) valuesInRange(other Range, fn func(float64) float64) []float64 {
	if !r.Overlap(other) {
		return []float64{}
	}

	if other.Infinite() && r.Infinite() {
		return []float64{fn(math.Inf(-1)), fn(math.Inf(1))}
	} else if other.Infinite() {
		return r.values(fn)
	} else if other.Start == math.Inf(-1) {
		if r.Start == math.Inf(-1) {
			return []float64{fn(math.Inf(-1)), fn(math.Min(other.End, r.End))}
		} else {
			values := []float64{}

			for i := r.Start; i <= math.Min(other.End, r.End); i++ {
				values = append(values, fn(i))
			}

			return values
		}
	} else if other.End == math.Inf(1) {
		if r.End == math.Inf(1) {
			return []float64{fn(math.Max(other.Start, r.Start)), fn(math.Inf(1))}
		} else {
			values := []float64{}

			for i := math.Max(other.Start, r.Start); i <= r.End; i++ {
				values = append(values, fn(i))
			}

			return values
		}
	} else {
		values := []float64{}

		for i := math.Max(other.Start, r.Start); i <= math.Min(other.End, r.End); i++ {
			values = append(values, fn(i))
		}

		return values
	}
}

// ValuesInRange returns the values in this range that overlap with the values in the supplied range
func (r Range) ValuesInRange(other Range) []float64 {
	return r.valuesInRange(other, func(v float64) float64 { return v })
}

// EachValueInRange executes function fn on each value in this range that overlaps with
// the supplied range.
func (r Range) EachValueInRange(other Range, fn func(float64) float64) {
	r.valuesInRange(other, fn)
}

// MapValueInRange executes function fn on each value in this range that overlaps with
// the supplied range, and returns the values modified by fn.
func (r Range) MapValueInRange(other Range, fn func(float64) float64) []float64 {
	return r.valuesInRange(other, fn)
}

// NewRange creates a new range. If end is less then start, it will return an error
func NewRange(start float64, end float64) (Range, error) {
	if start > end {
		return Range{}, errors.New(fmt.Sprintf("Start: %f is after End: %f", start, end))
	}
	return Range{Start: start, End: end}, nil
}

// ParseRange parses a range from a string. If the range is not in one of these forms
// (assuming the delimiter to be ":"), [":", "Num:", ":Num", "Num:Num"], ParseRange
// will return an error.
func ParseRange(srange string, delimiter string) (Range, error) {
	var prange Range

	if strings.Contains(srange, delimiter) {
		index := strings.Index(srange, delimiter)

		if srange == delimiter {
			return Range{Start: math.Inf(-1), End: math.Inf(1)}, nil
		}

		var float float64
		var err error = nil
		switch index {
		case 0:
			float, err = strconv.ParseFloat(srange[1:], 64)

			if err != nil {
				return prange, parsingError(srange, delimiter, err)
			}

			var err error = nil
			prange, err = NewRange(math.Inf(-1), float)
			if err != nil {
				return prange, err
			}
		case len(srange) - 1:
			float, err = strconv.ParseFloat(srange[:len(srange)-1], 64)

			if err != nil {
				return prange, parsingError(srange, delimiter, err)
			}

			var err error = nil
			prange, err = NewRange(float, math.Inf(1))
			if err != nil {
				return prange, err
			}
		default:
			ends := strings.Split(srange, delimiter)

			if len(ends) != 2 {
				return prange, parsingError(srange, delimiter, err)
			}

			start, err := strconv.ParseFloat(ends[0], 64)
			end, err := strconv.ParseFloat(ends[1], 64)

			if err != nil {
				return prange, parsingError(srange, delimiter, err)
			}

			err = nil
			prange, err = NewRange(start, end)
			if err != nil {
				return prange, err
			}
		}
	} else {
		float, err := strconv.ParseFloat(srange, 64)
		if err != nil {
			return prange, parsingError(srange, delimiter, err)
		}

		err = nil
		prange, err = NewRange(float, float)
	}

	return prange, nil
}

func parsingError(srange string, delimiter string, err error) error {
	return errors.New(fmt.Sprintf("Error parsing range (%s) with delimiter (%s): %v", srange, delimiter, err))
}
