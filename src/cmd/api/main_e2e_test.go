//+build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	authDomain "github.com/AngelVlc/todos/internal/api/auth/domain"
	authInfra "github.com/AngelVlc/todos/internal/api/auth/infrastructure"
	listsDomain "github.com/AngelVlc/todos/internal/api/lists/domain"
	"github.com/stretchr/testify/require"
)

func TestEndtoEnd(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	require.NotNil(t, baseURL)
	require.NotEmpty(t, baseURL)

	adminPass := os.Getenv("ADMIN_PASSWORD")
	require.NotNil(t, adminPass)
	require.NotEmpty(t, adminPass)

	client := &http.Client{}

	// Login
	loginBody := fmt.Sprintf("{\"username\": \"admin\",\"password\": \"%v\"}", adminPass)
	req := createRequest(t, "POST", baseURL+"/auth/login", strings.NewReader(loginBody), nil)
	req.Header.Set("Content-type", "application/json")
	res, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	loginRes := authDomain.LoginResponse{}
	err = objFromRes(res.Body, &loginRes)
	require.Nil(t, err)

	loginResCookies := res.Cookies()

	// Get users
	req = createRequest(t, "GET", baseURL+"/users", nil, loginResCookies)
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	usersRes := []authInfra.UserResponse{}
	err = objFromRes(res.Body, &usersRes)
	require.Nil(t, err)

	// Get user with the first id
	req = createRequest(t, "GET", fmt.Sprintf("%v/users/%v", baseURL, usersRes[0].ID), nil, loginResCookies)
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	userRes := authInfra.UserResponse{}
	err = objFromRes(res.Body, &userRes)
	require.Nil(t, err)

	// Creates a list
	listName := fmt.Sprintf("test %v", time.Now().Format("2006-01-02T15:04:05-0700"))
	listBody := fmt.Sprintf("{\"name\": \"%v\"}", listName)
	req = createRequest(t, "POST", baseURL+"/lists", strings.NewReader(listBody), loginResCookies)
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	createdRes := listsDomain.List{}
	err = objFromRes(res.Body, &createdRes)
	require.Nil(t, err)
	listID := fmt.Sprint(createdRes.ID)

	// Refreshes the token
	req = createRequest(t, "POST", baseURL+"/auth/refreshtoken", nil, loginResCookies)
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, 1, len(res.Cookies()))

	// Removes a list
	req = createRequest(t, "DELETE", baseURL+"/lists/"+string(listID), nil, loginResCookies)
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 204, res.StatusCode)
}

func bufferFromBody(body interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if body == nil {
		return buf, nil
	}

	err := json.NewEncoder(buf).Encode(body)
	return buf, err
}

func objFromRes(resBody io.Reader, obj interface{}) error {
	return json.NewDecoder(resBody).Decode(obj)
}

func createRequest(t *testing.T, method, url string, body io.Reader, cookies []*http.Cookie) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.Nil(t, err)

	if cookies != nil {
		for _, v := range cookies {
			req.AddCookie(v)
		}
	}

	return req
}
