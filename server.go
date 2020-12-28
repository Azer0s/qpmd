package main

import (
	"net"
	"time"
)

func handleClient(client net.Conn) {
	c := client.RemoteAddr().String()

	stdLog.Printf("Handling incoming connection from %s", c)

	defer func() {
		if r := recover(); r != nil {
			errLog.Print(r)
		}

		stdLog.Printf("Closing client %s", c)
		err := client.Close()
		if err != nil {
			errLog.Printf("Failed to close client %s", c)
			return
		}
		stdLog.Printf("Closed client %s", c)
	}()

	req, err := readRequest(client)

	if err != nil {
		return
	}

	handleRequest(client, req)
}

func handleRequest(client net.Conn, req Request) {
	switch req.RequestType {
	case REQUEST_REGISTER:
		handleRegister(client, req.Data[SYSTEM_NAME].(string), int(req.Data[PORT].(float64)), req.Data[ATTRIBUTES].(map[string]interface{}))
	case REQUEST_LOOKUP:
		handleLookup(client, req.Data[SYSTEM_NAME].(string))
	}
}

func handleLookup(client net.Conn, systemName string) {
	c := client.RemoteAddr().String()

	stdLog.Printf("Handling lookup request from client %s for system %s", c, systemName)

	s, err := getSystem(systemName)

	if err != nil {
		errLog.Printf("Error getting system %s, %s", systemName, err.Error())
		err = writeError(client, err)

		if err != nil {
			errLog.Printf("Error while sending error message to client %s, %s", c, err.Error())
		}
		return
	}

	err = writeOk(client, map[string]interface{}{
		SYSTEM_NAME: s.name,
		PORT: s.port,
		ATTRIBUTES: s.attributes,
	})

	if err != nil {
		errLog.Printf("Error while sending okay message to client %s, %s", c, err.Error())
	}
}

func handleRegister(client net.Conn, systemName string, port int, attributes map[string]interface{}) {
	c := client.RemoteAddr().String()

	stdLog.Printf("Handling registration request from client %s for system %s", c, systemName)

	stdLog.Printf("Adding mapping from system %s to port %d", systemName, port)

	addSystem(systemName, port, attributes)

	defer func() {
		stdLog.Printf("Removing system mapping from system %s to port %d", systemName, port)
		removeSystem(systemName)
	}()

	errChan := make(chan error)
	reqChan := make(chan Request)

	go func() {
		for {
			req, err := readRequest(client)

			if err != nil {
				errChan <- err
				return
			}

			reqChan <- req
		}
	}()

	for {
		timeoutChan := time.After(timeout)

		select {
		case req := <- reqChan:
			if req.RequestType != HEARTBEAT {
				errLog.Printf("Expected a heartbeat from %s, got a request of type %s", c, req.RequestType)
				return
			}
			stdLog.Printf("Got heartbeat from %s", c)

			err := writeOk(client, map[string]interface{}{})

			if err != nil {
				errLog.Printf("Error while sending okay message to client %s, %s", c, err.Error())
				return
			}
		case err := <- errChan:
			err = writeError(client, err)

			if err != nil {
				errLog.Printf("Error while sending error message to client %s, %s", c, err.Error())
			}

			return
		case <- timeoutChan:
			errLog.Printf("Didn't receive a heartbeat after 30 seconds from %s", c)
			err := writeTimeout(client)

			if err != nil {
				errLog.Printf("Error while sending timeout message to client %s, %s", c, err.Error())
			}

			return
		}
	}
}

