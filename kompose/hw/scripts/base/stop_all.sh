#!/usr/bin/env bash

source ./path_info.cfg

for i in ${service[*]}
do
  kubectl -n openim delete -f ../../${i}/service.yaml
  kubectl -n openim delete -f ../../${i}/deployment.yaml
done