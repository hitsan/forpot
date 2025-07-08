package ssh

const (
	opSet    = "set"
	opDelete = "delete"
	opGet    = "get"
	opGetAll = "getAll"
	opClose  = "close"
)

type SessionStore struct {
	sessions map[int]*ForwardSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[int]*ForwardSession),
	}
}

func (s *SessionStore) set(port int, session *ForwardSession) {
	s.sessions[port] = session
}

func (s *SessionStore) delete(port int) {
	if session, ok := s.sessions[port]; ok {
		session.Close()
		delete(s.sessions, port)
	}
}

func (s *SessionStore) get(port int) *ForwardSession {
	return s.sessions[port]
}

func (s *SessionStore) getAll() map[int]struct{} {
	pm := make(map[int]struct{})
	for port := range s.sessions {
		pm[port] = struct{}{}
	}
	return pm
}

func (s *SessionStore) close() {
	for _, session := range s.sessions {
		session.Close()
	}
}