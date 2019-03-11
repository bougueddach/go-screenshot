package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	"go-screenshot/chrome"
)

type Storage interface {
	Open() error
	SetHTTPData(data *chrome.HTTResponse)
}

// Storage handles the pointer to a buntdb instance
type FileStorage struct {
	Db   *buntdb.DB
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

// Open creates a new connection to a buntdb database
func (storage *FileStorage) Open() error {

	log.WithField("database-location", storage.path).Debug("Opening buntdb")

	db, err := buntdb.Open(storage.path)
	if err != nil {
		return err
	}

	// build some indexes
	db.CreateIndex("url", "*", buntdb.IndexJSON("url"))

	storage.Db = db

	return nil
}

// Close closes the connection to a buntdb connection
func (storage *FileStorage) Close() {

	log.Debug("Closing buntdb")
	storage.Db.Close()
}

// SetHTTPData stores HTTP information about a URL
func (storage *FileStorage) SetHTTPData(data *chrome.HTTResponse) {

	// marshal the data
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.WithField("err", err).Fatal("Error marshalling the HTTP response data to JSON")
	}

	// generate a key to use
	key := sha1.New()
	key.Write([]byte(data.URL))
	keyBytes := key.Sum(nil)
	keyString := hex.EncodeToString(keyBytes)
	log.WithFields(log.Fields{"url": data.URL, "key": keyString}).Debug("Calculated key for storage")

	// add the document
	err = storage.Db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(keyString, string(jsonData), nil)

		return err
	})

	if err != nil {
		log.WithField("err", err).Fatal("Error saving HTTP response data")
	}
}
