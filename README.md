# gobbq
gobbq is golang game server framework.

why gobbq? because I often go to barbecue with friends,  I'd rather everyone have time to barbecue than work all the time.

## components
gate
proxy - sidecar
game
client

# note!!! this repo is a baby now. everything is variable. but I hope to reach this feature:

# 理想的游戏后台框架能力

## 1、完备的基础能力
### 更加现代的编程语言：
* 可调试性：代码修改1分钟内完成重新部署和调试，调试过程可视化
* 安全：可靠的访存管理、异常保护机制
* 编码效率：能1行搞定的代码，为什么要写10行
* 良好的泛型：DRY
* 标准库：统一的标准库、完善的依赖管理
* 编译效率：减少的每一分钟都是生命
### 可靠的客户端与服务器通讯
#### 长连接保持
* 会话保持和重连能力
* 服务器指定一组客户端广播
* 支持后端任意服务下推消息
#### 接口级别网关
* 精确流控、熔断、接口保护
* 限频和限流
#### 通信安全和加密
#### 多协议支持，TCP/UDP/RUDP
#### 短链接支持，如HTTP/HTTPS支持
### 自由的服务器间通讯
* 服务注册和自动发现能力
* 支持消息指定地址单播、组合地址多播、受限广播等多种发送能力
* 负载均衡：被调方有负载自动迁移能力，发送方无需关注
* 可靠送达：在不发生网络分区或节点故障情况下，消息默认可靠、一次送达
* 消息保序：同一对互相通信的节点，多个连续包需要保证顺序达到

### 支持动态脚本： lua, python，简化逻辑编写或实现快速热更能力
### RPC：不管是gRPC, Spp, TAF, Tars, TRPC, IRPC(tsf4g2.0)，都需要有良好网络协议支持、IDL、多语言、性能和周边生态。
### 序列化：跨语言、高性能、自动代码生成或无IDL的自解析，良好的跨版本兼容能力。
### 定时器和事件循环：支持不同颗粒度的高性能定时器模块。
### 存储层抽象：需要一个解耦的存储接口层，并处理好实现的细节（如不同DB版本的特性兼容，分布式存储特性支持）
* Mysql
* Elastic Search
* Redis
* Etcd
* Zookeeper

### 可观测：几乎都是必备的了
* 调用链跟踪
* 日志接口
* 监控上报
* 共享内存支持（可选）
* 多线程模型（可选）
### 游戏业务通用能力
#### 并发编程模型: 多线程被证明不是一个很好的选择，CSP有局限性，相比之下，actor或许更加成熟，理解成本更低
#### 异步抽象（tsf4g里叫做transation）：80%的开发时间都在跟它打交道
* 异步回调：代码编写复杂，开发周期长，维护困难，BUG多
* 有限状态机（FSM）:本质是回调，但通过提前约定状态和流转，控制了严重BUG的出现几率
* Promise/Future
* 协程： 同步编程的诱惑实在是太大了，无法拒绝的特性。但许多细节需要考虑，如果是自己实现，有栈（独立or共享）vs无栈，对称vs非对称，更多(系统级Hook, 协程间* 交互数据，锁，定时器等)
#### 有状态服务支持：很重要的特性，需要考虑的点很多。
* 服务发现和路由：考虑路由切换式的请求低损
* 数据分区方式：如何分配数据（关键词 or hash or 多级索引）; 如何选主（主从，多主节点）
* 容灾方式：节点的新增、删除、数据恢复和数据迁移
#### 游戏资源管理： 游戏资源的热更和运营干预能力
#### 游戏配置管理：支持统一配置管理中心（如七彩石），实现配置的跨环境、多版本管理、权限管理、快速回滚/复制/灰度/校验、发布管理
* 事务一致性支持：从框架层面支持事务特性，能大大减少经济道具类的风险
### 完善的工具链和周边
#### 可测试性：符合测试左移的理念，保证提交质量可控，减少手工活
* 单测支持
* 代码覆盖
* 自动化测试
* 压力测试
* 静态检查
* 代码复杂度检测
* -
#### 完善的IDE支持：支持代码静态检测、编码规范等
#### 组件接入支持：适当抽象，并适配不同介入模式的组件（SDK、http、agent、udp上报等）
### 可持续运营能力
* 快速部署：1分钟快速构建开发环境和测试环境的能力
* 线上维护能力：TKE(k8s)基本都提供了解决方案
    - 自动扩缩容
    - 优雅启停
    - 原地热更
    - 滚动更新
* 更加智能的AI Ops： 如核心监控KPI曲线的异常检测、异常错误日志的自动告警
* 软件可定义网络（ SDN）：SDN是一个思想，其实这里讲的是业务层面，需要支持动态变化的配置（如依赖云、组件的地址配置、密钥）和版本镜像的解耦，真正实现“一次构* 建、多处可用”。
* ABTest能力：从框架层面支持用户/请求/功能系统/服务节点的AB测试/灰度能力
