package masterdatamanagement

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
)

// MasterData is the interface for managing master data
type MasterData interface {
	// Get retrieves a record from the master data store
	Get(key string) (interface{}, error)
	// Set stores a record in the master data store
	Set(key string, value interface{}) error
	// Delete removes a record from the master data store
	Delete(key string) error
	// List returns a list of all records in the master data store
	List() ([]interface{}, error)
}

// MasterDataManager is the implementation of the MasterData interface
type MasterDataManager struct {
	db     *sql.DB
	lock   sync.RWMutex
	tables map[string]string
}

// NewMasterDataManager creates a new instance of the MasterDataManager
func NewMasterDataManager(db *sql.DB, tables map[string]string) *MasterDataManager {
	return &MasterDataManager{
		db:     db,
		tables: tables,
	}
}

// Get retrieves a record from the master data store
func (m *MasterDataManager) Get(key string) (interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	table, ok := m.tables[key]
	if !ok {
		return nil, errors.New("invalid key")
	}

	row := m.db.QueryRow("SELECT * FROM " + table + " WHERE key = ?", key)
	var value interface{}
	err := row.Scan(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Set stores a record in the master data store
func (m *MasterDataManager) Set(key string, value interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	table, ok := m.tables[key]
	if !ok {
		return errors.New("invalid key")
	}

	_, err := m.db.Exec("INSERT INTO " + table + " (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a record from the master data store
func (m *MasterDataManager) Delete(key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	table, ok := m.tables[key]
	if !ok {
		return errors.New("invalid key")
	}

	_, err := m.db.Exec("DELETE FROM " + table + " WHERE key = ?", key)
	if err != nil {
		return err
	}

	return nil
}

// List returns a list of all records in the master data store
func (m *MasterDataManager) List() ([]interface{}, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	var values []interface{}
	for _, table := range m.tables {
		rows, err := m.db.Query("SELECT * FROM " + table)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var value interface{}
			err := rows.Scan(&value)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}
	}

	return values, nil
}

func main() {
	db, err := sql.Open("sqlite3", "./masterdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables := map[string]string{
		"users":    "users_table",
		"products": "products_table",
	}

	m := NewMasterDataManager(db, tables)

	// Set a record
	err = m.Set("users", "John Doe")
	if err != nil {
		log.Fatal(err)
	}

	// Get a record
	value, err := m.Get("users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)

	// List all records
	values, err := m.List()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(values)
}