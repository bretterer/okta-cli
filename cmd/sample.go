package cmd

import (
	"github.com/spf13/cobra"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)


// sampleCmd represents the log command
var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Install python-flask sample",

	Run: installSample,
}


func init() {
	RootCmd.AddCommand(sampleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}


// FIXME: So ugly ...

var base_url string
var ssws string

func readConfig() {
	usr, _ := user.Current()
	path :=  filepath.Join(usr.HomeDir, `.okta`, `keys.properties`)
	
	cfg, _ := ini.LooseLoad(path)
	org := cfg.Section("").Key("organization").String()
	domain := cfg.Section("").Key("domain").String()
	
	base_url = "https://" + org +"."+ domain
	ssws = cfg.Section("").Key("apiKey").String()
}

type HALLink struct {
	Href string `json:"href"`
}

type OAuthClient struct {
	ClientID string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	ClientIDIssuedAt int `json:"client_id_issued_at,omitempty"`
	ClientSecretExpiresAt int `json:"client_secret_expires_at,omitempty"`
	ClientName string `json:"client_name"`
	RedirectURIs []string `json:"redirect_uris"`
	ResponseTypes []string `json:"response_types"`
	GrantTypes []string `json:"grant_types"`
	TokenEndpointAuthMethod string `json:"token_endpoint_auth_method"`
	ApplicationType string `json:"application_type"`
	Links map[string]HALLink `json:"_links,omitempty"`
}

type OpenIdConfiguration struct {
	JwksURI string `json:"jwks_uri"`
}

type OktaUserProfile struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Login string `json:"login"`
}

type OktaUserCredentials struct {
	Password struct {
		Value string `json:"value"`
	} `json:"password"`
	RecoveryQuestion struct {
		Question string `json:"question"`
		Answer string `json:"answer"`
	} `json:"recovery_question"`
}

type OktaUser struct {
	Profile OktaUserProfile `json:"profile"`
	Credentials  OktaUserCredentials `json:"credentials"`
}

type OktaApplication struct {
	ID string `json:"id"`
	Label string `json:"label"`
	SignOnMode string `json:"signOnMode"`
}

type OktaResult struct {
	ID string `json:"id"`
}

type OktaApps struct {
	Applications []OktaApplication
}

func oktaReq(path string) []byte {
	client := &http.Client {}
	url := base_url + path
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "SSWS " + ssws)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return bodyBytes
}

func createUser(user_email string) string {
	fmt.Printf("User %s doesn't exist, creating it now\n", user_email)

	path := "/api/v1/users?activate=true"

	first_name := "Paul"
	last_name := "Cook"
	password := "L0rdn1k0n"
	recovery_question := "What's this one eat?"
	recovery_answer := "It nibbles. You see this?"
	
	m := OktaUser{
		Profile: OktaUserProfile{
			FirstName: first_name,
			LastName: last_name,
			Email: user_email,
			Login: user_email,
		},
		Credentials: OktaUserCredentials{
			Password: struct{
				Value string `json:"value"`
			}{password},
			RecoveryQuestion: struct {
				Question string `json:"question"`
				Answer string `json:"answer"`
			}{recovery_question, recovery_answer,},
		},
	}
	payload, err := json.Marshal(m)
	if err != nil {
		panic(err.Error())
	}

	payloadReader := bytes.NewReader(payload)

	url := base_url + path
	client := &http.Client {}
	req, _ := http.NewRequest("POST", url, payloadReader)
	req.Header.Add("Authorization", "SSWS " + ssws)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	
	var target OktaResult
	err = json.Unmarshal(bodyBytes, &target)
	if err != nil {
		panic(err.Error())
	}
	return target.ID
}

// FIXME: Change to "getUserID"
func getUserID(user string) string {
	// paul.cook
	var target []OktaResult
	path := "/api/v1/users?q=" + user
	body := oktaReq(path)
	err := json.Unmarshal(body, &target)
	if err != nil {
		panic(err.Error())
	}
	// fmt.Printf("%#v\n", target)
	if len(target) < 1 {
		return createUser(user)
	}
	
	if len(target) == 1 {
		return target[0].ID
	} else {
		panic("Returned more than 1 user")
	}
}

func getApps() {
	var target []OktaApplication

	body := oktaReq("/api/v1/apps?limit=2")
	err := json.Unmarshal(body, &target)
	if err != nil {
		panic(err.Error())
	}
	
	// fmt.Printf("%#v\n", target)
}

func createThing() {
}

