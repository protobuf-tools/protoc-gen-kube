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

// This file copied and edit from https://github.com/istio/tools/blob/1.11.0/cmd/kubetype-gen/generators/naming.go.

package kubetype

import (
	"strings"

	"k8s.io/gengo/namer"
)

// NameSystems used by the kubetype generator
func NameSystems(generatedPackage string, tracker namer.ImportTracker) namer.NameSystems {
	return namer.NameSystems{
		"public":       namer.NewPublicNamer(0),
		"raw":          namer.NewRawNamer(generatedPackage, tracker),
		"publicPlural": namer.NewPublicPluralNamer(map[string]string{}),
		"lower":        newLowerCaseNamer(0),
	}
}

// DefaultNameSystem to use if none is specified
func DefaultNameSystem() string {
	return "public"
}

func newLowerCaseNamer(prependPackageNames int, ignoreWords ...string) *namer.NameStrategy {
	n := &namer.NameStrategy{
		Join:                namer.Joiner(namer.IL, strings.ToLower),
		IgnoreWords:         map[string]bool{},
		PrependPackageNames: prependPackageNames,
	}
	for _, w := range ignoreWords {
		n.IgnoreWords[w] = true
	}
	return n
}
