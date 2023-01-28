package tool

import (
	"context"
	"fmt"
	"time"

	"github.com/0x00b/gobbq/cmd/bbq/proto/com/base"
	"github.com/spf13/cobra"
)

// TemplatePath 模版路径
const TemplatePath = "~/.gobbq/template"

// CmdInstall represents the new command.
var CmdInstall = &cobra.Command{
	Use:   "install",
	Short: "install tools",
	Long:  "install tools. Example: bbq install",
	Run:   install,
}

func install(cmd *cobra.Command, args []string) {
	err := base.GoInstall(
		"google.golang.org/protobuf/cmd/protoc-gen-go",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
		"github.com/envoyproxy/protoc-gen-validate",
		"github.com/0x00b/gobbq/cmd/protoc-gen-gobbq",
		// "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway",
		// "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",
	)
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	repo := base.NewRepo("https://github.com/0x00b/protoc-gen-gobbq.git", "")
	if err := repo.CopyFilesTo(ctx, "/usr/local/.gobbq/grpc-go-tpl",
		[]string{"grpc-go-tpl"}); err != nil {
		panic(err)
	}
	tplRepo := base.NewRepo("https://git.code.tencent.com/gobbq/gobbq-template.git", "")
	if err := tplRepo.CopyFilesTo(ctx, TemplatePath,
		[]string{""}); err != nil {
		panic(err)
	}

	giRepo := base.NewRepo("https://github.com/googleapis/googleapis.git", "")
	if err := giRepo.CopyFilesTo(ctx, "/usr/local/.gobbq/api/googleapis/google",
		[]string{"google"}); err != nil {
		panic(err)
	}
	pbRepo := base.NewRepo("https://github.com/protocolbuffers/protobuf.git", "")
	if err := pbRepo.CopyFilesTo(ctx, "/usr/local/.gobbq/api/protobuf",
		[]string{"src"}); err != nil {
		panic(err)
	}
	gwRepo := base.NewRepo("https://github.com/grpc-ecosystem/grpc-gateway.git", "")
	if err := gwRepo.CopyFilesTo(ctx, "/usr/local/.gobbq/api/grpc-gateway",
		[]string{""}); err != nil {
		panic(err)
	}
}
