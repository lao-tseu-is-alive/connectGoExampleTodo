#!/bin/bash
echo "let's create tasks via classical Connect routes, does not work with vanguard"
curl -s  -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo -H "Content-Type: application/json" -d '{"title": "Buy milk"}'
curl -s -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo -H "Content-Type: application/json" -d '{"title": "Test connectRPC"}'
echo "now let's list them (still usingConnect routes)"
curl -s -X POST -H "Content-Type: application/json" http://localhost:8080/todo.v1.TodoService/ListTodos -d '{"page_size":10}' |jq
echo "now let's use classic REST routes using vanguard"
curl -s  -X POST http://localhost:8080/v1/todos -H "Content-Type: application/json" -d '{"title": "Buy milk via openapi REST"}'
curl -s  -X POST http://localhost:8080/v1/todos -H "Content-Type: application/json" -d '{"title": "Buy milk via Try Connect RPC with vanguard"}'
curl -s http://localhost:8080/v1/todos?page_size=10 |jq

echo "now create another one with grpc via buf"
buf curl --schema ./proto/todo/v1/todo.proto --protocol grpc --http2-prior-knowledge --data '{"title": "Testing grpc connectRPC"}'   http://localhost:8080/todo.v1.TodoService/CreateTodo
echo "now call ListTodos with grpc via buf"
buf curl --schema ./proto/todo/v1/todo.proto --protocol grpc --http2-prior-knowledge --data '{"page_size": "30"}'   http://localhost:8080/todo.v1.TodoService/ListTodos |jq
