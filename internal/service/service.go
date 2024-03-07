package service

import (
	"github.com/WalkerChel/TO-DO/internal/entity"
	"github.com/WalkerChel/TO-DO/internal/repo"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list entity.TodoList) (int, error)
	GetAll(userId int) ([]entity.TodoList, error)
	GetById(userId, listId int) (entity.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input entity.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item entity.TodoItem) (int, error)
	GetAll(userId, listId int) ([]entity.TodoItem, error)
	GetById(userId, listId, itemId int) (entity.TodoItem, error)
	Delete(userId, listId, itemId int) error
	// DeleteAll(userId, listId int) error
	Update(userId, listId, itemId int, input entity.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repo.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
