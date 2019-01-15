package main

import (
	"fmt"

	"github.com/austindizzy/securitycenter-go"
)

func doAuth(s *sc.SC) {
	if len(username) == 0 || len(password) == 0 {
		for len(username) == 0 {
			user, err := readInput("Username:")
			if err != nil {
				panic(err)
			}
			if len(user) > 0 {
				username = user
			}
		}
		for len(password) == 0 {
			pass, err := readSecretInput("Password:", true)
			if err != nil {
				panic(err)
			}
			if len(pass) > 0 {
				password = pass
			}
		}
	}
	success, err := s.DoAuth(username, password)
	if err != nil {
		panic(err)
	}
	if success {
		fmt.Printf("auth:\n\ttoken: %s\n\tsession: %s", s.Keys["token"], s.Keys["session"])
	} else {
		fmt.Printf("auth: FAILED")
	}
}
