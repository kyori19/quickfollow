package main

import (
  "fmt"

  "github.com/spf13/viper"
)

type config struct {
  upstream string
  branch string
  target []string
}

func load(c *context, path string) config {
  viper.SetConfigName("quickfollow")
  viper.AddConfigPath(path)
  if err := viper.ReadInConfig(); err != nil {
    if _, ok := err.(viper.ConfigFileNotFoundError); ok {
      c.warn("Config file not found. Trying to checkout to master")
      checkout(c, path, "master")
      err = viper.ReadInConfig()
    }
    if err != nil {
      c.panic(err)
    }
  }

  checkIsSet(c, "upstream")
  viper.SetDefault("branch", "master")
  checkIsSet(c, "target")

  return config {
    upstream: viper.GetString("upstream"),
    branch: viper.GetString("branch"),
    target: viper.GetStringSlice("target"),
  }
}

func checkIsSet(c *context, key string) {
  if !viper.IsSet(key) {
    c.panic(fmt.Errorf("Config \"%s\" is not set", key))
  }
}
