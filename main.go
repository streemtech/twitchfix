package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommon "github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.SetLevel(gommon.DEBUG)
	e.POST("/oauth2/token", mutateToken)
	e.Logger.Fatal(e.Start(":8284"))
}

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

	bod := c.Request().Body

	defer bod.Close()
	//THIS MAY WORK, I Just need to set the content type.
	resp, err := http.Post(twitchURL, "application/x-www-form-urlencoded", bod)
	if err != nil {
		return c.String(500, "error on twitch post")
	}

	if resp.StatusCode != 200 {
		return c.String(500, "Unknown Twitch Status")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.String(500, "No Body On Twitch Post")
	}

	rfc, err := getRFCBody(body)
	if err != nil {
		return c.String(500, "Unable to parse the body from the twitch response: "+err.Error())
	}

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
