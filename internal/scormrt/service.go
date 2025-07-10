package scormrt

import (
	"sync"
)

// RuntimeService holds runtime session data.
type RuntimeService struct {
	mu        sync.RWMutex
	sessions  map[string]map[string]string
	lastError map[string]string
}

// NewService creates a new RuntimeService.
func NewService() *RuntimeService {
	return &RuntimeService{
		sessions:  make(map[string]map[string]string),
		lastError: make(map[string]string),
	}
}

var defaultService = NewService()

// Initialize starts a new session for the given id.
func (s *RuntimeService) Initialize(session string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[session]; !ok {
		s.sessions[session] = make(map[string]string)
	}
	s.lastError[session] = "0"
	return "true"
}

// Terminate ends an existing session.
func (s *RuntimeService) Terminate(session string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, session)
	s.lastError[session] = "0"
	return "true"
}

// GetValue retrieves a value for an element.
func (s *RuntimeService) GetValue(session, element string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if vals, ok := s.sessions[session]; ok {
		if v, ok := vals[element]; ok {
			s.lastError[session] = "0"
			return v
		}
	}
	s.lastError[session] = "101" // general error
	return ""
}

// SetValue stores a value for an element.
func (s *RuntimeService) SetValue(session, element, value string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[session]; !ok {
		s.sessions[session] = make(map[string]string)
	}
	s.sessions[session][element] = value
	s.lastError[session] = "0"
	return "true"
}

// Commit would persist the session values. Here we simply acknowledge.
func (s *RuntimeService) Commit(session string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[session]; ok {
		s.lastError[session] = "0"
		return "true"
	}
	s.lastError[session] = "101"
	return "false"
}

// GetLastError returns the last error code for a session.
func (s *RuntimeService) GetLastError(session string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if code, ok := s.lastError[session]; ok {
		return code
	}
	return "0"
}

// GetErrorString maps error codes to messages.
func (s *RuntimeService) GetErrorString(code string) string {
	switch code {
	case "0":
		return "No error"
	case "101":
		return "General exception"
	default:
		return "Unknown error"
	}
}

// GetDiagnostic returns diagnostic info for the code.
func (s *RuntimeService) GetDiagnostic(code string) string {
	return s.GetErrorString(code)
}

// exported helper functions using the default service
func Initialize(session string) string { return defaultService.Initialize(session) }
func Terminate(session string) string  { return defaultService.Terminate(session) }
func GetValue(session, element string) string {
	return defaultService.GetValue(session, element)
}
func SetValue(session, element, value string) string {
	return defaultService.SetValue(session, element, value)
}
func Commit(session string) string       { return defaultService.Commit(session) }
func GetLastError(session string) string { return defaultService.GetLastError(session) }
func GetErrorString(code string) string  { return defaultService.GetErrorString(code) }
func GetDiagnostic(code string) string   { return defaultService.GetDiagnostic(code) }
