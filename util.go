package main

import (
	"encoding/json"
	"net"
	"time"
)

func readRequest(client net.Conn) (Request, error) {
	c := client.RemoteAddr().String()

	buf := make([]byte, 4096)
	n, err := client.Read(buf)
	if n == 0 || err != nil {
		errLog.Printf("Failed to read from %s", c)
		return Request{}, err
	}
	stdLog.Printf("Read %d bytes from %s", n, c)

	req := Request{}
	err = json.Unmarshal(buf[:n], &req)
	if err != nil {
		errLog.Printf("Couldn't parse JSON request from %s", c)
		return req, err
	}

	return req, nil
}

func writeResponse(client net.Conn, response Response) error {
	response.Data[TIMESTAMP] = time.Now().Unix()

	b, err := json.Marshal(response)

	if err != nil {
		return err
	}

	n, err := client.Write(b)

	if err != nil {
		return err
	}

	stdLog.Printf("Wrote %d bytes to %s", n, client.RemoteAddr().String())

	return nil
}

func writeError(client net.Conn, err error) error {
	return writeResponse(client, Response{
		ResponseType: RESPONSE_ERROR,
		Data: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func writeTimeout(client net.Conn) error {
	return writeResponse(client, Response{
		ResponseType: RESPONSE_TIMEOUT,
		Data: map[string]interface{}{},
	})
}

func writeOk(client net.Conn, data map[string]interface{}) error {
	return writeResponse(client, Response{
		ResponseType: RESPONSE_OK,
		Data: data,
	})
}
