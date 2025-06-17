package handler

import (
	"errors"
	"time"
)

type session struct {
	scenario     scenario
	state        dialogState
	tempObjects  map[string]interface{}
	lastActivity time.Time
}

var (
	ErrObjectAlreadyExists = errors.New("can't save object: it already exist")
	ErrObjectNotExists     = errors.New("can't get object: it not exist")
)

func newSession(user string) *session {
	return &session{
		scenario:     none,
		state:        ready,
		tempObjects:  make(map[string]interface{}),
		lastActivity: time.Now(),
	}
}

func (session *session) setScenario(scenario scenario) {
	session.scenario = scenario
}

func (session *session) setState(state dialogState) {
	session.state = state
	session.updateLastActivity()
}

func (session *session) updateLastActivity() {
	session.lastActivity = time.Now()
}

func (session *session) saveObject(key string, tempObject interface{}) error {
	_, exists := session.tempObjects[key]
	if exists {
		return ErrObjectAlreadyExists
	}

	session.tempObjects[key] = tempObject
	return nil
}

func (session *session) updateObject(key string, tempObject interface{}) {
	session.tempObjects[key] = tempObject
}

func (session *session) getObject(key string) (interface{}, error) {

	tempObject, exists := session.tempObjects[key]
	if !exists {
		return nil, ErrObjectNotExists
	}

	return tempObject, nil
}

func (session *session) resetTempObjects() {
	session.tempObjects = make(map[string]interface{})
}

func (session *session) reset() {
	session.setScenario(none)
	session.setState(ready)
	session.resetTempObjects()
}
