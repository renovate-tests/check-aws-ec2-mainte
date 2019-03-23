package unit

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoErrorTestServer(t *testing.T) {

	expected := map[string]string{
		"/foo/bar":        "FooBar",
		"/fizz/buzz/piyo": "FizzBuzz",
	}

	server := StartTestServer(expected)
	defer server.Close()

	for uri, body := range expected {
		res, err := http.Get(server.URL + uri)
		defer res.Body.Close()

		actual, _ := ioutil.ReadAll(res.Body)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, body, string(actual))
	}
}

func TestNoErrorEventsMetadata(t *testing.T) {
	CreateEventsMetadata(t, "i-0dc818ea941b1ae18")
}
