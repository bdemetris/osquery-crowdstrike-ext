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
	GOOS=darwin go build -i -o build/${APP_NAME} -pkgdir ${PKGDIR_TMP}
