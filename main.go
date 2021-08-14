// Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: Apache-2.0

// Command protoc-gen-kube generates the Kubernetes controller APIs from Protocol Buffer schemas.
package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/multierr"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/protobuf-tools/protoc-gen-kube/pkg/generator"
	"github.com/protobuf-tools/protoc-gen-kube/pkg/version"
)

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stderr, "protoc-gen-kube %s\n", version.Version())
		os.Exit(0)
	}

	var flags flag.FlagSet
	g := protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: func(imp protogen.GoImportPath) protogen.GoImportPath { return "" },
	}

	g.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = generator.SupportedFeatures

		var errs error
		for _, files := range gen.Files {
			if files.Generate {
				if err := generator.Generate(gen, files); err != nil {
					errs = multierr.Append(errs, err)
				}
			}
		}

		if errs != nil {
			gen.Error(multierr.Combine(errs))
		}

		return nil
	})
}
