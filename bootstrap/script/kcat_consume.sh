#!/bin/bash
kubectl exec -it -n mytest kcat -- kcat -C -b mytest-kafka-kafka-bootstrap:9092 -t event