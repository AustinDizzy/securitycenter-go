package sc

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// Repository https://docs.tenable.com/sccv/api/Repository.html
type Repository struct {
	ID                string               `json:"id" sc:"id"`
	Name              string               `json:"name" sc:"name"`
	Description       string               `json:"description" sc:"description"`
	Type              string               `json:"type" sc:"type"`
	DataFormat        string               `json:"dataFormat" sc:"dataFormat"`
	VulnCount         int                  `json:"-" sc:"vulnCount"`
	VulnCountStr      string               `json:"vulnCount"`
	RemoteID          string               `json:"remoteID" sc:"remoteID"`
	RemoteIP          string               `json:"remoteIP" sc:"remoteIP"`
	IsRunning         bool                 `json:"-" sc:"running"`
	Running           string               `json:"running"`
	DownloadFormat    string               `json:"downloadFormat" sc:"downloadFormat"`
	LastSyncTime      time.Time            `json:"-" sc:"lastSyncTime"`
	LastSyncTimeStr   string               `json:"lastSyncTime"`
	LastVulnUpdate    time.Time            `json:"-" sc:"lastVulnUpdate"`
	LastVulnUpdateStr string               `json:"lastVulnUpdate"`
	CreatedTime       time.Time            `json:"-" sc:"createdTime"`
	CreatedTimeStr    string               `json:"createdTime"`
	ModifiedTime      time.Time            `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr   string               `json:"modifiedTime"`
	Organizations     []Organization       `json:"organization" sc:"organizations"`
	TypeFields        repositoryTypeFields `json:"typeFields" sc:"typeFields"`
}

type repositoryTypeFields struct {
	LastVulnUpdate            time.Time           `json:"-" sc:"lastVulnUpdate"`
	LastVulnUpdateStr         string              `json:"lastVulnUpdate"`
	VulnCount                 int                 `json:"-" sc:"vulnCount"`
	VulnCountStr              string              `json:"vulnCount"`
	NessusSchedule            map[string]string   `json:"nessusSchedule" sc:"nessusSchedule"`
	Correlation               []map[string]string `json:"correlation" sc:"correlation"`
	IPRange                   []*net.IPNet        `json:"-" sc:"ipRange"`
	IPRangeStr                string              `json:"ipRange"`
	IPCount                   int                 `json:"-" sc:"ipCount"`
	IPCountStr                string              `json:"ipCount"`
	RunningNessusStr          string              `json:"runningNessus"`
	RunningNessus             bool                `json:"-" sc:"runningNessus"`
	LastGenerateNessusTime    time.Time           `json:"-" sc:"lastGenerateNessusTime"`
	LastGenerateNessusTimeStr string              `json:"lastGenerateNessusTime"`
	TrendingDaysStr           string              `json:"trendingDays"`
	TrendingDays              int                 `json:"-" sc:"trendingDays"`
	TrendWithRawStr           string              `json:"trendWithRaw"`
	TrendWithRaw              bool                `json:"-" sc:"trendWithRaw"`
}

func (r *Repository) readAttr() error {
	var err error
	r.VulnCount, err = strconv.Atoi(r.VulnCountStr)
	if err != nil {
		return err
	}

	err = readBool(&r.Running, &r.IsRunning)
	if err != nil {
		return err
	}

	err = readTime(&r.LastSyncTimeStr, &r.LastSyncTime)
	if err != nil {
		return err
	}

	err = readTime(&r.LastVulnUpdateStr, &r.LastVulnUpdate)
	if err != nil {
		return err
	}

	err = readTime(&r.CreatedTimeStr, &r.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&r.ModifiedTimeStr, &r.ModifiedTime)
	if err != nil {
		return err
	}

	for _, o := range r.Organizations {
		err := o.readAttr()
		if err != nil {
			return err
		}
	}

	err = r.TypeFields.readAttr()
	if err != nil {
		return err
	}

	return nil
}

func (f *repositoryTypeFields) readAttr() error {
	// read time.Time
	err := readTime(&f.LastVulnUpdateStr, &f.LastVulnUpdate)
	if err != nil {
		return err
	}

	err = readTime(&f.LastGenerateNessusTimeStr, &f.LastGenerateNessusTime)
	if err != nil {
		return err
	}

	// read integers
	f.VulnCount, err = strconv.Atoi(f.VulnCountStr)
	if err != nil {
		return err
	}

	f.IPCount, err = strconv.Atoi(f.IPCountStr)
	if err != nil {
		return err
	}

	f.TrendingDays, err = strconv.Atoi(f.TrendingDaysStr)
	if err != nil {
		return err
	}

	// read Booleans
	err = readBool(&f.RunningNessusStr, &f.RunningNessus)
	if err != nil {
		return err
	}

	err = readBool(&f.TrendWithRawStr, &f.TrendWithRaw)
	if err != nil {
		return err
	}

	// read list of IP nets
	var nets []*net.IPNet
	for _, cidr := range strings.Split(f.IPRangeStr, ",") {
		_, n, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		nets = append(nets, n)
	}
	f.IPRange = nets

	return nil
}
