// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import (
	"fmt"
	"github.com/jackc/pgtype"
)

type GeoData struct {
	// Type specifies the geographic resolution of the match type.
	Type GeoType
	// Fullname is the full human-readable name of the matched location. In
	// the case of a country it is simply the country name. For US and World
	// states it is in the format of "State, Country Name", while for all other
	// matches it is in the format of "City/Landmark, State, Country".
	// This can be used to label locations when placing events on a map.
	Fullname string
	// CountryCode is the 2-character FIPS10-4 country code for the location.
	CountryCode string
	ADM1Code    string
	ADM2Code    string
	// Lat is the centroid latitude of the landmark for mapping.
	Lat NullableFloat64
	// Long is the centroid longitude of the landmark for mapping.
	Long      NullableFloat64
	FeatureID string
}

func (g *GeoData) PointCoordinates() pgtype.Point {
	if !g.Lat.Valid || !g.Long.Valid {
		return pgtype.Point{
			P:      pgtype.Vec2{X: 0, Y: 0},
			Status: pgtype.Null,
		}
	}
	return pgtype.Point{
		P:      pgtype.Vec2{X: g.Long.Float64, Y: g.Lat.Float64},
		Status: pgtype.Present,
	}
}

func (g *GeoData) CountryCodeISO31661() (string, error) {
	if len(g.CountryCode) == 0 {
		return "", nil
	}
	isoCode, ok := FIPS104ToISO31661[g.CountryCode]
	if !ok {
		return "", fmt.Errorf("gdelt: unknown FIPS 10-4 country code %#v", g.CountryCode)
	}
	return isoCode, nil
}
