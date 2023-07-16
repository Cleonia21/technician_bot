package db

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/mymmrac/telego"
)

type Data struct {
	db     *sql.DB
	logger telego.Logger
}

type Content struct {
	Path  string
	Text  string
	Photo []byte
	Child []string
}

func Init(logger telego.Logger) (d Data) {
	d = Data{}
	d.logger = logger
	connStr := "user=postgres password=2222 dbname=postgres sslmode=disable"
	var err error
	d.db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return
}

func (d *Data) Close() {
	_ = d.db.Close()
}

func (d *Data) Select(table string, path string) (*Content, error) {
	row := d.db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE path = $1", table), path)
	if row == nil {
		return nil, row.Err()
	}

	var c Content
	var id int
	err := row.Scan(&id, &c.Path, &c.Text, (*pq.StringArray)(&c.Child), &c.Photo)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

/*
	func (data *Data) Insert(model *Model) {
		if data.Select(model.OrderUid).OrderUid == model.OrderUid {
			return
		}

		b, err := json.Marshal(model)
		_, err = data.db.Exec("INSERT INTO orders (orderUid, attrs) VALUES($1, $2)", model.OrderUid, b)
		if err != nil {
			panic(err)
		}
	}



	func (data *Data) SelectAll() []Model {
		row := data.db.QueryRow("SELECT count(*) FROM orders")
		var num int
		err := row.Scan(&num)
		if err != nil {
			panic(err)
		}

		rows, err := data.db.Query("SELECT attrs FROM orders")
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		model := make([]Model, num)

		var b []byte
		for i := 0; rows.Next(); i++ {
			err = rows.Scan(&b)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(b, &(model[i]))
		}
		return model
	}
*/
