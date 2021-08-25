module k8s.io/gengo

go 1.17

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/google/go-cmp v0.5.6
	github.com/google/gofuzz v1.2.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/tools v0.1.6-0.20210825163236-3e0d083b858b
	k8s.io/klog/v2 v2.10.0
	sigs.k8s.io/yaml v1.2.0
)

require (
	github.com/go-logr/logr v1.1.0 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

// support go-logr/logr@v1.0.0
replace k8s.io/klog/v2 => k8s.io/klog/v2 v2.10.1-0.20210806124320-e1f317b53636

// replace
replace golang.org/x/sys => golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e
