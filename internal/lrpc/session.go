package lrpc

// SessionCtxtBase manages application-dependent session information that is shared
// by command handlers processing messages in the same LRPC session.
//
// Examples of such information are:
//
// - HTTP client handle
//
// - session context/state information
type SessionCtxtBase struct {
	prop map[string]interface{}
}

// Get returns the object that was previous saved in the session context
func (s *SessionCtxtBase) Get(key string) interface{} {
	if s.prop == nil {
		return nil
	}
	return s.prop[key]
}

// Set stores an object in the session context so that it can be retrieved later by Get()
func (s *SessionCtxtBase) Set(key string, val interface{}) {
	if s.prop == nil {
		s.prop = make(map[string]interface{})
	}
	s.prop[key] = val
}
