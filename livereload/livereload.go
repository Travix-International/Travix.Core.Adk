package livereload

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	hub *Hub
	// Appix always hosts the livereload server at 13221. This is the port the frontend has to try to use to connect to.
	livereloadAddress    = ":13221"
	livereloadAddressTLS = ":13222"
)

const (
	certContent = `-----BEGIN CERTIFICATE-----
MIICkzCCAfwCCQCbmnQ2PFatzzANBgkqhkiG9w0BAQsFADCBjTELMAkGA1UEBhMC
TkwxEzARBgNVBAgMClNvbWUtU3RhdGUxEjAQBgNVBAcMCUFtc3RlcmRhbTEPMA0G
A1UECgwGVHJhdml4MQ0wCwYDVQQLDARDb3JlMRIwEAYDVQQDDAlsb2NhbGhvc3Qx
ITAfBgkqhkiG9w0BCQEWEm12aW5jemVAdHJhdml4LmNvbTAeFw0xNjEwMTUxNzEy
NTVaFw0xOTA4MDUxNzEyNTVaMIGNMQswCQYDVQQGEwJOTDETMBEGA1UECAwKU29t
ZS1TdGF0ZTESMBAGA1UEBwwJQW1zdGVyZGFtMQ8wDQYDVQQKDAZUcmF2aXgxDTAL
BgNVBAsMBENvcmUxEjAQBgNVBAMMCWxvY2FsaG9zdDEhMB8GCSqGSIb3DQEJARYS
bXZpbmN6ZUB0cmF2aXguY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCf
OSB7LkkaRd6WXTplUHfD2k+EHoVi9flKcmbUlye9zHFzWVtCUQjhFjiZL1rNRQGn
9VMUqzpc55RyzTEy2KpyZ+7INR1ZAuqqXMxpNDzXeq+UQuAFnJrHnbwtiSYPiJ45
5EvysllYb5j6ihXEVZt+6QdMINFB+Gz0Xfrhug0+0QIDAQABMA0GCSqGSIb3DQEB
CwUAA4GBADrH8ibFye3iXHR6RkwVNBgeKyvL0kxs4C8785uYqjRJWVjAg2xJQyyZ
R3IHuvKqkmjs5i5d5CT9QT4t8Mlorg1XSnRz/HLf5zrRJlVzqrpd9N2+859TmTVD
9A91NtEwCNgBSGDGSCndjQ/dkPhbJFs28/ICujLySxbYswOGHGbK
-----END CERTIFICATE-----
`
	keyContent = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCfOSB7LkkaRd6WXTplUHfD2k+EHoVi9flKcmbUlye9zHFzWVtC
UQjhFjiZL1rNRQGn9VMUqzpc55RyzTEy2KpyZ+7INR1ZAuqqXMxpNDzXeq+UQuAF
nJrHnbwtiSYPiJ455EvysllYb5j6ihXEVZt+6QdMINFB+Gz0Xfrhug0+0QIDAQAB
AoGAIo9SxonwYhyCSN7peu4xYLh1A/df+m/rcUZNnZ1FigPjKCdgEI/oPnsFQ/Ks
Ydu1lVBBfT4BSAMYDKcPI7s1m5Hf++2TAWXuE/GiMmfmQq8QHVwdRERIzGo7BSIW
alA5tC4+dIe5gUKjR38MpG9VCEa3FBkNxlRQ2U1tIAoM9/ECQQDLWvbShPYpfKCM
8WlAGeWwgHJrjdmatMLsJepxFjGShxK1uhLy6mIMaVVCV0dFPk2Y81ACAirmev99
bqMd3sbtAkEAyHFgTZzQUrezQQhnfFcEDOaUrCwRBVERHFou6wHEwTLObJeedAuo
emRRpQkOp+wJq8y9eOI2pv0jpSI8pTKW9QJAdOuzOG1sX4Qhh4gSHOIG90mTABYK
BHJkFITkW+sHy5jQAB6hYHu0rjAt7jviZYSh9wwGd3Epm2Ui2sqvDLCXLQJBAKAk
NNTNXIM50TU8CbIFs267Kj0EV/Tvd8Q3KRUJLLFObi3EVQxR5CEk1TYNrm/q3S8t
PJO/5/oydLASUnGJoaECQGyPpJ6lVJb10yJKjcGtouwa+HFRJh9BxIQUHZRTbmHX
k7iRrF0Vcllo8k/Mos5PVPP0WIyS1l0lh4GZ+w8gA80=
-----END RSA PRIVATE KEY-----
`
)

func createCertFiles() (cert string, key string) {
	tempFolder, _ := ioutil.TempDir("", "appix")

	cert = tempFolder + "/livereload-cert.pem"
	key = tempFolder + "/livereload-key.pem"

	ioutil.WriteFile(cert, []byte(certContent), 0644)
	ioutil.WriteFile(key, []byte(keyContent), 0644)

	return cert, key
}

// StartServer starts the Websocket server listening for the websites that want to connect.
func StartServer() {
	hub = NewHub()
	go hub.Run()
	http.HandleFunc("/appixlivereload", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	go startLocalServer()
	go startLocalServerTLS()
	log.Println("Livereload server listening at", livereloadAddress)
}

func startLocalServer() {
	err := http.ListenAndServe(livereloadAddress, nil)

	if err != nil {
		log.Println("Failed to start up the Livereload server: ", err)
		return
	}
}

func startLocalServerTLS() {
	cert, key := createCertFiles()
	err := http.ListenAndServeTLS(livereloadAddressTLS, cert, key, nil)

	if err != nil {
		log.Println("Failed to start up the Livereload server with TLS: ", err)
		return
	}
}

// SendReload sends a message to the Websocket listeners to refresh the page.
func SendReload() {
	message := bytes.TrimSpace([]byte("reload"))
	hub.Broadcast <- message
}
