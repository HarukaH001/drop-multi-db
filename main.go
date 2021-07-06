package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DatabaseShow struct {
	Name string `json:"name" xml:"name"`
}

func GoDotEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can not connect to database : %v\n", err)
	}

	// ! change pattern here V
	db_pattern_name := GoDotEnv("NAME")
	// !

	sql_stm := `
		select 'drop database '||n.datname||';' as "name"
		FROM(select datname from pg_database where datname like '%%` + db_pattern_name + `%%') as n
	`
	rows, err := db.Query(sql_stm)
	data := []DatabaseShow{}
	for rows.Next() {
		temp := DatabaseShow{}
		err = rows.Scan(&temp.Name)
		if err != nil {
			log.Fatalf("Error from scan : %v\n", err)
		}
		data = append(data, temp)
	}

	for _, ele := range data {
		_, err := db.Exec(ele.Name)
		if err != nil {
			log.Fatalf("Error from drop database : %v\n", err)
		}
	}

	defer db.Close()
}
