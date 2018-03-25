#!/bin/bash

echo "Get the catalog:"
curl -X GET "http://localhost:3000/v2/catalog" -H "accept: application/json" -H "X-Broker-API-Version: 2.13"

echo "Provision a light:"
curl -X PUT "http://localhost:3000/v2/service_instances/aabbcc" -H "accept: application/json" \
  -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"location-bedroom\",\"plan_id\":\"location-bedroom-kind-green\"}"

echo "Bind to that light:"
curl -X PUT "http://localhost:3000/v2/service_instances/aabbcc/service_bindings/binding-aabbcc" \
  -H "accept: application/json" -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"location-bedroom\",\"plan_id\":\"location-bedroom-kind-green\"}"


# -d "{\"Intensity\":\"1.0\"}"

echo Now:  curl -X PUT -H "Content-Type: application/json" -d "{\"intensity\":1.0}" <url>