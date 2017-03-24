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
	"os/user"

	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"io/ioutil"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize the Okta CLI tool",
	Long: ``,
	Run: run,
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func run(cmd *cobra.Command, args []string) {

	usr, err := user.Current()


	if err != nil {
		panic(err.Error())
	}
	path := filepath.Join(usr.HomeDir, `.okta`, `keys.properties`)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		guidedInit()
	}

	fmt.Println("You are ready to use the Okta CLI!")

}

func guidedInit() bool {
	apiKey := askForApiKey()
	organization := askForOrganization()
	domain := askForDomain()

	confirmCreation := confirmConfigFileCreation(apiKey, organization, domain)


	if(confirmCreation) {
		usr, _ := user.Current()
		content := []byte(createContent(apiKey, organization, domain))
		os.MkdirAll(filepath.Join(usr.HomeDir, `.okta`), os.ModePerm)
		err := ioutil.WriteFile(filepath.Join(usr.HomeDir, `.okta`, `keys.properties`), content, 0644)

		if err != nil {
			panic(err.Error())
		}
	}

	return true

}
func confirmConfigFileCreation(apiKey string, organization string, domain string) bool {
	var response string
	fmt.Println("")
	fmt.Println("Generated File")
	fmt.Println("")

	contentToCreate := createContent(apiKey, organization, domain)

	fmt.Print(contentToCreate)

	fmt.Println("")

	fmt.Print("Is this information correct? [y]:  ")
	_, err := fmt.Scanln(&response)
	if err != nil {
		response = "y"
	}

	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}

	if containsString(okayResponses, response) {
		return true
	} else {
		return guidedInit()
	}
}
func createContent(apiKey string, organization string, domain string) string {
	return "" +
		"apiKey="+apiKey+"\n" +
		"organization="+organization+"\n" +
		"domain="+domain+"\n"
}

func askForApiKey() string {
	var response string

	fmt.Print("What is your api key? []:  ")

	_, err := fmt.Scanln(&response)
	if err != nil {
		panic(err.Error())
	}

	return response

}

func askForOrganization() string {
	var response string

	fmt.Print("What is your organization? (i.e. dev-28957684) []:  ")

	_, err := fmt.Scanln(&response)
	if err != nil {
		panic(err.Error())
	}

	return response

}

func askForDomain() string {
	var response string

	fmt.Print("What is your Okta domain? [oktapreview.com]:  ")

	fmt.Scanln(&response)

	if response == "" {
		response = "oktapreview.com"
	}

	return response
}


func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}