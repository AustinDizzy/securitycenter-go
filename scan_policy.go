package sc

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type policyResp struct {
	scResp
	Response Policy
}

type policiesResp struct {
	scResp
	Response []Policy
}

// Policy https://docs.tenable.com/sccv/api/Scan-Policy.html
type Policy struct {
	ID              string    `json:"id" sc:"id"`
	Name            string    `json:"name" sc:"name"`
	Status          string    `json:"status" sc:"status"`
	Description     string    `json:"description" sc:"description"`
	CanUseStr       string    `json:"canUse"`
	CanUse          bool      `json:"-" sc:"canUse"`
	CanManageStr    string    `json:"canManage"`
	CanManage       bool      `json:"-" sc:"canManage"`
	CreatedTime     time.Time `json:"-" sc:"createdTime"`
	CreatedTimeStr  string    `json:"createdTime"`
	ModifiedTime    time.Time `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr string    `json:"modifiedTime"`
	Owner           User      `json:"owner" sc:"owner"`
	OwnerGroup      Group     `json:"ownerGroup" sc:"ownerGroup"`
	Groups          []Group   `json:"groups" sc:"groups"`
}

func (p *Policy) readAttr() error {
	err := readBool(&p.CanUseStr, &p.CanUse)
	if err != nil {
		return err
	}

	err = readBool(&p.CanManageStr, &p.CanManage)
	if err != nil {
		return err
	}

	err = readTime(&p.CreatedTimeStr, &p.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&p.ModifiedTimeStr, &p.ModifiedTime)
	if err != nil {
		return err
	}

	err = p.Owner.readAttr()
	if err != nil {
		return err
	}

	err = p.OwnerGroup.readAttr()
	if err != nil {
		return err
	}

	for _, g := range p.Groups {
		err = g.readAttr()
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPolicies gets a list of policies with the optionally user-supplied fields
// and an error if there was one
func (sc *SC) GetPolicies(fields ...string) ([]Policy, error) {
	var (
		resp        policiesResp
		policies    []Policy
		req         = prepPolicyRequest(sc, -1, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return nil, err
	}

	policies = resp.Response

	for _, p := range policies {
		err := p.readAttr()
		if err != nil {
			return nil, err
		}
	}

	return policies, err
}

// GetPolicy gets the policy associated with the supplied id
func (sc *SC) GetPolicy(id int, fields ...string) (*Policy, error) {
	var (
		resp        policyResp
		policy      Policy
		req         = prepPolicyRequest(sc, id, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return nil, err
	}

	policy = resp.Response

	err = policy.readAttr()

	return &policy, err
}

func prepPolicyRequest(sc *SC, id int, fields ...string) *Request {
	var path string
	if id < 0 {
		path = "policy"
	} else {
		path = fmt.Sprintf("policy/%d", id)
	}

	req := sc.NewRequest("GET", path)

	if len(fields) > 0 {
		req.data["fields"] = strings.Join(fields, ",")
	}

	return req
}
