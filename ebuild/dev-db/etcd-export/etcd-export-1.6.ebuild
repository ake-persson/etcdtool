# Copyright 1999-2015 Gentoo Foundation
# Distributed under the terms of the GNU General Public License v2
# $Id$
# By Jean-Michel Smith, first created 9/21/15

EAPI=5

inherit user git-r3

DESCRIPTION="Expose hardware info using JSON/REST and provide a system HTML Front-End"
HOMEPAGE="https://github.com/mickep76/peekaboo.git"
SRC_URI=""

LICENSE="Apache-2.0"
SLOT="0"
KEYWORDS="amd64"
IUSE=""

DEPEND="dev-lang/go"

EGIT_REPO_URI="https://github.com/mickep76/etcd-export.git"
EGIT_COMMIT="${PV}"

GOPATH="${WORKDIR}/etcd-export-${PV}"

src_compile() {
	ebegin "Building etcd-export ${PV}"
	export GOPATH
	export PATH=${GOPATH}/bin:${PATH}
	cd ${GOPATH}
	./build
	cd
	eend ${?}
}

src_install() {
	ebegin "installing etcd-export ${PV}"
	dobin ${GOPATH}/bin/etcd-export
	dobin ${GOPATH}/bin/etcd-import
	dobin ${GOPATH}/bin/etcd-delete
	dobin ${GOPATH}/bin/etcd-tree
	dobin ${GOPATH}/bin/etcd-validate
	eend ${?}
}
