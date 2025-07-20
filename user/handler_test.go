package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	userpkg "local/chat/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	conn := setup(ctx, t)
	defer clean(ctx, t, conn)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	uh := userpkg.NewUserHandler(ctx, conn)
	uh.BindRoutes(router)

	pa := gin.H{
		"email":    "user1@gmail.com",
		"username": "USER1",
		"password": "password",
	}
	body, _ := json.Marshal(pa)
	req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, 201)

	var res gin.H
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Errorf("invalid JSON response")
	}

	user := res["user"].(map[string]any)
	t.Log(user)
	assert.Equal(t, user["id"], 1.0)
	assert.Equal(t, user["username"], pa["username"])
	assert.Equal(t, user["email"], pa["email"])
	assert.NotEmpty(t, user["createdAt"])
}
