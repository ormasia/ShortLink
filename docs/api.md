# 短链接服务 API 文档

## 基础信息

- 基础URL: `http://localhost:8080`
- 所有需要认证的接口都需要在请求头中携带 `Authorization: Bearer <token>`
- 响应格式统一为 JSON
- 所有接口都支持跨域访问

## 通用响应格式

```json
{
    "code": 200,       // 状态码
    "message": "成功",  // 响应消息
    "data": {}         // 响应数据
}
```

## 用户相关接口

### 注册用户

- **URL**: `/api/v1/users`
- **方法**: `POST`
- **描述**: 注册新用户
- **请求体**:
```json
{
    "username": "string",  // 用户名
    "password": "string",  // 密码
    "email": "string",     // 邮箱（可选）
    "nickname": "string"   // 昵称（可选）
}
```
- **响应**:
  - 成功: `200 OK`
  - 用户名已存在: `400 Bad Request`
  - 参数错误: `400 Bad Request`

### 用户登录

- **URL**: `/api/v1/users/login`
- **方法**: `POST`
- **描述**: 用户登录并获取token
- **请求体**:
```json
{
    "username": "string",  // 用户名
    "password": "string"   // 密码
}
```
- **响应**:
```json
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "string",  // JWT token
        "user": {
            "id": 1,
            "username": "string",
            "nickname": "string",
            "email": "string"
        }
    }
}
```

### 用户登出

- **URL**: `/api/v1/users/logout`
- **方法**: `POST`
- **描述**: 用户登出，使当前token失效
- **认证**: 需要
- **请求体**:
```json
{
    "token": "string"  // 当前token
}
```
- **响应**:
  - 成功: `200 OK`
  - 未认证: `401 Unauthorized`

## 短链接相关接口

### 创建短链接

- **URL**: `/api/v1/links`
- **方法**: `POST`
- **描述**: 创建新的短链接
- **认证**: 需要
- **限流**: 基于IP，每秒100个请求
- **请求体**:
```json
{
    "original_url": "string"  // 原始URL
}
```
- **响应**:
```json
{
    "code": 200,
    "message": "创建成功",
    "data": {
        "shortlink": "string"  // 生成的短链接
    }
}
```

### 批量创建短链接

- **URL**: `/api/v1/links/batch`
- **方法**: `POST`
- **描述**: 批量创建多个短链接
- **认证**: 需要
- **限流**: 基于用户ID，每分钟10个请求
- **请求体**:
```json
{
    "original_urls": ["string"],  // 原始URL列表
    "concurrency": 10             // 并发数（可选，默认10，最大50）
}
```
- **响应**:
```json
{
    "code": 200,
    "message": "批量创建成功",
    "data": {
        "results": [
            {
                "original_url": "string",
                "short_url": "string",
                "error": "string"  // 如果创建失败，这里会有错误信息
            }
        ],
        "total_count": 10,
        "success_count": 8,
        "elapsed_time": 1.5
    }
}
```

### 获取热门短链接

- **URL**: `/api/v1/links/top`
- **方法**: `GET`
- **描述**: 获取点击量最高的短链接列表
- **认证**: 需要
- **响应**:
```json
{
    "code": 200,
    "message": "获取成功",
    "data": {
        "top": [
            {
                "short_url": "string",
                "clicks": 100
            }
        ]
    }
}
```

### 访问短链接

- **URL**: `/api/v1/links/:short_url`
- **方法**: `GET`
- **描述**: 访问短链接并重定向到原始URL
- **限流**: 基于IP，每秒100个请求
- **响应**:
  - 成功: `302 Found` 重定向到原始URL
  - 短链接无效: `404 Not Found`

## 错误码说明

- `200`: 成功
- `400`: 请求参数错误
- `401`: 未认证或认证失败
- `404`: 资源不存在
- `429`: 请求过于频繁
- `500`: 服务器内部错误

## 限流说明

1. 全局限流（基于IP）:
   - 算法: 令牌桶
   - 速率: 每秒100个请求
   - 突发容量: 100个请求

2. 批量创建限流（基于用户）:
   - 算法: 滑动窗口
   - 时间窗口: 1分钟
   - 最大请求数: 10个

## 注意事项

1. 所有需要认证的接口必须在请求头中携带有效的JWT token
2. 批量创建接口的并发数建议不会超过50
3. 短链接访问接口会自动记录点击量
4. 系统会自动处理重复的URL，返回已存在的短链接 