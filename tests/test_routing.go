
package main

import (
  "regexp"
  "fmt"
)

func parseUrl(url string) {
  m := validPath.FindStringSubmatch(url)
  if m == nil {
    fmt.Printf("FAILED: %s\n\n", url)
    return
  }
  fmt.Printf("url: %s\n", m[0])
  fmt.Printf("dir: %s\n", m[1])
  fmt.Printf("res: %s\n", m[2])
  fmt.Println()
}

var regexString = "^(/|/view/|/edit/|/save/|/delete/)([a-zA-Z0-9_]*)$"

var validPath = regexp.MustCompile(regexString)

func main() {
  fmt.Printf("regex: %s\n\n", regexString)

  fmt.Println("should pass:")
  parseUrl("/view/Something")
  parseUrl("/edit/Something")
  parseUrl("/delete/Something")
  parseUrl("/delete/Something_Else")
  parseUrl("/Something")
  parseUrl("/save/something")
  parseUrl("/")

  fmt.Println("should fail:")
  parseUrl("/fake/Something")
  parseUrl("/view/fake/Something")
  parseUrl("//")
  parseUrl("")
  parseUrl("/delete/some thing")
  parseUrl("/view/delete/thing")
  parseUrl("/view/../file")
  parseUrl("../../etc/passwd")
  parseUrl("view/file")
  parseUrl("../view/file")
}

