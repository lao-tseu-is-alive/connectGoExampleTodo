#!/bin/bash
echo "let's create tasks"
curl -s  -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo -H "Content-Type: application/json" -d '{"title": "Buy milk"}'
curl -s -X POST http://localhost:8080/todo.v1.TodoService/CreateTodo -H "Content-Type: application/json" -d '{"title": "Test connectRPC"}'
curl -s -X POST -H "Content-Type: application/json" http://localhost:8080/todo.v1.TodoService/ListTodos -d '{"page_size":10}' |jq
