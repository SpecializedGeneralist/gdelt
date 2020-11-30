// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

// GeoType specifies the geographic resolution of the match type.
type GeoType uint8

const (
	NoGeoType GeoType = iota
	Country
	USState
	USCity
	WorldCity
	WorldState
)

func GeoTypeFromInt(value int) (GeoType, bool) {
	if value < 0 && value > 5 {
		return 0, false
	}
	return GeoType(value), true
}

func (g GeoType) String() string {
	switch g {
	case Country:
		return "COUNTRY"
	case USState:
		return "USSTATE"
	case USCity:
		return "USCITY"
	case WorldCity:
		return "WORLDCITY"
	case WorldState:
		return "WORLDSTATE"
	default:
		return ""
	}
}
