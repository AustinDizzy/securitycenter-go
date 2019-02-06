package sc

import (
	"time"
)

// Group is the application group located within SecurityCenter
type Group struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	UserCount       int          `json:"userCount"`
	Users           []User       `json:"users"`
	CreatedTimeStr  string       `json:"createdTime"`
	CreatedTime     time.Time    `json:"-" sc:"createdTime"`
	ModifiedTimeStr string       `json:"modifiedTime"`
	ModifiedTime    time.Time    `json:"-" sc:"modifiedTime"`
	LCEs            []LCE        `json:"lces" sc:"lces"`
	Repositories    []Repository `json:"repositories" sc:"repositories"`
	DefiningAssets  []Asset      `json:"definingAssets" sc:"definingAssets"`
	Assets          []Asset      `json:"assets" sc:"assets"`
	Policies        []Policy     `json:"policies" sc:"policies"`
}

func (g *Group) readAttr() error {
	for _, u := range g.Users {
		err := u.readAttr()
		if err != nil {
			return err
		}
	}

	for _, l := range g.LCEs {
		err := l.readAttr()
		if err != nil {
			return err
		}
	}

	for _, r := range g.Repositories {
		err := r.readAttr()
		if err != nil {
			return err
		}
	}

	for _, a := range g.DefiningAssets {
		err := a.readAttr()
		if err != nil {
			return err
		}
	}

	for _, a := range g.Assets {
		err := a.readAttr()
		if err != nil {
			return err
		}
	}

	for _, p := range g.Policies {
		err := p.readAttr()
		if err != nil {
			return err
		}
	}

	err := readTime(&g.CreatedTimeStr, &g.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&g.ModifiedTimeStr, &g.ModifiedTime)
	if err != nil {
		return err
	}

	return nil
}
