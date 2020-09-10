# go-config：配置文件读取库 #

`go-config` 是一个配置文件读取库，用来从约定好的配置文件里读取配置信息，并且反序列化成 Go 的数据结构。

当前只支持 TOML 文件，格式详见 [TOML v0.5.0](https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.5.0.md)。

## 使用方法 ##

`go-config` 可以解析约定位置的配置文件。

```go
import "github.com/altstory/go-config"

type MyConfig struct {
    Foo int    `config:"foo"`
    Bar string `config:"bar"`
}

func main() {
    // 假定配置文件 file.conf 内容如下：
    //
    // [my_config]
    // foo = 123
    // bar = "player"

    c, err := config.LoadFile("path/to/config/file.conf")

    if err != nil {
        // 处理错误……
        return
    }

    var myConf MyConfig
    c.Unmarshal("my_config", &myConf)
    fmt.Println(myConf.Foo, myConf.Bar) // 输出：123    player
}
```

业务代码推荐使用 `go-runner` 的自动加载机制来加载配置，不要直接使用 `go-config` 接口。
