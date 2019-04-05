# consulwatchexample

This example shows how consul watcher is broken in master as at commit `28c84845bf082f07724217cd198d5f119009417f` but https://github.com/hashicorp/consul/tree/v1.4.4 works correctly (as does commit `3d95636dd7ff354490523d6d5cae218e15a77393`).

With the following `go env`

```
GOARCH="amd64"
GOBIN=""
GOCACHE="/Users/blairmcmillan/Library/Caches/go-build"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="darwin"
GOOS="darwin"
GOPATH="/Users/blairmcmillan/Documents/dev/go"
GOPROXY=""
GORACE=""
GOROOT="/usr/local/Cellar/go/1.11.5/libexec"
GOTMPDIR=""
GOTOOLDIR="/usr/local/Cellar/go/1.11.5/libexec/pkg/tool/darwin_amd64"
GCCGO="gccgo"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD=""
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/tj/rr2xddq944s5xwm_gt4xgl0m0000gn/T/go-build675345840=/tmp/go-build -gno-record-gcc-switches -fno-common"
```

1. `go get` this package into your `$GOPATH`
1. Run a consul server at `127.0.0.1:8500`
1. Checkout `github.com/hashicorp/consul` onto the `master` branch (any commit after `28c84845bf082f07724217cd198d5f119009417f` inclusive)
1. Run `go build; ./consulwatchexample`
    1. Observe console output of `watch data was not of the expected type. Got []*api.ServiceEntry`
    1. This shows that we cannot correct type cast the watcher output
    1. Watch methods cannot work
1. Checkout `github.com/hashicorp/consol` onto the `tag/v1.4.4` branch (or any commit before `3d95636dd7ff354490523d6d5cae218e15a77393` inclusive)
1. Run `go build; ./consulwatchexample`
    1. Observe console out of `serviceEntry: &api.AgentService{...}`
    1. This shows that we can correctly type cast the watcher output
    1. Watch methods can be used
