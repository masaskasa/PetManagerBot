package handler

import (
	"errors"
	"time"
)

type Session struct {
	scenario     scenario
	state        dialogState
	tempObjects  map[string]interface{}
	lastActivity time.Time
}

var (
	ErrObjectAlreadyExists = errors.New("can't save object: it already exist")
	ErrObjectNotExists     = errors.New("can't get object: it not exist")
)

func newSession() *Session {
	return &Session{
		scenario:     none,
		state:        ready,
		tempObjects:  make(map[string]interface{}),
		lastActivity: time.Now(),
	}
}

func (session *Session) setScenario(scenario scenario) {
	session.scenario = scenario
}

func (session *Session) setState(state dialogState) {
	session.state = state
	session.updateLastActivity()
}

func (session *Session) updateLastActivity() {
	session.lastActivity = time.Now()
}

func (session *Session) saveObject(key string, tempObject interface{}) error {
	_, exists := session.tempObjects[key]
	if exists {
		return ErrObjectAlreadyExists
	}

	session.tempObjects[key] = tempObject
	return nil
}

func (session *Session) UpdateObject(key string, tempObject interface{}) {
	session.tempObjects[key] = tempObject
}

func (session *Session) GetObject(key string) (interface{}, error) {

	tempObject, exists := session.tempObjects[key]
	if !exists {
		return nil, ErrObjectNotExists
	}

	return tempObject, nil
}

func (session *Session) deleteTempObjects(keys ...string) {

	for _, key := range keys {
		delete(session.tempObjects, key)
	}
}

func (session *Session) resetTempObjects() {
	session.tempObjects = make(map[string]interface{})
}

func (session *Session) reset() {
	session.setScenario(none)
	session.setState(ready)
	session.resetTempObjects()
}

func (session *Session) isOldSession() bool {

	if time.Since(session.lastActivity) > 24*time.Hour {
		return true
	}

	return false
}
