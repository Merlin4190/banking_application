package database

type Repository interface {
	CreateTable(table string, data Entry) bool

	AddEntry(table string, data Entry)

	UpdateEntry(table string, id string, data Entry) error

	DeleteEntry(table string, id string) error

	GetEntries(table string) ([]Entry, error)

	GetEntry(table string, id string) (Entry, error)
}
