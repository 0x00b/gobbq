# usage
```shell
# 安装butc
go get -u github.com/0x00b/protoc-gen-gobbq/bbq
# 自动安装或升级依赖的工具
bbq install

# 在当前目录创建一个新的项目，go modulename：git.woa.com/test/test
bbq new git.woa.com/test/test

# 进入项目根目录
cd test

# 编译检查
make
# 添加一个新的pb协议文件
bbq proto add api/test.proto
# 根据项目项目下的pb协议生成代码
make gen

# 运行
go run main.go

```