// Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"runtime/debug"
	"strings"
	"text/template"
)

// version is a protoc-gen-kube vesion.
var version = "v0.0.0"

// gitCommit indicates which git commit hash the binary was built off of.
var gitCommit = "devel"

// buildInfo is the stub for moduleBuildInfo function which prints build information embedded
// in the running binary.
var buildInfo = func() string {
	return ""
}

func init() {
	buildInfo = moduleBuildInfo
}

var buildInfoTmpl = ` mod	{{ .Main.Path }}	{{ .Main.Version }}	{{ .Main.Sum }}
{{ range .Deps }} dep	{{ .Path }}		{{ .Version }}	{{ .Sum }}{{ if .Replace }}
	=> {{ .Replace.Path }}	{{ .Replace.Version }}	{{ .Replace.Sum }}{{ end }}
{{ end }}`

func moduleBuildInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "not built in module mode"
	}

	buf := new(strings.Builder)
	err := template.Must(template.New("buildinfo").Parse(buildInfoTmpl)).Execute(buf, info)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// Version returns the protoc-gen-kube current version and build informations.
func Version() string {
	return fmt.Sprintf("%s@%s\n\nBuildInfo:\n%s", version, gitCommit, buildInfo())
}
