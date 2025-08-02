package infrastructure

import (
	"net/http"
	"strings"
	"write_base/internal/domain"
	"github.com/gin-gonic/gin"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
    memory "github.com/ulule/limiter/v3/drivers/store/memory"
    "github.com/ulule/limiter/v3"
    "time"
)
type Middleware struct {
	tokenService domain.ITokenService
}

func NewMiddleware(ts domain.ITokenService) *Middleware{
	return &Middleware{tokenService: ts}
}
func (m *Middleware) Authmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrAuthorizationHeaderRequired.Error()})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenString := authParts[1]
		authclaim, err := m.tokenService.ValidateAccessToken(tokenString)
		if err != nil{
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidToken.Error()})
			c.Abort()
			return 
		}
		
		c.Set("user_id", authclaim.UserID)
		c.Set("role", authclaim.Role)

		c.Next()
	}
}


func RequireRole(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists{
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": domain.ErrUnauthorized})
			c.Abort()
			return
		}
		role := domain.UserRole(userRole.(string))
		for _, allowed := range roles{
			if role == allowed{
				c.Next()
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": domain.ErrUnauthorized})
		c.Abort()
	}
}

func SetupRateLimiter() gin.HandlerFunc {
    // 5 requests per minute
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  5,
    }

    store := memory.NewStore()
    middleware := ginlimiter.NewMiddleware(limiter.New(store, rate))

    return middleware
}