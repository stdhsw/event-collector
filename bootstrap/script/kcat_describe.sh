#!/bin/bash
kubectl exec -n mytest kcat -- kcat -L -b mytest-kafka-kafka-bootstrap:9092