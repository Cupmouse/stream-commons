package streamcommons

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-sql-driver/mysql"
	// needs to initialize mysql database driver
	_ "github.com/go-sql-driver/mysql"
)

var url = os.Getenv("DATABASE_URL")
var user = os.Getenv("DATABASE_USER")
var password = os.Getenv("DATABASE_PASSWORD")
var port = os.Getenv("DATABASE_PORT")

// ConnectDatabase tries to connect to the main database
func ConnectDatabase() (*sql.DB, error) {
	return connectDatabase(url, user, password, port)
}

// connectDatabase establishes connection to a database
func connectDatabase(url string, user string, password string, port string) (*sql.DB, error) {
	// TLS setup
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("./rds-ca-2019-root.pem")
	if err != nil {
		return nil, fmt.Errorf("io error when reading rds ca cert: %s", err.Error())
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, errors.New("failed to append pem")
	}
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})
	// connect to a database
	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?timeout=5s&tls=custom", user, password, url, port))
	if err != nil {
		return nil, fmt.Errorf("connecting to mysql server failed: %s", err.Error())
	}

	return database, nil
}
