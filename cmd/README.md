# `/cmd`
当前项目的可执行文件。cmd 目录下的每一个子目录名称都应该匹配可执行文件。比如果我们的项目是一个 grpc 服务，在 /cmd/myapp/main.go 中就包含了启动服务进程的代码，编译后生成的可执行文件就是 myapp。

不要在 /cmd 目录中放置太多的代码，我们应该将公有代码放置到 /pkg 中，将私有代码放置到 /internal 中并在 /cmd 中引入这些包，保证 main 函数中的代码尽可能简单和少。