package user

import (
	"fmt"
)

type Identity struct {
	Type      Type
	FirstName string
}
type Type string

const (
	FullAccess     Type = "full"
	RestrictAccess Type = "restrict"
)

// User struct holds ID and Type
type User struct {
	ID   int64
	Type Type
}

// GetRole returns if we are allowed to reply or not
func GetRole(id int64, firstName string, userRoles map[int64]Identity) (Type, error) {
	// if no rules, allow everyone not a bot
	if len(userRoles) == 0 {
		return FullAccess, nil
	}
	userData, ok := userRoles[id]
	if !ok {
		return RestrictAccess,
			fmt.Errorf("user %d is not allowed", id)
	}
	if len(userData.FirstName) > 0 && userData.FirstName != firstName {
		return RestrictAccess, fmt.Errorf(
			"user %d mistmatch firstname %v != %v",
			id, userData.FirstName, firstName)
	}
	return userData.Type, nil
}
