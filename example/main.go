// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/specializedgeneralist/gdelt"
)

func main() {
	fmt.Println("Getting latest events from GDELT...")
	events, err := gdelt.GetLatestEvents()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d new events.\n", len(events))
	for _, e := range events {
		fmt.Println(e.SourceURL)
	}
}
