package tenant

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Resolver interface {
	Resolve(c *gin.Context) (string, error)
}

type headerResolver struct {
	header string
}

func NewHeaderResolver(header string) Resolver {
	return &headerResolver{header: header}
}

var tenantRe = regexp.MustCompile(`^[a-z0-9_]{1,32}$`)

func (r *headerResolver) Resolve(c *gin.Context) (string, error) {
	tenant := c.GetHeader(r.header)
	if tenant == "" {
		return "defaulta", nil
		// return "", fmt.Errorf("missing tenant header %s", r.header)
	}
	if !tenantRe.MatchString(tenant) {
		return "", fmt.Errorf("invalid tenant")
	}
	return tenant, nil
}
