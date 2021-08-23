module github.com/protobuf-tools/protoc-gen-kube

go 1.17

require (
	github.com/go-logr/logr v1.1.0
	github.com/go-logr/zapr v1.0.1-0.20210809170106-a3325063a237 // v1.0.0+1 is not semver
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.19.0
	google.golang.org/protobuf v1.27.1
	k8s.io/apimachinery v0.22.1
	k8s.io/gengo v0.0.0-20201214224949-b6c5ce23f027
	k8s.io/klog/v2 v2.10.0
	sigs.k8s.io/controller-tools v0.0.0-00000000000000-000000000000
)

require (
	github.com/gobuffalo/flect v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/text v0.3.6 // indirect
	golang.org/x/tools v0.1.6-0.20210813165731-45389f592fe9 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apiextensions-apiserver v0.22.1 // indirect
	k8s.io/utils v0.0.0-20210707171843-4b05e18ac7d9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

// support go-logr/logr@v1.0.0
replace k8s.io/klog/v2 => k8s.io/klog/v2 v2.10.1-0.20210806124320-e1f317b53636

// pin k8s package
replace (
	k8s.io/api => k8s.io/api v0.22.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.1
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210527164424-3c818078ee3d
)

// forked
replace (
	// use forked
	k8s.io/gengo => ./pkg/internal/k8s.io/gengo

	// use forked
	sigs.k8s.io/controller-tools => ./pkg/internal/sigs.k8s.io/controller-tools
)

// CVE
exclude github.com/coreos/etcd v3.3.13+incompatible

// CVE
replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v3.2.2+incompatible
