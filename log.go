package main

import (
  "fmt"
  "log"
  "strings"

  "github.com/mattn/go-colorable"
)

var (
  colors = []string {
    "\x1b[31m",
    "\x1b[32m",
    "\x1b[33m",
    "\x1b[34m",
    "\x1b[35m",
    "\x1b[36m",
    "\x1b[37m",
  }
)

type colored struct {
  text string
  index int
}

func (c colored) build() string {
  return fmt.Sprintf("%s%s\x1b[0m", colors[c.index % len(colors)], c.text)
}

func (c colored) infect(s interface{}) string {
  return fmt.Sprintf("%s%s\x1b[0m", colors[c.index % len(colors)], s)
}

type context struct {
  value []colored
}

func initContext() *context {
  log.SetOutput(colorable.NewColorableStdout())

  return &context {
    value: []colored {
      colored {
        text: "I",
        index: 0,
      },
    },
  }
}

func pop(slice []colored) (colored, []colored) {
  ans := slice[len(slice) - 1]
  slice = slice[:len(slice) - 1]
  return ans, slice
}

func (c *context) next(s string) {
  last, value := pop(c.value)
  value = append(value, colored {
    text: s,
    index: last.index + 1,
  })
  c.value = value
}

func (c *context) step(s string) {
  value := append(c.value, colored {
    text: s,
    index: 0,
  })
  c.value = value
}

func (c *context) back() {
  _, value := pop(c.value)
  c.value = value
}

func (c context) build() string {
  var b strings.Builder
  for _, item := range c.value {
    b.WriteString(item.build())
  }
  return b.String()
}

func (c context) log(l colored, s string) {
  log.Printf("[%s%s] %s", l.build(), c.build(), l.infect(s))
}

func (c *context) debug(a ...interface{}) {
  s := fmt.Sprint(a...)
  c.log(colored {
    text: "D",
    index: 1,
  }, s)
}

func (c *context) info(format string, a ...interface{}) {
  s := fmt.Sprintf(format, a...)
  c.log(colored {
    text: "I",
    index: 6,
  }, s)
}

func (c *context) warn(format string, a ...interface{}) {
  s := fmt.Sprintf(format, a...)
  c.log(colored {
    text: "W",
    index: 2,
  }, s)
}

func (c *context) error(s string) {
  c.log(colored {
    text: "E",
    index: 0,
  }, s)
}

func (c *context) panic(e error) {
  c.error(e.Error())
  panic(e)
}
