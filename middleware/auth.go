package middleware

import (
	"hackathon/repositories"
	"hackathon/services"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func NewAuthMiddleware(userRepo repositories.UserRepository, secret string) echo.MiddlewareFunc {
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(services.JwtCustomClaims) },
		SigningKey:    []byte(secret),
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing or invalid token"})
		},
	}

	jwtMiddleware := echojwt.WithConfig(jwtConfig)

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return jwtMiddleware(func(c echo.Context) error {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
			}

			claims, ok := token.Claims.(*services.JwtCustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
			}

			userID, err := strconv.ParseUint(claims.ID, 10, 32)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid user ID in token"})
			}

			user, err := userRepo.FindByID(uint(userID))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User not found"})
			}

			issuedAt, err := claims.GetIssuedAt()
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token timestamp"})
			}

			if user.RevokeTokensBefore > 0 && issuedAt.Time.Unix() < user.RevokeTokensBefore {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token has been revoked, please login again"})
			}

			c.Set("user", user)

			return next(c)
		})
	}
}
