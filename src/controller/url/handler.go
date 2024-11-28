package url

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/middleware/error"
)

// @Summary URL Shortener
// @Description Shorten a URL using Bitly API
// @Tags URL
// @Accept json
// @Produce json
// @Param url body string true "URL to shorten"
// @Success 200 {string} string "https://bit.ly/example"
// @Failure 400 {object} error.ApiError
// @Failure 500 {object} error.ApiError
// @Router /url [post]
func ShortUrl(c *gin.Context) {
	var url_request struct {
		URL string `json:"url" example:"https://www.google.com" binding:"required"`
	}

	if err := c.ShouldBindJSON(&url_request); err != nil {
		api_error := error.NewBadRequestError(fmt.Sprintf("Incorrect request=%s\n", err.Error()))
		c.JSON(api_error.Code, api_error)
		return
	}

	group_guid := cmp.Or(os.Getenv("BITLY_GROUP_GUID"), "default")
	bitly_access_token := cmp.Or(os.Getenv("BITLY_ACCESS_TOKEN"), "")

	api_url := "https://api-ssl.bitly.com/v4/shorten"
	data, _ := json.Marshal(map[string]string{
		"group_guid": group_guid,
		"long_url":   url_request.URL,
	})

	request_body := bytes.NewBuffer(data)

	request, err := http.NewRequest("POST", api_url, request_body)
	if err != nil {
		api_error := error.NewInternalServerError(fmt.Sprintf("Failed to create request: %v", err))
		c.JSON(api_error.Code, api_error)
		return
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bitly_access_token))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		api_error := error.NewInternalServerError(fmt.Sprintf("Failed to send request: %v", err))
		c.JSON(api_error.Code, api_error)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		api_error := error.NewInternalServerError(fmt.Sprintf("Failed to read binary response: %v", err))
		c.JSON(api_error.Code, api_error)
		return
	}

	var response_data map[string]interface{}
	err = json.Unmarshal(body, &response_data)
	if err != nil {
		api_error := error.NewInternalServerError(fmt.Sprintf("Failed to parse response: %v", err))
		c.JSON(api_error.Code, api_error)
		return
	}

	if link, ok := response_data["link"].(string); ok {
		c.JSON(200, link)
	} else {
		for key, value := range response_data {
			fmt.Printf("%s: %d\n", key, value)
		}
		api_error := error.NewInternalServerError(fmt.Sprintf("Failed to get link from response: %v", response_data))
		c.JSON(api_error.Code, api_error)
		return
	}
}
