package handler

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleSignIn is oauth2 for Google
func GoogleSignIn(c echo.Context) error {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	oauthState := generateStateOuthCookie(c)
	callbackURL := googleOauthConfig.AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, callbackURL)
}

func generateStateOuthCookie(c echo.Context) string {
	expirationLength, err := strconv.Atoi(os.Getenv("OAUTH_EXPIRATION"))
	if err != nil {
		panic("OAUTH_EXPIRATION not presented")
	}
	var expiration = time.Now().Add(time.Duration(expirationLength) * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}

	c.SetCookie(&cookie)

	return state
}
