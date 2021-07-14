package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseShow struct {
	Name string `json:"name" xml:"name"`
}

func GoDotEnv(key string) string {
	path := os.Getenv("PDROP")
	err := godotenv.Load(fmt.Sprintf("%s.env", path))
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	arg := os.Args[1:]
	var db_pattern_name string
	if len(arg) == 0 {
		fmt.Printf("No PATTERN arguments. Using default in env\n")
		db_pattern_name = GoDotEnv("NAME")
	} else {
		fmt.Printf("Drop DB with pattern %s\n", arg[0])
		db_pattern_name = arg[0]
	}
	db, err := sql.Open("postgres", GoDotEnv("DB"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Can not connect to database : %v\n", err)
	}
	sql_stm := `
		select 'drop database '||n.datname||';' as "name"
		FROM(select datname from pg_database where datname like '%%` + db_pattern_name + `%%') as n
	`
	rows, _ := db.Query(sql_stm)
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
			fmt.Printf("Error from drop database : %v\n", err)
		}
	}

	defer rows.Close()
	defer db.Close()
}
