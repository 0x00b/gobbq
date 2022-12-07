package com

import (
	"path/filepath"
	"strings"
)

// TplSuffix eg main.grpc.go.tpl -> .grpc.go
//  @param name
//  @return string
func TplSuffix(name string) string {

	name = filepath.Base(name)
	start := strings.Index(name, ".")
	end := strings.LastIndex(name, ".")

	if start < end && start >= 0 {
		return name[start:end]
	}
	return ""
}
