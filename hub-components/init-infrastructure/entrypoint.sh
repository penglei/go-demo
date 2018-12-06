#!/bin/bash

#Example:
#
# docker run --rm \
#   -e GIT_CLONE_URL='https://github.com/qcloud2018/go-demo.git' \
#   -e GIT_REPO_DIR=/go/src/github.com/qcloud2018/go-demo \
#   hub.tencentyun.com/workshop/go-analysis
#

set -e
export SCRIPTDIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

source $SCRIPTDIR/component-base/libs.sh

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

create_mysql_deployment() {

	local mysql_root_password=`hub_var MYSQL_ROOT_PASSWORD true`

#cat <<EOF | kubectl create -f -
#EOF

cat > /mysql.yaml <<EOF
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
  - name: mysql
    port: 3306
  clusterIP: None
  selector:
    app: mysql
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  serviceName: mysql
  replicas: 1
  template:
    metadata:
      labels:
        app: mysql
    spec:
      volumes:
      - name: data
        hostPath:
          path: /data/go-demo/mysql
          type: DirectoryOrCreate
      nodeName: "172.27.100.24"
      containers:
      - name: mysql
        image: mysql:5.7
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "$mysql_root_password"
        ports:
        - name: mysql
          containerPort: 3306
        volumeMounts:
        - name: data
          mountPath: /var/lib/mysql
        livenessProbe:
          exec:
            command: ["mysqladmin", "ping"]
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
        readinessProbe:
          exec:
            # Check we can execute queries over TCP (skip-networking is off).
            command: ["mysql", "-h", "127.0.0.1", "-u", "root", "-p$mysql_root_password", "-e", "SELECT 1"]
          initialDelaySeconds: 5
          periodSeconds: 2
          timeoutSeconds: 1
EOF

kubectl apply -f /mysql.yaml
cat /mysql.yaml
}

do_task() {
	init_kubectl_config
	kubectl get nodes
	create_mysql_deployment
	echo "ok"
}

do_task
