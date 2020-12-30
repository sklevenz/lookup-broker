#!/usr/bin/env bash

curl -iLs 'http://localhost:5000/v2/service_instances/123/service_bindings/456?accepts_incomplete=true' -X GET -H "X-Broker-API-Version: 2.16" 