func createOktaApp() (string, string) {

	path := "/oauth2/v1/clients"
	m := OAuthClient{
		ClientName: "Okta CLI App",
		TokenEndpointAuthMethod: "client_secret_basic",
		ApplicationType: "web",
		RedirectURIs: []string{"http://localhost:3000/authorization-code/callback"},
		ResponseTypes: []string{"code", "token", "id_token"},
		GrantTypes: []string{"refresh_token", "authorization_code", "implicit"},
	}
	payload, err := json.Marshal(m)
	// fmt.Println(string(payload))
	if err != nil {
		panic(err.Error())
	}
	
	payloadReader := bytes.NewReader(payload)

	url := base_url + path
	
	client := &http.Client {}
	req, _ := http.NewRequest("POST", url, payloadReader)
	req.Header.Add("Authorization", "SSWS " + ssws)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var target OAuthClient
	err = json.Unmarshal(bodyBytes, &target)
	if err != nil {
		panic(err.Error())
	}
	// fmt.Println(target.ClientID)
	// fmt.Println(target.ClientSecret)
	app_href := target.Links["app"].Href

	everyone_group_id := groupIdForEveryoneGroup()
	app_to_group_href := app_href + "/groups/" + everyone_group_id
	// fmt.Println(app_to_group_href)
	
	// http://developer.okta.com/docs/api/resources/apps.html#assign-group-to-application

	url = app_to_group_href
	req, _ = http.NewRequest("PUT", url, strings.NewReader("{}"))
	req.Header.Add("Authorization", "SSWS " + ssws)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	
	return target.ClientID, target.ClientSecret
}

func resetOktaApp(app_id string) (string, string) {
	return "FAKE_CLIENT_ID", "FAKE_CLIENT_SECRET"
}

func groupIdForEveryoneGroup() string {
	path := "/api/v1/groups?q=Everyone"


	var target []OktaResult
	body := oktaReq(path)
	err := json.Unmarshal(body, &target)
	if err != nil {
		panic(err.Error())
	}
	// fmt.Printf("%#v\n", target)
	if len(target) == 1 {
		return target[0].ID
	} else {
		panic("Returned more than 1 user")
	}
}

func appsForUserID(user_id string) []string {
	path := fmt.Sprintf("/api/v1/apps?filter=user.id+eq+\"%s\"", user_id)
	
	var target []OktaApplication
	var found []string
	
	body := oktaReq(path)
	err := json.Unmarshal(body, &target)
	if err != nil {
		panic(err.Error())
	}

	for _, app := range target {
		if app.SignOnMode != "OPENID_CONNECT" {
			continue
		}
		if app.Label != "Okta CLI App" {
			continue
		}
		found = append(found, app.ID)
	}
	return found
}

func userSetup() {
	// Check for user, create user if they don't exist
	user_id := getUserID("joe.user@example.com")
	// Get apps assigned to user
	app_ids := appsForUserID(user_id)
	// fmt.Printf("%#v\n", app_ids)
	// If user has no apps assigned, then create one
	var client_id string
	var client_secret string
	
	if len(app_ids) < 1 {
		fmt.Println("Creating App")
		client_id, client_secret = createOktaApp()
	} else if len(app_ids) == 1 {
		fmt.Println("Resetting App")
		client_id, client_secret = resetOktaApp(app_ids[0])
	} else {
		panic("Unexpected error when finding apps")
	}

	oidc_config := fmt.Sprintf("clientId=%s\nclientSecret=%s\n", client_id, client_secret)
	// fmt.Println(oidc_config)
	usr, _ := user.Current()
	filename := filepath.Join(usr.HomeDir, `.okta`, `keys.properties`)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(oidc_config); err != nil {
		panic(err)
	}
}

func downloadAndUnzip() {
	curl, err := exec.Command("curl", "-OL", "https://github.com/okta/samples-python-flask/archive/2017.13-begin.zip").Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(curl)

	fmt.Println("Unzipping file")
	unzip, unzipErr := exec.Command("unzip", "-o", "2017.13-begin.zip").Output()
	// unzip, unzipErr := exec.Command("pwd").Output()
	if unzipErr != nil {
		panic(unzipErr)
	}
	fmt.Println(string(unzip))
}

func makeConfig(client_id string, client_secret string) string {
	out := fmt.Sprintf(`{
  "oktaSample": {
    "oidc": {
      "oktaUrl": "%s",
      "clientId": "%s",
      "clientSecret": "%s",
      "redirectUri": "http://localhost:3000/authorization-code/callback"
    }
  }
}`, base_url, client_id, client_secret)
	return out
}

func configApp(dirname string) {
	usr, _ := user.Current()
	path :=  filepath.Join(usr.HomeDir, `.okta`, `keys.properties`)

	ccfg, _ := ini.LooseLoad(path)
	
	client_id := ccfg.Section("").Key("clientId").String()
	client_secret := ccfg.Section("").Key("clientSecret").String()
	content := makeConfig(client_id, client_secret)

	filename := filepath.Join(dirname, ".samples.config.json")
	err := ioutil.WriteFile(filename, []byte(content), 0644)

	if err != nil {
		panic(err.Error())
	}

}

func installSample(cmd *cobra.Command, args []string) {
	readConfig()
	createUser("paul.cook@example.com")
	userSetup()
	downloadAndUnzip()
	configApp("samples-python-flask-2017.13-begin")
}


