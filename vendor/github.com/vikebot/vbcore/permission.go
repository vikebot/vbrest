package vbcore

const (
	// PermissionBanned defines that the user cannot access any vikebot
	// resources any more
	PermissionBanned = 0
	// PermissionDefault defines the initial permission-set given at the
	// users registration
	PermissionDefault = 1
	// PermissionVerified defines a state for profiles which where
	// controlled by a team member of vikebot manually
	PermissionVerified = 2
	// PermissionTeam defines that the user is working in any way for
	// vikebot and has supporting rights (e.g. set users to PermissionBanned,
	// PermissionDefault or PermissionVerified)
	PermissionTeam = 3
	// PermissionAdmin defines that the user is in vikebot's administration
	// board (e.g. director of ...)
	PermissionAdmin = 4
)

const (
	// PermissionBannedString is the string representation for PermissionBanned
	PermissionBannedString = "banned"
	// PermissionDefaultString is the string representation for PermissionDefault
	PermissionDefaultString = "default"
	// PermissionVerifiedString is the string representation for PermissionVerified
	PermissionVerifiedString = "verified"
	// PermissionTeamString is the string representation for PermissionTeam
	PermissionTeamString = "team"
	// PermissionAdminString is the string representation for PermissionAdmin
	PermissionAdminString = "admin"
)

// PermissionItoA converts an integer permission to value to it's equivalent string
// representation. If the value is lesser than PermissionBanned "banned" will
// be returned. If the value is bigger than PermissionAdmin "admin" will be
// returned.
func PermissionItoA(permission int) string {
	if permission < PermissionDefault {
		return PermissionBannedString
	} else if permission < PermissionVerified {
		return PermissionDefaultString
	} else if permission < PermissionTeam {
		return PermissionVerifiedString
	} else if permission < PermissionAdmin {
		return PermissionTeamString
	}
	return PermissionAdminString
}

// PermissionAtoI converts a permission string representation to it's equivalent
// real integer value. If the string value isn't defined PermissionAdmin
// will be returned.
func PermissionAtoI(permission string) int {
	switch permission {
	case PermissionBannedString:
		return PermissionBanned
	case PermissionDefaultString:
		return PermissionDefault
	case PermissionVerifiedString:
		return PermissionVerified
	case PermissionTeamString:
		return PermissionTeam
	default:
		return PermissionAdmin
	}
}
