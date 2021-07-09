# GDELT

A simple tool to get the latest events from http://data.gdeltproject.org/gdeltv2/lastupdate.txt

## Usage

Requirements:

* [Go 1.15](https://golang.org/dl/)
* [Go modules](https://blog.golang.org/using-go-modules)

Clone this repo or get the library:

```console
go get -u github.com/SpecializedGeneralist/gdelt
```

## Example

```go
package main

import (
	"fmt"
	"github.com/SpecializedGeneralist/gdelt"
)

func main()  {
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
```

