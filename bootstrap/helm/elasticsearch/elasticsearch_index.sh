#!/bin/bash
curl -u "elastic:elastic" -k -X PUT "$ES_HOST/event-000001"  -H "Content-Type: application/json" -d '{
    "aliases": {
        "event": {
            "is_write_index": true
        }
    }
}'
