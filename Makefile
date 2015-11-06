NAME=etcd-export
SRCDIR=src/github.com/mickep76
TMPDIR1=.build
VERSION:=$(shell awk -F '"' '/Version/ {print $$2}' ${SRCDIR}/common/version.go)
RELEASE:=$(shell date -u +%Y%m%d%H%M)
ARCH:=$(shell uname -p)

all: build

clean:
	rm -f *.rpm
	rm -rf pkg bin ${TMPDIR1}

test: clean
	gb test

build: test
	gb build all

update:
	gb vendor update --all

install:
	cp bin/* /usr/bin

docker-rpm: docker-clean
	docker pull mickep76/centos-golang:latest
	docker run --rm -it -v "$$PWD":/go/src/${SRCDIR} -w /go/src/${SRCDIR} mickep76/centos-golang:latest "make rpm"

rpm:	build
	mkdir -p ${TMPDIR1}/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	cp -r bin ${TMPDIR1}/SOURCES
	sed -e "s/%NAME%/${NAME}/g" -e "s/%VERSION%/${VERSION}/g" -e "s/%RELEASE%/${RELEASE}/g" \
		${NAME}.spec >${TMPDIR1}/SPECS/${NAME}.spec
	rpmbuild -vv -bb --target="${ARCH}" --clean --define "_topdir $$(pwd)/${TMPDIR1}" ${TMPDIR1}/SPECS/${NAME}.spec
	mv ${TMPDIR1}/RPMS/${ARCH}/*.rpm .
