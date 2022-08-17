<h2>1. 定义</h2>

全局服务节点（**Node GlobalService**）用于维护全服化共享数据。例如，中国大区的所有服务器共享同一套逻辑时间数据。

全局服务组启动和管理全局服务节点（**Node GlobalService**），并在全服化数据更新后通告所有服务器组。

<h2>2. 实现</h2>

<h3>2.1 全局服务节点</h3>

实现于`node/node_globalservice`，用于启动全服化系统功能。启动全局服务节点需要提供配置文件作为参数，配置文件模版位于`node/node_globalservice/config.toml.template`。

功能如下：

- 逻辑时间系统功能（主模式）。

<h3>2.2 生命周期</h3>

全局服务组必须先于其他所有服务器组启动，这样才能确保读取到正确的全服化共享数据。必须后于其他所有服务器组关闭，这样才能确保全服化数据的更新不会影响其他服务器的逻辑功能。

<h3>2.3 节点结构</h3>

全局服务组内涉及到3个服务节点：

- center： 用于管理全局服务组。
- globalservice：全局服务结点。
- harbor：在全局数据更新后向其他服务器组通告。

<h2>3. 使用</h2>

启动全局服务组:

```sh
// 假设当前处于 node 目录下
go run ../center/center.go ip:port
go run ../harbor/harbor.go centerip:centerport logicaddr ip:port
go run node_globalservice/main.go config.toml
```

<h2>4. Todo List</h2>

- [ ] 实现 harbor 向所有一直其他 harbor 转发消息的功能。
- [ ] 实现调整逻辑时间偏移的web接口。

