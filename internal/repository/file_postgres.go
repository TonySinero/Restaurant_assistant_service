package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var (
	tablesColumnsName = map[string][2]string{
		"/restaurant/image/:id": {"restaurants", "image"},
		"/dish/image/:id":       {"dishes", "image"},
	}
)

type FilePostgres struct {
	db *sqlx.DB
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}

func (r *FilePostgres) CheckUUID(uuid string, path string) error {
	tableColumn, ex := tablesColumnsName[path]
	if !ex {
		return errors.New("wrong path")
	}

	var result string
	checkRowQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = $1", tableColumn[0])

	err := r.db.Get(&result, checkRowQuery, uuid)
	if err != nil {
		log.Error().Err(err).Msg("uuid doesn't exist")
		return err
	}

	return nil
}

func (r *FilePostgres) Create(link string, uuid string, path string) error {
	tableColumn, ex := tablesColumnsName[path]
	if !ex {
		return errors.New("wrong path")
	}

	createTempImageRowQuery := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", tableColumn[0], tableColumn[1])

	_, err := r.db.Exec(createTempImageRowQuery, link, uuid)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while update image in restaurantservice SQL")
		return err
	}

	return nil
}
