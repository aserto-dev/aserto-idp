package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"gopkg.in/auth0.v5/management"
)

func main() {
	domain := flag.String("domain", "", "auth0 domain")
	clientId := flag.String("client_id", "", "auth0 client_id")
	clientSecret := flag.String("client_secret", "", "auth0 client_secret")
	deleteEmailDomain := flag.String("delete_email", "", "delete users with this email")
	flag.Parse()

	if *deleteEmailDomain == "" {
		log.Fatal("email domain required")
	}

	mgnt, err := management.New(
		*domain,
		management.WithClientCredentials(
			*clientId,
			*clientSecret,
		))
	if err != nil {
		log.Println(err)
		return
	}
	page := 0
	var users []string
	for {
		opts := management.Page(page)
		ul, err := mgnt.User.List(opts)
		if err != nil {
			log.Println(err)
		}

		for _, u := range ul.Users {
			if strings.Contains(*u.Email, *deleteEmailDomain) {
				users = append(users, *u.ID)
			}
		}
		if ul.Length < ul.Limit {
			break
		}
		page++
	}

	for _, u := range users {
		err := mgnt.User.Delete(u)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
