#!/bin/bash
kubectl exec -it -n mytest kcat -- kcat -P -b mytest-kafka-kafka-bootstrap:9092 -t test-topic