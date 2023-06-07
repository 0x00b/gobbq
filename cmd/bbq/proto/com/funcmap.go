package com

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/0x00b/gobbq/cmd/bbq/proto/com/gorewriter/rewrite"
	"github.com/0x00b/gobbq/proto/bbq"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/compiler/protogen"
	pb "google.golang.org/protobuf/proto"
)

// FuncMap 模版中使用的函数列表
var FuncMap = template.FuncMap{
	// for go
	"goComments":         GoComments,
	"goCamelcase":        GoCamelCase,
	"goCamelcaseType":    GoCamelCaseModel,
	"goGetMethodBody":    GoGetMethodBody,
	"goGetImportPaths":   GoGetImportPaths,
	"goExistImportPaths": GoExistImportPaths,
	"gopkg":              PBGoPackage,
	"gopkgSimple":        PBValidGoPackage,
	"goType":             PBGoType,
	"export":             GoExport,
	"simplify":           PBSimplifyGoType,
	"FileName":           FileName,

	// comm func
	"camelcase":      Camelcase,
	"lowerCamelcase": strcase.ToLowerCamel,
	"title":          Title,
	"untitle":        UnTitle,
	"trimright":      TrimRight,
	"trimleft":       TrimLeft,
	"splitList":      SplitList,
	"last":           Last,
	"hasPrefix":      HasPrefix,
	"hasSuffix":      HasSuffix,
	"contains":       strings.Contains,
	"add":            Add,
	"lower":          strings.ToLower,
	"snakecase":      strcase.ToSnake,
	"replace":        strings.ReplaceAll,
	"concat":         Concat,
	"isService":      IsService,
}

// GoGetMethodBody 获取函数已实现的函数体
//
//	@param rw
//	@param structname
//	@param methodName
//	@return string
func GoGetMethodBody(rw *rewrite.Rewriter, structname, methodName string) string {
	if rw == nil {
		return ""
	}
	return rw.GetMethodBody(structname, methodName)
}

// GoGetImportPaths 获取import列表
//
//	@param rw
//	@param fileName
//	@return []rewrite.Import
func GoGetImportPaths(rw *rewrite.Rewriter, fileName string) []rewrite.Import {
	if rw == nil {
		return nil
	}
	return rw.ExistingImports(fileName)
}

// GoExistImportPaths 判断import是不是存在
//
//	@param rw
//	@param fileName
//	@return []rewrite.Import
func GoExistImportPaths(is []rewrite.Import, alias, path string) bool {
	for _, i := range is {
		if i.Alias == alias && i.ImportPath == path {
			return true
		}
	}
	return false
}

// GoComments TODO
func GoComments(name string, comments protogen.CommentSet) string {
	comment := strings.TrimSpace(comments.Leading.String())
	if comment == "" {
		comment = strings.TrimSpace(string(comments.Trailing.String()))
	}
	if len(comment) > 2 && (comment[:2] == "//" || comment[:2] == "/*") {
		// rid of //
		comment = strings.TrimSpace(comment[2:])
	}
	return strings.TrimSpace(strings.TrimSpace(name) + " " + comment)
}

// GoCamelCaseModel TODO
func GoCamelCaseModel(s string) string {
	s = TrimLeft(s, ".")
	idx := strings.Index(s, ".")
	if idx > 0 {
		pkg := s[:idx]
		s = s[idx:]
		return pkg + "." + GoCamelCase(s)
	}
	return GoCamelCase(s)
}

// GoCamelCase camel-cases a protobuf name for use as a Go identifier.
//
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
func GoCamelCase(s string) string {
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '.' in ".{{lowercase}}".
		case c == '.':
			b = append(b, '_') // convert '.' to '_'
		case c == '_' && (i == 0 || s[i-1] == '.'):
			// Convert initial '_' to ensure we start with a capital letter.
			// Do the same for '_' after '.' to match historic behavior.
			b = append(b, 'X') // convert '_' to 'X'
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '_' in "_{{lowercase}}".
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if isASCIILower(c) {
				c -= 'a' - 'A' // convert lowercase to uppercase
			}
			b = append(b, c)

			// Accept lower case sequence that follows.
			for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ {
				b = append(b, s[i+1])
			}
		}
	}
	return string(b)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

//	func isASCIIUpper(c byte) bool {
//		return 'A' <= c && c <= 'Z'
//	}
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// PBSimplifyGoType determine whether to use fullyQualifiedPackageName or not,
// if the `fullTypeName` occur in code of `package goPackageName`, `package` part
// should be removed.
func PBSimplifyGoType(fullTypeName string, goPackageName string) string {
	idx := strings.LastIndex(fullTypeName, ".")
	if idx <= 0 {
		panic(fmt.Sprintf("invalid fullyQualifiedType: %s", fullTypeName))
	}

	pkg := fullTypeName[0:idx]
	typ := fullTypeName[idx+1:]

	if pkg == goPackageName {
		//fmt.Println("pkg:", pkg, "=", "gopkg:", goPackageName)
		return typ
	}
	//fmt.Println("pkg:", pkg, "!=", "gopkg:", goPackageName)
	return fullTypeName
}

func FileName(fullFileName string) string {
	idx := strings.LastIndex(fullFileName, ".")
	if idx <= 0 {
		return fullFileName
	}

	return fullFileName[0:idx]
}

