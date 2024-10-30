package db

import "technician_bot/temporary/entities"

type DB struct {
	db DataBase
}

type Value struct {
	Value string
}

type Child struct {
	Child map[string]string
}

type DataBase interface {
	GetValue(tableName string, key string) (Value, error)
	GetChild(tableName string, key string) (Child, error)
}

func NewDB(db DataBase) DB {
	return DB{db: db}
}

func (db *DB) GetResponse(request entities.Request) (resp entities.Response, err error) {
	resp.Id = request.Id

	tableName, key, err := db.parseData(request.Data)
	if err != nil {
		return entities.Response{}, err
	}

	dbValue, err := db.db.GetValue(tableName, key)
	if err != nil {
		return entities.Response{}, err
	}
	resp.Text = dbValue.Value

	dbChild, err := db.db.GetChild(tableName, key)
	if err != nil {
		return entities.Response{}, err
	}
	resp.Options = db.parseChild(dbChild)

	return
}

func (db *DB) parseData(data string) (tableName string, key string, err error) {
	return "", "", nil
}

func (db *DB) parseChild(child Child) (options []entities.Option) {
	return options
}
