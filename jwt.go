package main

import "github.com/golang-jwt/jwt/v5"
import "crypto/rand"
import "crypto/rsa"
import "encoding/base64"
import "os"
import "strings"
import "time"

type ScopeGroupClaims struct {
	Scope []string `json:"scope,omitempty"`
    Group []string `json:"group,omitempty"`
    jwt.RegisteredClaims
}

func JWTFromSiteDiscordUser(site SiteConfig, user DiscourceUser) (string, error) {
    var signKey *rsa.PrivateKey
    var signBytes []byte
    var err error

    if site.JwtConfig.PrivateKeyBytes == nil {
        signBytes, err = os.ReadFile(site.JwtConfig.PrivateKey)

        if err != nil {
            return "", err
        }

        site.JwtConfig.PrivateKeyBytes = signBytes

    } else {
        signBytes = site.JwtConfig.PrivateKeyBytes
    }

    signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)

    if err != nil {
        return "", err
    }

    subject := getSubject(user, site.JwtConfig.Subject)

    randomBytes := make([]byte, 15)

	if _, err = rand.Read(randomBytes); err != nil {
        return "", err
    }

    groups := make([]string,0)
    scopes := make([]string,0)

    if site.JwtConfig.AllGroups {
        groups = user.Groups 
    }

    if len(site.JwtConfig.Scope) > 0 {
        scopes = site.JwtConfig.Scope
    }

    token := getToken(randomBytes, subject, site.JwtConfig.Audience, scopes, groups, site.JwtConfig.Expiry)

    if site.JwtConfig.KeyId != "" {
        token.Header["kid"] = site.JwtConfig.KeyId
    }

    jwt, err := token.SignedString(signKey)

    return jwt, err
}

func getToken(id []byte, subject string, audience []string, scopes []string, groups []string, expiry int) *jwt.Token {

    
    expires := time.Now().Add(8 * time.Hour)

    if expiry > 0 {
        expires = time.Now().Add(time.Duration(expiry) * time.Second)
    }

    claims := ScopeGroupClaims{
        scopes,
        groups,
        jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expires),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer: "token-swap.portal.bristolhackspace.org",
            Subject: subject,
            ID: base64.URLEncoding.EncodeToString(id),
            Audience: audience,
        }, 
    }

    token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
    token.Header = map[string]any{
        "alg": "RS256",
        "typ": "JWT",
    }

    return token
}

func getSubject(user DiscourceUser, subject string) string {
    switch (strings.ToLower(subject)) {
        case "externalid":
            return user.ExternalID
        case "email":
            return user.Email
        case "username":
            return user.UserName
        default:
            return user.ExternalID
    }
}
