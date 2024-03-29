package main

import (
	"encoding/json"
	"errors"
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
	weight float32
}

func MakeFilterDataRecord(dr DataRecord, start time.Time, finish time.Time) FilteredDataRecord {

	var weight float32

	total_seconds := dr.Finish.Unix() - dr.Start.Unix()
	time_before := start.Unix() - dr.Start.Unix()
	if time_before < 0 {
		time_before = 0
	}
	time_after := dr.Finish.Unix() - finish.Unix()
	if time_after < 0 {
		time_after = 0
	}
	time_outside := time_before + time_after

	if total_seconds == 0 {
		if time_outside == 0 {
			weight = 1
		} else {
			panic("the passed data are outside this manager's day range")
		}
	} else {
		weight = float32(total_seconds-time_outside) / float32(total_seconds)
	}

	return FilteredDataRecord{
		data:   dr,
		weight: weight,
	}
}

type Manager struct {
	ID                  int
	Day                 Date
	Start               time.Time
	Finish              time.Time
	data                []FilteredDataRecord
	chan_in             chan DataRecord
	chan_out            chan float32
	computation_started bool
}

func MakeManager(id int, day Date, start time.Time, finish time.Time) Manager {
	return Manager{
		ID:                  id,
		Day:                 day,
		Start:               start,
		Finish:              finish,
		chan_in:             make(chan DataRecord),
		chan_out:            make(chan float32),
		computation_started: false,
	}
}

func (manager *Manager) PushRecord(record DataRecord) *FilteredDataRecord {
	filtered_record := MakeFilterDataRecord(record, manager.Start, manager.Finish)
	manager.data = append(manager.data, filtered_record)
	return &manager.data[len(manager.data)-1]
}

func (manager *Manager) TryPushRecordAsync(record DataRecord) error {
	if record.Start.After(manager.Finish) ||
		record.Finish.Before(manager.Start) {
		return errors.New("record outside this manager's duty")
	}
	manager.chan_in <- record
	return nil
}

// get the output (synchronous)
func (manager *Manager) GetOutput() float32 {
	s := float32(0)
	for _, v := range manager.data {
		s += v.weight * float32(v.data.X)
	}
	return s
}

// function intended to be called as a go routine to compute the data inside the manager
func (manager *Manager) Compute() {
	s := float32(0)
	for item := range manager.chan_in {
		v := manager.PushRecord(item)
		s += v.weight * float32(v.data.X)
	}
	manager.chan_out <- s
}

func (manager *Manager) TryStartCompute() error {
	if manager.computation_started {
		return errors.New("computation already started")
	}
	manager.computation_started = true
	go manager.Compute()
	return nil
}

// a light weight, hash-able representation of a date
type Date int64

func MakeDate(time *time.Time) Date {
	y, m, d := time.Date()
	m2 := int(m)
	v := y*13*32 + m2*32 + d
	return Date(v)
}

func (date Date) ToTime() time.Time {
	dateInt := int(date)
	y := dateInt / (13 * 32)
	dateInt %= 13 * 32
	m := dateInt / 32
	dateInt %= 32
	d := dateInt

	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

type ManagerKeeper struct {
	// keep a list of all the days
	days []Date
	// hashmap that get the ID (aka the order) from a date
	dayToIndex map[Date]int
	// nested dictionary of manager
	Managers map[int]*map[Date]*Manager
}

func MakeManagerKeeper(data []DataRecord) ManagerKeeper {

	// map the date to the index in the final array
	dayToIndex := make(map[Date]int)

	for _, e := range data {
		start := MakeDate(&e.Start)
		finish := MakeDate(&e.Finish)

		dayToIndex[start] = 0
		dayToIndex[finish] = 0
	}

	days := make([]Date, 0, len(dayToIndex))
	for k := range dayToIndex {
		days = append(days, k)
	}

	sort.Slice(days, func(i, j int) bool {
		return int(days[i]) < int(days[j])
	})

	for i, v := range days {
		dayToIndex[v] = i
	}

	managers := make(map[int]*map[Date]*Manager)

	return ManagerKeeper{
		days,
		dayToIndex,
		managers,
	}
}

// get the manager, AND create a new one if not exist
func (managerKeeper *ManagerKeeper) GetManager(id int, day Date) *Manager {

	val, ok := managerKeeper.Managers[id]

	if !ok {
		new_map := make(map[Date]*Manager)
		val = &new_map
		managerKeeper.Managers[id] = val
	}

	val2, ok := (*val)[day]

	if !ok {

		start := day.ToTime()

		finish := start.Add(ALMOST_A_DAY)

		new_manager := MakeManager(id, day, start, finish)

		val2 = &new_manager
		(*val)[day] = val2
	}
	return val2
}

func (managerKeeper *ManagerKeeper) GetNextDay(date Date) (Date, error) {
	index, ok := managerKeeper.dayToIndex[date]
	if !ok {
		return 0, errors.New("unable to find specified date")
	}
	index += 1
	if index >= len(managerKeeper.days) {
		return 0, errors.New("dey dose not have a successor")
	}
	return managerKeeper.days[index], nil
}

func (managerKeeper *ManagerKeeper) ProcessData(data []DataRecord) {

	for _, d := range data {
		id := d.ID
		day := MakeDate(&d.Start)

		// insert this record until it is no longer inside the time range of the manager
		for {
			manager := managerKeeper.GetManager(id, day)
			manager.TryStartCompute()
			err := manager.TryPushRecordAsync(d)
			if err != nil {
				break
			}
			day, err = managerKeeper.GetNextDay(day)
			if err != nil {
				break
			}
		}
	}
	managerKeeper.CloseAllChannels()
}

func (managerKeeper *ManagerKeeper) CloseAllChannels() {
	for _, v := range managerKeeper.Managers {
		for _, manager := range *v {
			close(manager.chan_in)
		}
	}
}

func (managerKeeper *ManagerKeeper) PrintData() {

	for id, v := range managerKeeper.Managers {
		for day, manager := range *v {

			day_fmt := day.ToTime()
			x := <-manager.chan_out

			fmt.Printf("Day: %v Id: %v X: %v\n", day_fmt, id, x)
		}
	}

}

var ALMOST_A_DAY time.Duration

func main() {

	var err error
	ALMOST_A_DAY, err = time.ParseDuration("23h59m59s")

	d := time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC)
	x := MakeDate(&d)
	x.ToTime()

	if err != nil {
		panic(err)
	}

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

	manager := MakeManagerKeeper(data)
	manager.ProcessData(data)
	manager.PrintData()

}
