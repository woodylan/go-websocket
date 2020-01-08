package define

const (
	MESSAGE_TYPE_CLIENT int = iota + 1
	MESSAGE_TYPE_GROUP
)

const REDIS_KEY_GROUP = "go-websocket-group:"

//redis clientId前缀
const REDIS_CLIENT_ID_PREFIX = "wsClientId:"

//redis 客户端ID过期时间
const REDIS_KEY_SURVIVAL_SECONDS = 172800 //2天
