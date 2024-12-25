package store

import (
	"fmt"
	"os"

	"github.com/marcos-brito/booklist/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Hate the global, but couldn't find another way.
// It gets initialized in cmd/graph. The graphql
// api should be the only service using it.
var DB *gorm.DB

type DSN struct {
	host     string
	user     string
	dbname   string
	password string
	port     string
	sslmode  bool
}

func (dsn *DSN) String() string {
	var sslmode string

	if dsn.sslmode {
		sslmode = "enable"
	} else {
		sslmode = "disable"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dsn.host, dsn.user, dsn.password, dsn.dbname, dsn.port, sslmode)
}

func With(db *gorm.DB) {
	DB = db
}

func Migrate(db *gorm.DB) error {
    err := db.AutoMigrate(&models.Book{}, &models.Author{}, &models.Publisher{}, &models.Profile{},
		&models.Settings{}, &models.List{}, &models.CollectionItem{})

    if err != nil {
        return err
    }

    return nil
}

func NewConnection() (*gorm.DB, error) {
	dsn := DSN{
		host:     os.Getenv("POSTGRES_HOST"),
		user:     os.Getenv("POSTGRES_USER"),
		dbname:   os.Getenv("POSTGRES_DB"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		port:     os.Getenv("POSTGRES_PORT"),
	}

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %s", err)
	}

	return db, nil
}
