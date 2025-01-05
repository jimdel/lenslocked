package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jimdel/lenslocked/models"
	"github.com/joho/godotenv"
)

func main() {
	// Grab env vars
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
	SMPT_HOST := os.Getenv("SMTP_HOST")
	SMPT_USERNAME := os.Getenv("SMTP_USERNAME")
	SMPT_PASSWORD := os.Getenv("SMTP_PASSWORD")
	SMPT_PORTSTR := os.Getenv("SMTP_PORT")
	SMTP_PORT, err := strconv.Atoi(SMPT_PORTSTR)
	if err != nil {
		panic(err)
	}
	//END - Grab env vars
	es, err := models.NewEmailService(models.SMTPConfig{
		Host:     SMPT_HOST,
		Port:     SMTP_PORT,
		Username: SMPT_USERNAME,
		Password: SMPT_PASSWORD,
	})
	if err != nil {
		panic(err)
	}

	err = es.ForgotPassword("jimdel@jim.com", "http://lenslocked.com/reset-pw?token=abc123")
	if err != nil {

		panic(err)
	}
	fmt.Println("sent email")
}
