package subscriptions

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	connection string
}

func NewRepository(connection string) IRepository {
	return &Repository{connection: connection}
}

func (r *Repository) GetAll() ([]*Subscription, error) {
	cnn := postgres.Open(r.connection)
	db, err := gorm.Open(cnn, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var subscriptions []*Subscription
	query := "SELECT Id, Email FROM subscription"
	db.Raw(query).Scan(&subscriptions)

	return subscriptions, nil
}

func (r *Repository) Add(email string) error {
	cnn := postgres.Open(r.connection)
	db, err := gorm.Open(cnn, &gorm.Config{})
	if err != nil {
		return err
	}

	insertQuery := "INSERT INTO subscription (Email) VALUES (?)"
	db.Exec(insertQuery, email)

	return nil
}

func (r *Repository) Exists(email string) (bool, error) {
	cnn := postgres.Open(r.connection)
	db, err := gorm.Open(cnn, &gorm.Config{})
	if err != nil {
		return false, err
	}

	var exists bool
	query := "SELECT 1 FROM subscription WHERE Email = ?"
	err = db.Raw(query, email).Scan(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}
