package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apcera/termtables"
	"github.com/austindizzy/securitycenter-go"
)

func doUsers(s *sc.SC) {
	var (
		userFields []string
		users, err = s.GetUsers()
		output     *termtables.Table
	)

	if err != nil {
		panic(err)
	}

	output = termtables.CreateTable()
	output.AddTitle("List of SecurityCenter Users")
	if len(fields) > 0 {
		var headers []interface{}
		for _, v := range strings.Split(fields, ",") {
			userFields = append(userFields, v)
			headers = append(headers, v)
		}
		output.AddHeaders(headers...)
	} else {
		output.AddHeaders("#", "username", "email", "role", "group", "name", "lastLogin", "lastLoginIP")
	}

	for i := 0; i < len(users); i++ {
		var (
			u   = users[i]
			ll  = "never"
			row []interface{}
		)
		if len(userFields) == 0 || strSliceContains(userFields, "lastLogin") {
			if u.LastLogin.Year() > 2001 {
				ll = u.LastLogin.Format(time.UnixDate)
			}
		}

		if len(userFields) == 0 {
			row = []interface{}{u.ID, u.Username, u.Email, u.Role.Name, u.Group.Name, fmt.Sprint(u.FirstName, " ", u.LastName), ll, u.LastLoginIP}
		} else {
			for _, name := range userFields {
				for i := 0; i < reflect.TypeOf(u).NumField(); i++ {
					field := reflect.TypeOf(u).Field(i)
					if alias, ok := field.Tag.Lookup("sc"); ok && alias == name {
						if !ok {
							row = append(row, nil)
						} else {
							val := reflect.ValueOf(u)
							row = append(row, reflect.Indirect(val).Field(i).Interface())
						}
					}
				}
			}
		}

		output.AddRow(row...)
	}

	output.AddSeparator()
	output.AddRow(fmt.Sprintf("%d total users", len(users)))

	fmt.Println(output.Render())
}

func doUser(s *sc.SC) {
	var (
		user       *sc.User
		id         int
		orgID      int
		err        error
		userFields []string
	)
	if len(fields) > 0 {
		userFields = strings.Split(fields, ",")
	}
	switch len(args) {
	default:
		fmt.Println("sc `user` subcommand\n\t* ID paramter (i.e. `sc user 4`)\n\t* OrgID and UserID (i.e. `sc user 2 4` <=> `sc user {orgID} {userID}`)")
		break
	case 2:
		id, err = strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		if len(userFields) > 0 {
			user, err = s.GetUser(id, userFields...)
		} else {
			user, err = s.GetUser(id)
		}
		break
	case 3:
		orgID, err = strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		id, err = strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}

		if len(userFields) > 0 {
			user, err = s.GetUserFromOrg(orgID, id, userFields...)
		} else {
			user, err = s.GetUserFromOrg(orgID, id)
		}

		break
	}

	if err != nil {
		panic(err)
	}

	if user != nil {
		fmt.Println(user)
	}
}
