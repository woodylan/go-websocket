package define

const (
	//每组客户端数量限制
	GROUP_CLIENT_LIMIT = 500;

	//redis 分组里的客户端列表key
	REDIS_KEY_GROUP = "ws-group-client-list:"

	//redis clientId前缀
	REDIS_CLIENT_ID_PREFIX = "ws-client-id:"

	//redis 客户端ID过期时间
	REDIS_KEY_SURVIVAL_SECONDS = 172800 //2天

	//redis 分组列表key
	REDIS_KEY_GROUP_LIST = "ws-group-list"

	//账号列表key
	REDIS_KEY_ACCOUNT_LIST = "ws-account-list"

	//服务器列表
	REDIS_KEY_SERVER_LIST = "ws-server-list"

	//账号前缀
	REDIS_PREFIX_ACCOUNT_INFO = "ws-account-info:"

	//token前缀
	REDIS_PREFIX_TOKEN = "ws-token:"

	//RPC消息类型
	RPC_MESSAGE_TYPE_GROUP  = 1; //分组消息
	RPC_MESSAGE_TYPE_SYSTEM = 2; //系统消息
)
