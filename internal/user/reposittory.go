package user

import "go-advanced/pkg/db"

type UserRepository struct {
	database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	res := repo.database.DB.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
	res := repo.database.DB.First(&user, "email = ?", email)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
