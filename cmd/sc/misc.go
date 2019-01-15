package main

import (
	"fmt"
	"time"

	"github.com/apcera/termtables"

	"github.com/austindizzy/securitycenter-go"
)

func doWhoAmI(s *sc.SC) {
	if !s.HasAuth() {
		fmt.Println("no valid authentication found")
		return
	}

	user, err := s.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n%s\n", user)
}

func doStatus(s *sc.SC) {
	if !s.HasAuth() {
		fmt.Println("no valid authentication found")
		return
	}

	req := s.NewRequest("GET", "status")
	res, err := req.Do()
	if err != nil {
		panic(err)
	}

	output := termtables.CreateTable()
	output.AddTitle("SecurityCenter Status")

	output.AddRow("Current Time", time.Now())
	output.AddSeparator()
	output.AddRow("Jobd Status", res.Data.Get("response").Get("jobd").MustString("Not Running"))
	output.AddRow("Plugin Subscription Status", res.Data.Get("response").Get("PluginSubscriptionStatus").MustString())
	output.AddRow("License Status", res.Data.Get("response").Get("licenseStatus").MustString())
	output.AddSeparator()
	output.AddRow("Licensed IPs", fmt.Sprintf("%s working / %s licensed", res.Data.GetPath("response", "activeIPs").MustString(), res.Data.GetPath("response", "licensedIPs").MustString()))
	t := time.Unix(res.Data.GetPath("response", "feedUpdates", "updateTime").MustInt64(), 0)
	output.AddRow("Last Feed Update", time.Since(t))
	output.AddRow("Feed Stale", res.Data.GetPath("response", "feedUpdates", "stale").MustBool(false))

	fmt.Println(output.Render())
	fmt.Println()

	zones := termtables.CreateTable()
	zones.AddTitle(fmt.Sprintf("SecurityCenter Zones (%d)", len(res.Data.GetPath("response", "zones").MustArray())))
	zones.AddHeaders("#", "Name", "Status", "Description")

	for i := range res.Data.GetPath("response", "zones").MustArray() {
		z := res.Data.GetPath("response", "zones").GetIndex(i)
		zones.AddRow(z.Get("id").MustString(), z.Get("name").MustString(), z.Get("status").MustString(), z.Get("description").MustString())
	}

	fmt.Println(zones.Render())
}
