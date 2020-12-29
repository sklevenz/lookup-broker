#!/usr/bin/env bash

curl -iLs 'http://localhost:5000/v2/service_instances/123' -X GET -H "X-Broker-API-Version: 2.16" 