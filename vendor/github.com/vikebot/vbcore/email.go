package vbcore

const (
	// EmailLinked defines that the user has added this email, but not verified
	// he has controll over it
	EmailLinked = 0
	// EmailVerified defines that the user has added this email and verfied that
	// he has controll over it
	EmailVerified = 1
	// EmailPrimary defines that the user has set this email as primary address
	// for his account. Also the one we show publicly.
	EmailPrimary = 2
)

const (
	// EmailLinkedString is the string representation for EmailLinked
	EmailLinkedString = "linked"
	// EmailVerifiedString is the string representation for EmailVerified
	EmailVerifiedString = "verified"
	// EmailPrimaryString is the string representation for EmailPrimary
	EmailPrimaryString = "primary"
)

// EmailItoA converts a status integer to it's equivalent string representation. If
// the value is lesser than EmailLinked "linked" will be returned. If the value
// is bigger than EmailPrimary "primary" will be returned.
func EmailItoA(status int) string {
	if status > EmailVerified {
		return EmailPrimaryString
	} else if status > EmailLinked {
		return EmailVerifiedString
	}
	return EmailLinkedString
}

// EmailAtoI converts a status string representation to it's equivalent real integer
// value. If the string value isn't defined EmailLinked will be returned.
func EmailAtoI(status string) int {
	switch status {
	case EmailPrimaryString:
		return EmailPrimary
	case EmailVerifiedString:
		return EmailVerified
	default:
		return EmailLinked
	}
}

// Email defines a single email entry of a user
type Email struct {
	Email  string
	Status int
	Public bool
}

// IsPrimary indicates if the `Status` is `EmailPrimary`
func (e Email) IsPrimary() bool {
	return e.Status == EmailPrimary
}

// IsVerified indicates if the `Status` is `EmailVerified`
func (e Email) IsVerified() bool {
	return e.Status == EmailVerified
}

// IsLinked indicates if the `Status` is `EmailLinked`
func (e Email) IsLinked() bool {
	return e.Status == EmailLinked
}
