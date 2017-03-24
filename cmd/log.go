// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"time"
	"net/http"
	"io/ioutil"
	"net/url"
	"github.com/go-ini/ini"
	"path/filepath"
	"os/user"
	"encoding/json"
	"github.com/fatih/color"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Tail follow your Okta event log",

	Run: logs,
}

var OktaEvent []struct {
	EventID   string `json:"eventId"`
	SessionID string `json:"sessionId"`
	RequestID string `json:"requestId"`
	Published string `json:"published"`
	Action    struct {
		Message    string `json:"message"`
		Categories []string `json:"categories"`
		ObjectType string `json:"objectType"`
		RequestURI string `json:"requestUri"`
	} `json:"action"`
	Actors    []struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		Login       string `json:"login,omitempty"`
		ObjectType  string `json:"objectType"`
		IPAddress   string `json:"ipAddress,omitempty"`
	} `json:"actors"`
	Targets   []struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		Login       string `json:"login"`
		ObjectType  string `json:"objectType"`
	} `json:"targets"`
}

func init() {
	RootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func logs(cmd *cobra.Command, args []string) {
	usr, _ := user.Current()
	path :=  filepath.Join(usr.HomeDir, `.okta`, `keys.properties`)

	cfg, _ := ini.LooseLoad(path)
	org := cfg.Section("").Key("organization").String()
	domain := cfg.Section("").Key("domain").String()

	OktaOrg := "https://" + org +"."+ domain
	OktaKey := cfg.Section("").Key("apiKey").String()

	lastEvent := ReturnTimeLastEvent(OktaOrg, OktaKey)

	i := 1
	for {
		i += 1
		duration := time.Second * 1
		time.Sleep(duration)
		events := GetOktaEvent(OktaOrg, OktaKey, "filter=published%20gt%20%22" + lastEvent + "%22")
		OktaEvent = nil
		json.Unmarshal([]byte (events), &OktaEvent)
		if (OktaEvent != nil && len (OktaEvent) !=0  ) {
			for v := range OktaEvent {
				fmt.Fprintln(
					os.Stderr,
					color.YellowString("["+OktaEvent[v].Published+"]") + " " +
					OktaEvent[v].Action.Message)

			}

			lastEvent = OktaEvent[len(OktaEvent) - 1].Published
			OktaEvent = nil

		}
	}


}


func ReturnTimeLastEvent(OktaOrg string, OktaKey string) string {

	url := OktaOrg + "/api/v1/events?limit=1&filter=published%20gt%20%222017-12-03T05%3A20%3A48.000Z%22"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "SSWS " + OktaKey)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "690b5379-d5f0-3cff-b1a9-a6a89bc40af4")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	date := string(res.Header.Get("Date"))

	t, err := time.Parse(time.RFC1123, date)

	if err != nil {
		fmt.Println("parse error", err.Error())
	}

	threeHours := time.Hour * 0
	newTime := t.Add(threeHours) // 7 hours actually

	returnString := newTime.Format("2006-01-02T15:04:05") + ".000Z"

	fmt.Fprintln(os.Stderr, color.GreenString("Wait for Events after this Published Date: " + returnString + ". " +
		"Events take " +
		"some time to hit the Event Log"))

	return returnString
}

func GetOktaEvent(OktaOrg string, OktaKey string, arguments string) []byte {

	url := OktaOrg + "/api/v1/events?" + arguments

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "SSWS " + OktaKey)
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "fcd54dc9-bd3b-bdbf-f99a-47272d773855")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
