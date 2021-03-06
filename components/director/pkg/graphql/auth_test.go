package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalBasicAuth(t *testing.T) {
	// GIVEN
	a := &Auth{}
	// WHEN
	err := a.UnmarshalJSON([]byte(`{
		"credential": {
			"username": "aaa",
			"password": "bbb"
		},
		"additionalHeaders": {
			"scopes": ["read", "write"]
		}
	}`))
	// THEN
	require.NoError(t, err)
	require.NotNil(t, a.AdditionalHeaders)
	scopes := (*a.AdditionalHeaders)["scopes"]
	assert.Len(t, scopes, 2)
	assert.Contains(t, scopes, "read")
	assert.Contains(t, scopes, "write")
	basic, ok := a.Credential.(*BasicCredentialData)
	require.True(t, ok)
	assert.Equal(t, "aaa", basic.Username)
	assert.Equal(t, "bbb", basic.Password)
}

func TestUnmarshalOAuth(t *testing.T) {
	// GIVEN
	a := &Auth{}
	// WHEN
	err := a.UnmarshalJSON([]byte(`{
  		"credential": {
			"url":"oauth.url",
			"clientId": "client-id",
			"clientSecret":"client-secret"
		},
		"additionalHeaders": {
			"scopes": ["read", "write"]
		}
	}`))
	// THEN
	require.NoError(t, err)
	require.NotNil(t, a.AdditionalHeaders)
	scopes := (*a.AdditionalHeaders)["scopes"]
	assert.Len(t, scopes, 2)
	assert.Contains(t, scopes, "read")
	assert.Contains(t, scopes, "write")
	oauth, ok := a.Credential.(*OAuthCredentialData)
	require.True(t, ok)
	assert.Equal(t, "oauth.url", oauth.URL)
	assert.Equal(t, "client-id", oauth.ClientID)
	assert.Equal(t, "client-secret", oauth.ClientSecret)
}

func TestUnmarshalCSRFBasicAuth(t *testing.T) {
	// GIVEN
	a := &CSRFTokenCredentialRequestAuth{}
	// WHEN
	err := a.UnmarshalJSON([]byte(`{
		"credential": {
			"username": "aaa",
			"password": "bbb"
		},
		"additionalHeaders": {
			"scopes": ["read", "write"]
		}
	}`))
	// THEN
	require.NoError(t, err)
	require.NotNil(t, a.AdditionalHeaders)
	scopes := (*a.AdditionalHeaders)["scopes"]
	assert.Len(t, scopes, 2)
	assert.Contains(t, scopes, "read")
	assert.Contains(t, scopes, "write")
	basic, ok := a.Credential.(*BasicCredentialData)
	require.True(t, ok)
	assert.Equal(t, "aaa", basic.Username)
	assert.Equal(t, "bbb", basic.Password)
}

func TestUnmarshalCSRFOAuth(t *testing.T) {
	// GIVEN
	a := &CSRFTokenCredentialRequestAuth{}
	// WHEN
	err := a.UnmarshalJSON([]byte(`{
  		"credential": {
			"url":"oauth.url",
			"clientId": "client-id",
			"clientSecret":"client-secret"
		},
		"additionalHeaders": {
			"scopes": ["read", "write"]
		}
	}`))
	// THEN
	require.NoError(t, err)
	require.NotNil(t, a.AdditionalHeaders)
	scopes := (*a.AdditionalHeaders)["scopes"]
	assert.Len(t, scopes, 2)
	assert.Contains(t, scopes, "read")
	assert.Contains(t, scopes, "write")
	oauth, ok := a.Credential.(*OAuthCredentialData)
	require.True(t, ok)
	assert.Equal(t, "oauth.url", oauth.URL)
	assert.Equal(t, "client-id", oauth.ClientID)
	assert.Equal(t, "client-secret", oauth.ClientSecret)
}
