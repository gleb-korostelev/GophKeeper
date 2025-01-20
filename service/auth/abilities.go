package auth

import (
	"github.com/gleb-korostelev/GophKeeper/models"
	auth "github.com/gleb-korostelev/GophKeeper/pkg/claims"
)

// getRole maps an account type to its corresponding role generator function.
//
// Parameters:
// - accType: The account type from the `models.AccountType` enum.
//
// Returns:
// - func(string) auth.Ability: A function that generates an `auth.Ability` for the given account type and username.
//
// Workflow:
// 1. Matches the account type to a specific role generator.
// 2. Returns the appropriate role generator function:
//   - `auth.AdminRole` for `AccountRoleAdmin`.
//   - `auth.SuperAdminRole` for `AccountRoleSuperAdmin`.
//   - `auth.RoleAuthorized` for all other account types.
//
// Example Usage:
//
//	accType := models.AccountRoleAdmin
//	roleGen := getRole(accType)
//	ability := roleGen("john_doe")
//
// Example Output:
//
//	If `accType` is `models.AccountRoleAdmin`, the returned `ability` will be:
//	Ability{Name: "admin", Scope: "john_doe"}.
//
// Supported Account Types:
// - `models.AccountRoleAdmin`: Generates an admin role.
// - `models.AccountRoleSuperAdmin`: Generates a superadmin role.
// - Other types (default): Generates an authorized user role.
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
