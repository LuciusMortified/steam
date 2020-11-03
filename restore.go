package steam

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type CookieData struct {
	TimezoneOffset     string `json:"timezone_offset"`
	SessionID          string `json:"session_id"`
	SteamCountry       string `json:"steam_country"`
	SteamLoginSecure   string `json:"steam_login_secure"`
	SteamMachineAuth   string `json:"steam_machine_auth"`
	SteamRememberLogin string `json:"steam_remember_login"`
}

type SessionData struct {
	SteamID     uint64     `json:"steam_id"`
	SessionID   string     `json:"session_id"`
	DeviceID    string     `json:"device_id"`
	UmqID       string     `json:"umq_id"`
	Token       string     `json:"token"`
	LoginSecure string     `json:"login_secure"`
	WebCookie   string     `json:"webcookie"`
	ApiKey      string     `json:"api_key"`
	ChatMessage int        `json:"chat_message"`
	Language    string     `json:"language"`
	Cookies     CookieData `json:"cookies"`
}

var (
	ErrRestoreClientNil = errors.New("restore client is nil")
	ErrRestoreDataNil   = errors.New("restore data is nil")
)

func (session *Session) Dump() (*SessionData, error) {
	steamUrl, err := url.Parse(httpBaseUrl)
	if err != nil {
		return nil, err
	}

	cookies := session.client.Jar.Cookies(steamUrl)
	cookieData := CookieData{}

	for _, cookie := range cookies {
		switch cookie.Name {
		case "sessionid":
			cookieData.SessionID = cookie.Value
		case "timezoneOffset":
			cookieData.TimezoneOffset = cookie.Value
		case "steamCountry":
			cookieData.SteamCountry = cookie.Value
		case "steamLoginSecure":
			cookieData.SteamLoginSecure = cookie.Value
		case fmt.Sprintf("steamMachineAuth%d", session.oauth.SteamID):
			cookieData.SteamMachineAuth = cookie.Value
		case "steamRememberLogin":
			cookieData.SteamRememberLogin = cookie.Value
		}
	}

	sessionData := SessionData{
		Cookies: cookieData,

		//OAuth
		SteamID:     uint64(session.oauth.SteamID),
		Token:       session.oauth.Token,
		LoginSecure: session.oauth.LoginSecure,
		WebCookie:   session.oauth.WebCookie,

		//Session
		SessionID:   session.sessionID,
		DeviceID:    session.deviceID,
		UmqID:       session.umqID,
		ApiKey:      session.apiKey,
		ChatMessage: session.chatMessage,
		Language:    session.language,
	}

	return &sessionData, nil
}

func RestoreSession(client *http.Client, data *SessionData, debug bool) (*Session, error) {
	if client == nil {
		return nil, ErrRestoreClientNil
	}

	if data == nil {
		return nil, ErrRestoreDataNil
	}

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	steamUrl, err := url.Parse(httpBaseUrl)
	if err != nil {
		return nil, err
	}

	cookieJar.SetCookies(steamUrl, []*http.Cookie{
		{Name: "sessionid", Value: data.Cookies.SessionID},
		{Name: "timezoneOffset", Value: data.Cookies.TimezoneOffset},
		{Name: "steamCountry", Value: data.Cookies.SteamCountry},
		{Name: "steamLoginSecure", Value: data.Cookies.SteamLoginSecure},
		{Name: fmt.Sprintf("steamMachineAuth%d", data.SteamID), Value: data.Cookies.SteamMachineAuth},
		{Name: "steamRememberLogin", Value: data.Cookies.SteamRememberLogin},
	})

	client.Jar = cookieJar

	oauth := OAuth{
		SteamID:     SteamID(data.SteamID),
		Token:       data.Token,
		LoginSecure: data.LoginSecure,
		WebCookie:   data.WebCookie,
	}

	session := Session{
		client:      client,
		oauth:       oauth,
		sessionID:   data.SessionID,
		apiKey:      data.ApiKey,
		deviceID:    data.DeviceID,
		umqID:       data.UmqID,
		chatMessage: data.ChatMessage,
		language:    data.Language,
		debug:       debug,
	}

	return &session, nil
}
