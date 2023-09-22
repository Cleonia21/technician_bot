package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mymmrac/telego"
	"os"
	"technician_bot/cmd/utils"
)

type DBinstance struct {
	db     *sql.DB
	logger telego.Logger
}

var DB DBinstance

func ConnectDB() {
	DB.logger, _ = utils.NewLogger("")

	var err error
	DB.db, err = sql.Open("postgres", connectString())

	if err != nil {
		DB.logger.Errorf("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	DB.logger.Debugf("data base connected")
}

func connectString() string {
	//connStr := "user=postgres password=2222 dbname=postgres port=5432 sslmode=disable"
	//connStr := "postgres://username:testpass@postgres:5432/mydb?sslmode=disable"
	//connStr := "postgres://cleonia:ek666x190@localhost:5432/tch_bot?sslmode=disable"
	var user, pass, host, dbName string
	var ok bool
	if user, ok = os.LookupEnv("POSTGRES_USER"); !ok {
		DB.logger.Errorf("env POSTGRES_USER not found")
	}
	if pass, ok = os.LookupEnv("POSTGRES_PASSWORD"); !ok {
		DB.logger.Errorf("env POSTGRES_PASSWORD not found")
	}
	if host, ok = os.LookupEnv("POSTGRES_HOST"); !ok {
		DB.logger.Errorf("env POSTGRES_HOST not found")
	}
	if dbName, ok = os.LookupEnv("POSTGRES_DB"); !ok {
		DB.logger.Errorf("env POSTGRES_DB not found")
	}
	connStr := fmt.Sprintf("postgres://%v:%v@%v:5432/%v?sslmode=disable",
		user,
		pass,
		host,
		dbName,
	)
	fmt.Println(connStr)
	return connStr
}

func DropTable(tableName string) error {
	_, err := DB.db.Exec(fmt.Sprintf("DROP TABLE %v", tableName))
	if err != nil {
		DB.logger.Errorf(err.Error())
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
	_, err := DB.db.Exec(query)
	if err != nil {
		DB.logger.Errorf(err.Error())
	}
	return err
}

func Close() {
	_ = DB.db.Close()
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
	_, err := DB.db.Exec(query)
	if err != nil {
		DB.logger.Errorf(err.Error())
	}
	return err
}

func GetKey(table string, value string) (key string, err error) {
	row := DB.db.QueryRow(fmt.Sprintf("select key from %v where value = '%v'", table, value))

	err = row.Scan(&key)
	return
}

func GetValue(table string, key string) (value string, err error) {
	row := DB.db.QueryRow(fmt.Sprintf("select value from %v where key = '%v'", table, key))

	err = row.Scan(&value)
	return
}

func GetChild(table string, key string) (child map[string]string, err error) {
	rows, err := DB.db.Query(fmt.Sprintf("select key, value from %v where parent = '%v'", table, key))
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
	row := DB.db.QueryRow(fmt.Sprintf("select parent from %v where key = '%v'", table, key))

	err = row.Scan(&parent)
	return
}

func GetTarget(table string, sourceKey string) (targetKey string, err error) {
	row := DB.db.QueryRow(fmt.Sprintf("select target from %v where source = '%v'", table, sourceKey))

	err = row.Scan(&targetKey)
	return
}

func GetSource(table string, targetKey string) (sourceKey string, err error) {
	row := DB.db.QueryRow(fmt.Sprintf("select source from %v where target = '%v'", table, targetKey))

	err = row.Scan(&sourceKey)
	return
}

func GetTables() (tables []string, err error) {
	rows, err := DB.db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return
}
