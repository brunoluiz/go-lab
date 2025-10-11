#!/bin/bash

grpcurl \
  -protoset <(buf build -o -) -plaintext \
  -d '{}' \
  -vv \
  'localhost:4000' 'acme.api.todo.v1.TodoService/ListTasks'
