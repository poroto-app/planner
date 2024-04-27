package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

const (
	sessionCookieExpiration = time.Hour * 24 * 5
	cookieKeySession        = "session"
)

type LoginRequest struct {
	IdToken string `json:"id_token"`
}

// SessionLoginHandler Firebase AuthによるセッションCookieを作成する
// https://firebase.google.com/docs/auth/admin/manage-cookies?hl=ja#go
func (s Server) SessionLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verify ID token passed in via HTTP json body
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		token, err := s.firebaseAuth.VerifyIdToken(c.Request.Context(), req.IdToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid ID token",
			})
			return
		}

		// Verify expiration time
		if token.Expires < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Recent login required",
			})
			return
		}

		cookie, err := s.firebaseAuth.CreateSessionCookie(
			c.Request.Context(),
			req.IdToken,
			sessionCookieExpiration,
		)

		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(
			cookieKeySession,
			cookie,
			int(sessionCookieExpiration.Seconds()),
			"/",
			os.Getenv("HOST"),
			true,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"message": "Session cookie created",
		})
	}
}

// SessionLogoutHandler Firebase AuthによるセッションCookieを削除する
func (s Server) SessionLogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie(
			cookieKeySession,
			"",
			-1,
			"/",
			os.Getenv("HOST"),
			true,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"message": "Session cookie deleted",
		})
	}
}

// Cookie情報からユーザー情報を取得する
func (s Server) SessionUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(cookieKeySession)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "No session cookie",
			})
			return
		}

		token, err := s.firebaseAuth.VerifySessionCookie(c.Request.Context(), cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session cookie",
			})
			return
		}

		user, err := s.firebaseAuth.GetUser(c.Request.Context(), token.UID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session cookie",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user": user.DisplayName,
		})
	}
}
