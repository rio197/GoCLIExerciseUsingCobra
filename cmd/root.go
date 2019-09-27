/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var name string
var greeting string
var preview bool
var prompt bool
var debug bool = false

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "motd",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Show usage if flags are invalid
		if prompt == false && (name == "" || greeting == "") {
			cmd.Usage()
			os.Exit(1)
		}

		// Optionally print flags and exit based on DEBUG environment variable
		if debug {
			fmt.Println("Name", name)
			fmt.Println("Greeting:", greeting)
			fmt.Println("Prompt:", prompt)
			fmt.Println("Preview:", preview)

			os.Exit(0)
		}

		// Conditionally read from stdin
		if prompt {
			name, greeting = renderPrompt()
		}

		// Generate message
		message := buildMessage(name, greeting)

		// Either preview the message or write to file
		if preview {
			fmt.Println(message)
		} else {
			// write content
			f, err := os.OpenFile("/etc/motd", os.O_WRONLY, 0644)

			if err != nil {
				fmt.Println("Error: Unable to open /etc/motd")
				os.Exit(1)
			}

			defer f.Close()

			_, err = f.Write([]byte(message))

			if err != nil {
				fmt.Println("Error: Failed to write to /etc/motd")
				os.Exit(1)
			}
		}
	},
}

func buildMessage(name, greeting string) string {
	return fmt.Sprintf("%s, %s", greeting, name)
}

func renderPrompt() (name, greeting string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your Greeting: ")
	greeting, _ = reader.ReadString('\n')
	greeting = strings.TrimSpace(greeting)

	fmt.Print("Your Name: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	return
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Name to use in this message")
	rootCmd.Flags().StringVarP(&greeting, "greeting", "g", "", "Greeting to use in this message")
	rootCmd.Flags().BoolVarP(&preview, "preview", "v", false, "Preview message instead of writing to /etc/motd")
	rootCmd.Flags().BoolVarP(&prompt, "prompt", "p", false, "Prompt for name and greeting")

	if os.Getenv("DEBUG") != "" {
		debug = true
	}
}
