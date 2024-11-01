package postgres

import "technician_bot/temporary/controller/db"

type Postgres struct {
}

func (p *Postgres) GetValue(tableName string, key string) (db.Value, error) {
	return db.Value{}, nil
}

func (p *Postgres) GetChild(tableName string, key string) (db.Child, error) {
	return db.Child{}, nil
}
