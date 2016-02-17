NAME=etcdtool
BUILDDIR=.build
SRCDIR=github.com/mickep76/$(NAME)
VERSION:=3.2
RELEASE:=$(shell date -u +%Y%m%d%H%M)
ARCH:=$(shell uname -p)

all: build

clean:
#	rm -f *.rpm
	rm -rf bin pkg ${NAME} ${BUILDDIR}

update:
	gb vendor update --all

deps:
	go get github.com/constabulary/gb/...

build: clean
	gb build

darwin:
	gb build
	mv bin/etcdtool etcdtool-${VERSION}-${RELEASE}.darwin.x86_64

rpm:
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/$(SRCDIR) -w /go/src/$(SRCDIR) mickep76/centos-golang:latest make build-rpm

binary:
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/$(SRCDIR) -w /go/src/$(SRCDIR) mickep76/centos-golang:latest make build-binary
	mv bin/etcdtool etcdtool-${VERSION}-${RELEASE}.linux.x86_64

release: clean darwin rpm binary

build-binary: deps build

build-rpm: deps build
	mkdir -p ${BUILDDIR}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	cp bin/${NAME} ${BUILDDIR}/SOURCES
	sed -e "s/%NAME%/${NAME}/g" -e "s/%VERSION%/${VERSION}/g" -e "s/%RELEASE%/${RELEASE}/g" \
		${NAME}.spec >${BUILDDIR}/SPECS/${NAME}.spec
	rpmbuild -vv -bb --target="${ARCH}" --clean --define "_topdir $$(pwd)/${BUILDDIR}" ${BUILDDIR}/SPECS/${NAME}.spec
	mv ${BUILDDIR}/RPMS/${ARCH}/*.rpm .
