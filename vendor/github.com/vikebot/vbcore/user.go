package vbcore

// User defines a user with all it's properties. All vikebot code
// should be able to access all information from the database through
// this struct
type User struct {
	ID               *int
	Permission       *int
	PermissionString *string
	Username         *string
	Name             *string
	Emails           []Email
	Bio              *string
	Location         *string
	Web              []string
	Company          *string
	Social           map[string]string
	OAuth            map[string]string
}

// Validate checks whether a `vbcore.User` struct, passed from clients, is in
// a valid state or not.
func (u *User) Validate() (safeUser *SafeUser, valid bool) {
	// It's not allowed to have the following attributes
	if u.ID != nil {
		return nil, false
	}
	if u.Permission != nil {
		return nil, false
	}
	if u.PermissionString != nil {
		return nil, false
	}
	if len(u.OAuth) != 0 {
		return nil, false
	}

	// You must have these attributes
	if u.Username == nil || len(*u.Username) > 32 {
		return nil, false
	}
	if u.Name == nil || len(*u.Name) > 32 {
		return nil, false
	}
	if len(u.Emails) == 0 {
		return nil, false
	}
	for _, i := range u.Emails {
		if len(i.Email) > 64 {
			return nil, false
		}
		if i.Status != EmailLinked && i.Status != EmailVerified && i.Status != EmailPrimary {
			return nil, false
		}
	}
	if u.Bio == nil || len(*u.Bio) > 1024 {
		return nil, false
	}
	if u.Location == nil || len(*u.Location) > 64 {
		return nil, false
	}
	if u.Company == nil || len(*u.Company) > 64 {
		return nil, false
	}
	for _, i := range u.Web {
		if len(i) > 128 {
			return nil, false
		}
	}

	// User is valid state -> convert and return
	return u.SafeUser(), true
}

// SafeUser converts a User to a SafeUser struct. All properties that
// aren't set in the User struct will have the type's default value
// in the SafeUser struct.
func (u *User) SafeUser() *SafeUser {
	su := SafeUser{}
	if u.ID != nil {
		su.ID = *u.ID
	}
	if u.Permission != nil {
		su.Permission = *u.Permission
	}
	if u.PermissionString != nil {
		su.PermissionString = *u.PermissionString
	}
	if u.Username != nil {
		su.Username = *u.Username
	}
	if u.Name != nil {
		su.Name = *u.Name
	}
	if u.Emails != nil {
		su.Emails = u.Emails
	}
	if u.Bio != nil {
		su.Bio = *u.Bio
	}
	if u.Location != nil {
		su.Location = *u.Location
	}
	if u.Web != nil {
		su.Web = u.Web
	}
	if u.Company != nil {
		su.Company = *u.Company
	}
	if u.Social != nil {
		su.Social = u.Social
	}
	if u.OAuth != nil {
		su.OAuth = u.OAuth
	}
	return &su
}

// SafeUser equals to User itself with one exception: all variables are
// already value fields (and not pointers) and there aren't any json
// declarations.
type SafeUser struct {
	ID               int               `json:"id"`
	Permission       int               `json:"permission"`
	PermissionString string            `json:"permission_string"`
	Username         string            `json:"username"`
	Name             string            `json:"name"`
	Emails           []Email           `json:"emails"`
	Bio              string            `json:"bio"`
	Location         string            `json:"location"`
	Web              []string          `json:"web"`
	Company          string            `json:"company"`
	Social           map[string]string `json:"social"`
	OAuth            map[string]string `json:"-"`
}

// User converts the SafeUser struct to a User struct.
func (su *SafeUser) User() *User {
	return &User{
		ID:               &su.ID,
		Permission:       &su.Permission,
		PermissionString: &su.PermissionString,
		Username:         &su.Username,
		Name:             &su.Name,
		Emails:           su.Emails,
		Bio:              &su.Bio,
		Location:         &su.Location,
		Web:              su.Web,
		Company:          &su.Company,
		Social:           su.Social,
		OAuth:            su.OAuth,
	}
}

// MakePublic converts the user into it's public version. Everyone is allowed
// to access this version without authentication. Currently only emails which
// are ether private or not verified are stripped away.
func (su *SafeUser) MakePublic() {
	emails := []Email{}
	for _, v := range su.Emails {
		if v.Public && (v.IsVerified() || v.IsPrimary()) {
			emails = append(emails, v)
		}
	}
	su.Emails = emails
}

// PrimaryEmail returns the primary email of this user. If there isn't anyone
// `nil` is returned.
func (su *SafeUser) PrimaryEmail() *Email {
	for _, v := range su.Emails {
		if v.IsPrimary() {
			return &v
		}
	}
	return nil
}
