package repo

import (
	"github.com/WalkerChel/TO-DO/internal/entity"
	"github.com/WalkerChel/TO-DO/internal/repo/pgdb"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username, password string) (entity.User, error)
}

type TodoList interface {
	Create(userId int, list entity.TodoList) (int, error)
	GetAll(useId int) ([]entity.TodoList, error)
	GetById(userId, listId int) (entity.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input entity.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item entity.TodoItem) (int, error)
	GetAll(userId, listId int) ([]entity.TodoItem, error)
	GetById(userId, listId, itemId int) (entity.TodoItem, error)
	Delete(userId, listId, itemId int) error
	Update(userId, listId, itemId int, input entity.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: pgdb.NewAuthPostgres(db),
		TodoList:      pgdb.NewTodoListPostgres(db),
		TodoItem:      pgdb.NewTodoItemPostgres(db),
	}
}
