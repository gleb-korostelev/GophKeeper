package auth

import (
	"github.com/gleb-korostelev/GophKeeper/models"
	auth "github.com/gleb-korostelev/GophKeeper/pkg/claims"
)

// getRole maps an account type to its corresponding role generator function.
func getRole(accType models.AccountType) func(string) auth.Ability {
	switch accType {
	case models.AccountRoleAdmin:
		return auth.AdminRole
	case models.AccountRoleSuperAdmin:
		return auth.SuperAdminRole
	default:
		return auth.RoleAuthorized
	}
}
