package config

type DBConfigExt struct {
	PublicSchema       string
	TenantSchemaPrefix string
}

// PublicSchema returns the configured shared schema/database name (default: "public").
func (c Config) PublicSchema() string {
	if c.DB.PublicSchema != "" {
		return c.DB.PublicSchema
	}
	return "public"
}

// TenantSchemaPrefix returns the configured tenant schema/database prefix (default: "tenant_").
func (c Config) TenantSchemaPrefix() string {
	if c.DB.TenantSchemaPrefix != "" {
		return c.DB.TenantSchemaPrefix
	}
	return "tenant_"
}
