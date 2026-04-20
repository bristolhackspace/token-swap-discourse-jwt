package main

import "os"
import "io"
import "encoding/json"

type JwtConfig struct {
    PublicKey string `json:"publicKey"`
    PublicKeyBytes []byte
    PrivateKey string `json:"privateKey"`
    PrivateKeyBytes []byte
    Subject string `json:"subject"`
    Audience []string `json:"audience"`
    KeyId string `json:"keyId"`
}

type SiteConfig struct {
    Domain string `json:"domain"`
    Alias []string `json:"alias"`
    JwtConfig JwtConfig `json:"jwt"`
    Redirect string `json:"redirect"`
    Cookie string `json:"cookie"`
}

func LoadSiteConfig(file string) (SiteConfig, error) {

    var result SiteConfig;

    jsonFile, err := os.Open(file)

    if err != nil {
        return result, err
    }

    defer jsonFile.Close()

    jsonBytes, err := io.ReadAll(jsonFile)

    if err != nil {
        return result, err
    }

    jsonErr := json.Unmarshal([]byte(jsonBytes), &result)

    if jsonErr != nil {
        return result, jsonErr
    }

    return result, nil

}