// PBGoType convert `t` to go style (like a.b.c.hello, it'll be changed to a_b_c.Hello)
func PBGoType(t string) string {
	var prefix string

	idx := strings.LastIndex(t, "/")
	if idx >= 0 {
		prefix = t[:idx]
		t = t[idx+1:]
	}

	idx = strings.LastIndex(t, ".")
	if idx <= 0 {
		panic(fmt.Sprintf("invalid go type: %s", t))
	}

	gopkg := PBGoPackage(t[0:idx])
	msg := t[idx+1:]

	return GoExport(prefix + gopkg + "." + msg)
}

// PBGoPackage convert a.b.c to a_b_c
func PBGoPackage(pkgName string) string {
	var (
		prefix string
		pkg    string
	)
	idx := strings.LastIndex(pkgName, "/")
	if idx < 0 {
		pkg = pkgName
	} else {
		prefix = pkgName[0:idx]
		pkg = pkgName[idx+1:]
	}

	pkg = strings.Replace(pkg, "-", "_", -1)
	gopkg := strings.Replace(pkg, ".", "_", -1)

	if len(prefix) == 0 {
		return gopkg
	}
	return prefix + "/" + gopkg
}

// GoExport export go type
func GoExport(typ string) string {
	idx := strings.LastIndex(typ, ".")
	if idx < 0 {
		return strings.Title(typ)
	}
	return typ[0:idx] + "." + strings.Title(typ[idx+1:])
}

// SplitList split string `str` via delimiter `sep` into a list of string
func SplitList(str, sep string) []string {
	return strings.Split(str, sep)
}

// TrimRight trim right substr starting at `sep`
func TrimRight(str, sep string) string {
	idx := strings.LastIndex(str, sep)
	if idx < 0 {
		return str
	}
	return str[:idx]
}

// TrimLeft trim left substr starting at `sep`
func TrimLeft(str, sep string) string {
	return strings.TrimPrefix(str, sep)
}

// Title uppercase the first character of `s`
func Title(s string) string {
	for k, v := range s {
		return string(unicode.ToUpper(v)) + s[k+1:]
	}
	return ""
}

// UnTitle make the first character of s lowercase
func UnTitle(s string) string {
	for k, v := range s {
		return string(unicode.ToLower(v)) + s[k+1:]
	}
	return ""
}

// PBValidGoPackage return valid go package
func PBValidGoPackage(pkgName string) string {
	var (
		pkg string
	)
	idx := strings.LastIndex(pkgName, "/")
	if idx < 0 {
		pkg = pkgName
	} else {
		pkg = pkgName[idx+1:]
	}

	pkg = strings.Replace(pkg, "-", "_", -1)
	pkg = strings.Replace(pkg, ".", "_", -1)
	return pkg
}

// Last returns the last element in `list`
func Last(list []string) string {
	idx := len(list) - 1
	return list[idx]
}

// HasPrefix test whether string `str` has prefix `prefix`
func HasPrefix(prefix, str string) bool {
	return strings.HasPrefix(str, prefix)
}

// HasSuffix test whether string `str` has suffix `suffix`
func HasSuffix(suffix, str string) bool {
	return strings.HasSuffix(str, suffix)
}

// Add add two number
func Add(num1, num2 int) int {
	return num1 + num2
}

// CheckSECVTpl 检查是否启用Validation特性，来决定导出的模板内容
func CheckSECVTpl(pkgMap map[string]string) bool {
	if _, isKeyFound := pkgMap["validate"]; isKeyFound {
		return true
	}
	return false
}

// Camelcase 驼峰处理，特殊case要兼容存量协议，否则转驼峰命名
func Camelcase(s string) string {
	if len(s) == 0 {
		return s
	}

	wordList := strings.Split(s, "_")
	if len(wordList) == 1 {
		if isAllUpper(wordList[0]) {
			return s
		}
		return strcase.ToCamel(s)
	}

	return CamelcaseList(wordList)
}

// CamelcaseList 驼峰处理字符串列表，最后输出一个字符串结果
func CamelcaseList(wordList []string) string {
	var camelWord string
	for i := 0; i < len(wordList); i++ {
		cur, seq := camelcaseListItem(wordList, i)
		camelWord = fmt.Sprintf("%s%s%s", camelWord, cur, seq)
	}

	return camelWord
}

func camelcaseListItem(wordList []string, i int) (string, string) {
	cur := wordList[i]
	seq := getCamelcaseSeq(wordList, i, cur)

	if !isAllUpper(cur) {
		cur = strcase.ToCamel(cur)
	}
	return cur, seq
}

func getCamelcaseSeq(wordList []string, i int, cur string) string {
	seq := ""
	// 如果当前和下一个词是全大写，需要用_拼接
	if i != len(wordList)-1 && len(wordList[i+1]) != 0 {
		curIsAllUpper := isAllUpper(cur)
		nextIsAllUpper := isAllUpper(wordList[i+1])
		if curIsAllUpper || nextIsAllUpper {
			seq = "_"
		}
	}
	return seq
}

// isAllUpper 是否全大写
func isAllUpper(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}
	return true
}

// Concat 连接字符串
func Concat(sep string, s ...string) string {
	ss := []string{}
	ss = append(ss, s...)
	return strings.Join(ss, sep)
}

func IsService(s *Service) bool {
	if s == nil {
		return false
	}

	if s.Options != nil {
		v := pb.GetExtension(s.Options, bbq.E_ServiceType)
		return v.(bbq.ServiceType) == bbq.ServiceType_Service
	}

	return false
}

// AStringContains 判断字符串数组中是否包含某个字符串
//
//	@param s 数组
//	@param n 字符串
//	@return bool
func AStringContains(s []string, n string) bool {
	for _, t := range s {
		if t == n {
			return true
		}
	}
	return false
}
