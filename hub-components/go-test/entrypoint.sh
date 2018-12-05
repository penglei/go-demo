#!/bin/bash

set -e

export SCRIPTDIR="$( cd "$( dirname "$0" )" && pwd )"

source $SCRIPTDIR/component-base/libs.sh

wait_service() {
	local time_out=$SVC_WAIT_TIMEOUT
	if [ -z $time_out ];then
		time_out=120
	fi

	echo "Wait MySQL..."
	for i in `seq 1 $time_out`;
	do
		nc -z localhost 3306 && echo Success && exit 0
		echo -n .
		sleep 1
	done
	echo Failed waiting for MySQL && exit 1

}

do_task() {
	runsvdir /etc/sv > /dev/null &

	(wait_service)

	hub_git_clone #clone and change work directory to git repo path

	# init database migrate
	vgo build -o /go-demo cmd/*.go
	/go-demo migrate up

	# run test!
	./test-service.sh

	echo "test completed successfully!"
}

do_task
