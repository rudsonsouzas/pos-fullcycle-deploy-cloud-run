package httpclient_test

import (
	"testing"
	"time"

	httpclient "api-server/pkg/http_client"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient(t *testing.T) {
	client := httpclient.NewHTTPClient(time.Minute)
	assert.NotNil(t, client)
}
