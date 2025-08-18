package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/otomaxv2/internal/config"
)
var oauthState = "apasih1788wwWW"
func GoogleLoginHandler(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL(oauthState)
	c.JSON(http.StatusOK, gin.H{
		"login_url": url,
	})
}

func GoogleCallbackHandler(c *gin.Context) {
	ctx := context.Background()
fmt.Println("CLIENT_ID:", config.GoogleOauthConfig.ClientID)
fmt.Println("CLIENT_SECRET:", config.GoogleOauthConfig.ClientSecret)
fmt.Println("REDIRECT_URL:", config.GoogleOauthConfig.RedirectURL)
	// validasi state (disarankan simpan di session)
	state := c.Query("state")
	if state != oauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	// ambil code
	code := c.Query("code")
	token, err := config.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed", "details": err.Error()})
		return
	}
	fmt.Println("STATE:", c.Query("state"))
fmt.Println("CODE:", c.Query("code"))

	// pakai token untuk akses API userinfo
	client := config.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user info", "details": err.Error()})
		return
	}

	// ✅ berhasil login
	c.JSON(http.StatusOK, gin.H{
		"message": "Google login success",
		"user":    userInfo,
	})
}