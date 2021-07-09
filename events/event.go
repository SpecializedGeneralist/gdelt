// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import (
	"fmt"
	"time"
)

type Event struct {
	// GlobalEventID is the globally unique identifier assigned to each event
	// record that uniquely identifies it in GDELT master dataset.
	GlobalEventID uint64
	Day           int
	MonthYear     int
	Year          int
	FractionDate  float64

	Actor1 ActorData
	Actor2 ActorData

	IsRootEvent int
	// EventCode is the raw CAMEO action code describing the action that Actor1
	// performed upon Actor2.
	EventCode string
	// EventBaseCode is the level two leaf root node category, when applicable.
	// CAMEO event codes are defined in a three-level taxonomy. For events at
	// level three in the taxonomy, this yields its level two leaf root node.
	// For example, code "0251" ("Appeal for easing of administrative
	// sanctions") would yield an EventBaseCode of "025" ("Appeal to yield").
	// This makes it possible to aggregate events at various resolutions of
	// specificity. For events at levels two or one, this field will be set
	// to EventCode.
	EventBaseCode string
	// EventRootCode is similar to EventBaseCode and defines the root-level
	// category the event code falls under. For example, code "0251" ("Appeal
	// for easing of administrative sanctions") has a root code of "02"
	// ("Appeal"). This makes it possible to aggregate events at various
	// resolutions of specificity. For events at levels two or one, this field
	// will be set to EventCode.
	EventRootCode  string
	QuadClass      int
	GoldsteinScale NullableFloat64
	NumMentions    int
	NumSources     int
	NumArticles    int
	AvgTone        float64

	Actor1Geo GeoData
	Actor2Geo GeoData
	// ActionGeo captures the location information closest to the point in the
	// event description that contains the actual statement of action and is
	// the best location to use for placing events on a map or in other spatial
	// context.
	ActionGeo GeoData

	// DateAdded stores the date the event was added to the master database in
	// "YYYYMMDDHHMMSS" format in the UTC timezone.
	DateAdded uint64
	// SourceURL records the URL or citation of the first news report it found
	// this event in. In most cases this is the first report it saw the article
	// in, but due to the timing and flow of news reports through the processing
	// pipeline, this may not always be the very first report, but is at least
	// in the first few reports.
	SourceURL string
}

var dateAddedTimeLayout = "20060102150405"

// DateAddedTime converts DateAdded int value to time.Time.
func (e *Event) DateAddedTime() (time.Time, error) {
	s := fmt.Sprintf("%014d", e.DateAdded)
	if len(s) != 14 {
		return time.Time{}, fmt.Errorf("unexpected DateAdded value %d", e.DateAdded)
	}
	return time.Parse(dateAddedTimeLayout, s)
}

// AllCameoEventCodes returns one or more CAMEO event codes from EventCode,
// EventBaseCode, and EventRootCode, keeping only one unique category code per
// level.
func (e *Event) AllCameoEventCodes() []string {
	s := make([]string, 0)
	if len(e.EventRootCode) == 0 {
		return s
	}
	s = append(s, e.EventRootCode)
	if e.EventBaseCode == e.EventRootCode || len(e.EventBaseCode) == 0 {
		return s
	}
	s = append(s, e.EventBaseCode)
	if e.EventCode == e.EventBaseCode || e.EventCode == e.EventRootCode || len(e.EventCode) == 0 {
		return s
	}
	s = append(s, e.EventCode)
	return s
}
