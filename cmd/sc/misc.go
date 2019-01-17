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

	status, err := s.GetStatus()
	if err != nil {
		panic(err)
	}

	output := termtables.CreateTable()
	output.AddTitle("SecurityCenter Status")

	output.AddRow("Current Time", time.Now())
	output.AddSeparator()
	output.AddRow("Jobd Status", status.Jobd)
	output.AddRow("Plugin Subscription Status", status.PluginSubscriptionStatus)
	output.AddRow("License Status", status.LicenseStatus)
	output.AddSeparator()
	output.AddRow("Licensed IPs", fmt.Sprintf("%d working / %d licensed", status.ActiveIPs, status.LicensedIPs))
	output.AddRow("Last Feed Update", time.Since(status.FeedUpdates.UpdateTime))
	output.AddRow("Feed Stale", status.FeedUpdates.IsStale)

	fmt.Println(output.Render())
	fmt.Println()

	zones := termtables.CreateTable()
	zones.AddTitle(fmt.Sprintf("SecurityCenter Zones (%d)", len(status.Zones)))
	zones.AddHeaders("#", "Name", "Status", "Description")

	for _, z := range status.Zones {
		zones.AddRow(z.ID, z.Name, z.Status, z.Description)
	}

	fmt.Println(zones.Render())
}
