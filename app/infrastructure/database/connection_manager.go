package database

import (
	"errors"
	"sync"

	"github.com/yanosea/jrp/pkg/proxy"
)

var (
	// gcm is a global connection manager.
	gcm ConnectionManager
	// gmutex is a global mutex.
	gmutex = &sync.Mutex{}
	// GetConnectionManagerFunc is a function to get the connection manager.
	GetConnectionManagerFunc = getConnectionManager
)

// ConnectionManager is an interface that manages database connections.
type ConnectionManager interface {
	CloseAllConnections() error
	CloseConnection(which DBName) error
	GetConnection(which DBName) (DBConnection, error)
	InitializeConnection(config ConnectionConfig) error
}

// connectionManager is a struct that implements the ConnectionManager interface.
type connectionManager struct {
	sql         proxy.Sql
	connections map[DBName]DBConnection
	mutex       *sync.RWMutex
}

// NewConnectionManager initializes the connection manager.
func NewConnectionManager(sql proxy.Sql) ConnectionManager {
	gmutex.Lock()
	defer gmutex.Unlock()

	if gcm == nil {
		gcm = &connectionManager{
			sql:         sql,
			connections: make(map[DBName]DBConnection),
			mutex:       &sync.RWMutex{},
		}
	}

	return gcm
}

// GetConnectionManager gets the connection manager.
func GetConnectionManager() ConnectionManager {
	return GetConnectionManagerFunc()
}

// getConnectionManager gets the connection manager.
func getConnectionManager() ConnectionManager {
	gmutex.Lock()
	defer gmutex.Unlock()

	if gcm == nil {
		return nil
	}

	return gcm
}

// ResetConnectionManager resets the connection manager.
func ResetConnectionManager() error {
	gmutex.Lock()
	defer gmutex.Unlock()

	if gcm == nil {
		return nil
	}

	if err := gcm.CloseAllConnections(); err != nil {
		return err
	} else {
		gcm = nil
	}

	return nil
}

// CloseAllConnections closes all database connections.
func (cm *connectionManager) CloseAllConnections() error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	for dbType, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			return err
		}
		delete(cm.connections, dbType)
	}

	return nil
}

// CloseConnection closes the database connection.
func (cm *connectionManager) CloseConnection(dbType DBName) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if conn, exists := cm.connections[dbType]; exists {
		if err := conn.Close(); err != nil {
			return err
		}
		delete(cm.connections, dbType)
	}

	return nil
}

// GetConnection gets the database connection.
func (cm *connectionManager) GetConnection(dbType DBName) (DBConnection, error) {
	cm.mutex.RLock()
	conn, exists := cm.connections[dbType]
	cm.mutex.RUnlock()

	if !exists {
		return nil, errors.New("connection not initialized")
	}

	return conn, nil
}

// InitializeConnectionManager initializes the connection manager.
func (cm *connectionManager) InitializeConnection(config ConnectionConfig) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if _, exists := cm.connections[config.DBName]; exists {
		return errors.New("connection already initialized")
	}

	cm.connections[config.DBName] = &dbConnection{
		sql:            cm.sql,
		db:             nil,
		driverName:     string(config.DBType),
		dataSourceName: config.DSN,
		mutex:          &sync.RWMutex{},
	}

	return nil
}
