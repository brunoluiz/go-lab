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

  *)
    echo 'command not found';;
esac
