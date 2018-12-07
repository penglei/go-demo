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
	local node_name=`hub_var TKE_MYSQL_NODE_NAME true`

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
      nodeName: "$node_name"
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
create_test_and_prod_database "$mysql_root_password"
}

create_test_and_prod_database() {
	local mysql_root_password="$1"
	kubectl delete pod wait-mysql-ready mysql-client > /dev/null 2>&1 || true
	kubectl run -i --tty --rm --image=busybox --restart=Never --command=true wait-mysql-ready -- sh -c 'for i in `seq 1 30`; do nc -z mysql 3306 && echo success && exit 0; echo -n .;  sleep 1; done' | grep 'success' || (echo error: mysql not ready. ; exit 1)
	# 创建测试和开发用的数据库
	kubectl run -i --tty --rm --image=mysql:5.7 --restart=Never --command=true mysql-client -- mysql -u root -h mysql -p"$mysql_root_password" -e "create database  if not exists go-demo-test; create database if not exists go-demo-prod; show databases;"
}

do_task() {
	init_kubectl_config
	kubectl get nodes
	create_mysql_deployment
	echo "ok"
}

do_task
