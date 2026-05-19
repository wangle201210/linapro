// This file defines the source-plugin visible authentication contract.

package contract

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// AuthService defines tenant token operations published to source plugins.
type AuthService interface {
	// SelectTenant consumes a pre-login token and issues a tenant-bound token.
	SelectTenant(ctx context.Context, in SelectTenantInput) (*TenantTokenOutput, error)
	// SwitchTenant validates membership, revokes the current token, and issues a new token.
	SwitchTenant(ctx context.Context, in SwitchTenantInput) (*TenantTokenOutput, error)
	// AuthenticateBearer validates one LinaPro access token, confirms the
	// online session is active, and returns the effective permission snapshot.
	AuthenticateBearer(ctx context.Context, bearerToken string) (*AuthenticatedIdentity, error)
}

// SelectTenantInput defines input for a pre-token tenant selection.
type SelectTenantInput struct {
	// PreToken is the short-lived pre-login token produced by host login.
	PreToken string
	// TenantID is the requested target tenant.
	TenantID int
}

// SwitchTenantInput defines input for authenticated tenant switching.
type SwitchTenantInput struct {
	// BearerToken is the current Authorization bearer token.
	BearerToken string
	// TenantID is the requested target tenant.
	TenantID int
}

// TenantTokenOutput contains one newly signed tenant-bound access token.
type TenantTokenOutput struct {
	// AccessToken is the host-compatible JWT.
	AccessToken string
	// RefreshToken is the host-compatible refresh JWT for the same session.
	RefreshToken string
}

// AuthenticatedIdentity is the plugin-visible LinaPro identity and access snapshot.
type AuthenticatedIdentity struct {
	// TokenID is the authenticated online-session token identifier.
	TokenID string
	// UserID is the authenticated LinaPro user ID.
	UserID int
	// Username is the authenticated LinaPro username.
	Username string
	// Status is the authenticated account status.
	Status int
	// TenantID is the effective host tenant for the request.
	TenantID int
	// ActingUserID is the real platform user ID during impersonation.
	ActingUserID int
	// ActingAsTenant reports whether the request acts through a tenant view.
	ActingAsTenant bool
	// IsImpersonation reports whether the current token represents impersonation.
	IsImpersonation bool
	// PlatformBypass reports whether the identity can bypass tenant filtering.
	PlatformBypass bool
	// DataScope is the effective role data scope.
	DataScope int
	// DataScopeUnsupported reports whether the effective role data scope is unsupported.
	DataScopeUnsupported bool
	// UnsupportedDataScope stores the first unsupported data-scope value.
	UnsupportedDataScope int
	// Permissions contains effective LinaPro permission strings.
	Permissions []string
	// IsSuperAdmin reports whether the user is the built-in administrator.
	IsSuperAdmin bool
}

// CurrentContext projects the authenticated identity into plugin-visible request context.
func (i *AuthenticatedIdentity) CurrentContext() CurrentContext {
	if i == nil {
		return CurrentContext{}
	}
	return CurrentContext{
		UserID:          i.UserID,
		Username:        i.Username,
		TenantID:        i.TenantID,
		ActingUserID:    i.ActingUserID,
		ActingAsTenant:  i.ActingAsTenant,
		IsImpersonation: i.IsImpersonation,
		PlatformBypass:  i.PlatformBypass,
	}
}

// HasPermissions applies LinaPro permission semantics to one required permission list.
func (i *AuthenticatedIdentity) HasPermissions(required []string) bool {
	if len(required) == 0 {
		return true
	}
	if i == nil {
		return false
	}
	if i.IsSuperAdmin {
		return true
	}
	granted := make(map[string]struct{}, len(i.Permissions))
	for _, permission := range i.Permissions {
		currentPermission := strings.TrimSpace(permission)
		if currentPermission == "" {
			continue
		}
		granted[currentPermission] = struct{}{}
	}
	if _, ok := granted["*:*:*"]; ok {
		return true
	}
	for _, permission := range required {
		normalizedPermission := strings.TrimSpace(permission)
		if normalizedPermission == "" {
			continue
		}
		if _, ok := granted[normalizedPermission]; !ok {
			return false
		}
	}
	return true
}

// BearerTokenFromContext extracts the bearer token from the current HTTP request.
func BearerTokenFromContext(ctx context.Context) (string, bool) {
	request := g.RequestFromCtx(ctx)
	if request == nil {
		return "", false
	}
	header := request.GetHeader("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")
	return token, token != "" && token != header
}
