package gogen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/0x00b/gobbq/cmd/bbq/proto/com"
	"github.com/0x00b/gobbq/cmd/bbq/proto/com/base"
	"github.com/0x00b/gobbq/cmd/bbq/proto/com/gorewriter/rewrite"
)

// GoGenerator generate go code
type GoGenerator struct {
	// Rewriter    map[*com.FileDescriptorProto]*rewrite.Rewriter
	RootPackage string
	com.RPC
}

// NewGoGenerator init a GoGenerator
//
//	@param pkg
//	@return *GoGenerator
//	@return error
func NewGoGenerator(rootPackage string, rpc com.RPC) (gg *GoGenerator, err error) {
	rootPackage, err = filepath.Abs(rootPackage)
	if err != nil {
		return nil, err
	}

	gg = &GoGenerator{
		RPC:         rpc,
		RootPackage: rootPackage,
	}

	// gg.Rewriter = make(map[*com.FileDescriptorProto]*rewrite.Rewriter)

	return gg, nil
}

// Generate TODO
func (g *GoGenerator) Generate(tplRoot string, proto *com.Proto) error {

	// fmt.Println(tplRoot)

	_ = filepath.Walk(tplRoot, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".tpl") && !strings.Contains(path, ".all.tpl") {
			tplInstance, err := template.New(info.Name()).Funcs(com.FuncMap).ParseFiles(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			for _, f := range proto.Files {
				if !com.AStringContains(proto.Plugin.Request.FileToGenerate, *f.Name) {
					continue
				}
				//需要更加go.mod找到文件路径
				name, _ := initGenerateDir(tplRoot, path, f)
				fmt.Printf("filename:[%s]\n", name)

				_ = g.setFileImpl(name, f)

				var b bytes.Buffer
				err = tplInstance.Execute(&b, f)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return err
				}
				g := proto.Plugin.NewGeneratedFile(name, "")
				g.P(b.String())
			}
		}
		return nil
	})

	_ = filepath.Walk(tplRoot, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".all.tpl") {
			tplInstance, err := template.New(info.Name()).Funcs(com.FuncMap).ParseFiles(path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			name, _ := initGenerateAllDir(tplRoot, path)
			proto.GoRewriter, err = rewrite.New(filepath.Dir(name))
			if err != nil {
				return err
			}

			var b bytes.Buffer
			err = tplInstance.Execute(&b, proto)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			g := proto.Plugin.NewGeneratedFile(name, "")
			g.P(b.String())
		}
		return nil
	})
	return nil
}

func (g *GoGenerator) setFileImpl(name string,
	f *com.File) error {
	var err error
	f.GoRewriter, err = rewrite.New(filepath.Dir(name))
	if err != nil {
		return err
	}

	f.GoImplImportPaths = f.GoRewriter.ExistingImports(name)

	return nil
}

func initGenerateAllDir(root, tplPath string) (fileName string, e error) {
	path := com.TrimRight(com.TrimLeft(com.TrimLeft(tplPath, root), "/"), ".all.tpl")

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0555)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}

func initGenerateDir(root, tplPath string, file *com.File) (fileName string, e error) {

	mod := base.ModName()
	// fmt.Println("mod:", mod)
	goPath := com.TrimLeft(file.GoImportPath, mod)
	tplRelativePath := com.TrimLeft(com.TrimLeft(filepath.Dir(tplPath), root), "/")
	file.GoImplPackage = mod + "/" + tplRelativePath + goPath

	path := tplRelativePath + goPath
	fmt.Printf("tplRelativePath:[%s]\n", tplRelativePath)
	fmt.Printf("goPath:[%s]\n", goPath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0555)
		if err != nil {
			fmt.Println("err:", err)
			return "", err
		}
	}

	name := filepath.Base(*file.Name)
	// fmt.Printf("name:[%s]\n", name)

	return path + "/" + com.TrimRight(name, filepath.Ext(name)) + com.TplSuffix(tplPath), nil
}
