package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mymmrac/telego"
	"log"
	"os"
	"technician_bot/cmd/utils"
)

type DBinstance struct {
	DB     *sql.DB
	logger telego.Logger
}

var DB DBinstance

func ConnectDB() {
	DB.logger, _ = utils.NewLogger("")

	connStr := "user=postgres password=2222 dbname=postgres port=5432 sslmode=disable"
	var err error
	DB.DB, err = sql.Open("postgres", connStr) //"postgres", "postgres://postgres:2222@postgres:5432/postgres?sslmode=disable")

	if err != nil {
		DB.logger.Errorf("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	DB.logger.Debugf("data base connected")
}

func DropTable(tableName string) error {
	_, err := DB.DB.Exec(fmt.Sprintf("DROP TABLE %v", tableName))
	if err != nil {
		log.Println(err)
	}
	return err
}

func CreateTable(name string) error {
	query := fmt.Sprintf(`create table %v
	(
		id     serial primary key,
		key    text,
		value  text,
		parent text,
		source text,
		target text
	);`, name)
	_, err := DB.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
	return err
}

func Close() {
	_ = DB.DB.Close()
}

type Line struct {
	Id     string `xml:"id,attr"`
	Value  string `xml:"value,attr"`
	Parent string `xml:"parent,attr"`
	Source string `xml:"source,attr"`
	Target string `xml:"target,attr"`
}

func InsertLines(table string, lines []Line) error {
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
	_, err := DB.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
	return err
}

func GetKey(table string, value string) (key string, err error) {
	row := DB.DB.QueryRow(fmt.Sprintf("select key from %v where value = '%v'", table, value))

	err = row.Scan(&key)
	return
}

func GetValue(table string, key string) (value string, err error) {
	row := DB.DB.QueryRow(fmt.Sprintf("select value from %v where key = '%v'", table, key))

	err = row.Scan(&value)
	return
}

func GetChild(table string, key string) (child map[string]string, err error) {
	rows, err := DB.DB.Query(fmt.Sprintf("select key, value from %v where parent = '%v'", table, key))
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

func GetParent(table string, key string) (parent string, err error) {
	row := DB.DB.QueryRow(fmt.Sprintf("select parent from %v where key = '%v'", table, key))

	err = row.Scan(&parent)
	return
}

func GetTarget(table string, sourceKey string) (targetKey string, err error) {
	row := DB.DB.QueryRow(fmt.Sprintf("select target from %v where source = '%v'", table, sourceKey))

	err = row.Scan(&targetKey)
	return
}

func GetSource(table string, targetKey string) (sourceKey string, err error) {
	row := DB.DB.QueryRow(fmt.Sprintf("select source from %v where target = '%v'", table, targetKey))

	err = row.Scan(&sourceKey)
	return
}
