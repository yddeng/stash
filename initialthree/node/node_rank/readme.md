## 排行榜

计算排名：排行榜中玩家数据一旦被更新，重新计算排名。
获取 top 数据：获取在当前排行榜中排名前 N 个的玩家数据。
获取当前玩家的排名：即使当前玩家没有名列前茅，也可以获取他在排行中的具体名次，以此了解到当前玩家和前排玩家的差距。
获取当前玩家附近排名的玩家：例如为当前玩家寻找水平相当的对手或好友。
数据重置：支持自动定期重置或手动重置数据，比如在一个赛季结束之时或在指定天、周、月后自动重置排行榜，或是在应用测试阶段、应用数据出现误差的情况下进行手动重置。

更新模式：提供两种分数更新模式，better 保留玩家的最好成绩，last 保留玩家的最新成绩。

### 数据库

    每一个排行榜实例，对应一张玩家数据表。由排行榜的生命周期 创建、删除玩家数据表
    
    ```
    rank_list (排行榜数据)
    id(榜实例ID), rank_info(排行榜基础数据)，settled_idx(结算时索引),logic_addr(创建该榜的服逻辑地址）
    
    rank_role_%d (根据已存在实例生成，删除。记录)
    key(玩家ID), score, role_info, role_idx
        
    ```
### 线程池 与 数据库连接池

  线程池，整个进程通用。异步处理耗时的 数据库操作，完成后返回对应的线程处理结果。
  连接池，整个进程通用。
    
    
### 脏标记存储

  定时存储 + 脏标记数量阀值 + 异步保存
  保存失败，将 数据添加脏标记（不存在添加，已经存在不添加。可能为第二次的 setSroce 标记的脏数据，则不能覆盖）

### 排行榜数据

  排行榜数据不一定为最新数据，请求消息可能在排队。
  进入 top 榜的玩家，数据立即存库，更新 top 榜数据。若 存库失败设置脏标记，由tick触发脏标记保存。
  没有进入 top 榜的玩家，设置脏标记，由tick触发保存。

### 结算
 
  获取排名数据，从头遍历下发奖励。邮件添加，数据库记录已经发送过的索引，每执行一次，数据库保存一次。

### 程序重启

  每个服只负责由其创建的榜，由逻辑地址区别。

  根据排行状态。从数据库 加载 玩家分数，重建排行榜实例、top榜数据。

  Begin，End，Setting，StatusWaitExpire ：状态需要重建排行榜及top榜数据。 StatusWaitExpire： 结算已经结束，仅生成 top 榜数据给玩家展示，清除排行榜。

  排行榜处于结算期：重建排行榜实例、top榜数据， 根据结算的索引开始， 继续向后下发奖励邮件。

### 程序停止

  拒绝所有请求
  等待异步调用队列执行完成
  等待各个 rank 脏数据落地；若在结算阶段，停止结算。
  释放数据库连接池
  
### rank 逻辑地址标记

  记录改榜由谁创建，并且重启后，该榜只能由对应的逻辑地址服重构运行。