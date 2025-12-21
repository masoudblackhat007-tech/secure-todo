package service

import "github.com/masoudblackhat007-tech/secure-todo/internal/domain"

type TodoService struct {
	// later: storage, logger, etc.
}

func NewTodoService() *TodoService {
	return &TodoService{}
}

func (s *TodoService) List() []domain.Todo {
	// TODO: implement with storage
	return nil
}
