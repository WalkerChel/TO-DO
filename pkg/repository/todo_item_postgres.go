package repository

import (
	"errors"
	"fmt"
	"strings"

	todo "github.com/WalkerChel/TO-DO"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)

	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listsItemsTable)

	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti 
							INNER JOIN %s li ON li.item_id = ti.id 
							INNER JOIN %s ul ON ul.list_id = li.list_id 
							WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, listId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done 
							FROM %s ti 
							INNER JOIN %s li ON ti.id = li.item_id 
							INNER JOIN %s ul ON ul.list_id = li.list_id 
							WHERE ul.user_id = $1 AND li.list_id = $2 AND ti.id =$3`,
		todoItemsTable, listsItemsTable, usersListsTable)

	err := r.db.Get(&item, query, userId, listId, itemId)

	return item, err
}

func (r *TodoItemPostgres) Delete(userId, listId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti 
							USING %s li, %s ul 
							WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND li.list_id = $2 AND ti.id = $3`,
		todoItemsTable, listsItemsTable, usersListsTable)

	res, err := r.db.Exec(query, userId, listId, itemId)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("no rows deleted")
	}

	return err
}

// func (r *TodoItemPostgres) DeleteAll(userId, listId int) error {
// 	query := fmt.Sprintf(`DELETE FROM %s ti 
// 							USING %s li, %s ul 
// 							WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND li.list_id = $2 `,
// 		todoItemsTable, listsItemsTable, usersListsTable)

// 	res, err := r.db.Exec(query, userId, listId)

// 	if rows, _ := res.RowsAffected(); rows == 0 {
// 		return errors.New("no rows deleted")
// 	}

// 	return err
// }

func (r *TodoItemPostgres) Update(userId, listId, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)

	args := make([]interface{}, 0)

	argId := 1

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ",")

	query := fmt.Sprintf(`UPDATE %s ti 
						SET %s FROM %s li, %s ul 
						WHERE ul.list_id = li.list_id AND ul.list_id = %d AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, listId, argId, argId+1)

	args = append(args, userId, itemId)

	res, err := r.db.Exec(query, args...)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("incorrect list or item id for update statement")
	}

	if err != nil {
		return err
	}

	return nil
}
