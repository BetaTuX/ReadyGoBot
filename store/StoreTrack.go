package store

import (
	"database/sql/driver"
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type PictureAttachment discordgo.MessageAttachment

type Track struct {
	gorm.Model
	TrackId string
	Name    string
	Picture PictureAttachment
}

func (msgAttachment PictureAttachment) Value() (driver.Value, error) {
	if value, error := json.Marshal(msgAttachment); error == nil {
		return string(value), nil
	} else {
		return nil, error
	}
}

func (msgAttachment *PictureAttachment) Scan(value interface{}) error {
	if value == nil {
		// set the value of the pointer yne to YesNoEnum(false)
		*msgAttachment = PictureAttachment{}
		return nil
	}
	return json.Unmarshal([]byte(value.(string)), msgAttachment)
}

type trackStore struct{}

var TrackStore = trackStore{}

func (store *trackStore) GetTrack(id string) Track {
	queryResult := make([]Track, 0, 1)

	res := DB.Limit(1).Find(&queryResult, Track{TrackId: id})

	if res.Error != nil {
		log.Printf("Error on DB Find:\n%v\n", res.Error)
	}
	return queryResult[0]
}

func (store *trackStore) GetTracks() []Track {
	queryResult := make([]Track, 0)

	res := DB.Find(&queryResult)

	if res.Error != nil {
		log.Printf("Error on DB Find:\n%v\n", res.Error)
	}
	return queryResult
}

func (store *trackStore) SetTrack(track Track) bool {

	queryResult := make([]Track, 0, 1)

	res := DB.Limit(1).Find(&queryResult, track)

	if res.Error != nil {
		log.Printf("Error on DB Find:\n%v\n", res.Error)
	}

	trackAlreadyExists := !(len(queryResult) <= 0)

	if !trackAlreadyExists {
		DB.Create(&track)
	} else {
		DB.Save(&track)
	}
	return trackAlreadyExists
}
