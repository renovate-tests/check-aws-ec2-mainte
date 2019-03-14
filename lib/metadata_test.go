package checkawsec2mainte

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func initTestServer(path string, resp string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != path {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Write([]byte(resp))
	}))
}

func TestGetInstanceId(t *testing.T) {
	server := initTestServer(
		"/latest/meta-data/instance-id",
		"i-09e032cce9ef71d84",
	)
	defer server.Close()

	cfg, _ := external.LoadDefaultAWSConfig()
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL + "/latest")

	fmt.Println(server.URL)
	fmt.Println(getInstanceIdFromMetadata(cfg))
}
