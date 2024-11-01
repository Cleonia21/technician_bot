package repository

import "technician_bot/temporary/entities"

type Repo struct {
	db DataBase
}

type Value struct {
	Key   string
	Value string
}

type Children struct {
	Children map[string]string
}

type DataBase interface {
	GetValue(tableName string, key string) (Value, error)
	GetChild(tableName string, key string) (Children, error)
}

func NewRepo(db DataBase) *Repo {
	return &Repo{db: db}
}

func (db *Repo) Get(request entities.Request) (resp entities.Response, err error) {
	resp.Id = request.Id

	dbValue, err := db.db.GetValue(tableName, key)
	if err != nil {
		return entities.Response{}, err
	}
	resp.Value = dbValue.Value

	dbChild, err := db.db.GetChild(tableName, key)
	if err != nil {
		return entities.Response{}, err
	}
	resp.Options = db.parseChild(dbChild)

	return
}

func (db *Repo) parseChild(children Children) (options []entities.Option) {
	return options
}
