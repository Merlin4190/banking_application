package database

import (
	"errors"
)

// Define a custom type for database entries
type Entry map[string]interface{}

// Initialize the database map
var db map[string][]Entry

// Initialize the database in memory
func InitDatabase() {
	db = make(map[string][]Entry)
}

func CreateTable(table string, data Entry) bool {
	db[table] = append(db[table], data)
	return true
}

// AddEntry adds a new entry to the specified table
func AddEntry(table string, data Entry) error {
	if _, exists := db[table]; !exists {
		return errors.New("table does not exist")
	}
	db[table] = append(db[table], data)
	return nil
}

// UpdateEntry updates an existing entry in the specified table
func UpdateIndexEntry(table string, index int, data Entry) error {
	if _, exists := db[table]; !exists {
		return errors.New("table does not exist")
	}
	if index < 0 || index >= len(db[table]) {
		return errors.New("invalid index")
	}
	db[table][index] = data
	return nil
}

func UpdateEntry(table string, id string, data Entry) error {
	if _, exists := db[table]; !exists {
		return errors.New("table does not exist")
	}

	for i, entry := range db[table] {
		entryID, ok := entry["ID"].(string)
		if !ok {
			return errors.New("ID is not a string")
		}

		if entryID == id {
			db[table][i] = data
			return nil
		}
	}

	return errors.New("entry not found")
}


// DeleteEntry deletes an existing entry from the specified table
func DeleteEntry(table string, index int) error {
	if _, exists := db[table]; !exists {
		return errors.New("table does not exist")
	}
	if index < 0 || index >= len(db[table]) {
		return errors.New("invalid index")
	}
	db[table] = append(db[table][:index], db[table][index+1:]...)
	return nil
}

// GetEntries retrieves all entries from the specified table
func GetEntries(table string) ([]Entry, error) {
	if _, exists := db[table]; !exists {
		return nil, errors.New("table does not exist")
	}
	return db[table], nil
}

// GetEntry retrieves a single entry from the specified table at the specified index
func GetIndexEntry(table string, index int) (Entry, error) {
	if _, exists := db[table]; !exists {
		return nil, errors.New("table does not exist")
	}
	if index < 0 || index >= len(db[table]) {
		return nil, errors.New("invalid index")
	}
	return db[table][index], nil
}

func GetEntry(table string, id string) (Entry, error) {
	if _, exists := db[table]; !exists {
		return nil, errors.New("table does not exist")
	}

	for _, entry := range db[table] {
		entryID, ok := entry["ID"].(string)
		if !ok {
			return nil, errors.New("ID is not a string")
		}

		if entryID == id {
			return entry, nil
		}
	}

	return nil, errors.New("entry not found")
}
