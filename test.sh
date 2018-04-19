#!/bin/bash

echo "Get the catalog:"
curl -X GET "http://localhost:3000/v2/catalog" -H "accept: application/json" -H "X-Broker-API-Version: 2.13"

echo "Provision a light:"
curl -X PUT "http://localhost:3000/v2/service_instances/aabbcc11" -H "accept: application/json" \
  -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"location-4a\",\"plan_id\":\"location-4a-kind-blue\"}"

echo "Bind to that light:"
curl -X PUT "http://localhost:3000/v2/service_instances/aabbcc11/service_bindings/binding-aabbcc11" \
  -H "accept: application/json" -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"location-4a\",\"plan_id\":\"location-4a-kind-blue\"}"


# -d "{\"Intensity\":\"1.0\"}"

echo Now:  curl -X PUT -H "Content-Type: application/json" -d "{\"intensity\":0.1}" <url>