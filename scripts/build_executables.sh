#!/bin/bash

cd $WORKSPACE

set -e
if [ -f $WORKSPACE/../TOGGLE ]; then
	echo "****************************************************"
    echo "data.stack.b2b.agent.watcher :: Toggle mode is on, terminating build"
    echo "data.stack.b2b.agent.watcher :: BUILD CANCLED"
    echo "****************************************************"
    exit 0
fi

echo "****************************************************"
echo "data.stack.b2b.agent.watcher :: Clearing exec folder"
echo "****************************************************"
rm exec/* || true
echo "****************************************************"
echo "data.stack.b2b.agent.watcher :: Building excutables"
echo "****************************************************"
# echo "env GOOS=android GOARCH=arm go build -ldflags="-s -w" -o ds-agent-android-arm ./v1"
# env GOOS=android GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-android-arm ./v1 || true

# echo "env GOOS=darwin GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-darwin-386 ./v1"
# env GOOS=darwin GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-darwin-386 ./v1

echo "env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-darwin-amd64 ./v1"
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-sentinel-darwin-amd64 main.go || true

# echo "env GOOS=darwin GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-darwin-arm ./v1"
# env GOOS=darwin GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-darwin-arm ./v1 || true

# echo "env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o exec/ds-agent-darwin-arm64 ./v1"
# env GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o exec/ds-agent-darwin-arm64 ./v1 || true

# echo "env GOOS=dragonfly GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-dragonfly-amd64 ./v1"
# env GOOS=dragonfly GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-dragonfly-amd64 ./v1 || true

# echo "env GOOS=freebsd GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-freebsd-386 ./v1"
# env GOOS=freebsd GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-freebsd-386 ./v1 || true

# echo "env GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-freebsd-amd64 ./v1"
# env GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-freebsd-amd64 ./v1 || true

# echo "env GOOS=freebsd GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-freebsd-arm ./v1"
# env GOOS=freebsd GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-freebsd-arm ./v1 || true

echo "env GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o exec/ds-sentinel-linux-386 main.go"
env GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o exec/ds-sentinel-linux-386 main.go

echo "env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-sentinel-linux-amd64 main.go || true"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-sentinel-linux-amd64 main.go || true

# echo "env GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-linux-arm ./v1"
# env GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-linux-arm ./v1 || true

# echo "env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o exec/ds-agent-linux-arm64 ./v1"
# env GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o exec/ds-agent-linux-arm64 ./v1 || true

# echo "env GOOS=linux GOARCH=ppc64 go build -ldflags="-s -w" -o exec/ds-agent-linux-ppc64 ./v1"
# env GOOS=linux GOARCH=ppc64 go build -ldflags="-s -w" -o exec/ds-agent-linux-ppc64 ./v1 || true

# echo "env GOOS=linux GOARCH=ppc64le go build -ldflags="-s -w" -o exec/ds-agent-linux-ppc64le ./v1"
# env GOOS=linux GOARCH=ppc64le go build -ldflags="-s -w" -o exec/ds-agent-linux-ppc64le ./v1 || true

# echo "env GOOS=linux GOARCH=mips go build -ldflags="-s -w" -o exec/ds-agent-linux-mips ./v1"
# env GOOS=linux GOARCH=mips go build -ldflags="-s -w" -o exec/ds-agent-linux-mips ./v1 || true

# echo "env GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o exec/ds-agent-linux-mipsle ./v1"
# env GOOS=linux GOARCH=mipsle go build -ldflags="-s -w" -o exec/ds-agent-linux-mipsle ./v1 || true

# echo "env GOOS=linux GOARCH=mips64 go build -ldflags="-s -w" -o exec/ds-agent-linux-mips64 ./v1"
# env GOOS=linux GOARCH=mips64 go build -ldflags="-s -w" -o exec/ds-agent-linux-mips64 ./v1 || true

# echo "env GOOS=linux GOARCH=mips64le go build -ldflags="-s -w" -o exec/ds-agent-linux-mips64le ./v1"
# env GOOS=linux GOARCH=mips64le go build -ldflags="-s -w" -o exec/ds-agent-linux-mips64le ./v1 || true

# echo "env GOOS=netbsd GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-netbsd-386 ./v1"
# env GOOS=netbsd GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-netbsd-386 ./v1 || true

# echo "env GOOS=netbsd GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-netbsd-amd64 ./v1"
# env GOOS=netbsd GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-netbsd-amd64 ./v1 || true

# echo "env GOOS=netbsd GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-netbsd-arm ./v1"
# env GOOS=netbsd GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-netbsd-arm ./v1 || true

# echo "env GOOS=openbsd GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-openbsd-386 ./v1"
# env GOOS=openbsd GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-openbsd-386 ./v1 || true

# echo "env GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-openbsd-amd64 ./v1"
# env GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-openbsd-amd64 ./v1 || true

# echo "env GOOS=openbsd GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-openbsd-arm ./v1"
# env GOOS=openbsd GOARCH=arm go build -ldflags="-s -w" -o exec/ds-agent-openbsd-arm ./v1 || true

# echo "env GOOS=plan9 GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-plan9-386 ./v1"
# env GOOS=plan9 GOARCH=386 go build -ldflags="-s -w" -o exec/ds-agent-plan9-386 ./v1 || true

# echo "env GOOS=plan9 GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-plan9-amd64 ./v1"
# env GOOS=plan9 GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-plan9-amd64 ./v1 || true

# echo "env GOOS=solaris GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-solaris-amd64 ./v1"
# env GOOS=solaris GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-agent-solaris-amd64 ./v1 || true

echo "env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o exec/ds-sentinel-windows-386.exe main.go"
env GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o exec/ds-sentinel-windows-386-unsigned.exe main.go
osslsigncode -h sha2 -certs certs/cd786349a667ff05-SHA2.pem -key certs/out.key -t http://timestamp.globalsign.com/scripts/timstamp.dll -in exec/ds-sentinel-windows-386-unsigned.exe -out exec/ds-sentinel-windows-386.exe

echo "env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-sentinel-windows-amd64.exe main.go"
env GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o exec/ds-sentinel-windows-amd64-unsigned.exe main.go
osslsigncode -h sha2 -certs certs/cd786349a667ff05-SHA2.pem -key certs/out.key -t http://timestamp.globalsign.com/scripts/timstamp.dll -in exec/ds-sentinel-windows-amd64-unsigned.exe -out exec/ds-sentinel-windows-amd64.exe
