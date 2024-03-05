package todo

import "errors"

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UserList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

type ListItem struct {
	Id     int
	ListId int
	ItemId int
}

func (s *UpdateListInput) Validate() error {
	if s.Title == nil && s.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

func (s *UpdateItemInput) Validate() error {
	if s.Title == nil && s.Description == nil && s.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
