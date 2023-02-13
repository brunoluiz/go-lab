#!/bin/bash

set +xae

content_type='Content-type: application/json'
endpoint=http://localhost:8080/api/v1

case $1 in
  "add_list")
    curl -sS -X POST \
      -H "$content_type" \
      -d '{"title":"Hello World"}' \
       "$endpoint/lists" | jq;;
  "get_list")
    curl -sS -X GET \
      -H "$content_type" \
       "$endpoint/lists/$2" | jq;;
  "update_list")
    curl -sS -X PUT \
      -H "$content_type" \
      -d '{"title":"Something"}' \
       "$endpoint/lists/$2" | jq;;
  "delete_list")
    curl -sS -X DELETE \
      -H "$content_type" \
       "$endpoint/lists/$2" | jq;;

  *)
    echo 'command not found';;
esac
