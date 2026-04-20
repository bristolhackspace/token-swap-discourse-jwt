package main
import "testing"
import "fmt"

func TestJwt(t *testing.T) {
    siteConfig, err := LoadSiteConfig("site-example.json")
    testUser := DiscourceUser{
        "01234-56778",
        "testing@testing.com",
        "Test Person",
        "test.person",
        "123456-asfaposi-01234abc",
    }

    jwt, err := JWTFromSiteDiscordUser(siteConfig, testUser)

    if jwt == "" || err != nil {
        t.Errorf("jwt == '' or error %v", err) 
    }

    fmt.Println(jwt)
}
