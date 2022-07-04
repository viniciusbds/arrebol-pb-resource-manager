#!/bin/bash

cd $1

VCPU=$2
RAM=$3
WORKER_ID=$4
QUEUE_ID=$5

vagrant ssh -c "mkdir ${WORKER_ID}"
vagrant ssh -c "git clone https://github.com/ufcg-lsd/arrebol-pb-worker ${WORKER_ID}"
vagrant ssh -c "touch ${WORKER_ID}/conf.json"
vagrant ssh -c "echo {vcpu:$VCPU, ram:${RAM}, id:${WORKER_ID}, queue_id:${QUEUE_ID}} > ${WORKER_ID}/conf.json"
