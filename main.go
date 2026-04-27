package main

import "net/http"
import "github.com/gin-gonic/gin"
import "log"
import "os"
import "strings"
import "time"
import "syscall"

func main() {
    router := gin.Default()
    discourceServer := os.Getenv("TOKEN_SWAP_DISCOURSE_SERVER")
    discourceKey := []byte(os.Getenv("TOKEN_SWAP_DISCOURSE_KEY"))
    port := os.Getenv("TOKEN_SWAP_PORT")
    //failRedirect := "http://localhost:5000/"

    if port == "" {
        port = "0.0.0.0:3400"
    }

    hostToSite := make(map[string]SiteConfig)

    entries, err := os.ReadDir("./sites")

    if err != nil {
        log.Fatal(err)
    }


    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(),".json") {
            site, err := LoadSiteConfig("./sites/" + e.Name()) 

            if err == nil {
                domain := strings.TrimSuffix(e.Name(), ".json")
                site.Domain = domain
                log.Printf("Loading: %s\n",domain)
                hostToSite[site.Domain] = site

                for _, alias := range site.Alias {
                    hostToSite[alias] = site
                    log.Printf("Alias: %s\n", alias)
                }
            } else {
                log.Println(err)
            }
        }
    }

    router.GET("/.well-known/token-swap-discourse-jwt/start", func(c *gin.Context) { 
        host := strings.ToLower(c.Request.Host)

        siteConfig, ok := hostToSite[host]

        if !ok {
            c.String(http.StatusNotFound, host)
            return
        }

        protocol := "https://" 

        if strings.HasPrefix(host, "localhost:") || strings.HasPrefix(siteConfig.Redirect, "http://")  {
            protocol = "http://"
        }

        siteHash := SiteHash(host)

        redirectUrl := BuildConnectUrl(protocol + host + "/.well-known/token-swap-discourse-jwt/end", discourceServer, discourceKey, siteHash)
        c.Redirect(http.StatusFound, redirectUrl)
    })

    router.GET("/.well-known/token-swap-discourse-jwt/end", func(c *gin.Context) {
        sso := c.Query("sso")
        sig := c.Query("sig")
        host := strings.ToLower(c.Request.Host)
        siteHash := SiteHash(host)
        siteConfig, ok := hostToSite[host]

        if !ok {
            c.String(http.StatusNotFound, host)
            return
        }

        valid, err := ValidateSSO(sso, sig, discourceKey)

        if err != nil {
            log.Println(err)
            c.String(http.StatusInternalServerError, "Validation error")
        }

        if valid {

            user, err := DecodeSSOParameter(sso)

            expiredOrAlien := !ValidateNumOnce(30, time.Now().Unix(), siteHash, user.Nonce)

            if expiredOrAlien {
                log.Println("Expired NumOnce")
                c.String(http.StatusInternalServerError, "Link expired")
                return
            }

            if err != nil {
                log.Println(err)
                c.String(http.StatusInternalServerError, "User decode error")
                return
            }

            jwt, err := JWTFromSiteDiscordUser(siteConfig, user)  

            if err != nil {
                log.Println(err)
                c.String(http.StatusInternalServerError, "JWT generate error")
                return
            }

            cookieSecure := true

            if strings.HasPrefix(host, "localhost:") {
                cookieSecure = false
            }

            if siteConfig.Cookie != "" {
                c.SetCookie(siteConfig.Cookie, jwt, 0, "/", "", cookieSecure, true)
            }

            if siteConfig.Redirect != "" {
                redirectUrl := strings.Replace(siteConfig.Redirect, "{{jwt}}", jwt, -1)
                c.Redirect(http.StatusFound, redirectUrl)
            } else {
                c.String(http.StatusOK, jwt)
            }

        } else {
            c.String(http.StatusUnauthorized,"Not OK")
        } 
    })

    router.GET("/.well-known/token-swap-discourse-jwt/health", func(c *gin.Context) {
        c.String(http.StatusOK, "OK")
    })


    if strings.HasPrefix(port, "/") {
        syscall.Umask(0o007)
        router.RunUnix(port)
    } else {
        router.Run(port);
    }
}
