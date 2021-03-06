package internals

import (
	"github.com/Azer0s/qpmd"
	"github.com/vmihailenco/msgpack/v5"
	"net"
	"time"
)

func readRequest(client net.Conn) (qpmd.Request, error) {
	c := client.RemoteAddr().String()

	buf := make([]byte, 4096)
	n, err := client.Read(buf)
	if n == 0 || err != nil {
		errLog.Printf("Failed to read from %s", c)
		return qpmd.Request{}, err
	}
	stdLog.Printf("Read %d bytes from %s", n, c)

	req := qpmd.Request{}
	err = msgpack.Unmarshal(buf[:n], &req)
	if err != nil {
		errLog.Printf("Couldn't parse JSON request from %s", c)
		return req, err
	}

	return req, nil
}

func writeResponse(client net.Conn, response qpmd.Response) error {
	response.Data[qpmd.TIMESTAMP] = time.Now().Unix()

	b, err := msgpack.Marshal(response)

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
	return writeResponse(client, qpmd.Response{
		ResponseType: qpmd.RESPONSE_ERROR,
		Data: map[string]interface{}{
			"error": err.Error(),
		},
	})
}

func writeTimeout(client net.Conn) error {
	return writeResponse(client, qpmd.Response{
		ResponseType: qpmd.RESPONSE_TIMEOUT,
		Data:         map[string]interface{}{},
	})
}

func writeOk(client net.Conn, data map[string]interface{}) error {
	return writeResponse(client, qpmd.Response{
		ResponseType: qpmd.RESPONSE_OK,
		Data:         data,
	})
}
