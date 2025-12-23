package todo

import (
	"context"
	"log/slog"

	todov1 "github.com/lao-tseu-is-alive/connectGoExampleTodo/gen/todo/v1"
)

type Service struct {
	store           *Store
	Logger          *slog.Logger
	ClientHeaderKey string
}

func NewTodoService(clientHeaderKey string, l *slog.Logger) *Service {
	return &Service{
		store:           NewTodoStore(),
		Logger:          l,
		ClientHeaderKey: clientHeaderKey,
	}
}

func (s *Service) CreateTodo(ctx context.Context, req *todov1.CreateTodoRequest) (*todov1.CreateTodoResponse, error) {
	s.Logger.Info("entering CreateTodo", "request", req)
	todo, err := s.store.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &todov1.CreateTodoResponse{Todo: todo}, nil
}

func (s *Service) GetTodo(ctx context.Context, req *todov1.GetTodoRequest) (*todov1.GetTodoResponse, error) {
	s.Logger.Info("entering GetTodo", "request", req)
	todo, err := s.store.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &todov1.GetTodoResponse{Todo: todo}, nil
}

func (s *Service) ListTodos(ctx context.Context, req *todov1.ListTodosRequest) (*todov1.ListTodosResponse, error) {
	s.Logger.Info("entering ListTodos", "request", req)
	todos, nextToken, err := s.store.List(ctx, req.PageSize, req.PageToken)
	if err != nil {
		return nil, err
	}
	return &todov1.ListTodosResponse{Todos: todos, NextPageToken: nextToken}, nil
}

func (s *Service) UpdateTodo(ctx context.Context, req *todov1.UpdateTodoRequest) (*todov1.UpdateTodoResponse, error) {
	s.Logger.Info("entering UpdateTodo", "request", req)
	todo, err := s.store.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return &todov1.UpdateTodoResponse{Todo: todo}, nil
}

func (s *Service) DeleteTodo(ctx context.Context, req *todov1.DeleteTodoRequest) (*todov1.DeleteTodoResponse, error) {
	s.Logger.Info("entering DeleteTodo", "request", req)
	if err := s.store.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &todov1.DeleteTodoResponse{}, nil
}
