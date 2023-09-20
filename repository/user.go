package repository

import (
	"context"
	"github.com/restore/user/entity"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

func (u *User) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	result := u.db.Create(user)
	if result.Error != nil {
		return 0, result.Error
	}
	return user.ID, nil
}

func (u *User) CreateProfile(ctx context.Context, profile *entity.Profile) error {
	return u.db.Create(profile).Error
}

func (u *User) CreateStore(ctx context.Context, store *entity.Store) error {
	return u.db.Create(store).Error
}

func (u *User) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var result entity.User
	res := u.db.Where("email = ?", email).First(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return &result, nil
}

func (u *User) GetProfileByID(ctx context.Context, id int) (*entity.Profile, error) {
	result := entity.Profile{ID: id}
	res := u.db.First(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return &result, nil
}

func (u *User) GetStoreByID(ctx context.Context, id int) (*entity.Store, error) {
	result := entity.Store{ID: id}
	res := u.db.First(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return &result, nil
}

func (u *User) SearchStore(ctx context.Context, name string) ([]entity.Store, error) {
	var result []entity.Store
	res := u.db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (u *User) UpdateProfile(ctx context.Context, id int, profile *entity.Profile) error {
	result := entity.Profile{ID: id}
	res := u.db.First(&result)
	if res.Error != nil {
		return res.Error
	}

	result.Address = profile.Address
	result.City = profile.City
	result.Block = profile.Block
	result.State = profile.State
	result.ZipCode = profile.ZipCode
	result.Name = profile.Name

	res = u.db.Save(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
