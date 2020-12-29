#!/usr/bin/env bash

curl -iLs 'http://localhost:5000/v2/service_instances/123?accepts_incomplete=true' -d '{
  "service_id": "1",
  "plan_id": "1.1",
  "context": {
    "platform": "cloudfoundry",
    "some_field": "some-contextual-data"
  },
  "organization_guid": "org-guid-here",
  "space_guid": "space-guid-here",
  "parameters": {
    "parameter1": 1,
    "parameter2": "foo"
  }
}' -X PUT -H "X-Broker-API-Version: 2.16" -H "Content-Type: application/json"
