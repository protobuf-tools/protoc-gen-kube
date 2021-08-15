module k8s.io/gengo

go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/go-cmp v0.5.6
	github.com/google/gofuzz v1.2.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.0.0-20200505023115-26f46d2f7ef8
	k8s.io/klog/v2 v2.10.0
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/go-logr/logr v1.0.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	// support typeparams build tag
	golang.org/x/tools => github.com/zchee/golang-tools v0.0.0-20210814085923-cf63d8262102

	// support go-logr/logr@v1.0.0
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.10.1-0.20210806124320-e1f317b53636
)

// replace
replace golang.org/x/sys => golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e
