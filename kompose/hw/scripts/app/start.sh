#!/usr/bin/env bash

kubectl -n openim apply -f ../../openimserver/deployment.yaml
kubectl -n openim apply -f ../../openimserver/service.yaml