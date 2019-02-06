package sc

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"
)

type zonesResp struct {
	scResp
	Response []Zone
}

type zoneResp struct {
	scResp
	Response Zone
}

// Zone https://docs.tenable.com/sccv/api/Scan-Zone.html
type Zone struct {
	ID              string       `json:"id" sc:"id"`
	Name            string       `json:"name" sc:"name"`
	Description     string       `json:"description" sc:"description"`
	Status          string       `json:"status" sc:"status"`
	IPList          []*net.IPNet `json:"-" sc:"ipList"`
	IPListStr       string       `json:"ipList"`
	CreatedTime     time.Time    `json:"-" sc:"createdTime"`
	CreatedTimeStr  string       `json:"createdTime"`
	ModifiedTime    time.Time    `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr string       `json:"modifiedTime"`
	Scanners        []Scanner    `json:"scanners" sc:"scanners"`
	ActiveScanners  int          `json:"activeScanners" sc:"activeScaners"`
	TotalScanners   int          `json:"totalScanners" sc:"totalScanners"`
}

// todo: finish Zones, Scanner, Scan, Repository, other key objects, and basic editing

// GetZones retrieves a slice of Zones with either the user-supplied fields or their default values
func (sc *SC) GetZones(fields ...string) ([]Zone, error) {
	var (
		resp   zonesResp
		zones  []Zone
		req    = sc.NewRequest("GET", "zone")
		scResp *Response
		err    error
	)

	if len(fields) > 0 {
		req.data["fields"] = strings.Join(fields, ",")
	}

	scResp, err = req.Do()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return resp.Response, err
	}

	for _, z := range resp.Response {
		err = z.readAttr()
		if err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}

	return zones, err
}

func (z *Zone) readAttr() error {
	// read list of IP networks
	for _, ipNetStr := range strings.Split(z.IPListStr, ",") {
		_, ipNet, err := net.ParseCIDR(ipNetStr)
		if err != nil {
			return err
		}
		z.IPList = append(z.IPList, ipNet)
	}

	createdTime, err := strconv.ParseInt(z.CreatedTimeStr, 10, 64)
	if err != nil {
		return err
	}
	z.CreatedTime = time.Unix(createdTime, 0)

	modifiedTime, err := strconv.ParseInt(z.ModifiedTimeStr, 10, 64)
	if err != nil {
		return err
	}
	z.ModifiedTime = time.Unix(modifiedTime, 0)

	for _, s := range z.Scanners {
		if err := (&s).readAttr(); err != nil {
			return err
		}
	}

	return nil
}
