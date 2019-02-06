package sc

import (
	"strconv"
	"time"
)

// Organization https://docs.tenable.com/sccv/api/Organization.html
type Organization struct {
	ID                string              `json:"id" sc:"id"`
	Name              string              `json:"name" sc:"name"`
	Description       string              `json:"description" sc:"description"`
	Email             string              `json:"email" sc:"email"`
	Address           string              `json:"address" sc:"address"`
	City              string              `json:"city" sc:"city"`
	State             string              `json:"state" sc:"state"`
	Country           string              `json:"country" sc:"country"`
	Phone             string              `json:"phone" sc:"phone"`
	Fax               string              `json:"fax" sc:"fax"`
	IPinfoLinks       []map[string]string `json:"ipInfoLinks" sc:"ipInfoLinks"`
	ZoneSelection     string              `json:"zoneSelection" sc:"zoneSelection"`
	RestrictedIPs     string              `json:"restrictedIPs" sc:"restrictedIPs"`
	VulnScoreLow      string              `json:"vulnScoreLow" sc:"vulnScoreLow"`
	VulnScoreMedium   string              `json:"vulnScoreMedium" sc:"vulnScoreMedium"`
	VulnScoreHigh     string              `json:"vulnScoreHigh" sc:"vulnScoreHigh"`
	VulnScoreCritical string              `json:"vulnScoreCritical" sc:"vulnScoreCritical"`
	CreatedTime       time.Time           `json:"-" sc:"createdTime"`
	CreatedTimeStr    string              `json:"createdTime"`
	ModifiedTime      time.Time           `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr   string              `json:"modifiedTime"`
	UserCount         int                 `json:"-" sc:"userCount"`
	UserCountStr      string              `json:"userCount"`
	LCEs              []LCE               `json:"lces" sc:"lces"`
}

func (o *Organization) readAttr() error {
	err := readTime(&o.CreatedTimeStr, &o.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&o.ModifiedTimeStr, &o.ModifiedTime)
	if err != nil {
		return err
	}

	o.UserCount, err = strconv.Atoi(o.UserCountStr)
	if err != nil {
		return err
	}

	return nil
}
