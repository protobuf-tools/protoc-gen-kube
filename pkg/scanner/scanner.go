// Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2019 Istio Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// This file copied and edit from https://github.com/istio/tools/blob/1.11.0/cmd/kubetype-gen/scanner/scanner.go.

package scanner

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"

	"github.com/protobuf-tools/protoc-gen-kube/pkg/generator/kubetype"
	"github.com/protobuf-tools/protoc-gen-kube/pkg/metadata"
)

const (
	// enabledTagName is the root tag used to identify types that need a corresponding kube type generated.
	enabledTagName = "kubetype-gen"

	// groupVersionTagName is the tag used to identify the k8s group/version associated with the generated types.
	groupVersionTagName = enabledTagName + ":groupVersion"

	// kubeTypeTagName is used to identify the name(s) of the types to be generated from the type with this tag.
	// If this tag is not present, the k8s type will have the same name as the source type.  If this tag is specified
	// multiple times, a k8s type will be generated for each value.
	kubeTypeTagName = enabledTagName + ":kubeType"

	// kubeTagsTagTemplate is used to identify a comment tag that should be added to the generated kubeType.  This
	// allows different sets of tags to be used when a single type is the source for multiple kube types (e.g. where one
	// is namespaced and another is not).  The tag should not be prefixed with '+', as this will be added by the
	// generator.  This may be specified multiple times, once for each tag to be added to the generated type.
	kubeTagsTagTemplate = enabledTagName + ":%s:tag"
)

// Scanner is used to scan input packages for types with kubetype-gen tags.
type Scanner struct {
	ctx context.Context

	arguments *args.GeneratorArgs
	gctxt     *generator.Context
}

// WithContext sets ctx to Scanner.
//
// If not set, use context.Background instead.
func (s *Scanner) WithContext(ctx context.Context) *Scanner {
	s.ctx = ctx

	return s
}

// Scan the input packages for types with kubetype-gen tags.
func (s *Scanner) Scan(c *generator.Context, arguments *args.GeneratorArgs) generator.Packages {
	if s.ctx != nil {
		s.ctx = context.Background()
	}
	logf := logr.FromContextOrDiscard(s.ctx)

	s.arguments = arguments
	s.gctxt = c

	boilerplate, err := arguments.LoadGoBoilerplate()
	if err != nil {
		panic(fmt.Errorf("failed loading boilerplate: %w", err))
	}

	// scan input packages for kubetype-gen
	metadataStore := metadata.NewMetadataStore(s.ctx, s.getBaseOutputPackage(), &c.Universe)
	fail := false

	logf.V(1).Info("scanning input packages")
	for _, input := range c.Inputs {
		logf.V(1).Info("scanning package", "input", input)

		pkg := c.Universe[input]
		if pkg == nil {
			logf.V(1).Info("package not found, continue", "input", input)
			continue
		}

		if strings.HasPrefix(arguments.OutputPackagePath, pkg.Path) {
			logf.V(1).Info("ignoring package because located in the output package, continue", "ignored package", pkg.Path, "output package", arguments.OutputPackagePath)
			continue
		}

		pkgTags := types.ExtractCommentTags("+", pkg.DocComments)

		// group/version for generated types from this package
		defaultGV, err := s.getGroupVersion(pkgTags, nil)
		switch {
		case err != nil:
			logf.Error(err, "could not calculate Group/Version for package", "package path", pkg.Path)
			fail = true
		case defaultGV != nil && defaultGV.Group == "":
			logf.Error(errors.New("invalid Group/Version"), "invalid Group/Version for package because Group not specified for Group/Version", "package path", pkg.Path, "Group/Version", defaultGV)
			fail = true
		default:
			logf.V(1).Info("use default Group/Version for package", "Group/Version", defaultGV)
		}

		// scan package for types that need kube types generated
		for _, t := range pkg.Types {
			comments := make([]string, len(t.CommentLines)+len(t.SecondClosestCommentLines))
			for i, line := range append(t.CommentLines, t.SecondClosestCommentLines...) {
				comments[i] = line
			}

			typeTags := types.ExtractCommentTags("+", comments)
			if _, ok := typeTags[enabledTagName]; ok {
				var gv *schema.GroupVersion
				gv, err = s.getGroupVersion(typeTags, defaultGV)
				if err != nil {
					logf.Error(err, "could not calculate Group/Version for type", "type", t)
					fail = true
					continue
				}

				if gv == nil || gv.Group == "" {
					logf.Error(errors.New("invalid Group/Version"), "invalid Group/Version for type", "type", t, "Group/Version", gv)
					fail = true
					continue
				}

				packageMetadata := metadataStore.MetadataForGV(gv)
				if packageMetadata == nil {
					logf.Error(errors.New("create metadata"), "could not create metadata for type", "type", t)
					fail = true
					continue
				}

				kubeTypes := s.createKubeTypesForType(t, packageMetadata.TargetPackage())
				logf.V(1).Info("Kube types will be generated with Group/Version for raw type", "kubeTypes", kubeTypes, "Group/Version", gv, "type", t)

				if err := packageMetadata.AddMetadataForType(t, kubeTypes...); err != nil {
					logf.Error(err, "error adding metadata source", "type", t)
					fail = true
				}
			}
		}
	}
	logf.V(1).Info("finished scanning input packages")

	validationErrors := metadataStore.Validate()
	if len(validationErrors) > 0 {
		for _, validationErr := range validationErrors {
			logf.Error(validationErr, "failed to validate metadata", "metadata store", metadataStore)
		}
		fail = true
	}

	if fail {
		panic(errors.New("errors occurred while scanning input, see previous output for details"))
	}

	mds := metadataStore.AllMetadata()
	pkgs := make([]generator.Package, 0, len(mds))
	for _, source := range mds {
		if len(source.RawTypes()) == 0 {
			logf.V(1).Info("skipping generation, no types to generate", "Group/Version", source.GroupVersion(), "source", source)
			continue
		}

		logf.V(1).Info("adding package generator", "group/version", source.GroupVersion())
		pkgs = append(pkgs, kubetype.NewPackageGenerator(s.ctx, source, boilerplate))
	}

	return pkgs
}

