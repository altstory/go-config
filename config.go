package config

import (
	"github.com/altstory/go-data"
	"github.com/pelletier/go-toml"
)

const (
	tagName = "config"
)

// Config 代表一个配置解析器。
type Config struct {
	data data.Data
}

// LoadFile 解析指定路径的文件并返回 Config 结构。
// 当前只支持 toml 格式的配置文件，未来如果有需要再增加其他格式。
func LoadFile(path string) (c *Config, err error) {
	tree, err := toml.LoadFile(path)

	if err != nil {
		return
	}

	m := tree.ToMap()
	d := data.Make(m)
	c = &Config{
		data: d,
	}
	return
}

// Unmarshal 反序列化 p 中指定的配置文件。
// 如果 section 为空，则将整个配置文件反序列化到 v。
//
// 需要注意，section 名字中如果含有“.”，会被当做分隔符使用，
// 例如 http.server，会访问配置文件的 http -> server 的配置内容。
func (c *Config) Unmarshal(section string, v interface{}) error {
	dec := &data.Decoder{
		TagName: tagName,
	}
	return dec.DecodeQuery(c.data, section, v)
}

const deletesKey = "_deletes"

type deletesConfig struct {
	Deletes []string `toml:"_deletes"`
}

// LoadExt 加载额外的配置文件，用来覆盖主配置文件，
// 额外配置与主配置相同的配置项会采用覆盖策略。
//
// 额外配置文件有一个特殊字段 `_deletes` 用来删除任意字段。
//
//     # 删除配置项 foo 和 bar.player
//     _deletes = ['foo', 'bar.player']
//
// 需要注意，默认覆盖策略是深度合并各种选项，对于数组和对象类型的配置而言，
// 不会简单替换，而是合并，数组会追加到原数组后，对象会追加更多 key。
// 如果需要完全替换原来的数组、对象，应该使用 `_deletes` 来删除原来的字段，
// 再添加新的字段。
func (c *Config) LoadExt(path string) error {
	tree, err := toml.LoadFile(path)

	if err != nil {
		return err
	}

	var del deletesConfig
	tree.Unmarshal(&del)
	tree.Delete(deletesKey)
	m := tree.ToMap()

	d := data.Make(m)
	patch := data.NewPatch()
	patch.Add(del.Deletes, map[string]data.Data{
		"": d,
	})
	err = patch.ApplyTo(&c.data)
	return err
}
