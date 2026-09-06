package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// oidcDiscoveryDocument mirrors the subset of an OpenID Connect discovery
// document (OIDC Discovery 1.0 / RFC 8414) this CLI needs.
type oidcDiscoveryDocument struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

// DiscoverEndpoints fetches the OIDC discovery document at
// issuer + "/.well-known/openid-configuration" and returns its authorization
// and token endpoints. This lets callers configure just an issuer/base URL,
// as most standards-compliant identity providers support discovery, instead
// of every individual endpoint.
func DiscoverEndpoints(ctx context.Context, httpClient *http.Client, issuer string) (authURL, tokenURL string, err error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	discoveryURL := strings.TrimRight(issuer, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create discovery request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to reach discovery endpoint %s: %w", discoveryURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read discovery response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("discovery endpoint %s returned status %d: %s", discoveryURL, resp.StatusCode, string(body))
	}

	var doc oidcDiscoveryDocument
	if err := json.Unmarshal(body, &doc); err != nil {
		return "", "", fmt.Errorf("failed to parse discovery document from %s: %w", discoveryURL, err)
	}
	if doc.AuthorizationEndpoint == "" || doc.TokenEndpoint == "" {
		return "", "", fmt.Errorf("discovery document from %s is missing authorization_endpoint or token_endpoint", discoveryURL)
	}
	return doc.AuthorizationEndpoint, doc.TokenEndpoint, nil
}
