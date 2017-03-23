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

package main

import (
	"github.com/bretterer/okta-cli/cmd"
	"fmt"
)

func main() {
	fmt.Println(`                                              `)
	fmt.Println(`==============================================`)
	fmt.Println(`                                              `)
	fmt.Println(`  ______    __  ___ .___________.    ___      `)
	fmt.Println(` /  __  \  |  |/  / |           |   /   \     `)
	fmt.Println("|  |  |  | |  '  /  `---|  |----`  /  ^  \\    ")
	fmt.Println(`|  |  |  | |    <       |  |      /  /_\  \   `)
	fmt.Println("|  `--'  | |  .  \\      |  |     /  _____  \\  ")
	fmt.Println(` \______/  |__|\__\     |__|    /__/     \__\ `)
	fmt.Println(`                                              `)
	fmt.Println(`==============================================`)
	fmt.Println(`                                              `)

	cmd.Execute()
}
