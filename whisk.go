package main

// https://github.com/apache/openwhisk-client-go/blob/13fc65f65684e04f401fee67b231b370c53b3dcd/whisk/shared.go#L96-L101
type Limits struct {
	Timeout     int `yaml:"timeout,omitempty"`     // in seconds
	Memory      int `yaml:"memory,omitempty"`      // in MB
	Logs        int `yaml:"logs,omitempty"`        // in MB
	Concurrency int `yaml:"concurrency,omitempty"` // number of concurrent invocations allowed
}
