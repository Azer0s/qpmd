package qpmd

type RequestType string

const REQUEST_REGISTER RequestType = "register"
const REQUEST_LOOKUP RequestType = "lookup"
const REQUEST_HELLO RequestType = "hello"
const HEARTBEAT RequestType = "heartbeat"

const SYSTEM_NAME = "system"
const PORT = "port"
const MACHINE_ID = "machine_id"

type Request struct {
	RequestType RequestType
	Data        map[string]interface{}
}
