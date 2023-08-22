package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mymmrac/telego"
	"log"
	"main/utils"
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

func Init() (d *Data) {
	d = new(Data)
	d.logger, _ = utils.NewLogger("")
	connStr := "user=postgres password=2222 dbname=postgres sslmode=disable"
	var err error
	d.db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return
}

func (d *Data) Exec(tableName string, lines []Line) {
	log.Println(d.dropTable(tableName))
	log.Println(d.createTable(tableName))
	log.Println(d.InsertLines(tableName, lines))
}

func (d *Data) dropTable(tableName string) error {
	_, err := d.db.Exec(fmt.Sprintf("DROP TABLE %v", tableName))
	return err
}

func (d *Data) createTable(name string) error {
	query := fmt.Sprintf(`create table %v
	(
		id     serial primary key,
		key    text,
		value  text,
		parent text,
		source text,
		target text
	);`, name)
	_, err := d.db.Exec(query)
	return err
}

func (d *Data) Close() {
	_ = d.db.Close()
}

//insert into table values (1,1), (1,2), (1,3), (2,1);

type Line struct {
	Id     string `xml:"id,attr"`
	Value  string `xml:"value,attr"`
	Parent string `xml:"parent,attr"`
	Source string `xml:"source,attr"`
	Target string `xml:"target,attr"`
}

func (d *Data) InsertLines(table string, lines []Line) error {
	query := fmt.Sprintf("insert into %v (key, value, parent, source, target) values", table)
	for _, line := range lines {
		query = fmt.Sprintf("%v ('%v','%v','%v','%v','%v'),",
			query,
			line.Id,
			line.Value,
			line.Parent,
			line.Source,
			line.Target)
	}
	query = query[:len(query)-1] + ";"
	_, err := d.db.Exec(query)
	return err
}

func (d *Data) GetKey(table string, value string) (key string, err error) {
	row := d.db.QueryRow(fmt.Sprintf("select key from %v where value = '%v'", table, value))

	err = row.Scan(&key)
	return
}

func (d *Data) GetValue(table string, key string) (value string, err error) {
	row := d.db.QueryRow(fmt.Sprintf("select value from %v where key = '%v'", table, key))

	err = row.Scan(&value)
	return
}

func (d *Data) GetChild(table string, key string) (child map[string]string, err error) {
	rows, err := d.db.Query(fmt.Sprintf("select key, value from %v where parent = '%v'", table, key))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	child = make(map[string]string)
	for rows.Next() {
		var resultKey, resultValue string
		if err = rows.Scan(&resultKey, &resultValue); err != nil {
			return nil, err
		}
		child[resultKey] = resultValue
	}
	return
}

func (d *Data) GetTarget(table string, sourceKey string) (targetKey string, err error) {
	row := d.db.QueryRow(fmt.Sprintf("select target from %v where source = '%v'", table, sourceKey))

	err = row.Scan(&targetKey)
	return
}
