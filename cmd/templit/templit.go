package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"maps"
	"os"

	"github.com/euforic/templit"
	"github.com/spf13/cobra"
)

// flagValues stores the values of command-line flags
var flagValues = struct {
	token  string
	branch string
}{}

// templitCmd represents the templit command
var templitCmd = &cobra.Command{
	Use:   "templit <command>",
	Short: "A CLI tool for rendering templates from remote repositories",
	Long:  `templit is a CLI tool for rendering templates from remote repositories.`,
	Run: func(cmd *cobra.Command, args []string) {
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
	Run: func(cmd *cobra.Command, args []string) {
		// Check for the correct number of command-line arguments
		if len(args) < 3 {
			if err := cmd.Help(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}

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
		var funcMap = template.FuncMap{
			// return an error for the embed and import functions if no GitHub token is provided
			"embed": func(repoAndPath string, ctx interface{}) (string, error) {
				return "", fmt.Errorf("embed function requires a GitHub token")
			},
			"import": func(repoAndPath string, destPath string, ctx interface{}) (string, error) {
				return "", fmt.Errorf("import function requires a GitHub token")
			},
		}

		if flagValues.token != "" {
			funcMap["embed"] = executor.EmbedFunc
			funcMap["import"] = executor.ImportFunc(outputPath)
		}

		// Copy the default function map from the templit package
		maps.Copy(funcMap, templit.DefaultFuncMap)
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
}

// main is the entrypoint of the application
func main() {
	if err := templitCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
