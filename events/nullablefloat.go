// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import "strconv"

// NullableFloat64 represents a float64 value that may be null.
type NullableFloat64 struct {
	Float64 float64
	// Valid is true if Float64 is not NULL
	Valid bool
}

var nullNullableFloat64 = NullableFloat64{Float64: 0, Valid: false}

// ParseNullableFloat64 parses a string value, converting it to NullableFloat64.
func ParseNullableFloat64(value string) (NullableFloat64, error) {
	if len(value) == 0 {
		return nullNullableFloat64, nil
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nullNullableFloat64, err
	}
	return NullableFloat64{Float64: f, Valid: true}, nil
}
