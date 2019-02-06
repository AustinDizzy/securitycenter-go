package sc

import (
	"net"
	"time"
)

// LCE https://docs.tenable.com/sccv/api/LCE.html
type LCE struct {
	ID                string         `json:"id" sc:"id"`
	Name              string         `json:"name" sc:"name"`
	Description       string         `json:"description" sc:"description"`
	IP                net.IP         `json:"-" sc:"ip"`
	IPStr             string         `json:"ip"`
	NTPIP             net.IP         `json:"-" sc:"ntpIP"`
	NTPIPStr          string         `json:"ntpIP"`
	Port              string         `json:"port" sc:"port"`
	Username          string         `json:"username" sc:"username"`
	Password          string         `json:"password" sc:"password"`
	ManagedRanges     string         `json:"managedRanges" sc:"managedRanges"`
	Version           string         `json:"version" sc:"version"`
	DownloadVulnsStr  string         `json:"downloadVulns"`
	DownloadVulns     bool           `json:"-" sc:"downloadVulns"`
	Status            string         `json:"status" sc:"status"`
	VulnStatus        string         `json:"vulnStatus" sc:"vulnStatus"`
	LastReportTime    time.Time      `json:"-" sc:"lastReportTime"`
	LastReportTimeStr string         `json:"lastReportTime"`
	CreatedTime       time.Time      `json:"-" sc:"createdTime"`
	CreatedTimeStr    string         `json:"createdTime"`
	ModifiedTime      time.Time      `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr   string         `json:"modifiedTime"`
	Organizations     []Organization `json:"organizations" sc:"organizations"`
	CanUse            bool           `json:"-" sc:"canUse"`
	CanUseStr         string         `json:"canUse"`
	CanManage         bool           `json:"-" sc:"canManage"`
	CanManageStr      string         `json:"canManage"`
}

func (l *LCE) readAttr() error {
	// read net.IPs
	l.IP = net.ParseIP(l.IPStr)

	l.NTPIP = net.ParseIP(l.NTPIPStr)

	// read Booleans
	err := readBool(&l.DownloadVulnsStr, &l.DownloadVulns)
	if err != nil {
		return err
	}

	err = readBool(&l.CanUseStr, &l.CanUse)
	if err != nil {
		return err
	}

	err = readBool(&l.CanManageStr, &l.CanManage)
	if err != nil {
		return err
	}

	// read time.Time
	err = readTime(&l.LastReportTimeStr, &l.LastReportTime)
	if err != nil {
		return err
	}

	err = readTime(&l.CreatedTimeStr, &l.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&l.ModifiedTimeStr, &l.ModifiedTime)
	if err != nil {
		return err
	}

	// read nested types
	for _, o := range l.Organizations {
		err := o.readAttr()
		if err != nil {
			return err
		}
	}

	return nil
}
