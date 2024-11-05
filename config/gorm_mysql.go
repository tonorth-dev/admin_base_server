package config

import "fmt"

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", "root", "123456Abc", "127.0.0.1", "3306", "hongshi_interview", "charset=utf8mb4&parseTime=True&loc=Local")
}
