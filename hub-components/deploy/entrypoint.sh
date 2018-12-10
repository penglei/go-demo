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

create_configmap() {
	local namespace="$1"
	kubectl delete configmap go-demo-config -n "$namespace" || true

	replace_template k8s/config/go-demo.yaml db_name "$2"

	#XXX 可以在这里替换配置
	kubectl create configmap go-demo-config --from-file=k8s/config/ -n "$namespace"
}

create_deployment() {
	local namespace="$1"
	local tke_docker_image=`hub_var TKE_DOCKER_IMAGE`

	replace_template k8s/config/deployment.yaml docker_image "$tke_docker_image"
	# upgrade deployment
	kubectl apply -f k8s/deployment.yaml -n "$namespace"
}

do_database_migrate() {
	local namespace="$1"
	local tke_docker_image=`hub_var TKE_DOCKER_IMAGE`
	# create configmap
	local yaml=$(cat <<EOF
{ "apiVersion": "v1", "spec": {
	"imagePullSecrets": [{"name": "myhubsecret"}],
	"volumes":[{
		"name": "config",
		"configMap": {
			"name": "go-demo-config"
		}
	}],
	"containers": [{
		"name": "migrate-database",
		"image": "$tke_docker_image",
		"stdin": true,
		"stdinOnce": true,
		"tty": true,
		"command": ["/go-demo"],
		"args": ["-c", "/go-demo-config/go-demo.yaml", "migrate", "up"],
		"workingDir": "/go/src/github.com/qcloud2018/go-demo",
		"volumeMounts": [{
			"mountPath": "/go-demo-config",
			"name": "config"
		}]
	}]
} }

EOF
)
	kubectl run -it --rm --image="$tke_docker_image" --restart=Never --command=true --pod-running-timeout=2m migrate-database -n "$namespace" --overrides="$yaml"
}

do_task() {
	init_kubectl_config
	local tke_namespace=`hub_var TKE_CLUSTER_NAMESPACE true`

	local action=`hub_var TASK_ACTION true`

	create_configmap "$tke_namespace" "go-demo-$tke_namespace"

	case "$action" in
	migrate_database)
		do_database_migrate "$tke_namespace"
		;;
	upgrade_service)
		# 创建service
		create_service "$tke_namespace"
		# 创建deploy
		create_deployment "$tke_namespace"
		;;

	esac
	echo "ok"
}

do_task

