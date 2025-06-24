package handler

import "time"

type SessionsMap map[string]*session

func NewSessionsMap() SessionsMap {
	return make(map[string]*session)
}

func (sessions SessionsMap) getSession(userName string) *session {

	session, exists := sessions[userName]
	if exists {
		return session
	}

	sessions[userName] = newSession(userName)
	return sessions[userName]
}

func (sessions SessionsMap) cleanOldSessions() {

	for {
		time.Sleep(1 * time.Hour)
		for userName, session := range sessions {
			if session.isOldSession() {
				delete(sessions, userName)
			}
		}
	}
}
