#!/bin/bash
curl -u "elastic:elastic" -k -X PUT "$ES_HOST/_template/event_template"  -H "Content-Type: application/json" -d '{
    "index_patterns": ["event-*"],
    "settings": {
        "index": {
            "number_of_shards": 3,
            "number_of_replicas": 1,
            "refresh_interval": "5s",
            "lifecycle": {
                "name": "event_policy",
                "rollover_alias": "event"
            }
        }
    },
    "mappings": {
        "dynamic": "false",
        "properties": {
            "metadata": {
                "properties": {
                    "name": { "type": "keyword" },
                    "namespace": { "type": "keyword" },
                    "uid": { "type": "keyword" },
                    "resourceVersion": { "type": "keyword" },
                    "creationTimestamp": { "type": "date" }
                }
            },
            "eventTime": { "type": "date" },
            "reportingController": { "type": "keyword" },
            "reason": { "type": "keyword" },
            "regarding": {
                "properties": {
                    "kind": { "type": "keyword" },
                    "namespace": { "type": "keyword" },
                    "name": { "type": "keyword" },
                    "uid": { "type": "keyword" },
                    "apiVersion": { "type": "keyword" },
                    "resourceVersion": { "type": "keyword" }
                }
            },
            "note": { "type": "text" },
            "type": { "type": "keyword" },
            "deprecatedFirstTimestamp": { "type": "date" },
            "deprecatedLastTimestamp": { "type": "date" },
            "deprecatedCount": { "type": "integer" }
        }
    }
}'
