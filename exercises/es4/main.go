package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type RawDataRecord struct {
	Start  string `json:"start,omitempty"`
	Finish string `json:"finish,omitempty"`
	ID     string `json:"id,omitempty"`
	X      string `json:"x,omitempty"`
}

type DataRecord struct {
	Start  time.Time `json:"start,omitempty"`
	Finish time.Time `json:"finish,omitempty"`
	ID     int       `json:"id,omitempty"`
	X      int       `json:"x,omitempty"`
}

const TIME_LAYOUT = "2006-01-02 15:04:05.999999"

func MakeDataRecord(rdr *RawDataRecord) (DataRecord, error) {

	start, err := time.Parse(TIME_LAYOUT, rdr.Start)
	if err != nil {
		return DataRecord{}, err
	}
	finish, err := time.Parse(TIME_LAYOUT, rdr.Finish)
	if err != nil {
		return DataRecord{}, err
	}
	id, err := strconv.Atoi(rdr.ID)
	if err != nil {
		return DataRecord{}, err
	}
	x, err := strconv.Atoi(rdr.X)
	if err != nil {
		return DataRecord{}, err
	}
	return DataRecord{start, finish, id, x}, nil
}
func MakeDataRecordList(rdr *[]RawDataRecord) ([]DataRecord, error) {
	to_return := make([]DataRecord, len(*rdr))
	for i, e := range *rdr {
		value, err := MakeDataRecord(&e)
		if err != nil {
			return to_return, err
		}
		to_return[i] = value
	}
	return to_return, nil
}

func GetData() ([]DataRecord, error) {
	dat, err := os.ReadFile("./dataset.json")
	if err != nil {
		return nil, err
	}

	var raw_record []RawDataRecord

	err = json.Unmarshal([]byte(dat), &raw_record)
	if err != nil {
		panic(err)
	}

	record, err := MakeDataRecordList(&raw_record)
	if err != nil {
		panic(err)
	}

	return record, nil
}

type FilteredDataRecord struct {
	data   DataRecord
	waight float32
}

type Manager struct {
	ID       int
	Day      int
	Start    time.Time
	Finish   time.Time
	data     []FilteredDataRecord
	chan_in  chan DataRecord
	chan_out chan float32
}

func main() {
	// get the data
	data, err := GetData()
	if err != nil {
		panic(err)
	}

	// sort by time start
	sort.Slice(data, func(i, j int) bool {
		return data[i].Start.Before(data[j].Start)
	})

	for _, e := range data {
		fmt.Println(e.Start)
	}

}
