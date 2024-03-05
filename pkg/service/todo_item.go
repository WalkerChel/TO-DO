package service

import (
	todo "github.com/WalkerChel/TO-DO"
	"github.com/WalkerChel/TO-DO/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		//list does't exist or does't belong to user
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, listId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, listId, itemId)
}

func (s *TodoItemService) Delete(userId, listId, itemId int) error {
	return s.repo.Delete(userId, listId, itemId)
}

// func (s *TodoItemService) DeleteAll(userId, listId int) error {
// 	return s.repo.DeleteAll(userId, listId)
// }

func (s *TodoItemService) Update(userId, listId, itemId int, input todo.UpdateItemInput) error {

	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, listId, itemId, input)
}
