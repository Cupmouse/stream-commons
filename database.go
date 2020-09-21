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

// ConnectDatabase tries to connect to the main database
func ConnectDatabase() (*sql.DB, error) {
	return connectDatabase(
		os.Getenv("DATABASE_URL"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_DISABLE_CERT") == "1",
	)
}

// connectDatabase establishes connection to a database
func connectDatabase(url string, user string, password string, port string, disableCert bool) (*sql.DB, error) {
	if !disableCert {
		// TLS setup
		rootCertPool := x509.NewCertPool()
		pem, serr := ioutil.ReadFile("./rds-ca-2019-root.pem")
		if serr != nil {
			return nil, fmt.Errorf("rds ca cert: %v", serr)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, errors.New("failed to append pem")
		}
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs: rootCertPool,
		})
		// connect to a database
		database, serr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?timeout=5s&tls=custom", user, password, url, port))
		if serr != nil {
			return nil, fmt.Errorf("mysql connect: %v", serr)
		}
		return database, nil
	}
	database, serr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?timeout=5s", user, password, url, port))
	if serr != nil {
		return nil, fmt.Errorf("mysql connect: %v", serr)
	}
	return database, nil
}
