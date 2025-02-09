// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/file/converter": {
            "post": {
                "description": "Convert a word file into pdf",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "multipart/form-data"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Document converter",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File that will be converted",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Expected type that will be converted",
                        "name": "expected_type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    }
                }
            }
        },
        "/file/image/resizer": {
            "post": {
                "description": "Edit an image based on the parameters in the request",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "image/jpeg"
                ],
                "tags": [
                    "File"
                ],
                "summary": "Image Editor",
                "parameters": [
                    {
                        "type": "file",
                        "description": "JPEG image file to edit (only accepts .jpeg files)",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Image width",
                        "name": "width",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Image height",
                        "name": "height",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Converted file",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    }
                }
            }
        },
        "/secret/generator": {
            "post": {
                "description": "Generate secret based in the params",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Secret"
                ],
                "summary": "Secret Generator",
                "parameters": [
                    {
                        "description": "Lenght of the secret that'll be generated",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/secret.GenerateSecret.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "sl5=wv_X/OK/",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    }
                }
            }
        },
        "/url/shortener": {
            "post": {
                "description": "Shorten a URL using Bitly API",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "URL"
                ],
                "summary": "URL Shortener",
                "parameters": [
                    {
                        "description": "URL to shorten",
                        "name": "url",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "https://bit.ly/example",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/error.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "error.ApiError": {
            "type": "object",
            "properties": {
                "causes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/error.Causes"
                    }
                },
                "code": {
                    "type": "integer"
                },
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "error.Causes": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "secret.GenerateSecret.Request": {
            "type": "object",
            "properties": {
                "length": {
                    "type": "integer"
                },
                "lowcase_characters": {
                    "type": "boolean"
                },
                "numbers": {
                    "type": "boolean"
                },
                "special_characters": {
                    "type": "boolean"
                },
                "uppercase_characters": {
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
