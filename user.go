package sc

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const defaultUserFields = "id,username,firstname,lastname,status,role,title,email,address,city,state,country,phone,fax,createdTime,modifiedTime,lastLogin,lastLoginIP,mustChangePassword,locked,failedLogins,authType,fingerprint,password,description,canUse,canManage,managedUsersGroups,managedObjectsGroups,preferences,ldaps,ldapUsername,group,responsibleAsset"

type usersResp struct {
	scResp
	Response []User
}

type userResp struct {
	scResp
	Response User
}

// userProp is a single record containing an application user's profile property value
// for now, userProp is mostly used in cases where other objects (i.e. assets, groups) go,
// but is used as a bandaid until full native object permeability is implemeted
type userPropStr struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type userPropInt struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// userPref is a single record containing an application user's preference value
type userPref struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

// User https://docs.tenable.com/sccv/api/User.html
type User struct {
	ID                    string                 `json:"id" sc:"id"`
	Status                string                 `json:"status" sc:"status"`
	Username              string                 `json:"username" sc:"username"`
	LdapUsername          string                 `json:"ldapUsername" sc:"ldapUsername"`
	FirstName             string                 `json:"firstname" sc:"firstname"`
	LastName              string                 `json:"lastname" sc:"lastname"`
	Title                 string                 `json:"title" sc:"title"`
	Email                 string                 `json:"email" sc:"email"`
	Address               string                 `json:"address" sc:"address"`
	City                  string                 `json:"city" sc:"city"`
	State                 string                 `json:"state" sc:"state"`
	Country               string                 `json:"country" sc:"country"`
	Phone                 string                 `json:"phone" sc:"phone"`
	Fax                   string                 `json:"fax" sc:"fax"`
	CreatedTime           time.Time              `json:"-" sc:"createdTime"`
	CreatedTimeStr        string                 `json:"createdTime"`
	ModifiedTime          time.Time              `json:"-" sc:"modifiedTime"`
	ModifiedTimeStr       string                 `json:"modifiedTime"`
	LastLogin             time.Time              `json:"-" sc:"lastLogin"`
	LastLoginStr          string                 `json:"lastLogin"`
	LastLoginIP           net.IP                 `json:"-" sc:"lastLoginIP"`
	LastLoginIPStr        string                 `json:"lastLoginIP"`
	MustChangePassword    bool                   `json:"-" sc:"mustChangePassword"`
	MustChangePasswordStr string                 `json:"mustChangePassword"`
	IsLocked              bool                   `json:"-" sc:"locked"`
	Locked                string                 `json:"locked"`
	FailedLogins          string                 `json:"failedLogins" sc:"failedLogins"`
	AuthType              string                 `json:"authType" sc:"authType"`
	Fingerprint           string                 `json:"fingerprint" sc:"fingerprint"`
	Password              string                 `json:"password" sc:"password"`
	ManagedUsersGroups    []userPropStr          `json:"managedUsersGroups" sc:"managedUsersGroups"`
	ManagedObjectsGroups  []userPropStr          `json:"managedObjectsGroups" sc:"managedObjectsGroups"`
	UserPreferences       []userPref             `json:"userPrefs" sc:"userPrefs"`
	Preferences           []userPref             `json:"preferences" sc:"preferences"`
	CanUse                bool                   `json:"canUse" sc:"canUse"`
	CanManage             bool                   `json:"canManage" sc:"canManage"`
	Role                  userPropStr            `json:"role" sc:"role"`
	ResponsibleAsset      userPropInt            `json:"responsibleAsset" sc:"responsibleAsset"`
	Organization          userPropStr            `json:"organization" sc:"organization"`
	Group                 userPropStr            `json:"group" sc:"group"`
	Ldap                  map[string]interface{} `json:"ldap" sc:"ldap"`
}

func (u *User) String() string {
	msg := "#%s %s (%s <%s>)\n\t%s of %s\n\tLast Login: "
	if u.LastLogin.Year() < 2001 {
		msg += "never"
	} else {
		msg += fmt.Sprintf("%s (%s) from %s", u.LastLogin.String(), time.Since(u.LastLogin), u.LastLoginIP.String())
	}

	email := "nil"
	if len(u.Email) > 0 {
		email = u.Email
	}

	return fmt.Sprintf(msg, u.ID, u.Username, fmt.Sprintf("%s %s", u.FirstName, u.LastName), email, u.Role.Name, u.Group.Name)
}

