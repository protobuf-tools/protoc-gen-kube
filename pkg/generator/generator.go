// Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// SupportedFeatures reports the set of supported protobuf language features.
const SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// Generate generates the Kubernetes controller APIs from Protocol Buffer schemas.
func Generate(gen *protogen.Plugin, file *protogen.File) error { return nil }
