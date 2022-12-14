package project

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a service template",
	Long:  "Create a service project using the repository template. Example: bbq new helloworld",
	Run:   run,
}

// var repoURL string
// var branch string

func init() {
	// if repoURL = os.Getenv("BUTC_LAYOUT_REPO"); repoURL == "" {
	// 	repoURL = "https://git.code.tencent.com/butchery/butchery-template.git"
	// }
	// CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	// CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		_ = survey.AskOne(prompt, &name)
		if name == "" {
			return
		}
	} else {
		name = args[0]
	}

	p := &Project{Name: path.Base(name), Path: name}
	if err := p.New(ctx, wd); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		return
	}
}
