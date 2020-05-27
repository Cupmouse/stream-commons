package streamcommons

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// APIKey is the object of API-key
type APIKey struct {
	Key  []byte
	Demo bool
}

// CheckAvalability returns error with detailed message
func (a *APIKey) CheckAvalability(db *sql.DB) error {
	if a.Demo {
		// availability of test apikey can not be checked
		// this is a fail safe machanism so as to prevent bugs that allow
		// unlimited access to the API
		return errors.New("availabiliy of demo API-key can not be checked")
	}

	// check if API key is valid
	var apikeyAvailable int
	row := db.QueryRow("SELECT exchangedataset.apikey_available(?)", a.Key)
	err := row.Scan(&apikeyAvailable)
	if err != nil {
		return fmt.Errorf("apikey_available procedure call failed: %s", err.Error())
	}
	if apikeyAvailable != 1 {
		return errors.New("API key does not exist or reached the quota or is not enabled")
	}
	return nil
}

// IncrementUsed tries to increment used count tied to API key
func (a *APIKey) IncrementUsed(db *sql.DB, amount int64) (incremented int64, err error) {
	if a.Demo {
		err = errors.New("this is demo test apikey, can not perform quota increment")
		return
	}
	incremented = CalcQuotaUsed(amount)

	// increase api-key's quota used bytes
	res, serr := db.Exec("CALL exchangedataset.increment_apikey_used_now(?, ?)", a.Key, incremented)
	if serr != nil {
		err = fmt.Errorf("Call failed: %s", err.Error())
		return
	}
	rows, serr := res.RowsAffected()
	if serr != nil {
		err = fmt.Errorf("RowsAffected returned error: %s", err.Error())
		return
	}
	if rows != 1 {
		err = errors.New("Too many or less rows affected: API key might not exist")
		return
	}
	return
}

// CalcQuotaUsed returns how much quota will be decremented if a request of specified amount were processed
func CalcQuotaUsed(amount int64) (incremented int64) {
	// request always consume 1kb of quota at minimun for spamming counter method
	if amount < 1024 {
		incremented = 1024
	} else {
		incremented = amount
	}
	return
}

// NewAPIKey creates new instance of APIKey from headers
func NewAPIKey(event events.APIGatewayProxyRequest) (apikey *APIKey, err error) {
	headerAuthorization, ok := event.Headers["Authorization"]
	if !ok {
		// if Authroization header was not present, reject
		err = errors.New("Authorization header is not present")
		return
	}
	if !strings.HasPrefix(headerAuthorization, "Bearer") {
		// does not have Bearer as prefix
		return nil, errors.New("you must add 'Bearer' as a prefix to API-key")
	}
	apikeyString := strings.TrimSpace(strings.TrimPrefix(headerAuthorization, "Bearer"))
	if apikeyString == APIKeyDemo {
		// demo apikey
		return &APIKey{[]byte(APIKeyDemo), true}, nil
	}
	// normal apikey
	var key []byte
	key, err = base64.RawURLEncoding.DecodeString(apikeyString)
	if err != nil {
		err = fmt.Errorf("base64 decoding failed: %v", err)
		return
	}
	return &APIKey{key, false}, nil
}
