package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// CmdGen represents the source command.
	CmdGen = &cobra.Command{
		Use:   "gen",
		Short: "Generate the proto code",
		Long:  "Generate the proto code. Example: bbq proto gen helloworld.proto",
		Run:   run,
	}
)

var protoPath []string

func init() {
	CmdGen.Flags().StringArrayVarP(&protoPath, "proto_path", "I", protoPath, "proto path")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please enter the proto file or directory")
		return
	}
	var (
		err   error
		proto = strings.TrimSpace(args[0])
	)
	if err = look("protoc-gen-go", "protoc-gen-go-grpc", "protoc-gen-gobbq", "protoc-gen-validate"); err != nil {
		// update the bbq plugins
		cmd := exec.Command("bbq", "install")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}
	if strings.HasSuffix(proto, ".proto") {
		err = generateRPC(proto, args)
	} else {
		err = walk(proto, args)
	}
	if err != nil {
		fmt.Println(err)
	}
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}

	var protoFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if ext := filepath.Ext(path); ext != ".proto" {
			return nil
		}
		protoFiles = append(protoFiles, path)
		return generateRPC(path, args)
	})
	if err != nil {
		return err
	}
	return generateGobbq(args, protoFiles)
}

// generate is used to execute the generate command for the specified proto file
func generateRPC(proto string, args []string) error {
	input := []string{
		"--proto_path=.",
	}
	for _, pp := range protoPath {
		if pathExists(pp) {
			input = append(input, "--proto_path="+pp)
		}
	}

	for _, pp := range shouldIncludePath(proto) {
		if pathExists(pp) {
			input = append(input, "--proto_path="+pp)
		}
	}

	inputExt := []string{
		// "--proto_path=" + base.ButcMod(),
		"--go_out=paths=source_relative:.",
		// "--go-grpc_out=paths=source_relative:.",
		// "--grpc-gateway_out=paths=source_relative:.",
		// "--gobbq_out=plugins=grpc,tpl_dir=/root/code/protoc-gen/protoc-gen-gobbq/gogen/tpl:.",
	}
	input = append(input, inputExt...)
	protoBytes, err := ioutil.ReadFile(proto)
	if err == nil && len(protoBytes) > 0 {
		if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
			input = append(input, "--validate_out=lang=go,paths=source_relative:.")
		}
	}
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}

// generate is used to execute the generate command for the specified proto file
func generateGobbq(args []string, protoFiles []string) error {
	input := []string{
		"--proto_path=.",
	}
	for _, pp := range protoPath {
		if pathExists(pp) {
			input = append(input, "--proto_path="+pp)
		}
	}

	for _, pp := range shouldIncludePath(protoFiles...) {
		if pathExists(pp) {
			input = append(input, "--proto_path="+pp)
		}
	}

	inputExt := []string{
		// "--proto_path=" + base.ButcMod(),
		// "--go_out=paths=source_relative:.",
		// "--go-grpc_out=paths=source_relative:.",
		// "--grpc-gateway_out=paths=source_relative:.",
		"--gobbq_out=plugins=grpc,tpl_dir=/usr/local/.gobbq/grpc-go-tpl:.",
	}
	input = append(input, inputExt...)
	// protoBytes, err := ioutil.ReadFile(proto)
	// if err == nil && len(protoBytes) > 0 {
	// 	if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
	// 		input = append(input, "--validate_out=lang=go,paths=source_relative:.")
	// 	}
	// }
	input = append(input, protoFiles...)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}

	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func shouldIncludePath(protoFiles ...string) []string {

	paths := make(map[string]bool)

	regex, err := regexp.Compile(`\n[^/]*(import)\s+"google/api/`)
	if err != nil {
		return nil
	}
	for _, protoFile := range protoFiles {
		protoBytes, err := ioutil.ReadFile(protoFile)
		if err == nil && len(protoBytes) > 0 {
			if ok := regex.Match(protoBytes); ok {
				paths["/usr/local/.gobbq/api/googleapis"] = true
				paths["/usr/local/.gobbq/api/googleapis/google"] = true
				paths["/usr/local/.gobbq/api/protobuf"] = true
				// paths["/usr/local/.gobbq/api/grpc-gateway"] = true
			}
		}
	}

	var res []string
	for p := range paths {
		res = append(res, p)
	}

	return res
}
