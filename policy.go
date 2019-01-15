package sc

import (
	"fmt"
	"strings"
	"time"
)

const defaultPolicyFields = "id,name,description,status,policyTemplate,policyProfileName,creator,tags,createdTime,modifiedTime,context,generateXCCDFResults,auditFiles,preferences,targetGroup,owner,ownerGroup,groups,families"

type Policy struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Status       int       `json:"status"`
	Description  string    `json:"description"`
	CanUse       bool      `json:"canUse"`
	CanManage    bool      `json:"canManage"`
	CreatedTime  time.Time `json:"createdTime"`
	ModifiedTime time.Time `json:"modifiedTime"`
	Owner        User      `json:"_"`
	OwnerGroup   Group     `json:"ownerGroup"`
	Groups       []Group   `json:"groups"`
}

func (sc *SC) GetPolicy(id int, fields ...string) (Policy, error) {
	var (
		policy Policy
		path   = "policy"
	)

	if id >= 0 {
		path = fmt.Sprintf("policy/%d", id)
	}
	req := sc.NewRequest("GET", path)
	if len(fields) > 0 {
		req.data["fields"] = strings.Join(fields, ",")
	} else {
		req.data["fields"] = defaultPolicyFields
	}

	// resp, err := req.Do()
	// parseJSONToUser(resp.Data.GetPath("response", "owner")

	return policy, nil
}
