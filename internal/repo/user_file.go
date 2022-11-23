package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"sync"
	"users/internal/entity"
)

type userDAO struct {
	Increment int      `json:"increment"`
	List      userList `json:"list"`
}

type userList map[string]entity.User

// User implements usecases.UserRepo interface
type User struct {
	f     os.File
	cache userList
	mu    sync.RWMutex
	inc   int
}

// New is a constructor for User
func New(f os.File) (*User, error) {
	user := &User{
		f:     f,
		cache: map[string]entity.User{},
		mu:    sync.RWMutex{},
	}
	err := user.getCache()
	if err != nil {
		return nil, err
	}
	max := 0
	for k := range user.cache {
		v, _ := strconv.Atoi(k)
		if v > max {
			max = v
		}
	}
	user.inc = max
	return user, nil
}

func (r *User) getCache() error {
	return json.NewDecoder(&r.f).Decode(&r.cache)
}

// GetByID return entity.User by its id, entity.ErrNoID if there is no one
func (r *User) GetByID(ctx context.Context, id int) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	k := strconv.Itoa(id)
	user, ok := r.cache[k]
	if !ok {
		return entity.User{}, entity.ErrNoID
	}
	return user, nil
}

// Create make new user, return its id
func (r *User) Create(ctx context.Context, user entity.User) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.inc++
	r.cache[strconv.Itoa(r.inc)] = user
	data := userDAO{Increment: r.inc, List: r.cache}
	buf, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("repo - Create: %w", err)
	}
	err = os.WriteFile(r.f.Name(), buf, fs.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("repo - Create: %w", err)
	}
	return r.inc, nil
}

// Update change user info
func (r *User) Update(ctx context.Context, user entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[strconv.Itoa(user.ID)] = user
	data := userDAO{Increment: r.inc, List: r.cache}
	buf, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("repo - Update: %w", err)
	}
	err = os.WriteFile(r.f.Name(), buf, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("repo - Update: %w", err)
	}
	return nil
}

// Delete remove user from repository
func (r *User) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cache, strconv.Itoa(id))
	data := userDAO{Increment: r.inc, List: r.cache}
	buf, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("repo - Delete: %w", err)
	}
	err = os.WriteFile(r.f.Name(), buf, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("repo - Delete: %w", err)
	}
	return nil
}
