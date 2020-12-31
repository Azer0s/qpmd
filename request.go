package qpmd

type RequestType string

const REQUEST_REGISTER RequestType = "register"
const REQUEST_LOOKUP RequestType = "lookup"
const REQUEST_HELLO RequestType = "hello"
const HEARTBEAT RequestType = "heartbeat"

const SYSTEM_NAME = "system"
const PORT = "port"
const MACHINE = "machine"
const MACHINE_ID = "machine_id"
const MESSAGE_GATEWAY_PORT = "message_gateway_port"
const GP_GATEWAY_PORT = "gp_gateway_port"

type Request struct {
	RequestType RequestType
	Data        map[string]interface{}
}
