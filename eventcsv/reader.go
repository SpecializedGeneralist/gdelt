// Copyright 2020 SpecializedGeneralist. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eventcsv

import (
	"encoding/csv"
	"fmt"
	"github.com/SpecializedGeneralist/gdelt/events"
	"io"
	"strconv"
)

type Reader struct {
	csvReader *csv.Reader
}

func NewReader(r io.Reader) *Reader {
	csvReader := csv.NewReader(r)
	csvReader.Comma = '\t'
	return &Reader{csvReader: csvReader}
}

func (r *Reader) Read() (*events.Event, error) {
	csvRecord, err := r.csvReader.Read()
	if err != nil {
		// This includes io.EOF
		return nil, err
	}

	if len(csvRecord) != 61 {
		return nil, fmt.Errorf("expected 61 CSV columns, actual %d", len(csvRecord))
	}

	event := &events.Event{}

	event.GlobalEventID, err = strconv.ParseUint(csvRecord[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse GlobalEventID %#v", csvRecord[0])
	}

	event.Day, err = strconv.Atoi(csvRecord[1])
	if err != nil {
		return nil, fmt.Errorf("parse Day %#v", csvRecord[1])
	}

	event.MonthYear, err = strconv.Atoi(csvRecord[2])
	if err != nil {
		return nil, fmt.Errorf("parse MonthYear %#v", csvRecord[2])
	}

	event.Year, err = strconv.Atoi(csvRecord[3])
	if err != nil {
		return nil, fmt.Errorf("parse Year %#v", csvRecord[3])
	}

	event.FractionDate, err = strconv.ParseFloat(csvRecord[4], 64)
	if err != nil {
		return nil, fmt.Errorf("parse FractionDate %#v", csvRecord[4])
	}

	readActorData(&event.Actor1, csvRecord[5:15])
	readActorData(&event.Actor2, csvRecord[15:25])

	event.IsRootEvent, err = strconv.Atoi(csvRecord[25])
	if err != nil {
		return nil, fmt.Errorf("parse IsRootEvent %#v", csvRecord[25])
	}

	event.EventCode = csvRecord[26]
	event.EventBaseCode = csvRecord[27]
	event.EventRootCode = csvRecord[28]

	event.QuadClass, err = strconv.Atoi(csvRecord[29])
	if err != nil {
		return nil, fmt.Errorf("parse QuadClass %#v", csvRecord[29])
	}

	event.GoldsteinScale, err = strconv.ParseFloat(csvRecord[30], 64)
	if err != nil {
		return nil, fmt.Errorf("parse GoldsteinScale %#v", csvRecord[30])
	}

	event.NumMentions, err = strconv.Atoi(csvRecord[31])
	if err != nil {
		return nil, fmt.Errorf("parse NumMentions %#v", csvRecord[31])
	}

	event.NumSources, err = strconv.Atoi(csvRecord[32])
	if err != nil {
		return nil, fmt.Errorf("parse NumSources %#v", csvRecord[32])
	}

	event.NumArticles, err = strconv.Atoi(csvRecord[33])
	if err != nil {
		return nil, fmt.Errorf("parse NumArticles %#v", csvRecord[33])
	}

	event.AvgTone, err = strconv.ParseFloat(csvRecord[34], 64)
	if err != nil {
		return nil, fmt.Errorf("parse AvgTone %#v", csvRecord[34])
	}

	err = readGeoData(&event.Actor1Geo, csvRecord[35:43])
	if err != nil {
		return nil, fmt.Errorf("reading Actor1Geo: %v", err)
	}
	err = readGeoData(&event.Actor2Geo, csvRecord[43:51])
	if err != nil {
		return nil, fmt.Errorf("reading Actor2Geo: %v", err)
	}
	err = readGeoData(&event.ActionGeo, csvRecord[51:59])
	if err != nil {
		return nil, fmt.Errorf("reading ActionGeo: %v", err)
	}

	event.DateAdded, err = strconv.ParseUint(csvRecord[59], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse DATEADDED %#v", csvRecord[59])
	}

	event.SourceURL = csvRecord[60]

	return event, nil
}

func readActorData(a *events.ActorData, csvFields []string) {
	a.Code = csvFields[0]
	a.Name = csvFields[1]
	a.CountryCode = csvFields[2]
	a.KnownGroupCode = csvFields[3]
	a.EthnicCode = csvFields[4]
	a.Religion1Code = csvFields[5]
	a.Religion2Code = csvFields[6]
	a.Type1Code = csvFields[7]
	a.Type2Code = csvFields[8]
	a.Type3Code = csvFields[9]
}

func readGeoData(g *events.GeoData, csvFields []string) error {
	var err error
	intGeoType, err := strconv.Atoi(csvFields[0])
	if err != nil {
		return fmt.Errorf("parse Type %#v", csvFields[0])
	}
	var geoTypeOk bool
	g.Type, geoTypeOk = events.GeoTypeFromInt(intGeoType)
	if !geoTypeOk {
		return fmt.Errorf("unexpected GeoType value %d", intGeoType)
	}

	g.Fullname = csvFields[1]
	g.CountryCode = csvFields[2]
	g.ADM1Code = csvFields[3]
	g.ADM2Code = csvFields[4]

	if len(csvFields[5]) > 0 {
		g.Lat, err = events.ParseNullableFloat64(csvFields[5])
		if err != nil {
			return fmt.Errorf("parse Lat %#v", csvFields[5])
		}
	}

	if len(csvFields[6]) > 0 {
		g.Long, err = events.ParseNullableFloat64(csvFields[6])
		if err != nil {
			return fmt.Errorf("parse Long %#v", csvFields[6])
		}
	}

	g.FeatureID = csvFields[7]

	return nil
}
