// Copyright 2021 The protobuf-tools Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/davecgh/go-spew/spew"
	_ "sigs.k8s.io/controller-tools/pkg/crd"
)

func init() {
	spew.Config = spew.ConfigState{
		Indent:                  "  ",
		SortKeys:                true, // maps should be spewed in a deterministic order
		DisablePointerAddresses: true, // don't spew the addresses of pointers
		DisableCapacities:       true, // don't spew capacities of collections
		ContinueOnMethod:        true, // recursion should continue once a custom error or Stringer interface is invoked
		SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
		MaxDepth:                4,    // maximum number of levels to descend into nested data structures.
	}
}
