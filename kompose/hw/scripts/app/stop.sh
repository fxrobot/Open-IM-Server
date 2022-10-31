#!/usr/bin/env bash

kubectl -n openim delete -f ../../openimserver/service.yaml
kubectl -n openim delete -f ../../openimserver/deployment.yaml