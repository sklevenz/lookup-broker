#!/usr/bin/env bash

curl 'http://localhost:5000/v2/service_instances/123?accepts_incomplete=true&service_id=1&plan_id=1.1' -X DELETE -H "X-Broker-API-Version: 2.16"