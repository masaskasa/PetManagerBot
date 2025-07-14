package handler

import "time"

type SessionsMap map[string]*Session

func NewSessionsMap() SessionsMap {
	return make(map[string]*Session)
}

func (sessions SessionsMap) GetSession(userName string) *Session {

	session, exists := sessions[userName]
	if exists {
		return session
	}

	sessions[userName] = newSession()
	return sessions[userName]
}

func (sessions SessionsMap) CleanOldSessions() {
	for {
		time.Sleep(1 * time.Hour)
		for userName, session := range sessions {
			if session.isOldSession() {
				delete(sessions, userName)
			}
		}
	}
}
