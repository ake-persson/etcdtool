NAME=dcrank
SRCDIR=github.com/mickep76/dcrank
BUILDDIR=.build
VERSION:=2.5
RELEASE:=$(shell date -u +%Y%m%d%H%M)
ARCH:=$(shell uname -p)

all: build

clean:
	rm -f *.rpm
	rm -rf ${NAME} ${BUILDDIR}

update:
	glide up

build: clean
	go build

rpm:
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/${SRCDIR} -w /go/src/${SRCDIR} mickep76/centos-golang:latest make build-rpm

build-rpm: build
	mkdir -p ${BUILDDIR}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	cp -r ${NAME} ${BUILDDIR}/SOURCES
	sed -e "s/%NAME%/${NAME}/g" -e "s/%VERSION%/${VERSION}/g" -e "s/%RELEASE%/${RELEASE}/g" \
		${NAME}.spec >${BUILDDIR}/SPECS/${NAME}.spec
	rpmbuild -vv -bb --target="${ARCH}" --clean --define "_topdir $$(pwd)/${BUILDDIR}" ${BUILDDIR}/SPECS/${NAME}.spec
	mv ${BUILDDIR}/RPMS/${ARCH}/*.rpm .
