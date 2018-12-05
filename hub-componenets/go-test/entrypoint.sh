#!/bin/bash

export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

source $SCRIPTDIR/component-base/libs.sh

init_mysql() {
	source $SCRIPTDIR/mysql-entrypoint.sh
}

do_task() {
	init_mysql()
	hub_git_clone

	#TODO vgo test

	echo "test completed successfully!"
}

do_task
