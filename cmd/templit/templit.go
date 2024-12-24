//nolint:gochecknoglobals,gochecknoinits
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"maps"
	"os"

	"github.com/euforic/templit"
	"github.com/spf13/cobra"
)

var ErrMissingToken = errors.New("embed and import functions requires a GitHub token")

// flagValues stores the values of command-line flags
var flagValues = struct {
	token  string
	branch string
	remote string
}{}

// templitCmd represents the templit command
var templitCmd = &cobra.Command{
	Use:   "templit <command>",
	Short: "A CLI tool for rendering templates from remote repositories",
	Long:  `templit is a CLI tool for rendering templates from remote repositories.`,
	Args:  cobra.MinimumNArgs(3), //nolint:mnd
	Run: func(cmd *cobra.Command, _ []string) {
		if err := cmd.Help(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render <inputPath> <outputPath> <jsonData>",
	Short: "A CLI tool for rendering templates from remote repositories",
	Long:  `generate is for rendering templates.`,
	Run: func(_ *cobra.Command, args []string) {
		if flagValues.token == "" {
			flagValues.token = os.Getenv("GIT_TOKEN")
		}

		// Extract the command-line arguments
		inputPath, outputPath, inputData := args[0], args[1], args[2]

		// values stores the JSON data
		var values map[string]interface{}

		// Parse the JSON data from the command-line argument
		if err := json.Unmarshal([]byte(inputData), &values); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON data: %s\n", err)

			return
		}

		// executor is the template executor
		executor := templit.NewExecutor(templit.NewDefaultGitClient(flagValues.branch, flagValues.token))

		// funcMap defines the custom functions that can be used in templates
		funcMap := template.FuncMap{
			// return an error for the embed and import functions if no GitHub token is provided
			"embed": func(string, interface{}) (string, error) {
				return "", ErrMissingToken
			},
			"import": func(string, string, interface{}) (string, error) {
				return "", ErrMissingToken
			},
		}

		if flagValues.token != "" {
			funcMap["embed"] = executor.EmbedFunc
			funcMap["import"] = executor.ImportFunc(outputPath)
		}

		// Copy the default function map from the templit package
		executor.Funcs(templit.DefaultFuncMap())
		executor.Funcs(funcMap)

		// If a remote repository is specified, process the template and write it to the output directory
		if flagValues.remote != "" {
			importParts, err := templit.ParseDepURL(flagValues.remote)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing remote: %s\n", err)

				return
			}

			importParts.Path = inputPath

			if _, err := executor.ImportFunc(outputPath)(importParts.String(), "./", values); err != nil {
				fmt.Fprintf(os.Stderr, "Error processing template: %s\n", err)
			}

			return
		}

		// Copy the default function map from the templit package
		maps.Copy(funcMap, templit.DefaultFuncMap())
		executor.Funcs(funcMap)

		// Process the templates in the input directory and write them to the output directory
		if err := executor.WalkAndProcessDir(inputPath, outputPath, values); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing template: %s\n", err)
		}
	},
}

func init() {
	templitCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVarP(&flagValues.token, "git_token", "t", "", "GitHub token")
	renderCmd.Flags().StringVarP(&flagValues.branch, "branch", "b", "main", "GitHub branch")
	renderCmd.Flags().StringVarP(&flagValues.remote, "remote", "r", "", "remote repository to use. (example: github.com/owner/repo@ref)")
}

// main is the entrypoint of the application
func main() {
	if err := templitCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
