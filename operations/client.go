package operations

// DATABASE_URL=postgres://{user}:{password}@{hostname}:{port}/{database-name}
// https://pkg.go.dev/github.com/jackc/pgx/v5/pgtype

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbClient *pgxpool.Pool
)

// clientStartup - entry point for creating the pooling client and verifying a connection can be acquired
func clientStartup(cfg map[string]string) error {
	if cfg == nil {
		return errors.New("map configuration is nil")
	}
	url, ok := cfg[uriConfigKey]
	if !ok {
		return errors.New("database URL is empty")
	}
	// Create connection string with credentials
	s, err := connectString(url, cfg)
	if err != nil {
		return err
	}
	// Create pooled client and acquire connection
	dbClient, err = pgxpool.New(context.Background(), s)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to create connection pool: %v\n", err))
	}
	conn, err1 := dbClient.Acquire(context.Background())
	if err1 != nil {
		clientShutdown()
		return errors.New(fmt.Sprintf("unable to acquire connection from pool: %v\n", err1))
	}
	conn.Release()
	return nil
}

func clientShutdown() {
	if dbClient != nil {
		dbClient.Close()
		dbClient = nil
	}
}

func connectString(url string, cfg map[string]string) (string, error) {
	user, ok := cfg[userConfigKey]
	pswd, ok0 := cfg[pswdConfigKey]
	// Username and password can be in the connect string Url
	if !ok && !ok0 {
		return url, nil
	}
	if !ok {
		return "", errors.New("error: user is not configured")
	}
	if !ok0 {
		return "", errors.New("error: password is not configured")
	}
	return fmt.Sprintf(url, user, pswd), nil
}

/*
// accessCredentials - access function for Credentials in a message
func accessCredentials(msg *messaging.Message) startupCredentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(func() (user string, pswd string, err error)); ok {
			return fn
		}
	}
	return nil
}

// accessResource - access function for a test in a message
func accessResource(msg *messaging.Message) startupResource {
	if msg == nil || msg.Content == nil {
		return startupResource{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(struct{ Uri string }); ok {
			return url
		}
	}
	return startupResource{}
}


*/
