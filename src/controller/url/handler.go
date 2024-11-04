package url

import (
	"bytes"
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
// @Success 200 {string} Shortened URL
// @Failures 404 {object} httputil.HTTPError
// Failures 500 {object} httputil.HTTPError
// @Router /url [post]
func ShortUrl(c *gin.Context) {
	var url_request struct {
		URL string `json:"url" example:"https://www.google.com"`
	}
	
	if err := c.ShouldBindJSON(&url_request) ; err != nil {
		api_error := error.NewBadRequestError(fmt.Sprintf("Incorrect request=%s\n", err.Error()))
		c.JSON(api_error.Code, api_error)
		return
	}

	api_url := "https://api-ssl.bitly.com/v4/shorten"
	data, _ := json.Marshal(map[string]string{
		"group_guid": os.Getenv("BITLY_GROUP_GUID"),
		"long_url": url_request.URL,
	})

	request_body := bytes.NewBuffer(data)

	request, err := http.NewRequest("POST", api_url, request_body)
	if err != nil {
		error_message := fmt.Sprintf("Failed to create request: %v", err)
		c.JSON(500, error_message)
		return
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BITLY_ACCESS_TOKEN")))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		error_message := fmt.Sprintf("Failed to read response: %v", err)
		c.JSON(500, error_message)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		error_message := fmt.Sprintf("Failed to read binary response: %v", err)
		c.JSON(500, error_message)
		return
	}

	var response_data map[string]interface{}
	err = json.Unmarshal(body, &response_data)
	if err != nil {
		error_message := fmt.Sprintf("Failed to parse response: %v", err)
		c.JSON(500, error_message)
		return
	}

	if link, ok := response_data["link"].(string); ok {
		c.JSON(200, link)
	} else {
		error_message := fmt.Sprintf("Failed to get link from response: %v", response_data)
		c.JSON(500, error_message)
		return
	}
}