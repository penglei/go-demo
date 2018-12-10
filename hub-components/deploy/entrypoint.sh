#!/bin/bash

#Example:
#
# docker run --rm \
#   hub.tencentyun.com/workshop/deploy
#

set -e
export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

source $SCRIPTDIR/component-base/libs.sh


replace_template() {
	local file_path="$1"
	local key="$2"
	local val="$3"
	sed -i.bak -E "s/\\{\\{$key\\}\\}/$val/g" $file_path
	rm $file_path.bak
}


init_kubectl_config() {
	local tke_server=`hub_var TKE_SERVER true`
	local tke_username=`hub_default_var TKE_USERNAME admin`
	local tke_password=`hub_var TKE_PASSWORD true`

	hub_var TKE_CERTIFICATE true | tee -a /tke-cluster-ca.crt

	kubectl config set-credentials default-admin --username="$tke_username" --password="$tke_password"
	kubectl config set-cluster default-cluster --server="$tke_server" --certificate-authority=/tke-cluster-ca.crt
	kubectl config set-context default-system --cluster=default-cluster --user=default-admin
	kubectl config use-context default-system
}

create_service() {
	local namespace="$1"
	cat k8s/service.yaml | kubectl apply -n "$namespace" -f -
}

create_deployment() {
	local namespace="$1"
	local tke_docker_image=`hub_var TKE_DOCKER_IMAGE`
	replace_template config/deployment.yaml docker_image "$tke_docker_image"
	# update configmap
	kubectl delete configmap go-demo-config || true
	kubectl create configmap go-demo-config --from-file=config/ -n "$namespace"

	# upgrade deployment
	kubectl apply -f k8s/deployment.yaml -n "$namespace"
}

do_task() {
	init_kubectl_config

	local tke_namespace=`hub_var TKE_CLUSTER_NAMESPACE true`

	# 创建service
	create_service "$tke_namespace"
	# 创建deploy
	create_deployment "$tke_namespace"

	echo "ok"
}

do_task

