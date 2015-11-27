# Copyright 1999-2015 Gentoo Foundation
# Distributed under the terms of the GNU General Public License v2
# $Id$
# By Jean-Michel Smith, first created 9/21/15

EAPI=5

inherit user git-r3

DESCRIPTION="Export/Import/Edit etcd directory as JSON/YAML/TOML and validate directory using JSON schema"
HOMEPAGE="https://github.com/mickep76/etcdtool.git"
SRC_URI=""

LICENSE="Apache-2.0"
SLOT="0"
KEYWORDS="amd64"
IUSE=""

DEPEND="dev-lang/go"

EGIT_REPO_URI="https://github.com/mickep76/etcdtool.git"
EGIT_COMMIT="${PV}"

GOPATH="${WORKDIR}/etcdtool-${PV}"

src_compile() {
	ebegin "Building etcdtool ${PV}"
	export GOPATH
	export PATH=${GOPATH}/bin:${PATH}
	cd ${GOPATH}
	./build
	cd
	eend ${?}
}

src_install() {
	ebegin "installing etcdtool ${PV}"
	dobin ${GOPATH}/bin/etcdtool
	eend ${?}
}
