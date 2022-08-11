all: build

APP_NAME = crowdstrike.ext
PKGDIR_TMP = ${TMPDIR}golang

.pre-build:
	mkdir -p build

download:
	go mod download

clean:
	rm -rf build/
	rm -rf ${PKGDIR_TMP}_darwin

build: .pre-build
	GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}-amd64 -pkgdir ${PKGDIR_TMP}
	GOOS=darwin GOARCH=arm64 go build -o build/${APP_NAME}-arm64 -pkgdir ${PKGDIR_TMP}
	lipo -create -output build/${APP_NAME} build/${APP_NAME}-amd64 build/${APP_NAME}-arm64
