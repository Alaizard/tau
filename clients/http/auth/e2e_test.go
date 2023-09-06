// This file implements testing functionality for sending client requests on the test environment as a means for end-to-end testing with the aim of catching bugs.
package client // functionality for interacting with API

import (
	"context"
	"os"
	"testing"

	"github.com/taubyte/tau/clients/http"
	"gotest.tools/v3/assert"
)

var (
	authNodeUrl = "https://auth.tau.sandbox.taubyte.com"
	testToken   = testGitToken()
) // Setting the Url to ping and the auth token for access to the API

func testGitToken() string { // Making sure the test environment is being used for test-specifc API requests
	token := os.Getenv("TEST_GIT_TOKEN")

	if token == "" {
		panic("TEST_GIT_TOKEN not set") // Panic is Go's Error handling throw keyword
	}

	return token
}

func newTestClient() (*Client, error) {
	ctx := context.Background() // Contexts are necessary for sending HTTP requests. They dictate behavior for Cancelling, Timing out, and closing channels when a deadline is reached. context.Background is initalizing a blank context.
	client, err := New(ctx, http.URL(authNodeUrl), http.Auth(testToken), http.Provider(http.Github)) // Initialize a new client with the previously created test token, context, and url
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newTestUnsecureClient() (*Client, error) { // A new test client to be used when Certification is not required for the request to go through
	ctx := context.Background()
	client, err := New(ctx, http.URL(authNodeUrl), http.Auth(testToken), http.Provider(http.Github), http.Unsecure())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestConnectionToProdNodeWithoutCheckingCertificates(t *testing.T) { // A simple request to the production server using a valid auth token but an unsecure client
	t.Skip("test needs to be redone")
	t.Run("Given an Unsecure Client with a valid token", func(t *testing.T) {
		client, err := newTestUnsecureClient()
		assert.NilError(t, err)

		t.Run("Getting /me", func(t *testing.T) { // path '/me' is a GET request for the current user
			me := client.User()
			_, err := me.Get()
			assert.NilError(t, err)
		})
	})
}

func TestConnectionToProdNode(t *testing.T) {
	t.Skip("tests need to be redone")
	t.Run("Given a Client with a valid token", func(t *testing.T) {
		client, err := newTestClient()
		assert.NilError(t, err)

		t.Run("Getting /me", func(t *testing.T) {
			me := client.User()
			_, err := me.Get()
			assert.NilError(t, err)
		})
	})
}

// I'd reccomend DRYing the code by having a ''newTestClient(secure)'' function that takes in a boolean 'secure' value to indicate whether or not certification is required or not. 
// We can then have a ''TestConnectionToProdNode(t *testing.T, secure) function utilizing the same boolean effectively reducing the total amountt of code by half.
