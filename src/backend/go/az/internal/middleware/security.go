package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds comprehensive security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' https://login.microsoftonline.com; " +
			"style-src 'self' 'unsafe-inline'; " +
			"connect-src 'self' https://login.microsoftonline.com https://management.azure.com; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' https: data:; " +
			"frame-src https://login.microsoftonline.com; " +
			"object-src 'none'; " +
			"base-uri 'self'"

		c.Header("Content-Security-Policy", csp)

		// Strict Transport Security (HSTS)
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Frame-Options
		c.Header("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		c.Header("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=(), magnetometer=(), gyroscope=(), accelerometer=()")

		// Cross-Origin Embedder Policy
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")

		// Cross-Origin Opener Policy
		c.Header("Cross-Origin-Opener-Policy", "same-origin")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// CSRFProtection validates CSRF tokens
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip GET, HEAD, OPTIONS, and TRACE methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || 
		   c.Request.Method == "OPTIONS" || c.Request.Method == "TRACE" {
			c.Next()
			return
		}

		// Get CSRF token from header
		headerToken := c.GetHeader("X-CSRF-Token")
		
		// Get CSRF token from cookie
		cookieToken, err := c.Cookie("csrf_token")
		
		// Validate tokens exist and match (double-submit cookie pattern)
		if err != nil || headerToken == "" || cookieToken == "" || headerToken != cookieToken {
			c.Header("X-CSRF-Error", "true")
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Invalid or missing CSRF token. Please refresh the page and try again.",
				"code":    "CSRF_TOKEN_MISMATCH",
				"action":  "refresh_page",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateCSRFToken creates a new CSRF token
func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SetCSRFToken sets CSRF token in cookie
func SetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if CSRF token already exists
		if _, err := c.Cookie("csrf_token"); err == nil {
			c.Next()
			return
		}

		// Generate new CSRF token
		token, err := GenerateCSRFToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate CSRF token",
			})
			c.Abort()
			return
		}

		// Set CSRF token in cookie
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("csrf_token", token, 3600, "/", "", true, false) // httpOnly=false so JS can read it

		c.Next()
	}
}