#!/bin/bash

set -e

export SCRIPTDIR="$( cd "$( dirname "$0" )" && pwd )"

source $SCRIPTDIR/component-base/libs.sh

wait_service() {
	echo "Wait MySQL..."
	for i in `seq 1 30`;
	do
		nc -z localhost 3306 && echo Success && exit 0
		echo -n .
		sleep 1
	done
	echo Failed waiting for MySQL && exit 1

}

do_task() {
	runsvdir /etc/sv > /dev/null 2>&1 &

	(wait_service)

	hub_git_clone

	#TODO vgo test

	echo "test completed successfully!"
}

do_task
