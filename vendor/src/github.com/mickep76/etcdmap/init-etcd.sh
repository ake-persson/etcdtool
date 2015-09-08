#!/bin/bash

IMAGE='quay.io/coreos/etcd:v0.4.6'
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
    docker run -d -p 8001:8001 -p 5001:5001 --name ${NAME}1 ${IMAGE} -name ${NAME}1 \
        -peer-addr ${IP}:8001 -addr ${IP}:5001
    docker run -d -p 8002:8002 -p 5002:5002 --name ${NAME}2 ${IMAGE} -name ${NAME}2 \
        -peer-addr ${IP}:8002 -addr ${IP}:5002 -peers ${IP}:8001,${IP}:8002,${IP}:8003
    docker run -d -p 8003:8003 -p 5003:5003 --name ${NAME}3 ${IMAGE} -name ${NAME}3 \
        -peer-addr ${IP}:8003 -addr ${IP}:5003 -peers ${IP}:8001,${IP}:8002,${IP}:8003
}

stop() {
    docker stop ${NAME}3
    docker rm ${NAME}3
    docker stop ${NAME}2
    docker rm ${NAME}2
    docker stop ${NAME}1
    docker rm ${NAME}1
}

status() {
    echo NODE1: ${NAME}1 IP: ${IP} PORT: 5001
    echo NODE2: ${NAME}2 IP: ${IP} PORT: 5002
    echo NODE3: ${NAME}3 IP: ${IP} PORT: 5003
    echo
    curl -s -L ${IP}:5001/v2/stats/leader | python -mjson.tool
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
	echo "export ETCD_CONN=\"http://${IP}:5001\""
	echo "# Run the following to export the environment"
	echo "# eval \"\$(./init-etcd.sh env)\""

	;;
    *)
        echo "Usage: $0 {start|stop|restart|status|env}"
        ;;
esac
