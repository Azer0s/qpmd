package qpmd

type ResponseType string

const RESPONSE_OK ResponseType = "ok"
const RESPONSE_TIMEOUT ResponseType = "timeout"
const RESPONSE_ERROR ResponseType = "error"

const TIMESTAMP = "timestamp"

type Response struct {
	ResponseType ResponseType
	Data         map[string]interface{}
}
