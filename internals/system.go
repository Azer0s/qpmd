package internals

import (
	"errors"
	"fmt"
	"sync"
)

var systems = map[string]system{}
var systemsMu = &sync.RWMutex{}

var machines = map[string]machine{}
var machinesMu = &sync.RWMutex{}

func addMachine(machineId string, gatewayPort, gpPort uint16) {
	machinesMu.Lock()
	defer machinesMu.Unlock()

	machines[machineId] = machine{
		machineId:   machineId,
		gatewayPort: gatewayPort,
		gpPort:      gpPort,
	}
}

func getMachine(machineId string) (machine, error) {
	machinesMu.RLock()
	defer machinesMu.RUnlock()

	m, ok := machines[machineId]

	if ok {
		return m, nil
	}

	return machine{}, errors.New(fmt.Sprintf("Couldn't find machine %s in machines map!", machineId))
}

func addSystem(systemName string, port uint16, machineId string) {
	systemsMu.Lock()
	defer systemsMu.Unlock()

	systems[systemName] = system{
		name:      systemName,
		port:      port,
		machineId: machineId,
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
	name      string
	port      uint16
	machineId string
}
