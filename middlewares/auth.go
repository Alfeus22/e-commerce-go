package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Header Authorization diperlukan"})
			c.Abort()
			return
		}

		// Hapus "Bearer " dengan benar (termasuk spasinya)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		// Validasi token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Method Signing salah: %v", t.Header["alg"])
			}
			return []byte("secretbos"), nil
		})

		// Cek apakah ada error saat parse (misal token kadaluarsa)
		if err != nil || !token.Valid {
			fmt.Println("DEBUG ERROR JWT:", err) // Lihat di terminal kenapa gagal
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau kadaluarsa"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			userID, ok := claims["_id"].(string)

			if !ok {
				fmt.Print("Gagal mengambil id user dari token")
				c.JSON(400, gin.H{"error": "tken cacat : ID user tidak ditemukan"})
				c.Abort()
				return
			}

			c.Set("currentuser", userID)
			c.Set("role", claims["role"])
			c.Next()
		}
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")

		if !exists {
			ctx.JSON(400, gin.H{"error": "Role salah"})
			ctx.Abort()
			return
		}

		if role != requiredRole {
			fmt.Print(role)
			fmt.Print(requiredRole)
			ctx.JSON(403, gin.H{"error": "Akses ditolak: Anda bukan " + requiredRole})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
