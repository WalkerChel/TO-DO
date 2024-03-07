package service

import (
	"github.com/WalkerChel/TO-DO/internal/entity"
	"github.com/WalkerChel/TO-DO/internal/repo"
)

type TodoItemService struct {
	repo     repo.TodoItem
	listRepo repo.TodoList
}

func NewTodoItemService(repo repo.TodoItem, listRepo repo.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(userId, listId int, item entity.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		//list does't exist or does't belong to user
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]entity.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, listId, itemId int) (entity.TodoItem, error) {
	return s.repo.GetById(userId, listId, itemId)
}

func (s *TodoItemService) Delete(userId, listId, itemId int) error {
	return s.repo.Delete(userId, listId, itemId)
}

func (s *TodoItemService) Update(userId, listId, itemId int, input entity.UpdateItemInput) error {

	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, listId, itemId, input)
}
