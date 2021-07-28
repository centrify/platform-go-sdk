package utils

const (
	productDir  = "/var/centrify/"
	prodDataDir = productDir + "cloud/"

	// LRPC service endpoints.
	// Note: should be in secure directory (own by root, not writable by world etc)
	endpointDMC = prodDataDir + "daemon2"
)
