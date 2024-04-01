package utils

import (
	"github.com/les-cours/auth-service/api/users"
	"github.com/les-cours/auth-service/types"
)

func PermissionsRPCToCode(permissions *users.Permissions) types.Permissions {
	return types.Permissions{
		Bots:     permissions.Bots,
		Triggers: permissions.Triggers,
		Tickets:  permissions.Tickets,
		Profiles: permissions.Profiles,
		Kbas:     permissions.Kbas,
		Settings: permissions.Settings,
	}
}
