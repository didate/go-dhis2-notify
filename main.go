package main

import (
	"fmt"
	"os"
	"time"

	"github.com/didate/go-dhis2-notify/api"
	"github.com/didate/go-dhis2-notify/dao"
	"github.com/didate/go-dhis2-notify/mail"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	ctLayout = "2006-01-02T15:04:05"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic("failed to load env variables")
	}

	users := api.Fetch(os.Getenv("DHIS2_URL"), os.Getenv("DHIS2_AUTH"))

	db, err := gorm.Open(sqlite.Open(os.Getenv("BD_FILE")), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	dao.InitDb(users, db, os.Getenv("BD_FILE"))

	c := gocron.NewScheduler(time.UTC)

	c.Every(30).Seconds().Do(func() {

		users = api.Fetch(os.Getenv("DHIS2_URL"), os.Getenv("DHIS2_AUTH"))
		dao.UpdateDb(users, db)
		notifier(db)

	})

	c.StartBlocking()
}

func notifier(db *gorm.DB) {

	tx := db.Begin()

	mBody := "Nouveaux comptes\n\n"

	users := dao.NewerUsers(db)
	createdAccount := len(users)
	if createdAccount > 0 {
		for _, u := range users {
			mBody += "Login : " + u.Username + "\n"
			mBody += "Nom : " + u.DisplayName + "\n"
			mBody += "Telphone : " + u.PhoneNumber + "\n"
			mBody += "Crée Par : " + u.CreatedBy + "\n"
			mBody += "Crée Le : " + u.Created.Format(ctLayout) + "\n"
			mBody += "______________________\n\n"

			dao.UpdateTrack(u, "track_newer", false, tx)
		}
	}

	mBody += "Comptes modifiés\n\n"

	users = dao.UpdatedUsers(db)
	updatedAccount := len(users)
	if updatedAccount > 0 {
		for _, u := range users {
			mBody += "Login : " + u.Username + "\n"
			mBody += "Nom : " + u.DisplayName + "\n"
			mBody += "Telphone : " + u.PhoneNumber + "\n"
			mBody += "Modifiée Par : " + u.LastUpdatedBy + "\n"
			mBody += "Modifiée Le : " + u.LastUpdated.Format(ctLayout) + "\n"
			mBody += "______________________\n\n"
			dao.UpdateTrack(u, "track_updated", false, tx)

		}
	}
	fmt.Printf("%d updated users \n", updatedAccount+createdAccount)
	if updatedAccount+createdAccount > 0 {
		err := mail.Send("Situation des comptes DHIS2", mBody)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Error when sending email %v", err)
		} else {
			tx.Commit()
		}
	}
}
