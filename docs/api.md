## 接口

#### 连接接口

**请求地址：**/ws?systemId=xxx

**协议：** websocket

**请求参数**：systemId 系统ID

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "clientId": "9fa54bdbbf2778cb"
  }
}
```

#### 注册系统

**请求地址：**/api/register

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```

#### 发送信息给指定客户端

**请求地址：**/api/send_to_client

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| clientId | string | 是       | 客户端ID |
| sendUserId | string | 是       | 发送者ID |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "messageId": "5b4646dd8328f4b1"
    }
}
```

#### 批量发送信息给指定客户端

**请求地址：**/api/send_to_clients

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| clientIds | array | 是       | 客户端ID列表 |
| sendUserId | string | 是       | 发送者ID |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "messageId": "5b4646dd8328f4b1"
    }
}
```

#### 绑定客户端到分组

**请求地址：**/api/bind_to_group

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| sendUserId | string | 是       | 发送者ID |
| clientId | string | 是       | 客户端ID |
| groupName | string | 是       | 分组名 |

**响应示例：**

```json
{
  "code": 0,
  "msg": "success",
  "data": []
}
```

#### 发送信息给指定分组

**请求地址：**/api/send_to_group

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| sendUserId | string | 是       | 发送者ID |
| groupName | string | 是       | 分组名 |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "messageId": "5b4646dd8328f4b1"
    }
}
```

#### 获取在线的客户端列表

**请求地址：**/api/get_online_list

**请求方式：** POST

**Content-Type：** application/json; charset=UTF-8

**请求头Header**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| systemId | string | 是       | 系统ID |

**请求头Body**

| 字段     | 类型   | 是否必须 | 说明     |
| -------- | ------ | -------- | -------- |
| groupName | string | 是       | 分组名 |
| code | integer | 是       | 自定义的状态码 |
| msg | string | 是       | 自定义的状态消息 |
| data | sring、array、object | 是       | 消息内容 |

**响应示例：**

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "count": 2,
        "list": [
            "WQReWw6m+wct+eKk/2rDiWcU4maU8JRTRZEX8c7Te6LzCa//VCXr/0KeVyO0sdNt",
            "j6YdsGFH4rfbYN/vS6UavJ5fVclWIB9W+Gqg9R/92cLJqgAp2ZPkvMbQiwQBJmDc"
        ]
    }
}
```