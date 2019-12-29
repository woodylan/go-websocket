package define

const (
	MESSAGE_TYPE_CLIENT int = iota + 1
	MESSAGE_TYPE_GROUP
)

const REDIS_KEY_GROUP = "go-websocket-group"
