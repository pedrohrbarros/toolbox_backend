package secret

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/pedrohrbarros/toolbox_backend/src/middleware/error"
)

func GenerateSecret(c *gin.Context) {
	type Request struct {
		SpecialCharacters bool `json:"special_characters"`
		UpperCaseCharacters bool `json:"uppercase_characters"`
		Letters bool `json:"letters"`
		Numbers bool `json:"numbers"`
		Length int `json:"length"`
	}

	var request Request

	if err := c.ShouldBindJSON(&request) ; err != nil {
		api_error := error.NewBadRequestError(fmt.Sprintf("Incorrect request=%s\n", err.Error()))
		c.JSON(api_error.Code, api_error)
		return
	}

	const (
		letters        = "abcdefghijklmnopqrstuvwxyz"
		uppercase_letters   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers        = "0123456789"
		special_characters   = "!@#$%^&*()-_=+[]{}<>?/|"
	)
	var all_characters string

	if request.SpecialCharacters {
		all_characters += special_characters
	}

	if request.UpperCaseCharacters {
		all_characters += uppercase_letters
	}

	if request.Letters {
		all_characters += letters
	}

	if request.Numbers {
		all_characters += numbers
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