NAME=etcdtool
BUILDDIR=.build
SRCDIR=github.com/mickep76/$(NAME)
VERSION:=$(shell git describe --abbrev=0 --tags)
RELEASE:=$(shell date -u +%Y%m%d%H%M)
ARCH:=$(shell uname -p)

all: build

clean:
	rm -rf bin pkg ${NAME} ${BUILDDIR} release

build: clean
	go build

darwin:
	gb build
	mkdir release || true
	mv bin/etcdtool release/etcdtool-${VERSION}-${RELEASE}.darwin.x86_64

rpm:
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/$(SRCDIR) -w /go/src/$(SRCDIR) mickep76/centos-golang:latest make build-rpm

binary:
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/$(SRCDIR) -w /go/src/$(SRCDIR) mickep76/centos-golang:latest make build-binary
	mkdir release || true
	mv bin/etcdtool release/etcdtool-${VERSION}-${RELEASE}.linux.x86_64

set-version:
	sed -i .tmp "s/const Version =.*/const Version = \"${VERSION}\"/" src/${SRCDIR}/version.go
	rm -f src/${SRCDIR}/version.go.tmp

release: clean set-version darwin rpm binary

build-binary:
	gb build

build-rpm:
	gb build
	mkdir -p ${BUILDDIR}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	cp bin/${NAME} ${BUILDDIR}/SOURCES
	sed -e "s/%NAME%/${NAME}/g" -e "s/%VERSION%/${VERSION}/g" -e "s/%RELEASE%/${RELEASE}/g" \
		${NAME}.spec >${BUILDDIR}/SPECS/${NAME}.spec
	rpmbuild -vv -bb --target="${ARCH}" --clean --define "_topdir $$(pwd)/${BUILDDIR}" ${BUILDDIR}/SPECS/${NAME}.spec
	mkdir release || true
	mv ${BUILDDIR}/RPMS/${ARCH}/*.rpm release
