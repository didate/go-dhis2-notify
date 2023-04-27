package dao

import (
	"fmt"

	"github.com/didate/go-dhis2-notify/api"
	"github.com/didate/go-dhis2-notify/tools"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserDb struct {
	ID            string `gorm:"primarykey"`
	Username      string
	DisplayName   string
	PhoneNumber   string
	LastUpdated   tools.CustomTime
	Created       tools.CustomTime
	CreatedBy     string
	LastUpdatedBy string
	Roles         string
	TrackNewer    bool
	TrackUpdated  bool
}

func UpdateDb(users []api.User, db *gorm.DB) error {
	fmt.Println("Updating database")
	for _, u := range users {
		userDb := UserDb{
			ID:            u.ID,
			Username:      u.UserCredential.Username,
			DisplayName:   u.DisplayName,
			PhoneNumber:   u.PhoneNumber,
			LastUpdated:   u.LastUpdated,
			LastUpdatedBy: u.UserCredential.LastUpdatedBy.Username,
			Created:       u.Created,
			CreatedBy:     u.UserCredential.CreatedBy.Username,
		}

		var userFromDb UserDb
		db.Find(&userFromDb, "ID =?", userDb.ID)

		if userFromDb.ID == "" {
			userDb.TrackNewer = true
			db.Create(&userDb)
		} else {
			userDb.TrackUpdated = true
			db.Model(&userDb).Where("last_updated < ?", userDb.LastUpdated).Updates(&userDb)
		}

	}
	fmt.Println("Updating database finished")
	return nil
}

func InitDb(users []api.User, db *gorm.DB, pDb string) int64 {
	db, err := gorm.Open(sqlite.Open(pDb), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&UserDb{})

	result := db.Find(&UserDb{})
	if result.RowsAffected == 0 {

		for _, u := range users {
			userDb := UserDb{
				ID:            u.ID,
				Username:      u.UserCredential.Username,
				DisplayName:   u.DisplayName,
				PhoneNumber:   u.PhoneNumber,
				LastUpdated:   u.LastUpdated,
				LastUpdatedBy: u.UserCredential.LastUpdatedBy.Username,
				Created:       u.Created,
				CreatedBy:     u.UserCredential.CreatedBy.Username,
				TrackNewer:    false,
				TrackUpdated:  false,
			}
			result = db.Create(userDb)
		}

		return result.RowsAffected
	}

	return 0
}

func UpdateTrack(u UserDb, cUpdate string, v bool, db *gorm.DB) int64 {
	result := db.Model(&u).Where("id = ? ", u.ID).Update(cUpdate, v)
	return result.RowsAffected
}

func NewerUsers(db *gorm.DB) []UserDb {
	var users []UserDb
	db.Find(&users, "track_newer = ? ", 1)
	return users
}

func UpdatedUsers(db *gorm.DB) []UserDb {
	var users []UserDb
	db.Find(&users, "track_updated = ? ", 1)
	return users
}
