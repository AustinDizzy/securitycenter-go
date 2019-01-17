package sc

import (
	"encoding/json"
	"strconv"
	"time"
)

type statusResp struct {
	scResp
	Response Status
}

type statusUpdate struct {
	UpdateTime          time.Time `sc:"updateTime"`
	UpdateTimeStr       string    `json:"updateTime"`
	IsStale             bool      `sc:"stale"`
	Stale               string    `json:"stale"`
	PluginCurrentSet    time.Time `sc:"pluginCurrentSet" json:"-"`
	PluginCurrentSetStr string    `json:"pluginCurrentSet,omitempty"`
}

// Status https://docs.tenable.com/sccv/api/Status.html
type Status struct {
	Jobd                            string                   `json:"jobd" sc:"jobd"`
	LicenseStatus                   string                   `json:"licenseStatus" sc:"licenseStatus"`
	PluginSubscriptionStatus        string                   `json:"pluginSubscriptionStatus" sc:"pluginSubscriptionStatus"`
	LCEPluginSubscriptionStatus     string                   `json:"LCEPluginSubscriptionStatus" sc:"LCEPluginSubscriptionStatus"`
	PassivePluginSubscriptionStatus string                   `json:"passivePluginSubscriptionStatus" sc:"passivePluginSubscriptionStatus"`
	PluginUpdates                   map[string]*statusUpdate `json:"pluginUpdates" sc:"pluginUpdates"`
	FeedUpdates                     statusUpdate             `json:"feedUpdates" sc:"feedUpdates"`
	ActiveIPs                       int                      `json:"-" sc:"activeIPs"`
	ActiveIPsStr                    string                   `json:"activeIPs"`
	LicensedIPs                     int                      `json:"-" sc:"licensedIPs"`
	LicensedIPsStr                  string                   `json:"licensedIPs"`
	NoLCEs                          bool                     `json:"-" sc:"noLCEs"`
	NoLCEsStr                       string                   `json:"noLCEs"`
	NoReps                          bool                     `json:"-" sc:"noReps"`
	NoRepsStr                       string                   `json:"noReps"`
	Zones                           []Zone                   `json:"zones" sc:"zones"`
}

// GetStatus returns the current SecurityCenter status
// including license information
func (sc *SC) GetStatus() (*Status, error) {
	var (
		resp        statusResp
		status      *Status
		req         = sc.NewRequest("GET", "status")
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return &resp.Response, err
	}

	status = &resp.Response

	err = status.readAttr()

	return status, err
}

func (status *Status) readAttr() error {
	for _, set := range status.PluginUpdates {
		err := set.readAttr()
		if err != nil {
			return err
		}
	}

	err := status.FeedUpdates.readAttr()
	if err != nil {
		return err
	}

	status.ActiveIPs, err = strconv.Atoi(status.ActiveIPsStr)
	if err != nil {
		return err
	}

	status.LicensedIPs, err = strconv.Atoi(status.LicensedIPsStr)
	if err != nil {
		return err
	}

	err = readBool(&status.NoLCEsStr, &status.NoLCEs)
	if err != nil {
		return err
	}

	err = readBool(&status.NoRepsStr, &status.NoReps)

	return nil
}

func (s *statusUpdate) readAttr() error {
	updateTime, err := strconv.ParseInt(s.UpdateTimeStr, 10, 64)
	if err != nil {
		return err
	}
	s.UpdateTime = time.Unix(updateTime, 0)

	err = readBool(&s.Stale, &s.IsStale)
	if err != nil {
		return err
	}

	if len(s.PluginCurrentSetStr) > 0 {
		currentSet, err := strconv.ParseInt(s.PluginCurrentSetStr, 10, 64)
		if err != nil {
			return err
		}
		s.PluginCurrentSet = time.Unix(currentSet, 0)
	}
	return err
}
