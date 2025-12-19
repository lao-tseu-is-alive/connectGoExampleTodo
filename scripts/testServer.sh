#!/bin/bash
echo "let's create tasks"
curl -s  -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo -H "Content-Type: application/json" -d '{"title": "Buy milk"}'
curl -s -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo -H "Content-Type: application/json" -d '{"title": "Test connectRPC"}'
echo "now let's list them"
curl -s -X POST -H "Content-Type: application/json" http://localhost:8080/todo.v1.TodoService/ListTodos -d '{"page_size":10}' |jq
echo "now create another one with grpc via buf"
buf curl --schema ./proto/todo/v1/todo.proto --protocol grpc --http2-prior-knowledge --data '{"title": "Testing grpc connectRPC"}'   http://localhost:8080/todo.v1.TodoService/CreateTodo
