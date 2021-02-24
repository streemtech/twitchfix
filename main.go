package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommon "github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger = gommon.New("test")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(gommon.INFO)
	e.POST("/oauth2/token", mutateToken)
	e.GET("/oauth2/userinfo", userinfo) //Created for debug purposes.
	e.GET("/oauth2/authorize", auth)    //not possible because it returns html
	e.Logger.Fatal(e.Start(":8284"))
}
func auth(c echo.Context) error {
	req, _ := httputil.DumpRequest(c.Request(), true)
	fmt.Printf("%s\n", string(req))

	bod := c.Request().Body

	defer bod.Close()
	authU := authURL + "?" + c.QueryString()
	//THIS MAY WORK, I Just need to set the content type.
	resp, err := http.Get(authU)
	if err != nil {
		c.Logger().Warnf("Error on post to twitch: %s", err.Error())
		return c.String(500, "error on twitch post")
	}

	if resp.StatusCode != 200 {
		c.Logger().Warnf("Unknown Twitch Status: %s", resp.Status)
		return c.String(500, "Unknown Twitch Status")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Warnf("error on body readAll: %s", err.Error())
		return c.String(500, "No Body On Twitch Post")
	}

	c.Logger().Info("Auth Returned")
	return c.Blob(200, "text/html", body)
}

func userinfo(c echo.Context) error {

	bod := c.Request().Body

	defer bod.Close()
	//THIS MAY WORK, I Just need to set the content type.
	resp, err := http.Get(userinfoURL + c.QueryString())
	if err != nil {
		c.Logger().Warnf("Error on post to twitch: %s", err.Error())
		return c.String(500, "error on twitch post")
	}

	if resp.StatusCode != 200 {
		c.Logger().Warnf("Unknown Twitch Status: %s", err.Error())
		return c.String(500, "Unknown Twitch Status")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Warnf("error on body readAll: %s", err.Error())
		return c.String(500, "No Body On Twitch Post")
	}

	c.Logger().Info("Userinfo Returned")

	return c.Blob(200, "application/x-www-form-urlencoded", body)
}

var userinfoURL = "https://id.twitch.tv/oauth2/userinfo"
var authURL = "https://id.twitch.tv/oauth2/authorize"
var twitchURL = "https://id.twitch.tv/oauth2/token"

type twitchResp struct {
	TokenType    string   `json:"token_type"`
	IDToken      string   `json:"id_token"`
	RefreshToken string   `json:"refresh_token"`
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	Scope        []string `json:"scope"`
}

type rfcResp struct {
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func mutateToken(c echo.Context) error {

	if c.Logger().Level() == gommon.DEBUG {
		req, _ := httputil.DumpRequest(c.Request(), true)
		fmt.Printf("%s\n", string(req))
	}
	bod := c.Request().Body

	defer bod.Close()
	//THIS MAY WORK, I Just need to set the content type.
	resp, err := http.Post(twitchURL, "application/x-www-form-urlencoded", bod)
	if err != nil {
		c.Logger().Warnf("Error on post to twitch: %s", err.Error())
		return c.String(500, "error on twitch post")
	}

	if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return c.String(500, "Unknown Twitch Status")
		}

		c.Logger().Warnf("Unknown Twitch Status: %s. DATA: %s", resp.Status, bodyBytes)
		return c.String(500, "Unknown Twitch Status, failed to read body.")
	}

	if resp.Body == nil {
		c.Logger().Warnf("Null Body")
		return c.String(500, "No Body On Twitch Post")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Warnf("error on body readAll: %s", err.Error())
		return c.String(500, "No Body On Twitch Post")
	}

	rfc, err := getRFCBody(body)
	if err != nil {
		c.Logger().Warnf("Unable to parse the body from the twitch response: %s", err.Error())
		return c.String(500, "Unable to parse the body from the twitch response: "+err.Error())
	}

	c.Logger().Info("RFC Returned")
	return c.JSON(200, rfc)

}

func getRFCBody(body []byte) (rfcResp, error) {

	twitch := twitchResp{}
	rfc := rfcResp{}

	err := json.Unmarshal(body, &twitch)
	if err != nil {
		return rfc, err
	}

	rfc.AccessToken = twitch.AccessToken
	rfc.ExpiresIn = twitch.ExpiresIn
	rfc.IDToken = twitch.IDToken
	rfc.RefreshToken = twitch.RefreshToken
	rfc.TokenType = twitch.TokenType
	rfc.Scope = makeScopes(twitch.Scope)

	return rfc, nil

}

func makeScopes(scopes []string) string {

	if len(scopes) == 0 {
		return ""
	}
	if len(scopes) == 1 {
		return scopes[0]
	}

	str := ""
	for i := 0; i < len(scopes)-1; i++ {
		str += scopes[i] + " "
	}
	str += scopes[len(scopes)-1]
	return str
}
