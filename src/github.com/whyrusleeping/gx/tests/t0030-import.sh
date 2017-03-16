#!/bin/sh
#
# Copyright (c) 2015 Jeromy Johnson
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="test package importing"

. lib/test-lib.sh

check_package_import() {
	pkg=$1
	imphash=$2
	name=$3

	test_expect_success "dir was created" '
		stat $pkg/vendor/gx/ipfs/$imphash > /dev/null
	'

	test_expect_success "dep set in package.json" '
		jq -r ".gxDependencies[] | select(.hash == \"$imphash\") | .name" $pkg/package.json > name
	'

	test_expect_success "name looks good" '
		echo "$name" > exp_name &&
		test_cmp exp_name name
	'
}

test_init_ipfs
test_launch_ipfs_daemon

test_expect_success "setup test packages" '
	make_package a none
	make_package b none
	make_package c none
	make_package d none
'

test_expect_success "publish the packages a and b" '
	pkgA=$(publish_package a) &&
	pkgB=$(publish_package b)
'

test_expect_success "import package a succeeds" '
	pkg_run c gx import $pkgA
'

check_package_import c $pkgA a

test_expect_success "import package b succeeds" '
	pkg_run c gx import $pkgB
'

check_package_import c $pkgB b

test_expect_success "publish c succeeds" '
	pkgC=$(publish_package c)
'

test_expect_success "d imports c succeeds" '
	pkg_run d gx import $pkgC
'

check_package_import d $pkgC c

test_expect_success "importing c brought along a and b" '
	stat d/vendor/gx/ipfs/$pkgA/a/package.json > /dev/null &&
	stat d/vendor/gx/ipfs/$pkgB/b/package.json > /dev/null
'

test_expect_success "install d works" '
	pkg_run d gx --verbose install > install_out
'

test_expect_success "install output looks good" '
	grep "installing package: d-0.0.0" install_out &&
	grep "installing package: c-0.0.0" install_out &&
	grep "installing package: a-0.0.0" install_out &&
	grep "installing package: b-0.0.0" install_out &&
	grep "installation of dep a complete!" install_out &&
	grep "installation of dep b complete!" install_out &&
	grep "installation of dep c complete!" install_out &&
	grep "installation of dep d complete!" install_out
'

test_expect_success "deps look correct" '
	pkg_run d gx deps --tree > deps_out
'

test_expect_success "deps tree looks right" '
	printf "└─ \033[1mc\033[0m    %s 0.0.0\n" "$pkgC" > deps_exp &&
	printf "   ├─ \033[1ma\033[0m %s 0.0.0\n" "$pkgA" >> deps_exp &&
	printf "   └─ \033[1mb\033[0m %s 0.0.0\n" "$pkgB" >> deps_exp &&
	test_cmp deps_exp deps_out
'

test_kill_ipfs_daemon

test_done
