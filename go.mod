module github.com/protobuf-tools/protoc-gen-kube

go 1.17

require (
	github.com/go-logr/logr v1.0.0
	github.com/go-logr/zapr v1.0.0
	go.uber.org/zap v1.19.0
	google.golang.org/protobuf v1.27.1
	k8s.io/apimachinery v0.22.0
	k8s.io/gengo v0.0.0-20210813121822-485abfe95c7c
	k8s.io/klog/v2 v2.10.0
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/tools v0.1.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace (
	// support typeparams build tag
	golang.org/x/tools => github.com/zchee/golang-tools v0.0.0-20210814085923-cf63d8262102
	// support go-logr/logr@v1.0.0
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.10.1-0.20210806124320-e1f317b53636
)
