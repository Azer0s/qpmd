package main

import (
	"errors"
	"fmt"
	"sync"
)

var systems = make(map[string]system)
var systemsMu = &sync.RWMutex{}

func addSystem(systemName string, port int, attributes map[string]interface{}) {
	systemsMu.Lock()
	defer systemsMu.Unlock()

	systems[systemName] = system{
		name:       systemName,
		port:       port,
		attributes: attributes,
	}
}

func removeSystem(systemName string) {
	systemsMu.Lock()
	defer systemsMu.Unlock()

	delete(systems, systemName)
}

func getSystem(systemName string) (system, error) {
	systemsMu.RLock()
	defer systemsMu.RUnlock()

	s, ok := systems[systemName]

	if ok {
		return s, nil
	}

	return system{}, errors.New(fmt.Sprintf("Couldn't find system %s in systems map", systemName))
}

type system struct {
	name string
	port int
	attributes map[string]interface{}
}
