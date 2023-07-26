
-- 1、InstID 是自己生成的，同一个proxy下会重复
-- 2、ID 需要看看怎么实现
-- 3、proxy service/entity的 ProxyID + 0 + 0 
-- 4、gate的ID可以细分一下，自己的entity和client区分开
*** 5、nets.conn需要重构一下，层次不够清晰，还有超时处理等，断连通知等 ***
6、entity manager 启动之后不能再注册service
7、packet重构一下，buf包含长度，write不需要append
-- 8、entity的主动释放
-- 9、entity watch/unwatch
*** 10、帧同步 ***
11、log
12、长连接断开