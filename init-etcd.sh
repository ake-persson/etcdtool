#!/bin/bash

IMAGE='quay.io/coreos/etcd:latest'
NAME='etcd'

case $(uname -s) in
    'Darwin')
        HOST_PORT=${DOCKER_HOST#tcp://}
        IP=${HOST_PORT%:[0-9]*}
        ;;
    'Linux')
        IP=$(ifconfig eth0 | awk '/inet / {print $2}')
        ;;
esac

start() {
    docker run -d -v /usr/share/ca-certificates/:/etc/ssl/certs -p 4001:4001 -p 2380:2380 -p 2379:2379 \
        --name etcd quay.io/coreos/etcd:v2.0.8 \
        -name etcd0 \
        -advertise-client-urls http://${IP}:2379,http://${IP}:4001 \
        -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
        -initial-advertise-peer-urls http://${IP}:2380 \
        -listen-peer-urls http://0.0.0.0:2380 \
        -initial-cluster-token etcd-cluster-1 \
        -initial-cluster etcd0=http://${IP}:2380 \
        -initial-cluster-state new
}

stop() {
    docker stop ${NAME}0
    docker rm ${NAME}0
}

status() {
    echo NODE1: ${NAME}1 IP: ${IP} PORT: 4001
}

CMD=$1
case ${CMD} in
    'start')
    start
    sleep 1
    status
    echo "# Run the following to export the environment"
    echo "# eval \"\$(./init-etcd.sh env)\""
        ;;
    'stop')
    stop
        ;;
    'restart')
    stop
    start
    sleep 1
    status
    ;;
    'status')
    status
        ;;
    'env')
    echo "export ETCDCTL_PEERS=\"http://${IP}:4001\""
    echo "# Run the following to export the environment"
    echo "# eval \"\$(./init-etcd.sh env)\""

    ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|env}"
        ;;
esac
