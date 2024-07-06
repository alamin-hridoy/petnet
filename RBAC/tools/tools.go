//go:build tools
// +build tools

package tools

import (
	_ "cirello.io/openapigen"
	_ "cuelang.org/go/cmd/cue"
	_ "github.com/brankas/git-buildnumber"
	_ "github.com/go-bindata/go-bindata"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/grpc-ecosystem/grpc-health-probe"
	_ "github.com/gunk/gunk"
	_ "github.com/gunk/opt/openapiv2"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/cmd/stringer"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "mvdan.cc/gofumpt"
	_ "sigs.k8s.io/kustomize/kustomize/v3"
)
