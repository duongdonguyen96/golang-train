package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"golang-train/internal/shared/database"
	"golang-train/internal/shared/tenant"
)

const (
	ginDBKey     = "tenant_db"
	ginTenantKey = "tenant"
)

func TenantDBMiddleware(resolver tenant.Resolver, factory database.Factory) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("222222222222222222222222222222222222222222")
		tenantID, err := resolver.Resolve(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("000000000000000000000000000000000000000000")
		db, err := factory.ForTenant(c.Request.Context(), tenantID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("111111111111111111111111111111111111111111")

		c.Set(ginTenantKey, tenantID)
		c.Set(ginDBKey, db)
		c.Request = c.Request.WithContext(database.WithDB(c.Request.Context(), db))
		c.Next()
	}
}

func TenantFromGin(c *gin.Context) (string, bool) {
	v, ok := c.Get(ginTenantKey)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
