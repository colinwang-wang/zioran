package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
	"gorm.io/gorm"
)

// adminRoles lists roles that are allowed to access the admin panel.
var adminRoles = map[string]bool{
	"super_admin": true,
	"admin":       true,
	"operator":    true,
	"support":     true,
}

// AdminRequired verifies the authenticated user has an admin-level role.
// Must be placed after JWTAuth middleware.
func AdminRequired(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		var role string
		err := db.Table("users").Select("role").Where("id = ?", userID).Row().Scan(&role)
		if err != nil || !adminRoles[role] {
			response.Error(c, errcode.ErrForbidden)
			c.Abort()
			return
		}

		c.Set("admin_role", role)
		c.Next()
	}
}

// RoleRequired restricts access to specific roles. super_admin always passes.
func RoleRequired(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("admin_role")
		roleStr, _ := role.(string)
		if roleStr == "super_admin" {
			c.Next()
			return
		}
		if !allowed[roleStr] {
			response.Error(c, errcode.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
