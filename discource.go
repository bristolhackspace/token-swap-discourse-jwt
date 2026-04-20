package main

import "net/url"
import "encoding/base64"
import "encoding/hex"
import "crypto/hmac"
import "crypto/sha256"
import "crypto/rand"
import "time"
import "strconv"
import "strings"

type DiscourceUser struct {
    ExternalID string
    Email string
    Name string
    UserName string
    Nonce string
}

func SiteHash(host string) string {
    sh := sha256.Sum256([]byte(host))
    siteHash := hex.EncodeToString(sh[:])
    return siteHash[:10]
}

func EncodeSSO(retUrl string, nonce string) string {
    v := url.Values{}

    v.Set("return_sso_url", retUrl)
    v.Set("nonce", nonce)

    urlEncoded := v.Encode()
    return base64.URLEncoding.EncodeToString([]byte(urlEncoded))
}

func ValidateSSO(payload string, sig string, key[] byte) (bool, error) {
    mac := hmac.New(sha256.New, key)
    mac.Write([]byte(payload))
    sigBytes := mac.Sum(nil)

    decodeBytes, err := hex.DecodeString(sig)

    if err != nil {
        return false, err
    }

    return hmac.Equal(sigBytes, decodeBytes), nil 
}

func SignSSO(sso string, key []byte) string {
    mac := hmac.New(sha256.New, key)
    mac.Write([]byte(sso))
    sigBytes := mac.Sum(nil)
    return hex.EncodeToString(sigBytes)
}

func BuildConnectUrl(returnUrl string, discourceConnectUrl string, key []byte, siteHash string) string {
   numOnce := strconv.FormatInt(time.Now().Unix(), 10) + "-" + shortID(10) + "-" + siteHash
   sso := EncodeSSO(returnUrl, numOnce)
   sig := SignSSO(sso, key)

    v := url.Values{}
    v.Set("sso", sso)
    v.Set("sig", sig)
    urlEncoded := v.Encode()

    return discourceConnectUrl + "?" + urlEncoded
}

func ValidateNumOnce(maxAge int64, now int64, siteHash string, numOnce string) bool {
    parts := strings.Split(numOnce, "-")
    
    if len(parts) < 3 {
        return false
    }

    issued, err := strconv.ParseInt(parts[0], 10, 64)

    if err != nil || (now - issued) > maxAge {
        return false
    }

    if parts[2] != siteHash {
        return false
    }

    return true
}

func DecodeSSOParameter(sso string) (DiscourceUser, error) {
    queryBytes, err := base64.URLEncoding.DecodeString(sso)

    if err != nil {
        return DiscourceUser{}, err
    }

    return DecodeDiscourceUser(string(queryBytes)) 
}

func DecodeDiscourceUser(query string) (DiscourceUser, error) {
    var dcUser DiscourceUser 

    queryValues, err := url.ParseQuery(query)

    if err != nil {
        return dcUser, nil
    }

    if queryValues.Get("email") != "" {
        dcUser.Email = queryValues.Get("email")
    }

    if queryValues.Get("external_id") != "" {
        dcUser.ExternalID = queryValues.Get("external_id")
    }

    if queryValues.Get("name") != "" {
        dcUser.Name = queryValues.Get("name")
    }

    if queryValues.Get("username") != "" {
        dcUser.UserName = queryValues.Get("username")
    }

    if queryValues.Get("nonce") != "" {
        dcUser.Nonce = queryValues.Get("nonce")
    }

    return dcUser, nil
}



func shortID(length int) string {
    var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
    ll := len(chars)
    b := make([]byte, length)
    rand.Read(b) // generates len(b) random bytes
    for i := 0; i < length; i++ {
        b[i] = chars[int(b[i])%ll]
    }
    return string(b)
}
