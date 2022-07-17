#!/bin/bash

VAGRANTFILE_PATH=$1
VCPU=$2
RAM=$3
WORKER_ID=$4
QUEUE_ID=$5

cd $VAGRANTFILE_PATH
vagrant ssh -c "mkdir ${WORKER_ID}"
vagrant ssh -c "git clone https://github.com/viniciusbds/arrebol-pb-worker ${WORKER_ID}"
vagrant ssh -c "mv ${WORKER_ID}/.env.example2 ${WORKER_ID}/.env"
vagrant ssh -c "cp server.pub ${WORKER_ID}/certs"
vagrant ssh -c "cd ${WORKER_ID} && ./create_worker_conf.sh  ${VCPU}  ${RAM}  ${WORKER_ID}  ${QUEUE_ID}"
vagrant ssh -c "cd ${WORKER_ID} && /usr/local/go/bin/go build"
vagrant ssh -c "cd ${WORKER_ID} && ./arrebol-pb-worker"
