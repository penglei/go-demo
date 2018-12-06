FROM ubuntu:stretch

LABEL TencentHubComponent='{"description": "init the infra"}'

WORKDIR /

ENTRYPOINT /entrypoint.sh

RUN apt-get update -y && apt-get install -y apt-transport-https && \
(curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -) && \
(echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list) \
apt-get update && apt-get install -y kubectl

ADD ./*.sh /

