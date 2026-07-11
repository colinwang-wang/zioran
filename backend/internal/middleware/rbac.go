package middleware

import (
	"sync"
	"time"

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

// permCache caches role->permissions with TTL to avoid per-request DB queries.
type permCache struct {
	mu    sync.RWMutex
	data  map[string]map[string]bool // role -> permission set
	expAt time.Time
	ttl   time.Duration
	db    *gorm.DB
}

func newPermCache(db *gorm.DB, ttl time.Duration) *permCache {
	return &permCache{db: db, ttl: ttl, data: make(map[string]map[string]bool)}
}

func (pc *permCache) hasPermission(role, perm string) bool {
	pc.mu.RLock()
	if time.Now().Before(pc.expAt) {
		perms, ok := pc.data[role]
		pc.mu.RUnlock()
		return ok && perms[perm]
	}
	pc.mu.RUnlock()

	// Refresh cache
	pc.mu.Lock()
	defer pc.mu.Unlock()
	// Double-check after lock
	if time.Now().Before(pc.expAt) {
		perms, ok := pc.data[role]
		return ok && perms[perm]
	}

	type row struct {
		Role       string
		Permission string
	}
	var rows []row
	pc.db.Table("role_permissions").Select("role, permission").Find(&rows)
	newData := make(map[string]map[string]bool)
	for _, r := range rows {
		if newData[r.Role] == nil {
			newData[r.Role] = make(map[string]bool)
		}
		newData[r.Role][r.Permission] = true
	}
	pc.data = newData
	pc.expAt = time.Now().Add(pc.ttl)

	perms, ok := pc.data[role]
	return ok && perms[perm]
}

var globalPermCache *permCache

// InitPermCache must be called once at startup with the db connection.
func InitPermCache(db *gorm.DB) {
	globalPermCache = newPermCache(db, 30*time.Second)
}

// PermissionRequired checks if the user's role has the specified permission.
// super_admin always passes. If permission cache is not initialized, falls back to pass-through.
func PermissionRequired(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("admin_role")
		if !exists {
			// AdminRequired not applied (test environment) — pass through
			c.Next()
			return
		}
		roleStr, _ := role.(string)
		if roleStr == "super_admin" {
			c.Next()
			return
		}
		if globalPermCache == nil {
			// Cache not initialized — pass through (shouldn't happen in production)
			c.Next()
			return
		}
		for _, perm := range permissions {
			if globalPermCache.hasPermission(roleStr, perm) {
				c.Next()
				return
			}
		}
		response.Error(c, errcode.ErrForbidden)
		c.Abort()
	}
}

// RoleRequired restricts access to specific roles. super_admin always passes.
// If admin_role is not set in context (e.g. AdminRequired not applied), passes through.
// DEPRECATED: Use PermissionRequired for dynamic DB-based checks.
func RoleRequired(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, exists := c.Get("admin_role")
		if !exists {
			// AdminRequired not applied (e.g. test environment) — pass through
			c.Next()
			return
		}
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
