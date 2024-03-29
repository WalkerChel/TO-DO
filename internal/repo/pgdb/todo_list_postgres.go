package pgdb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/WalkerChel/TO-DO/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list entity.TodoList) (int, error) {
	// checkTitle := fmt.Sprintf(`SELECT tl.id FROM %s tl
	// 							INNER JOIN %s ul ON ul.list_id = tl.id
	// 							WHERE ul.user_id = $1 AND LOWER(tl.title) = TRIM(LOWER($2))`,
	// 	todoListsTable, usersListsTable)

	var id int

	// err := r.db.Get(&id, checkTitle, userId, list.Title)

	// if err == nil {
	// 	return 0, fmt.Errorf("list with such title '%s' already exists", strings.ToLower(list.Title))
	// }

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING ID", todoListsTable)

	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)

	_, err = tx.Exec(createUsersListQuery, userId, id)

	if err != nil {
		tx.Rollback()
		return 0, nil
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]entity.TodoList, error) {
	var lists []entity.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)

	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (entity.TodoList, error) {
	var list entity.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2",
		todoListsTable, usersListsTable)

	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	deleteAllItemsQuery := fmt.Sprintf(`DELETE FROM %s ti 
										USING %s li, %s ul 
										WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND li.list_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	_, err = tx.Exec(deleteAllItemsQuery, userId, listId)

	if err != nil {
		// tx.Rollback()
		return fmt.Errorf("1)%s; 2)%s", err, tx.Rollback())
	}

	deleteListQuery := fmt.Sprintf(`DELETE FROM %s tl 
							USING %s ul 
							WHERE tl.id=ul.list_id AND ul.user_id=$1 AND ul.list_id=$2`,
		todoListsTable, usersListsTable)

	res, err := tx.Exec(deleteListQuery, userId, listId)

	if err != nil {
		// tx.Rollback()
		return fmt.Errorf("1)%s; 2)%s", err, tx.Rollback())
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("no rows deleted")
	}

	return tx.Commit()
}

func (r *TodoListPostgres) Update(userId, listId int, input entity.UpdateListInput) error {
	setValues := make([]string, 0)
	// args = [title, description, listId, userId]
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	res, err := r.db.Exec(query, args...)

	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("incorrect list or item id for update statement")
	}

	return err
}
