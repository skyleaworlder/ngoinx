# Roadmap

## Target

### 1. 负载均衡

根据某一路由下对接的 `target` 数量区分均衡算法：

* 1~4：加权轮询法
* 4+：一致性哈希

### 2. 反向代理

通过维护映射表转发请求。

### 3. 动静分离

通过 `getContent` 获取返回：

1. 维护 `ngoinx cache`；
2. 对接静态资源服务器；