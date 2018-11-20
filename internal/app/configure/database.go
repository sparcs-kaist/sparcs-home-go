package configure

import (
	"fmt"
	"log"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	// SQL : mysql connection
	SQL *sqlx.DB
)

// DatabaseInfo : database connection info
type DatabaseInfo struct {
	User     string
	Password string
	DBname   string
}

// ConnectDB to the database
func ConnectDB(i DatabaseInfo) {
	var err error
	addr := fmt.Sprintf("%s:%s@/%s", i.User, i.Password, i.DBname)
	log.Println("Addr: ", addr)
	if SQL, err = sqlx.Connect("mysql", addr); err != nil {
		log.Println("Failed to connect to mysql", err)
	}

	if err = SQL.Ping(); err != nil {
		log.Println("Database Error", err)
	}

}
