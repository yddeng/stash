### new 架构


### login

对客户端进行登陆验证，分配gate服

### gate

主要功能 承载客户端流量，分发到对应内部服务器

缓存各服务器的对象，各服务器能使用的常用数据（玩家外观）。

gate将玩家的心跳包转发到各内部服务

### game

处理玩家业务逻辑

### 地图（world，map）

处理玩家视野，地图对象

### 队伍 （team）

处理玩家队伍


## 玩家上线

gate - game， rpc请求保证对象创建成功

- 用户登陆 设置数据库标记成功。到game上创建业务逻辑对象，创建成功返回到gate。gate返回用户成功，创建与game对象关系。
- 用户登陆竞争 在gate端直接能判断是否已经在游戏中。向gate发起踢人请求 
    1）gate正在异步创建game对象，这时设置下线标记，待响应时处理。
    2）正常游戏中，直接执行下线流程

## 玩家下线

gate - 各服务器 post请求，不需要保证能摧毁对象。各服务器有对象心跳。
gate与玩家保持心跳，超时切换到短线状态。 短线超时后 gate特销毁对象，并发送销毁协议到各服务。
各服务器未收到gate消息（心跳），超时销毁对象，并通告给gate。



### 规则

创建对象的服务器 需向数据库添加标记，保证唯一性。

gate -> 各服务器
创建对象，消息投递 为rpc。保证成功


各服务器 -> gate
销毁对象，消息投递 为post。不保证成功


### 宕机

1. gate宕机，各服务器运行
    1）gate对象销毁，各服务器不能收到来自gate的消息，对象心跳超时销毁。不能在其他gate上登陆。
    2）gate 重启，玩家登陆，根据数据库标记，将各服务器对象剔除。重登陆流程
    
2. 各服务器宕机，gate运行
    同步消息到各服务器，发现对象不存在。由gate执行创建对象流程。若创建失败按照对象的重要程度选择是否踢玩家下线。
    
### 创建对象

创建对象的服务器需向数据库添加标记

登陆时，数据库标记需全部为空。保证对象在整个游戏集群只在一处修改。

gate标记必定为空。如不为空表示玩家在其他gate登陆。到对应gate踢人，待标记为空时才能登陆。
    
内部服务器创建 -> gate创建 -> 完成
    
### 销毁对象

玩家下线，由gate发起。向各服务器发送下线消息，待各标记清理完成，gate清理标记，删除对象。

若内部服务器宕机，导致下线消息不能处理，这时标记不能清理。 gate对象也不能清理 ？？？？？？？ 需要整个游戏进程都稳定。

清理内部服务器标记 -> gate标记 -> 完成

### 踢人

由gate发起，向各服务器发起下线消息。处理流程通销毁对象。

## 处理

分服清理，标记间不存在依赖关系。

- 登陆要保证gate和game都为空
- map gate已经登陆完成后才能的执行的，然后再判断map标记。