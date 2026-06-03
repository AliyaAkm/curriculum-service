package middleware

import (
	"curriculum-service/internal/domain"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleManager = "manager"
	RoleAdmin   = "admin"
)
const (
	contextUserIDKey = "auth.user_id"
	contextRolesKey  = "auth.roles"
)

var rolePriority = map[string]int{
	RoleStudent: 1,
	RoleTeacher: 2,
	RoleManager: 3,
	RoleAdmin:   4,
}

type Claims struct {
	Role     string   `json:"role"`
	Roles    []string `json:"roles,omitempty"`
	Login    string   `json:"login,omitempty"`
	IsActive bool     `json:"is_active"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	secret   []byte
	issuer   string
	audience string
	ttl      time.Duration
}

func New(secret []byte, issuer, audience string, ttl time.Duration) *Manager {
	return &Manager{secret: secret, issuer: issuer, audience: audience, ttl: ttl}
}

func GetUserID(jwtMgr *Manager, c *gin.Context) *uuid.UUID {
	claims := GetClaims(jwtMgr, c)
	if claims == nil {
		return nil
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil
	}

	return &userID
}

func GetClaims(jwtMgr *Manager, c *gin.Context) *Claims {
	tokenStr := bearerToken(c.GetHeader("Authorization"))
	if tokenStr == "" {
		return nil
	}

	claims, err := jwtMgr.VerifyAccessToken(tokenStr)
	if err != nil {
		return nil
	}

	return claims
}

func ClaimsHasRole(claims *Claims, role string) bool {
	if claims == nil {
		return false
	}
	return slices.Contains(claims.Roles, role) || claims.Role == role
}

func (m *Manager) VerifyAccessToken(tokenStr string) (*Claims, error) {
	claims, err := m.Verify(tokenStr)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}

func (m *Manager) Verify(tokenStr string) (*Claims, error) {
	parser := jwtlib.NewParser(jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}))

	tok, err := parser.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtlib.Token) (any, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, domain.ErrInvalidToken
	}

	if claims.Issuer != m.issuer {
		return nil, domain.ErrInvalidToken
	}

	if !audienceHas(claims.Audience, m.audience) {
		return nil, domain.ErrInvalidToken
	}

	claims.Roles = normalizeRoleClaims(claims.Role, claims.Roles)
	if claims.Role == "" && len(claims.Roles) > 0 {
		claims.Role = claims.Roles[0]
	}

	return claims, nil
}

func bearerToken(header string) string {
	if header == "" {
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
func audienceHas(auds jwtlib.ClaimStrings, want string) bool {
	return slices.Contains(auds, want)
}

func normalizeRoleClaims(primaryRole string, roles []string) []string {
	normalized := make([]string, 0, len(roles)+1)

	add := func(role string) {
		if role == "" || !IsValidRoleCode(role) {
			return
		}
		if slices.Contains(normalized, role) {
			return
		}
		normalized = append(normalized, role)
	}

	add(primaryRole)
	for _, role := range roles {
		add(role)
	}

	return normalized
}
func IsValidRoleCode(code string) bool {
	_, ok := rolePriority[code]
	return ok
}
