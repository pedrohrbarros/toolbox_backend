package secret

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/middleware/error"
)

// @Summary Secret Generator
// @Description Generate secret based in the params
// @Tags Secret
// @Accept json
// @Produce json
// @Param request body secret.GenerateSecret.Request true "Lenght of the secret that'll be generated"
// @Success 200 {string} string "sl5=wv_X/OK/"
// @Failure 400 {object} error.ApiError
// @Failure 500 {object} error.ApiError
// @Router /secret-generator [post]
func GenerateSecret(c *gin.Context) {
	type Request struct {
		SpecialCharacters   bool `json:"special_characters"`
		UpperCaseCharacters bool `json:"uppercase_characters"`
		LowCaseCharacters   bool `json:"lowcase_characters"`
		Numbers             bool `json:"numbers"`
		Length              int  `json:"length"`
	}

	var request Request

	if err := c.ShouldBindJSON(&request); err != nil {
		api_error := error.NewBadRequestError(fmt.Sprintf("Incorrect request=%s\n", err.Error()))
		c.JSON(api_error.Code, api_error)
		return
	}

	const (
		lowcase_characters   = "abcdefghijklmnopqrstuvwxyz"
		uppercase_characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers              = "0123456789"
		special_characters   = "!@#$%^&*()-_=+[]{}<>?/|"
	)
	var all_characters string

	if request.SpecialCharacters {
		all_characters += special_characters
	}

	if request.UpperCaseCharacters {
		all_characters += uppercase_characters
	}

	if request.LowCaseCharacters {
		all_characters += lowcase_characters
	}

	if request.Numbers {
		all_characters += numbers
	}

	if (!request.Numbers && !request.LowCaseCharacters && !request.UpperCaseCharacters && !request.SpecialCharacters) {
		error_message := error.NewBadRequestError("At least one type of secret must be selected")
		c.JSON(error_message.Code, error_message)
		return
	}

	password := make([]byte, request.Length)

	for i := 0; i < request.Length; i++ {
		random_index, err := rand.Int(rand.Reader, big.NewInt(int64(len(all_characters))))
		if err != nil {
			error_message := error.NewInternalServerError(fmt.Sprintf("Error while selecting characters to compose password: %s", err.Error()))
			c.JSON(error_message.Code, error_message)
			return
		}
		password[i] = all_characters[random_index.Int64()]
	}

	c.JSON(200, string(password))
}