// GetCurrentUser returns the current user authenticated with the SC
// instance, using /rest/currentUser endpoint response, with the optionally
// user-supplied or default user fields
func (sc *SC) GetCurrentUser(fields ...string) (*User, error) {
	var (
		resp        userResp
		user        *User
		req         = prepUserRequest(sc, 0, -1, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return &resp.Response, err
	}

	user = &resp.Response

	err = user.readAttr()

	return user, err
}

// GetUser returns the user in the default organization with the
// supplied user id and the optionally user-supplied or default user fields
func (sc *SC) GetUser(id int, fields ...string) (*User, error) {
	var (
		resp        userResp
		user        *User
		req         = prepUserRequest(sc, id, -1, fields...)
		scResp, err = req.Do()
	)

	fmt.Println("get user", id)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return &resp.Response, err
	}

	user = &resp.Response

	err = user.readAttr()

	return user, err
}

// GetUserFromOrg returns the user in the organization with the
// user-supplied orgID and user id and the optionally user-supplied or default user fields
func (sc *SC) GetUserFromOrg(id, orgID int, fields ...string) (*User, error) {
	var (
		resp        userResp
		user        *User
		req         = prepUserRequest(sc, id, orgID, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return &resp.Response, err
	}

	user = &resp.Response

	err = user.readAttr()

	return user, err
}

// GetUsersFromOrg returns the users in the organization with the
// user-supplied orgID and the optionally user-supplied or default user fields
func (sc *SC) GetUsersFromOrg(orgID int, fields ...string) ([]User, error) {
	var (
		resp        usersResp
		users       []User
		req         = prepUserRequest(sc, -1, orgID, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return nil, err
	}

	users = resp.Response

	for _, u := range users {
		err = u.readAttr()
		if err != nil {
			return nil, err
		}
	}

	return users, err
}

// GetUsers returns the users in the current organization with the
// optionally user-supplied or default user fields
func (sc *SC) GetUsers(fields ...string) ([]User, error) {
	var (
		resp        usersResp
		users       []User
		req         = prepUserRequest(sc, -1, -1, fields...)
		scResp, err = req.Do()
	)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(scResp.RespBody, &resp)
	if err != nil {
		return nil, err
	}

	for _, u := range resp.Response {
		err = u.readAttr()
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, err
}

func prepUserRequest(sc *SC, id, orgID int, fields ...string) *Request {
	var path string
	if id < 0 {
		path = "user"
	} else if id == 0 {
		path = "currentUser"
	} else {
		path = fmt.Sprintf("user/%d", id)
	}

	req := sc.NewRequest("GET", path)

	if orgID >= 0 {
		req.data["orgID"] = orgID
	}

	if len(fields) > 0 {
		req.data["fields"] = strings.Join(fields, ",")
	} else {
		req.data["fields"] = defaultUserFields
	}

	return req
}

func parseRespToUser(data []byte) (User, error) {
	var (
		user User
		err  = json.Unmarshal(data, &user)
	)

	if err != nil {
		return user, err
	}

	user.LastLoginIP = net.ParseIP(user.LastLoginIPStr)
	n, err := strconv.ParseInt(user.CreatedTimeStr, 10, 64)
	if err != nil {
		return user, err
	}
	user.CreatedTime = time.Unix(n, 0)

	return user, err
}

func (u *User) readAttr() error {
	// read net.IPs
	u.LastLoginIP = net.ParseIP(u.LastLoginIPStr)

	// read Booleans
	err := readBool(&u.Locked, &u.IsLocked)
	if err != nil {
		return err
	}

	err = readBool(&u.MustChangePasswordStr, &u.MustChangePassword)
	if err != nil {
		return err
	}

	// read Unix timestamps
	err = readTime(&u.CreatedTimeStr, &u.CreatedTime)
	if err != nil {
		return err
	}

	err = readTime(&u.ModifiedTimeStr, &u.ModifiedTime)
	if err != nil {
		return err
	}

	err = readTime(&u.LastLoginStr, &u.LastLogin)
	if err != nil {
		return err
	}

	return nil
}
