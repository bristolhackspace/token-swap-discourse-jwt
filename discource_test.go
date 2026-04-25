package main

import "fmt"
import "testing"

func TestEncodeSSO(t *testing.T) {
    sso := EncodeSSO("http://www.testing.com/login", "testing12345")
    fmt.Println("SSO: " + sso)
}

func TestSignSSO(t *testing.T) {
    sig := SignSSO("a-test-string", []byte("G0o1RdIzhusaJgDy4OR8CfX"))
    fmt.Println(sig)
}

func TestBuildConnectUrl(t *testing.T) {
    siteHash := SiteHash("grafana.bristolhackspace.org")
    url := BuildConnectUrl("http://grafana.bristolhackspace.org/path-to/end", "https://bristolhackspace.discourse.group/path-to/sso", []byte("G0o1RdIzhusaJgDy4OR8CfX"), siteHash) 

    fmt.Println(url)
}

func TestDecodeSSOParameter(t *testing.T) {
    sso := "bm9uY2U9MTc3NzE1NzIyMC1pOWxvNXRmV1ZZLWQ2ZTEwMmJkYWYmZW1haWw9YWFyb25kcyU0MGdtYWlsLmNvbSZleHRlcm5hbF9pZD0xMCZuYW1lPUFhcm9uK1NocmltcHRvbiZncm91cHM9YWRtaW4lMkNjb21taXR0ZWUmdXNlcm5hbWU9YWFyb25kcw=="

    user, err := DecodeSSOParameter(sso)

    if err != nil {
        t.Errorf("SSO Decode error: %v", err) 
    }

    if user.ExternalID == "" {
        t.Errorf("No external_id")
    }

    if user.Email == "" {
        t.Errorf("No email")
    }

    if user.Name == "" {
        t.Errorf("No name")
    }

    if user.UserName == "" {
        t.Errorf("No UserName")
    }

    if len(user.Groups) < 2 {
        t.Errorf("Groups missing")
    }
}
