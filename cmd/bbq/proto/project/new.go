package project

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/0x00b/gobbq/cmd/bbq/proto/com/base"
	"github.com/0x00b/gobbq/cmd/bbq/proto/tool"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// Project is a project template.
type Project struct {
	Name string
	Path string
}

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string) error {

	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		_ = survey.AskOne(prompt, &override)
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	// fmt.Printf("🚀 Creating service %s, template repo is %s, please wait a moment.\n\n", p.Name, repoURL)
	// repo := base.NewRepo(repoURL, branch)
	// if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
	// 	return err
	// }

	fmt.Printf("🚀 Creating service %s, template repo is %s, please wait a moment.\n\n", p.Name, tool.TemplatePath)
	mod, err := base.ModulePath(path.Join(tool.TemplatePath, "go.mod"))
	if err != nil {
		return err
	}
	err = base.CopyDir(tool.TemplatePath, to, []string{mod, p.Path}, []string{".git", ".github"})
	if err != nil {
		return err
	}

	// os.Rename(
	// 	path.Join(to, "cmd", "server"),
	// 	path.Join(to, "cmd", p.Name),
	// )
	base.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))
	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ make "))
	fmt.Println("			🤝 Thanks for using bbq")
	fmt.Println("	📚 Tutorial: https://iwiki.woa.com/pages/viewpage.action?pageId=966642249")
	return nil
}
