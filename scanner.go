package sc

import (
	"net"
	"time"
)

// Scanner https://docs.tenable.com/sccv/api/Scanner.html
type Scanner struct {
	ID                string    `json:"id" sc:"id"`
	Name              string    `json:"name" sc:"name"`
	Description       string    `json:"description" sc:"description"`
	Status            string    `json:"status" sc:"status"`
	IP                net.IP    `json:"-" sc:"ip"`
	IPStr             string    `json:"ip"`
	Port              int       `json:"-" sc:"port"`
	PortStr           string    `json:"port"`
	UseProxy          bool      `json:"-" sc:"useProxy"`
	UseProxyStr       string    `json:"useProxy"`
	Enabled           bool      `json:"-" sc:"enabled"`
	EnabledStr        string    `json:"enabled"`
	VerifyHost        bool      `json:"-" sc:"verifyHost"`
	VerifyHostStr     string    `json:"verifyHost"`
	ManagePlugins     bool      `json:"-" sc:"managePlugins"`
	ManagePluginsStr  string    `json:"managePlugins"`
	AuthType          string    `json:"authType" sc:"authType"`
	Username          string    `json:"username" sc:"username"`
	Password          string    `json:"password" sc:"password"`
	Admin             bool      `json:"-" sc:"admin"`
	AdminStr          string    `json:"admin"`
	MSP               bool      `json:"-" sc:"msp"`
	MSPStr            string    `json:"msp"`
	NumScans          int       `json:"-" sc:"numScans"`
	NumScansStr       string    `json:"numScans"`
	NumHosts          int       `json:"-" sc:"numHosts"`
	NumHostsStr       string    `json:"numHosts"`
	NumSessions       int       `json:"-" sc:"numSessions"`
	NumSessionsStr    string    `json:"numSessions"`
	NumTCPSessions    int       `json:"-" sc:"numTCPSessions"`
	NumTCPSessionsStr string    `json:"numTCPSessions"`
	CreatedTime       time.Time `json:"-" sc:"createdTime"`
	CreatedTimeStr    string    `json:"createdTime"`
	ModifedTime       time.Time `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr   string    `json:"modifiedTime"`
	Zones             []Zone    `json:"zones" sc:"zones"`
}

func (s *Scanner) readAttr() error {
	// read net.IPs
	s.IP = net.ParseIP(s.IPStr)

	// read Booleans
	err := readBool(&s.UseProxyStr, &s.UseProxy)
	if err != nil {
		return err
	}

	err = readBool(&s.EnabledStr, &s.Enabled)
	if err != nil {
		return err
	}

	err = readBool(&s.VerifyHostStr, &s.VerifyHost)
	if err != nil {
		return err
	}

	err = readBool(&s.ManagePluginsStr, &s.ManagePlugins)
	if err != nil {
		return err
	}

	err = readBool(&s.AdminStr, &s.Admin)
	if err != nil {
		return err
	}

	err = readBool(&s.MSPStr, &s.MSP)
	if err != nil {
		return err
	}

	return nil
}
