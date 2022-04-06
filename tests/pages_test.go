package tests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

//func TestHomePage(t *testing.T) {
//	baseURL := "http://localhost:8088"
//
//	var (
//		resp *http.Response
//		err error
//	)
//
//	resp, err = http.Get(baseURL + "/")
//
//	assert.NoError(t, err, "有错误，错误不为空")
//	assert.Equal(t, 200, resp.StatusCode, "返回码错误，应为200")
//}
//
//func TestAboutPage(t *testing.T) {
//	baseURL := "http://localhost:8088"
//
//	var (
//		resp *http.Response
//		err error
//	)
//
//	resp, err = http.Get(baseURL + "/about")
//
//	assert.NoError(t, err, "有错误，错误不为空")
//	assert.Equal(t, 200, resp.StatusCode, "返回码错误，应为200")
//}

func TestAllPages(t *testing.T) {
	baseURL := "http://localhost:8088"

	var tests = []struct{
		method string
		url string
		expected int
	}{
		{"GET", "/", 200},
		{"GET", "/about", 200},
		{"GET", "/notfound", 404},
		{"GET", "/articles", 200},
		{"GET", "/articles/create", 200},
		{"GET", "/articles/3", 200},
		{"GET", "/articles/3/edit", 200},
		{"POST", "/articles/3", 200},
		{"POST", "/articles", 200},
		{"POST", "/articles/1/delete", 404},
	}

	for _, test := range tests {
		t.Logf("当前请求URL: %v\n", test.url)

		var (
			resp *http.Response
			err error
		)

		switch {
		case test.method == "POST":
			data := make(map[string][]string)
			resp, err = http.PostForm(baseURL + test.url, data)
		default:
			resp, err = http.Get(baseURL + test.url)
		}

		assert.NoError(t, err, "存在错误")
		assert.Equal(t, test.expected, resp.StatusCode, "返回码与预期" + strconv.Itoa(test.expected) + "不符，返回：" + strconv.Itoa(resp.StatusCode))
	}
}