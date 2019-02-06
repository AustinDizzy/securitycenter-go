package sc

import (
	"strconv"
	"time"
)

// AssetTemplate https://docs.tenable.com/sccv/api/Asset-Template.html
type AssetTemplate struct {
	ID          int              `json:"id" sc:"id"`
	Name        string           `json:"name" sc:"name"`
	Description string           `json:"description" sc:"description"`
	Summary     string           `json:"summary" sc:"summary"`
	Type        string           `json:"type" sc:"type"`
	Category    TemplateCategory `json:"category" sc:"category"`
	Definition  struct {
		Rules           assetTemplateDefinitionRules `json:"rules" sc:"rules"`
		AssetDataFields []map[string]interface{}
	} `json:"definition" sc:"definition"`
	AssetType             string              `json:"assetType" sc:"assetType"`
	Enabled               bool                `json:"-" sc:"enabled"`
	EnabledStr            string              `json:"enabled"`
	MinUpgradeVersion     string              `json:"minUpgradeVersion"`
	TemplatePubTime       time.Time           `json:"-" sc:"templatePubTime"`
	TemplatePubTimeStr    string              `json:"templatePubTime"`
	TemplateModTime       time.Time           `json:"-" sc:"templateModTime"`
	TemplateModTimeStr    string              `json:"templateModTime"`
	TemplateDefModTime    time.Time           `json:"-" sc:"templateDefModTime"`
	TemplateDefModTimeStr string              `json:"templateDefModTime"`
	DefinitionModTime     time.Time           `json:"-" sc:"definitionModTime"`
	DefinitionModTimeStr  string              `json:"definitionModTime"`
	CreatedTime           time.Time           `json:"-" sc:"createdTime"`
	CreatedTimeStr        string              `json:"createdTime"`
	ModifiedTime          time.Time           `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr       string              `json:"modifiedTime"`
	Tags                  []string            `json:"tags" sc:"tags"`
	Requirements          []map[string]string `json:"requirements" sc:"requirements"`
}

type assetTemplateDefinitionRules struct {
	Operator string              `json:"operator" sc:"operator"`
	Children []map[string]string `json:"children" sc:"children"`
	Type     string              `json:"type" sc:"type"`
}

// TemplateCategory https://docs.tenable.com/sccv/api/Asset-Template.html#AssetTemplateRESTReference-/assetTemplate/categories
type TemplateCategory struct {
	ID          string `json:"id" sc:"id"`
	Name        string `json:"name" sc:"name"`
	Description string `json:"description" sc:"description"`
	CountStr    string `json:"count"`
	Count       int    `json:"-" sc:"count"`
	Status      string `json:"status" sc:"status"`
}

func (c *TemplateCategory) readAttr() error {
	var err error

	c.Count, err = strconv.Atoi(c.CountStr)

	return err
}

func (t *AssetTemplate) readAttr() error {
	err := readBool(&t.EnabledStr, &t.Enabled)
	if err != nil {
		return err
	}

	err = readTime(&t.TemplatePubTimeStr, &t.TemplatePubTime)
	if err != nil {
		return err
	}

	err = readTime(&t.TemplateModTimeStr, &t.TemplateModTime)
	if err != nil {
		return err
	}

	err = readTime(&t.TemplateDefModTimeStr, &t.TemplateDefModTime)
	if err != nil {
		return err
	}

	err = readTime(&t.DefinitionModTimeStr, &t.DefinitionModTime)
	if err != nil {
		return err
	}

	err = readTime(&t.CreatedTimeStr, &t.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&t.ModifiedTimeStr, &t.ModifiedTime)
	if err != nil {
		return err
	}

	return t.Category.readAttr()
}
