//+build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/stretchr/testify/require"
)

func TestEndtoEnd(t *testing.T) {
	baseURL := os.Getenv("BASE_URL")
	require.NotNil(t, baseURL)

	adminPass := os.Getenv("ADMIN_PASSWORD")
	require.NotNil(t, adminPass)

	client := &http.Client{}

	tokenDtoBody, _ := bufferFromBody(dtos.TokenDto{UserName: "admin", Password: adminPass})
	req := createRequest(t, "POST", baseURL+"/auth/token", tokenDtoBody)
	req.Header.Set("Content-type", "application/json")
	res, err := client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	tokenRes := dtos.TokenResponseDto{}
	err = objFromRes(res.Body, &tokenRes)
	require.Nil(t, err)

	authHeaderContent := fmt.Sprintf("Bearer %v", tokenRes.Token)

	req = createRequest(t, "GET", baseURL+"/users/1", nil)
	req.Header.Set("Authorization", authHeaderContent)
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)
	userRes := dtos.UserResponseDto{}
	err = objFromRes(res.Body, &userRes)
	require.Nil(t, err)
	require.True(t, userRes.IsAdmin)

	listName := fmt.Sprintf("test %v", time.Now().Format("2006-01-02T15:04:05-0700"))
	listDtoBody, _ := bufferFromBody(dtos.ListDto{Name: listName})
	req = createRequest(t, "POST", baseURL+"/lists", listDtoBody)
	req.Header.Set("Authorization", authHeaderContent)
	req.Header.Set("Content-type", "application/json")
	res, err = client.Do(req)
	require.Nil(t, err)
	require.Equal(t, 201, res.StatusCode)
	createdRes := int32(0)
	err = objFromRes(res.Body, &createdRes)
	require.Nil(t, err)
	listID := fmt.Sprint(createdRes)

	req = createRequest(t, "DELETE", baseURL+"/lists/"+string(listID), nil)
	req.Header.Set("Authorization", authHeaderContent)
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
	err := json.NewDecoder(resBody).Decode(obj)
	return err
}

func createRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	require.Nil(t, err)

	return req
}
