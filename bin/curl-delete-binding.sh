#!/usr/bin/env bash

curl -iLs 'http://localhost:5000/v2/service_instances/123/service_bindings/456?service_id=1&plan_id=1.1&accepts_incomplete=true' -X DELETE -H "X-Broker-API-Version: 2.16"

