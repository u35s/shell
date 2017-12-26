## shell
* shell  是一个内网穿透工具,提供了当网络失效时自动恢复特性
* shelld 是公网守护进程,shellc webc 都是内网客户端
* shellc 可以把自己的端口映射到公网的某个端口上
* webc   本身可以启动一个web客户端并映射到公网指定端口
## 部署
* 服务器端 go get github.com/u35s/shell/shelld
* 拷贝 example.shell.conf 到某个目录 修改配置
* 执行 shelld
* 客户端 go get github.com/u35s/shell/shelld
* 拷贝 example.shell.conf 到某个目录 修改配置
* 执行 shellc
