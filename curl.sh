#!/bin/bash

set +xae

content_type='Content-type: application/json'
endpoint=http://localhost:8080/api

case $1 in
  "add_radar")
    curl -sS -X POST \
      -H "$content_type" \
      -d '{"title":"Hello World"}' \
       "$endpoint/radars" | jq;;
  "add_radar_item")
    curl -sS -X POST \
      -H "$content_type" \
      -d '{"name":"Hello World", "quadrant_id": "'$3'", "description": "Foo Bar"}' \
       "$endpoint/radars/$2/items" | jq;;
  "get_radar")
    curl -sS -X GET \
      -H "$content_type" \
       "$endpoint/radars/$2" | jq;;
  "update_radar")
    curl -sS -X PUT \
      -H "$content_type" \
      -d '{"title":"Something"}' \
       "$endpoint/radars/$2" | jq;;
  "delete_radar")
    curl -sS -X DELETE \
      -H "$content_type" \
       "$endpoint/radars/$2" | jq;;

  *)
    echo 'command not found';;
esac
