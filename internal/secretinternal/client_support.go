package secretinternal

// SetDebug enables/disables debug mode
func (c *APIClient) SetDebug(value bool) {
	c.cfg.setDebug(value)
}

// AddDefaultHeaders sets the default HTTP header
func (c *APIClient) AddDefaultHeaders(hdrs map[string]string) {
	if hdrs == nil {
		return
	}
	for key, value := range hdrs {
		c.cfg.AddDefaultHeader(key, value)
	}
}

// SetUserAgent sets the user agent
func (c *APIClient) SetUserAgent(agent string) {
	c.cfg.setUserAgent(agent)
}

// setDebug enables/disables debug mode
func (c *Configuration) setDebug(value bool) {
	c.Debug = value
}

// setUserAgent sets UserAgent in HTTP header
func (c *Configuration) setUserAgent(agent string) {
	c.UserAgent = agent
}
