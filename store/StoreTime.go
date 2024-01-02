package store

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Hotlap struct {
	gorm.Model
	/** User's Discord id */
	DriverUid string
	TrackId   string
	Time      time.Duration
}

type hotlapStore struct {
	listMutex sync.Mutex
	list      []Hotlap
}

var HotlapStore = hotlapStore{
	list: make([]Hotlap, 0),
}

/*
Returns bool true if a new hotlap as been saved (added or updated), false if driver didn't beat his hotlap
*/
func (store *hotlapStore) SetHotlap(trackId, driverUid string, pTime time.Duration) (bool, error) {
	var hotlap Hotlap

	if err := DB.Find(&hotlap, Hotlap{
		DriverUid: driverUid,
		TrackId:   trackId,
	}).Error; err != nil {
		log.Printf("WARN: Finding track resulted in error: %v\n", err)
		return false, fmt.Errorf("finding track resulted in error: %v", err)
	}
	for _, v := range store.list {
		timeExists := v.DriverUid == driverUid && v.TrackId == trackId

		if !timeExists {
			continue
		}
		if v.Time > pTime {
			hotlap.Time = pTime
		} else {
			// Didn't beat his previous time
			return false, nil
		}
	}
	if (hotlap == Hotlap{}) {
		hotlap = Hotlap{
			TrackId:   trackId,
			DriverUid: driverUid,
			Time:      pTime,
		}
		DB.Create(&hotlap)
		store.list = append(store.list, hotlap)
	} else {
		hotlap.Time = pTime
		DB.Save(&hotlap)
	}
	fmt.Println("debug: StoreTime: added new time:", TrackStore.GetTrack(trackId).Name, driverUid, pTime.String())
	return true, nil
}

type ByDuration []Hotlap

func (laps ByDuration) Len() int           { return len(laps) }
func (laps ByDuration) Swap(i, j int)      { laps[i], laps[j] = laps[j], laps[i] }
func (laps ByDuration) Less(i, j int) bool { return laps[i].Time < laps[j].Time }

type ByTime []Hotlap

func (laps ByTime) Len() int           { return len(laps) }
func (laps ByTime) Swap(i, j int)      { laps[i], laps[j] = laps[j], laps[i] }
func (laps ByTime) Less(i, j int) bool { return laps[i].UpdatedAt.Before(laps[j].UpdatedAt) }

func (store *hotlapStore) GetTrackHotlapList(trackId string, limit int) []Hotlap {
	queryResult := make([]Hotlap, 0, limit)

	res := DB.Order("Time ASC").Limit(limit).Find(&queryResult, Hotlap{TrackId: trackId})

	if res.Error != nil {
		log.Printf("Error on DB Find:\n%v\n", res.Error)
	}
	return queryResult
}

func (store *hotlapStore) GetDriverHotlapList(pDriverUid string) []Hotlap {
	queryResult := make([]Hotlap, 0)

	res := DB.Order("UpdatedAt ASC").Find(&queryResult, Hotlap{DriverUid: pDriverUid})

	if res.Error != nil {
		log.Printf("Error on DB Find:\n%v\n", res.Error)
	}
	return queryResult
}

func (store *hotlapStore) GetDriverHottestLap(pDriverUid, pTrackId string) (Hotlap, error) {
	queryResult := make([]Hotlap, 0)

	res := DB.Order("Time ASC").Find(&queryResult, Hotlap{TrackId: pTrackId, DriverUid: pDriverUid})

	if res.Error != nil {
		log.Printf("Error on DB Find:\n%v\n", res.Error)
	}

	if len(queryResult) <= 0 {
		return Hotlap{}, fmt.Errorf("No hotlap set on track (%s) for driver (%s)", pTrackId, pDriverUid)
	}
	return queryResult[0], nil
}
