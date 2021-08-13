// Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: Apache-2.0

//go:build tools
// +build tools

// Package tools manages the project dependency tools using during development.
package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "gotest.tools/gotestsum"
	_ "mvdan.cc/gofumpt"
)
