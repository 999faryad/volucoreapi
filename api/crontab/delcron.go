package crontab

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func DeleteCron() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 * * * * *", func() {
		fmt.Printf("Hello, World! %s\n", time.Now().Format(time.RFC3339))
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()

	// Läuft für immer, oder du kannst auch eine bestimmte Dauer einstellen.
	select {}
}
