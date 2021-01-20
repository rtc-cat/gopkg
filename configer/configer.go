/*
Package configer 按照优先级读取配置
1. 配置文件
2. 环境变量
*/
package configer

import (
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"

	"github.com/krizsx/gopkg/errorx"
)

const (
	DefaultConfigFile = "app.yaml"
	DefaultEnvPrefix  = "APP_"
)

var (
	_configer *Configer // 默认
)

func init() {
	_configer = NewConfiger()
}

// Configer 多种方式读取配置
type Configer struct {
	file      string
	envPrefix string
}

// NewConfiger 构造函数
func NewConfiger(opts ...Option) *Configer {
	c := &Configer{
		file:      DefaultConfigFile,
		envPrefix: DefaultEnvPrefix,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(*Configer)

// WithFile 指定配置文件
func WithFile(file string) Option {
	return func(c *Configer) {
		c.file = file
	}
}

// WithEnvPrefix 指定环境前缀
func WithEnvPrefix(envPrefix string) Option {
	return func(c *Configer) {
		c.envPrefix = strings.ToUpper(envPrefix)
	}
}

func Load(config interface{}) error {
	return _configer.Load(config)
}

// Load 读取文件,以及环境变量
func (c *Configer) Load(config interface{}) error {
	const op = errorx.ErrOperation("configer.Load")

	f, err := os.Open(c.file)
	if err != nil {
		return errorx.E(err, op)
	}
	defer f.Close()
	if err = yaml.NewDecoder(f).Decode(config); err != nil {
		return errorx.E(err, op)
	}
	if err = envconfig.Process(c.envPrefix, config); err != nil {
		return errorx.E(err, op)
	}
	return nil
}
