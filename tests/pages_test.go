package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHomePage(t *testing.T) {
	baseURL := "http://localhost:8088"

	var (
		resp *http.Response
		err error
	)

	resp, err = http.Get(baseURL + "/")

	assert.NoError(t, err, "有错误，错误不为空")
	assert.Equal(t, 200, resp.StatusCode, "返回码错误，应为200")
}
