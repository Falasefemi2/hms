package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/falasefemi2/hms/internal/utils"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}
		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, utils.RoleKey, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleValue := r.Context().Value(utils.RoleKey)
		role, ok := roleValue.(string)
		if !ok || role != "ADMIN" {
			utils.WriteError(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func DoctorOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleValue := r.Context().Value(utils.RoleKey)
		role, ok := roleValue.(string)

		if !ok || role != "DOCTOR" {
			utils.WriteError(w, http.StatusForbidden, "doctor access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func NurseOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleValue := r.Context().Value(utils.RoleKey)
		role, ok := roleValue.(string)

		if !ok || role != "NURSE" {
			utils.WriteError(w, http.StatusForbidden, "nurse access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func PatientOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleValue := r.Context().Value(utils.RoleKey)
		role, ok := roleValue.(string)

		if !ok || role != "PATIENT" {
			utils.WriteError(w, http.StatusForbidden, "patient access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HasAnyRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleValue := r.Context().Value(utils.RoleKey)
			role, ok := roleValue.(string)

			if !ok {
				utils.WriteError(w, http.StatusForbidden, "user role not found in context")
				return
			}

			// Check if user's role is in allowed roles
			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				utils.WriteError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(utils.UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func GetRoleFromContext(ctx context.Context) string {
	role, ok := ctx.Value(utils.RoleKey).(string)
	if !ok {
		return ""
	}
	return role
}
