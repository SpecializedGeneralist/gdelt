// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdelt

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/SpecializedGeneralist/gdelt/eventcsv"
	"github.com/SpecializedGeneralist/gdelt/events"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// DataExportCSVZipReference is a reference to GDELT Export CSV Zip-compressed file.
type DataExportCSVZipReference struct {
	Size   int
	MD5Sum string
	URL    string
}

var reExportCSVZip = regexp.MustCompile(
	`(?P<size>\d+) (?P<md5sum>[0-9a-f]{32}) (?P<url>http://data.gdeltproject.org/gdeltv2/\d{14}.export.CSV.zip)`)

var defaultCSVDataFileListURL = "http://data.gdeltproject.org/gdeltv2/lastupdate.txt"

func GetLatestEvents() ([]*events.Event, error) {
	zipRef, err := FetchDataExportCSVZipReference(defaultCSVDataFileListURL)
	if err != nil {
		return nil, fmt.Errorf("gdelt: fetching data export CSV Zip reference: %v", err)
	}
	eventRecords, err := zipRef.GetLatestEvents()
	if err != nil {
		return nil, fmt.Errorf("gdelt: fetching events from CSV Zip reference: %v", err)
	}
	return eventRecords, nil
}

// FetchDataExportCSVZipReference fetches a reference to GDELT Export CSV
// Zip-compressed file from the given file list URL.
func FetchDataExportCSVZipReference(CSVDataFileListURL string) (*DataExportCSVZipReference, error) {
	content, err := httpGetString(CSVDataFileListURL)
	if err != nil {
		return nil, err
	}

	matches := reExportCSVZip.FindAllStringSubmatch(content, -1)
	if len(matches) != 1 {
		return nil, fmt.Errorf("unexpected GDELT Export CSV Zip %d matches in content %#v",
			len(matches), content)
	}
	match := matches[0]

	size, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, fmt.Errorf("convert file size %#v to int: %v", match[0], err)
	}

	return &DataExportCSVZipReference{
		Size:   size,
		MD5Sum: match[2],
		URL:    match[3],
	}, nil
}

func (d *DataExportCSVZipReference) GetLatestEvents() ([]*events.Event, error) {
	content, err := httpGet(d.URL)
	if err != nil {
		return nil, err
	}

	err = d.checkMD5Sum(content)
	if err != nil {
		return nil, err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return nil, fmt.Errorf("new zip reader: %v", err)
	}

	var events []*events.Event = nil

	for _, zipFile := range zipReader.File {
		if !strings.HasSuffix(zipFile.Name, ".export.CSV") {
			continue
		}
		if events != nil {
			return nil, fmt.Errorf("multiple export CSV files found in Zip archive")
		}
		events, err = processCSVZipFile(zipFile)
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func processCSVZipFile(zipFile *zip.File) (records []*events.Event, err error) {
	f, err := zipFile.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		e := f.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	records = make([]*events.Event, 0)

	r := eventcsv.NewReader(f)
	for {
		event, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, event)
	}

	return records, nil
}

func (d *DataExportCSVZipReference) checkMD5Sum(content []byte) error {
	actual := fmt.Sprintf("%x", md5.Sum(content))
	if actual != d.MD5Sum {
		return fmt.Errorf("md5 sum: expected %s, actual %s", d.MD5Sum, actual)
	}
	return nil
}

func httpGetString(url string) (string, error) {
	content, err := httpGet(url)
	return string(content), err
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET %s: %v", url, err)
	}
	defer func() {
		e := resp.Body.Close()
		if e != nil && err == nil {
			err = e
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET %s returned status code %d", url, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read whole body from HTTP GET %s: %v", url, err)
	}
	return body, nil
}
