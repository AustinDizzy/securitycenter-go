package sc

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type roleResp struct {
	scResp
	Response Role
}

type rolesResp struct {
	scResp
	Response []Role
}

// Role https://docs.tenable.com/sccv/api/Role.html
type Role struct {
	ID                           string              `json:"id" sc:"id"`
	Name                         string              `json:"name" sc:"name"`
	Description                  string              `json:"description" sc:"description"`
	Creator                      User                `json:"creator" sc:"creator"`
	CreatedTime                  time.Time           `json:"-" sc:"createdTime"`
	CreatedTimeStr               string              `json:"createdTime"`
	ModifiedTime                 time.Time           `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr              string              `json:"modifiedTime"`
	PermManageApp                string              `json:"permManageApp"`
	PermManageGroups             string              `json:"permManageGroups"`
	PermManageRoles              string              `json:"permManageRoles"`
	PermManageImages             string              `json:"permManageImages"`
	PermManageGroupRelationships string              `json:"permManageGroupRelationships"`
	PermManageBlackoutWindows    string              `json:"permManageBlackoutWindows"`
	PermManageAttributeSets      string              `json:"permManageAttributeSets"`
	PermCreateTickets            string              `json:"permCreateTickets"`
	PermCreateAlerts             string              `json:"permCreateAlerts"`
	PermCreateAuditFiles         string              `json:"permCreateAuditFiles"`
	PermCreateLDAPAssets         string              `json:"permCreateLDAPAssets"`
	PermCreatePolicies           string              `json:"permCreatePolicies"`
	PermPurgeTickets             string              `json:"permPurgeTickets"`
	PermPurgeScanResults         string              `json:"permPurgeScanResults"`
	PermPurgeReportResults       string              `json:"permPurgeReportResults"`
	PermScan                     string              `json:"permScan"`
	PermAgentScan                string              `json:"permAgentScan"`
	PermShareObjects             string              `json:"permShareObjects"`
	PermUpdateFeeds              string              `json:"permUpdateFeeds"`
	PermUploadNessusResults      string              `json:"permUploadNessusResults"`
	PermViewOrgLogs              string              `json:"permViewOrgLogs"`
	PermManageAcceptRiskRules    string              `json:"permManageAcceptRiskRules"`
	PermManageRecastRiskRules    string              `json:"permManageRecastRiskRules"`
	OrganizationCounts           []map[string]string `json:"organizationCounts" sc:"organizationCounts"`
	CanManageApp                 bool                `json:"-" sc:"permManageApp"`
	CanManageGroups              bool                `json:"-" sc:"permManageGroups"`
	CanManageRoles               bool                `json:"-" sc:"permManageRoles"`
	CanManageImages              bool                `json:"-" sc:"permManageImages"`
	CanManageGroupRelationships  bool                `json:"-" sc:"permManageGroupRelationships"`
	CanManageBlackoutWindows     bool                `json:"-" sc:"permManageBlackoutWindows"`
	CanManageAttributeSets       bool                `json:"-" sc:"permManageAttributeSets"`
	CanCreateTickets             bool                `json:"-" sc:"permCreateTickets"`
	CanCreateAlerts              bool                `json:"-" sc:"permCreateAlerts"`
	CanCreateAuditFiles          bool                `json:"-" sc:"permCreateAuditFiles"`
	CanCreateLDAPAssets          bool                `json:"-" sc:"permCreateLDAPAssets"`
	CanCreatePolicies            bool                `json:"-" sc:"permCreatePolicies"`
	CanPurgeTickets              bool                `json:"-" sc:"permPurgeTickets"`
	CanPurgeScanResults          bool                `json:"-" sc:"permPurgeScanResults"`
	CanPurgeReportResults        bool                `json:"-" sc:"permPurgeReportResults"`
	CanScan                      bool                `json:"-" sc:"permScan"`
	CanAgentScan                 bool                `json:"-" sc:"permAgentScan"`
	CanShareObjects              bool                `json:"-" sc:"permShareObjects"`
	CanUpdateFeeds               bool                `json:"-" sc:"permUpdateFeeds"`
	CanUploadNessusResults       bool                `json:"-" sc:"permUploadNessusResults"`
	CanViewOrgLogs               bool                `json:"-" sc:"permViewOrgLogs"`
	CanManageAcceptRiskRules     bool                `json:"-" sc:"permManageAcceptRiskRules"`
	CanManageRecastRiskRules     bool                `json:"-" sc:"permManageRecastRiskRules"`
}

// GetRoles gets the list of roles with the optionally user-supplied fields
// and an error if there was one
func (sc *SC) GetRoles(fields ...string) ([]Role, error) {
	var (
		resp        rolesResp
		roles       []Role
		req         = prepRoleRequest(sc, -1, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return nil, err
	}

	for _, r := range resp.Response {
		err = r.readAttr()
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, err
}

// GetRole returns the role in associated with the supplied id and
// the optionally user-supplied fields, including an error if there was one
func (sc *SC) GetRole(id int, fields ...string) (*Role, error) {
	var (
		resp        roleResp
		role        *Role
		req         = prepRoleRequest(sc, id, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return &resp.Response, err
	}

	role = &resp.Response

	err = role.readAttr()

	return role, err
}

func prepRoleRequest(sc *SC, id int, fields ...string) *Request {
	var path string
	if id <= 0 {
		path = "role"
	} else {
		path = fmt.Sprintf("role/%d", id)
	}

	req := sc.NewRequest("GET", path)

	if len(fields) > 0 {
		req.data["fields"] = strings.Join(fields, ",")
	}

	return req
}

func (r *Role) readAttr() error {
	// read nested types
	err := r.Creator.readAttr()
	if err != nil {
		return err
	}

	// read Unix timestamps
	err = readTime(&r.CreatedTimeStr, &r.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&r.ModifiedTimeStr, &r.ModifiedTime)
	if err != nil {
		return err
	}

	// read Booleans
	err = readBool(&r.PermManageApp, &r.CanManageApp)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageGroups, &r.CanManageGroups)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageRoles, &r.CanManageRoles)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageImages, &r.CanManageImages)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageGroupRelationships, &r.CanManageGroupRelationships)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageBlackoutWindows, &r.CanManageBlackoutWindows)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageAttributeSets, &r.CanManageAttributeSets)
	if err != nil {
		return err
	}

	err = readBool(&r.PermCreateTickets, &r.CanCreateTickets)
	if err != nil {
		return err
	}

	err = readBool(&r.PermCreateAlerts, &r.CanCreateAlerts)
	if err != nil {
		return err
	}

	err = readBool(&r.PermCreateAuditFiles, &r.CanCreateAuditFiles)
	if err != nil {
		return err
	}

	err = readBool(&r.PermCreateLDAPAssets, &r.CanCreateLDAPAssets)
	if err != nil {
		return err
	}

	err = readBool(&r.PermCreatePolicies, &r.CanCreatePolicies)
	if err != nil {
		return err
	}

	err = readBool(&r.PermPurgeTickets, &r.CanPurgeTickets)
	if err != nil {
		return err
	}

	err = readBool(&r.PermPurgeScanResults, &r.CanPurgeScanResults)
	if err != nil {
		return err
	}

	err = readBool(&r.PermPurgeReportResults, &r.CanPurgeReportResults)
	if err != nil {
		return err
	}

	err = readBool(&r.PermScan, &r.CanScan)
	if err != nil {
		return err
	}

	err = readBool(&r.PermAgentScan, &r.CanAgentScan)
	if err != nil {
		return err
	}

	err = readBool(&r.PermShareObjects, &r.CanShareObjects)
	if err != nil {
		return err
	}

	err = readBool(&r.PermUpdateFeeds, &r.CanUpdateFeeds)
	if err != nil {
		return err
	}

	err = readBool(&r.PermUploadNessusResults, &r.CanUploadNessusResults)
	if err != nil {
		return err
	}

	err = readBool(&r.PermViewOrgLogs, &r.CanViewOrgLogs)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageAcceptRiskRules, &r.CanManageAcceptRiskRules)
	if err != nil {
		return err
	}

	err = readBool(&r.PermManageRecastRiskRules, &r.CanManageRecastRiskRules)
	if err != nil {
		return err
	}

	return nil
}
