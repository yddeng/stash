# TODO List
- [x] 调整结构，开放自定义节点创建功能（用户自定义节点并在行为树系统中注册）
- [x] 行为树系统可以注册用户自定义行为
- [x] Task 类型定义能否优化设计，类型检查过程能否优化？Task只有软类型（TaskType，帮助debug），没有硬类型（interface），统一了接口，避免因类型断言带来的开销。
- [x] 通过 Task 创建 Agent 的过程能否优化？agent定义统一，不再根据 TaskType 针对性的设计。
- [x] 生命周期短，创建频率高的对象（agent，task），需要通过对象池进行管理。行为树系统统一管理
- [x] Env 改名。Runtime 不太合适，运行时虽然也携带数据，具备上下文的属性，但概念更广。直接使用 Context 吧。
- [x] DataContext 改名，DataCenter？挪回根包下。改为dataSet，类型不导出，导出的接口继承到 Context 中。
- [x] 实现了测试用subtree结点
- [ ] Context 之间是否需要共享 DataSet ？。例如 subtree 与 父树（完整的行为树）共享 DataSet