func (s *Scanner) getGroupVersion(tags map[string][]string, defaultGV *schema.GroupVersion) (*schema.GroupVersion, error) {
	if value, ok := tags[groupVersionTagName]; ok && len(value) > 0 {
		gv, err := schema.ParseGroupVersion(value[0])
		if err != nil {
			return nil, fmt.Errorf("invalid group version %q specified: %w", value[0], err)
		}

		return &gv, nil
	}

	return defaultGV, nil
}

func (s *Scanner) getBaseOutputPackage() *types.Package {
	return s.gctxt.Universe.Package(s.arguments.OutputPackagePath)
}

func (s *Scanner) createKubeTypesForType(t *types.Type, outputPackage *types.Package) []metadata.KubeType {
	namesForType := s.kubeTypeNamesForType(t)
	newKubeTypes := make([]metadata.KubeType, 0, len(namesForType))

	for _, name := range namesForType {
		tags := s.getTagsForKubeType(t, name)
		newKubeTypes = append(newKubeTypes, metadata.NewKubeType(t, s.gctxt.Universe.Type(types.Name{Name: name, Package: outputPackage.Path}), tags))
	}

	return newKubeTypes
}

func (s *Scanner) kubeTypeNamesForType(t *types.Type) []string {
	comments := make([]string, 0, len(t.CommentLines)+len(t.SecondClosestCommentLines))
	comments = append(comments, t.CommentLines...)
	comments = append(comments, t.SecondClosestCommentLines...)
	tags := types.ExtractCommentTags("+", comments)
	names := make([]string, 0, len(tags[kubeTypeTagName]))

	if value, ok := tags[kubeTypeTagName]; ok {
		for _, name := range value {
			if name != "" {
				names = append(names, name)
			}
		}
	} else {
		if len(value) == 0 || value[0] == "" {
			logr.FromContextOrDiscard(s.ctx).Error(errors.New("invalid value specified"), "using default name", "kubeTypeTagName", kubeTypeTagName, "type", t, "name", t.Name.Name)
		}
		names = append(names, t.Name.Name)
	}

	return names
}

func (s *Scanner) getTagsForKubeType(t *types.Type, name string) []string {
	tagName := fmt.Sprintf(kubeTagsTagTemplate, name)
	comments := make([]string, 0, len(t.CommentLines)+len(t.SecondClosestCommentLines))
	comments = append(comments, t.CommentLines...)
	comments = append(comments, t.SecondClosestCommentLines...)
	tags := types.ExtractCommentTags("+", comments)

	if value, ok := tags[tagName]; ok {
		return value
	}

	return []string{}
}
