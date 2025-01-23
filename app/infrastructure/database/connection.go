package database

import (
	"sync"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

// DBConnection is an interface that contains the database connection.
type DBConnection interface {
	Close() error
	Open() (proxy.DB, error)
}

// dbConnection is a struct that contains the database connection.
type dbConnection struct {
	sql            proxy.Sql
	db             proxy.DB
	driverName     string
	dataSourceName string
	mutex          *sync.RWMutex
}

// Close closes the database connection.
func (c *dbConnection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			return err
		}
		c.db = nil
	}
	return nil
}

// Open opens the database connection.
func (c *dbConnection) Open() (proxy.DB, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.db != nil {
		return c.db, nil
	}

	db, err := c.sql.Open(c.driverName, c.dataSourceName)
	if err != nil {
		return nil, err
	}

	c.db = db
	return c.db, nil
}
