package main

import (
	"errors"
	"sync"

	ct "github.com/flynn/flynn-controller/types"
)

type AppRepo struct {
	appNames map[string]*ct.App
	appIDs   map[string]*ct.App
	apps     []*ct.App
	mtx      sync.RWMutex
}

func NewAppRepo() *AppRepo {
	return &AppRepo{
		appNames: make(map[string]*ct.App),
		appIDs:   make(map[string]*ct.App),
	}
}

// - validate
// - set id
// - check name doesn't exist
// - persist
func (r *AppRepo) Add(data interface{}) error {
	app := data.(*ct.App)
	// TODO: actually validate
	if app.Name == "" {
		return errors.New("controller: app name must not be blank")
	}
	app.ID = uuid()
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, exists := r.appNames[app.Name]; exists {
		return errors.New("controller: app name is already in use")
	}

	r.appNames[app.Name] = app
	r.appIDs[app.ID] = app
	r.apps = append(r.apps, app)

	return nil
}

var ErrNotFound = errors.New("controller: resource not found")

func (r *AppRepo) Get(id string) (interface{}, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	app := r.appIDs[id]
	if app == nil {
		return nil, ErrNotFound
	}
	return app, nil
}
