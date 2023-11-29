package store

import (
	"sync"
)

type Track struct {
	Id      string
	Name    string
	Picture any
}

type trackStore struct {
	listMutex sync.Mutex
	list      map[string]Track
}

var TrackStore = trackStore{list: make(map[string]Track)}

func (store *trackStore) GetTrack(id string) Track {
	store.listMutex.Lock()
	defer store.listMutex.Unlock()
	return store.list[id]
}

func (store *trackStore) GetTracks() []Track {
	store.listMutex.Lock()
	defer store.listMutex.Unlock()
	trackList := make([]Track, 0, len(store.list))

	for _, v := range store.list {
		trackList = append(trackList, v)
	}
	return trackList
}

func (store *trackStore) SetTrack(track Track) bool {
	store.listMutex.Lock()
	defer store.listMutex.Unlock()
	_, trackExists := store.list[track.Id]
	store.list[track.Id] = track

	return trackExists
}
