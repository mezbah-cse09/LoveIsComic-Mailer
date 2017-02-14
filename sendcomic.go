package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jordan-wright/email"
)

const (
	loveIsComicHostName = "http://loveiscomix.com"
)

func main() {

	fromEmail := flag.String("from", "", "gmail address of Sender")
	toEmail := flag.String("to", "", "gmail address of Recipient")
	password := flag.String("password", "", "password for Sender account")

	flag.Parse()
	if *fromEmail == "" {
		handleError(errors.New("Please provide the gmail address of the Sender"))
	}
	if *toEmail == "" {
		handleError(errors.New("Please provide the gmail address of the Recipient"))
	}
	if *password == "" {
		handleError(errors.New("Please provide the password for the Sender account"))
	}

	//HTTP GET for a Random Image from Love Is Comic
	resp, err := http.Get(loveIsComicHostName + "/random/")
	handleError(err)

	//Build out a jQuery like object
	doc, err := goquery.NewDocumentFromResponse(resp)
	handleError(err)

	//Search for the image
	doc.Find(".comiccell .comicbox a img").Each(func(i int, s *goquery.Selection) {
		val, b := s.Attr("src")
		if b != false {
			t := time.Now()
			loveIsURL := loveIsComicHostName + val
			e := email.NewEmail()
			e.From = *fromEmail
			e.To = []string{*toEmail}
			e.Subject = "Love Is Comic - " + t.Format("Mon Jan _2 2006")
			e.HTML = []byte("<img src='" + loveIsURL + "'/><hr/>Powered by <a href='" + loveIsComicHostName + "'>Love Is Comic</a>")
			auth := smtp.PlainAuth("", *fromEmail, *password, "smtp.gmail.com")
			err := e.Send("smtp.gmail.com:587", auth)
			if err != nil {
				handleError(err)
			}
		}
	})
}

//Error Handler. Currently just does prints out the error message
func handleError(err error) {
	if err != nil {
		log.Println(err.Error())
		return
	}
}
