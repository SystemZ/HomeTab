package web

import (
	"crypto/tls"
	"net/http"
)

func WebhookArgocd(w http.ResponseWriter, r *http.Request) {
	// TODO rate limiting
	RefreshArgoCd()
	w.WriteHeader(http.StatusOK)
	return
}

func RefreshArgoCd() {
	req, err := http.NewRequest("POST", "https://argocd-server.argocd/api/webhook", nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
