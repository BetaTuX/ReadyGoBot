package store

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Hotlap struct {
	Id int64
	/** User's Discord id */
	DriverUid string
	TrackId   string
	Time      time.Duration
	UpdatedAt time.Time
}

type hotlapStore struct {
	listMutex sync.Mutex
	list      []Hotlap
}

var HotlapStore = hotlapStore{
	list: make([]Hotlap, 0),
}

func generateId() int64 {
	if len(HotlapStore.list) > 0 {
		return HotlapStore.list[len(HotlapStore.list)-1].Id + 1
	}
	return 1
}

/*
Returns bool true if a new hotlap as been saved (added or updated), false if driver didn't beat his hotlap
*/
func (store *hotlapStore) SetHotlap(trackId, driverUid string, pTime time.Duration) bool {
	var hotlap Hotlap

	store.listMutex.Lock()
	defer store.listMutex.Unlock()
	for _, v := range store.list {
		timeExists := v.DriverUid == driverUid && v.TrackId == trackId

		if !timeExists {
			continue
		}
		if v.Time > pTime {
			hotlap.Time = pTime
			hotlap.UpdatedAt = time.Now()
		} else {
			// Didn't beat his previous time
			return false
		}
	}
	if (hotlap == Hotlap{}) {
		hotlap = Hotlap{
			Id:        generateId(),
			TrackId:   trackId,
			DriverUid: driverUid,
			Time:      pTime,
			UpdatedAt: time.Now(),
		}
		store.list = append(store.list, hotlap)
	}
	fmt.Println("debug: StoreTime: added new time:", TrackStore.GetTrack(trackId).Name, driverUid, pTime.String())
	return true
}

type ByDuration []Hotlap

func (laps ByDuration) Len() int           { return len(laps) }
func (laps ByDuration) Swap(i, j int)      { laps[i], laps[j] = laps[j], laps[i] }
func (laps ByDuration) Less(i, j int) bool { return laps[i].Time < laps[j].Time }

type ByTime []Hotlap

func (laps ByTime) Len() int           { return len(laps) }
func (laps ByTime) Swap(i, j int)      { laps[i], laps[j] = laps[j], laps[i] }
func (laps ByTime) Less(i, j int) bool { return laps[i].UpdatedAt.Before(laps[j].UpdatedAt) }

func (store *hotlapStore) GetTrackHotlapList(trackId string) []Hotlap {
	store.listMutex.Lock()
	defer store.listMutex.Unlock()
	lapList := make([]Hotlap, 0)

	for _, hotlap := range store.list {
		if hotlap.TrackId == trackId {
			lapList = append(lapList, hotlap)
		}
	}
	sort.Sort(ByDuration(lapList))
	return lapList
}

func (store *hotlapStore) GetDriverHotlapList(pDriverUid string) []Hotlap {
	store.listMutex.Lock()
	defer store.listMutex.Unlock()
	lapList := make([]Hotlap, 0)

	for _, hotlap := range store.list {
		if hotlap.DriverUid == pDriverUid {
			lapList = append(lapList, hotlap)
		}
	}
	sort.Sort(ByTime(lapList))
	return lapList
}

func (store *hotlapStore) GetDriverHottestLap(pDriverUid, pTrackId string) (Hotlap, error) {
	store.listMutex.Lock()
	defer store.listMutex.Unlock()

	for _, v := range store.list {
		if v.DriverUid == pDriverUid && v.TrackId == pTrackId {
			return v, nil
		}
	}
	return Hotlap{}, fmt.Errorf("No hotlap set on track (%s) for driver (%s)", pTrackId, pDriverUid)
}
