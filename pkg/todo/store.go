package todo

import (
	"context"
	"errors"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	todov1 "github.com/lao-tseu-is-alive/connectGoExampleTodo/gen/todo/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Store struct {
	mu    sync.RWMutex
	todos map[string]*todov1.Todo
}

func NewTodoStore() *Store {
	return &Store{todos: make(map[string]*todov1.Todo)}
}

func (s *Store) Create(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.Todo, error) {
	now := time.Now()
	id := uuid.NewString()
	todo := &todov1.Todo{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}
	s.mu.Lock()
	s.todos[id] = todo
	s.mu.Unlock()
	return todo, nil
}

func (s *Store) Get(ctx context.Context, id string) (*todov1.Todo, error) {
	s.mu.RLock()
	todo, ok := s.todos[id]
	s.mu.RUnlock()
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	return todo, nil
}

func (s *Store) List(ctx context.Context, pageSize int32, pageToken string) ([]*todov1.Todo, string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var todos []*todov1.Todo
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	// Simple pagination (implement properly for large datasets)
	if pageSize > 0 && int(pageSize) < len(todos) {
		todos = todos[:pageSize]
	}
	return todos, "", nil
}

func (s *Store) Update(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo, ok := s.todos[req.Id]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = timestamppb.New(time.Now())
	return todo, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return connect.NewError(connect.CodeNotFound, errors.New("todo not found"))
	}
	delete(s.todos, id)
	return nil
}
