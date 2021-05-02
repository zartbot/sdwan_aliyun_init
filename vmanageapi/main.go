package vmanageapi

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/sirupsen/logrus"
)

type VManageAPI struct {
	Client   *http.Client
	Username string
	Password string
	BaseURL  string
}

type VmanageInfo struct {
	BaseURL  string `yaml:baseurl`
	Username string `yaml:username`
	Password string `yaml:password`
}

func NewVmanage(baseurl string, username string, password string) *VManageAPI {
	jar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	app := VManageAPI{
		Client:   &http.Client{Jar: jar, Transport: tr},
		Username: username,
		Password: password,
		BaseURL:  "https://" + baseurl,
	}
	app.login()
	return &app

}
func (app *VManageAPI) login() {
	client := app.Client
	loginURL := app.BaseURL + "/j_security_check"
	data := url.Values{
		"j_username": {app.Username},
		"j_password": {app.Password},
	}
	response, err := client.PostForm(loginURL, data)
	if err != nil {
		logrus.Fatal(err)
	}
	defer response.Body.Close()

	/*	fmt.Println("HTTP Response Status:", response.StatusCode, http.StatusText(response.StatusCode))
		if response.StatusCode >= 200 && response.StatusCode <= 299 {
			logrus.Warn("HTTP Status is in the 2xx range")
		} else {
			fmt.Println("Argh! Broken")
		}*/
}
