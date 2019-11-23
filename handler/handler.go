package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{}

const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// GoogleUserInfo is user info will be passed after successfully signed in at Google
type GoogleUserInfo struct {
	ID    string
	Email string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleSignIn is oauth2 for Google
func GoogleSignIn(c echo.Context) error {
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

// GoogleCallbackAfterSuccess is to get user info after sign in success
func GoogleCallbackAfterSuccess(c echo.Context) error {
	oauthState, err := c.Cookie("oauthstate")
	if err != nil {
		panic(err)
	}

	if c.FormValue("state") != oauthState.Value {
		return echo.NewHTTPError(http.StatusUnauthorized, "User is invalid")
	}

	data, err := getUserDataFromGoogle(c.FormValue("code"))

	if err != nil {
		return err
	}

	return c.JSONPretty(http.StatusOK, data, " ")
}

func getUserDataFromGoogle(code string) (*GoogleUserInfo, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Code exchange wrong")
	}

	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed getting user info")
	}
	defer response.Body.Close()

	var data GoogleUserInfo
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed read response")
	}

	err = json.Unmarshal(contents, &data)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Failed to parse user info")
	}

	return &data, nil
}
