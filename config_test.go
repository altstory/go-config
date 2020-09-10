package config

import (
	"testing"

	"github.com/huandu/go-assert"
)

type testConfig struct {
	FooConfig testFoo `config:"foo"`
	BarConfig testBar `config:"bar"`
}

type testFoo struct {
	Foo1 int      `config:"foo1"`
	Foo2 string   `config:"foo2"`
	Foo3 []string `config:"foo3"`
}

type testBar struct {
	Bar1 string  `config:"bar1"`
	Bar2 float64 `config:"bar2"`
}

func TestConfig(t *testing.T) {
	a := assert.New(t)

	// 测试加载文件。
	c, err := LoadFile("internal/testdata/service.conf")
	a.NilError(err)
	a.Assert(c)

	// 测试反序列化。
	var tc testConfig
	a.NilError(c.Unmarshal("", &tc))
	a.Equal(tc, testConfig{
		FooConfig: testFoo{
			Foo1: 123,
			Foo2: "abc",
			Foo3: []string{"a:b", "c:d", "e:f"},
		},
		BarConfig: testBar{
			Bar1: "single quote",
			Bar2: 456.78,
		},
	})

	// 测试部分反序列化。
	var tf *testFoo
	a.NilError(c.Unmarshal("foo", &tf))
	a.Equal(tf, &testFoo{
		Foo1: 123,
		Foo2: "abc",
		Foo3: []string{"a:b", "c:d", "e:f"},
	})

	var bar1 string
	a.NilError(c.Unmarshal("bar.bar1", &bar1))
	a.Equal(bar1, "single quote")

	var notExist *testBar
	a.NilError(c.Unmarshal("not_exist", &notExist))
	a.Equal(notExist, nil)

	// 测试加载覆盖文件。
	var tcExt testConfig
	a.NilError(c.LoadExt("internal/testdata/service-ext.conf"))
	a.NilError(c.Unmarshal("", &tcExt))
	a.Equal(tcExt, testConfig{
		FooConfig: testFoo{
			Foo1: 456,
			Foo2: "",
			Foo3: []string{"a:b", "e:f", "g:h"},
		},
		BarConfig: testBar{
			Bar1: "another",
			Bar2: 456.78,
		},
	})
}

func TestConfigFailure(t *testing.T) {
	a := assert.New(t)
	a.NonNilError(LoadFile("do/not/exist"))
}
