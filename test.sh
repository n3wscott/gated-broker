#!/bin/bash

printf "\n\nGet the catalog:\n\n"

curl -X GET "http://localhost:3000/v2/catalog" -H "accept: application/json" -H "X-Broker-API-Version: 2.13"

printf "\n\nProvision the light registry:\n\n"

curl -X PUT "http://localhost:3000/v2/service_instances/lightreg" -H "accept: application/json" \
  -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"light-registry\",\"plan_id\":\"default\"}"

printf "\n\nProvision a light:\n\n"

curl -X PUT "http://localhost:3000/v2/service_instances/aabbcc11" -H "accept: application/json" \
  -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"location-4a\",\"plan_id\":\"location-4a-kind-blue\"}"

printf "\n\nBind to that light:\n\n"

curl -X PUT "http://localhost:3000/v2/service_instances/aabbcc11/service_bindings/binding-aabbcc11" \
  -H "accept: application/json" -H "X-Broker-API-Version: 2.13" \
  -H "X-Broker-API-Originating-Identity: petstore eyJoZWxsbyI6MH0=" \
  -H "Content-Type: application/json" \
  -d "{\"service_id\":\"location-4a\",\"plan_id\":\"location-4a-kind-blue\"}"


# -d "{\"Intensity\":\"1.0\"}"

printf '\n\n'
echo "Now:  curl -X PUT -H \"Content-Type: application/json\" -d '{\"intensity\":0.1}' <url>"
printf '\n\n'