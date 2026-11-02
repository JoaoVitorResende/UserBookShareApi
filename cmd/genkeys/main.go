package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "log"
    "os"
)

func main() {
    key, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        log.Fatal(err)
    }

    privPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(key),
    })
    if err := os.WriteFile("private.pem", privPEM, 0o600); err != nil {
        log.Fatal(err)
    }

    pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
    if err != nil {
        log.Fatal(err)
    }
    pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
    if err := os.WriteFile("public.pem", pubPEM, 0o644); err != nil {
        log.Fatal(err)
    }
}
