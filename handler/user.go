package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/restore/user/config"
	"github.com/restore/user/entity"
	"net/http"
)

type controller interface {
	Register(ctx context.Context, profile *entity.Profile) (string, error)
	RegisterStore(ctx context.Context, store *entity.Store) (string, error)
	Login(ctx context.Context, user *entity.User) (string, bool, error)
	GetProfile(ctx context.Context, id string) (*entity.Profile, error)
	GetStore(ctx context.Context, id string) (*entity.Store, error)
	SearchStore(ctx context.Context, name string) ([]entity.Store, error)
	UpdateProfile(ctx context.Context, id string, profile *entity.Profile) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetSelfProfile(ctx context.Context) (*entity.Profile, error)
	GetSelfStore(ctx context.Context) (*entity.Store, error)
}

type User struct {
	controller controller
}

func NewUser(c controller) *User {
	return &User{
		controller: c,
	}
}

// Register creates a new User.
func (u *User) Register(c *gin.Context) {
	var profile entity.Profile
	if err := c.BindJSON(&profile); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	result, err := u.controller.Register(c.Request.Context(), &profile)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		JWT string
	}{
		result,
	})
}

// RegisterStore creates a new Store.
func (u *User) RegisterStore(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	var store entity.Store
	if err := c.BindJSON(&store); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	result, err := u.controller.RegisterStore(ctx, &store)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		JWT string
	}{
		result,
	})
}

// Login creates a new User session.
func (u *User) Login(c *gin.Context) {
	var user entity.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	result, admin, err := u.controller.Login(c.Request.Context(), &user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		JWT     string
		IsAdmin bool
	}{
		result,
		admin,
	})
}

// GetProfile finds a Profile.
func (u *User) GetProfile(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	result, err := u.controller.GetProfile(ctx, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

// GetStore finds a Store.
func (u *User) GetStore(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	result, err := u.controller.GetStore(ctx, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

// SearchStore search a Store.
func (u *User) SearchStore(c *gin.Context) {
	name := c.Param("name")
	result, err := u.controller.SearchStore(c.Request.Context(), name)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

// SearchAdminStore search a Store.
func (u *User) SearchAdminStore(c *gin.Context) {
	name := c.Query("name")
	result, err := u.controller.SearchStore(c.Request.Context(), name)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}

// UpdateProfile updates a Profile.
func (u *User) UpdateProfile(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	var profile entity.Profile
	if err := c.BindJSON(&profile); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	err := u.controller.UpdateProfile(ctx, id, &profile)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, profile)
}

// GetSelfProfile finds the user Profile.
func (u *User) GetSelfProfile(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	result, err := u.controller.GetSelfProfile(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// GetSelfStore finds the user Store.
func (u *User) GetSelfStore(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	result, err := u.controller.GetSelfStore(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}
