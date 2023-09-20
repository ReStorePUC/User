package controller

import (
	"context"
	"errors"
	"github.com/restore/user/config"
	"github.com/restore/user/entity"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type repository interface {
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	CreateProfile(ctx context.Context, profile *entity.Profile) error
	CreateStore(ctx context.Context, store *entity.Store) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetProfileByID(ctx context.Context, id int) (*entity.Profile, error)
	GetStoreByID(ctx context.Context, id int) (*entity.Store, error)
	SearchStore(ctx context.Context, name string) ([]entity.Store, error)
	UpdateProfile(ctx context.Context, id int, profile *entity.Profile) error
}

type kong interface {
	CreateCustomer(email string) error
	CreateCredentials(email string) (string, error)
}

type User struct {
	repo repository
	kong kong
}

func NewUser(r repository, k kong) *User {
	return &User{
		repo: r,
		kong: k,
	}
}

func (u *User) Register(ctx context.Context, profile *entity.Profile) (string, error) {
	log := zap.NewNop()
	pass, err := crypt(profile.User.Password)
	if err != nil {
		log.Error(
			"error to crypt password",
			zap.Error(err),
		)
		return "", err
	}
	profile.User.Password = pass

	id, err := u.repo.CreateUser(ctx, &profile.User)
	if err != nil {
		log.Error(
			"error to register user",
			zap.Error(err),
		)
		return "", err
	}

	profile.UserID = id
	err = u.repo.CreateProfile(ctx, profile)
	if err != nil {
		log.Error(
			"error to create profile",
			zap.Error(err),
		)
		return "", err
	}

	err = u.kong.CreateCustomer(profile.User.Email)
	if err != nil {
		log.Error(
			"error creating kong consumer",
			zap.Error(err),
		)
		return "", err
	}

	jwt, err := u.kong.CreateCredentials(profile.User.Email)
	if err != nil {
		log.Error(
			"error creating kong credentials",
			zap.Error(err),
		)
		return "", err
	}

	return jwt, nil
}

func (u *User) RegisterStore(ctx context.Context, store *entity.Store) (string, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	result, err := u.repo.GetUserByEmail(ctx, admin.(string))
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return "", err
	}
	if !result.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return "", errors.New("unauthorized action")
	}

	pass, err := crypt(store.User.Password)
	if err != nil {
		log.Error(
			"error to crypt password",
			zap.Error(err),
		)
		return "", err
	}
	store.User.Password = pass

	id, err := u.repo.CreateUser(ctx, &store.User)
	if err != nil {
		log.Error(
			"error to register user",
			zap.Error(err),
		)
		return "", err
	}

	store.UserID = id
	err = u.repo.CreateStore(ctx, store)
	if err != nil {
		log.Error(
			"error to create store",
			zap.Error(err),
		)
		return "", err
	}

	err = u.kong.CreateCustomer(store.User.Email)
	if err != nil {
		log.Error(
			"error creating kong consumer",
			zap.Error(err),
		)
		return "", err
	}

	jwt, err := u.kong.CreateCredentials(store.User.Email)
	if err != nil {
		log.Error(
			"error creating kong credentials",
			zap.Error(err),
		)
		return "", err
	}

	return jwt, nil
}

func (u *User) Login(ctx context.Context, user *entity.User) (string, bool, error) {
	isAdmin, err := u.validate(ctx, user)
	if err != nil {
		return "", false, err
	}

	jwt, err := u.kong.CreateCredentials(user.Email)
	if err != nil {
		return "", false, err
	}

	return jwt, isAdmin, nil
}

func (u *User) GetProfile(ctx context.Context, id string) (*entity.Profile, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	result, err := u.repo.GetUserByEmail(ctx, admin.(string))
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return nil, err
	}
	if !result.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return nil, errors.New("unauthorized action")
	}

	profileID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return nil, err
	}

	profile, err := u.repo.GetProfileByID(ctx, profileID)
	if err != nil {
		log.Error(
			"error to get profile",
			zap.Error(err),
		)
		return nil, err
	}
	profile.User.Password = ""

	return profile, nil
}

func (u *User) GetStore(ctx context.Context, id string) (*entity.Store, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	result, err := u.repo.GetUserByEmail(ctx, admin.(string))
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return nil, err
	}
	if !result.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return nil, errors.New("unauthorized action")
	}

	storeID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return nil, err
	}

	store, err := u.repo.GetStoreByID(ctx, storeID)
	if err != nil {
		log.Error(
			"error to get store",
			zap.Error(err),
		)
		return nil, err
	}
	store.User.Password = ""

	return store, nil
}

func (u *User) SearchStore(ctx context.Context, name string) ([]entity.Store, error) {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	result, err := u.repo.GetUserByEmail(ctx, admin.(string))
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return nil, err
	}
	if !result.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return nil, errors.New("unauthorized action")
	}

	stores, err := u.repo.SearchStore(ctx, name)
	if err != nil {
		log.Error(
			"error to get store",
			zap.Error(err),
		)
		return nil, err
	}

	return stores, nil
}

func (u *User) UpdateProfile(ctx context.Context, id string, profile *entity.Profile) error {
	log := zap.NewNop()

	admin := ctx.Value(config.EmailHeader)
	result, err := u.repo.GetUserByEmail(ctx, admin.(string))
	if err != nil {
		log.Error(
			"error getting admin",
			zap.Error(err),
		)
		return err
	}
	if !result.IsAdmin {
		log.Error(
			"unauthorized action",
		)
		return errors.New("unauthorized action")
	}

	profileID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return err
	}

	err = u.repo.UpdateProfile(ctx, profileID, profile)
	if err != nil {
		log.Error(
			"error to get profile",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (u *User) validate(ctx context.Context, user *entity.User) (bool, error) {
	result, err := u.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		return false, err
	}

	return result.IsAdmin, nil
}

func crypt(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
