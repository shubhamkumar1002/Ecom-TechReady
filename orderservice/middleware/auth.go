package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func ValidateTokenMiddleware(ctx iris.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{"error": "Authorization header required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{"error": "Invalid Authorization header format"})
		return
	}
	token := parts[1]

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service.default.svc.cluster.local/auth/validate"
	}

	requestBody, err := json.Marshal(map[string]string{
		"access_token": token,
	})

	if err != nil {
		log.Printf("Error creating auth request body: %v", err)
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}

	resp, err := http.Post(authServiceURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error contacting auth service: %v", err)
		ctx.StopWithJSON(iris.StatusServiceUnavailable, iris.Map{"error": "Could not validate token"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Token validation successful")
		ctx.Next()
	} else {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("Failed to read error response body from auth service: %v", readErr)
			// Send a generic error if we can't even read the response.
			ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{"error": "Invalid token (and failed to parse auth response)"})
			return
		}

		var errorResponse map[string]interface{}
		json.Unmarshal(bodyBytes, &errorResponse)

		//ctx.StopWithJSON(resp.StatusCode, errorResponse)

		fmt.Println("Token validation failed")
		ctx.StopWithJSON(iris.StatusUnauthorized, iris.Map{"error": errorResponse["message"]})
	}
}
