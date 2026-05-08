package schema

// Identity carries the caller's tenant, user, and use-case context,
// normalized from request headers by the ingress layer.
type Identity struct {
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`
	UseCaseID string `json:"use_case_id"`
}
