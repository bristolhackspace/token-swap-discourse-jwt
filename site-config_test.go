package main

import "fmt"
import "testing"

func TestLoadSiteConfig(t *testing.T) {
    siteConfig, err := LoadSiteConfig("site-example.json")

    fmt.Println(siteConfig)
    fmt.Println(err)
}
