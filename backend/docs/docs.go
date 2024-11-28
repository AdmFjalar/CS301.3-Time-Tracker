package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/authentication/token": {
            "post": {
                "description": "Creates a token for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Creates a token",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.CreateUserTokenPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Token",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/authentication/user": {
            "post": {
                "description": "Registers a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authentication"
                ],
                "summary": "Registers a user",
                "parameters": [
                    {
                        "description": "User credentials",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.RegisterUserPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User registered",
                        "schema": {
                            "$ref": "#/definitions/main.UserWithToken"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Healthcheck endpoint",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ops"
                ],
                "summary": "Healthcheck",
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/users/activate/{token}": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Activates/Register a user by invitation token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Activates/Register a user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Invitation token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User activated",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/users/feed": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetches the user feed",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "feed"
                ],
                "summary": "Fetches the user feed",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Since",
                        "name": "since",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Until",
                        "name": "until",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Search",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/store.Timestamp"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/users/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetches a user profile by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Fetches a user profile",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/users/{userID}/follow": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Follows a user by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Follows a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User followed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "User payload missing",
                        "schema": {}
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {}
                    }
                }
            }
        },
        "/users/{userID}/unfollow": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Unfollow a user by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Unfollow a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "User unfollowed",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "User payload missing",
                        "schema": {}
                    },
                    "404": {
                        "description": "User not found",
                        "schema": {}
                    }
                }
            }
        },
        "/timestamps": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Creates a timestamp for a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timestamps"
                ],
                "summary": "Creates a timestamp",
                "parameters": [
                    {
                        "description": "Timestamp information",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.CreateTimestampPayload"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Timestamp created",
                        "schema": {
                            "$ref": "#/definitions/store.Timestamp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/timestamps/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetches a timestamp by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timestamps"
                ],
                "summary": "Fetches a timestamp",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Timestamp ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Timestamp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Updates a timestamp by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timestamps"
                ],
                "summary": "Updates a timestamp",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Timestamp ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated timestamp information",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.UpdateTimestampPayload"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Timestamp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Deletes a timestamp by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timestamps"
                ],
                "summary": "Deletes a timestamp",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Timestamp ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Timestamp deleted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/timestamps/latest": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetches the most recent timestamp for a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "timestamps"
                ],
                "summary": "Fetches the latest timestamp",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Timestamp"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/shifts": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetches finished shifts for a user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shifts"
                ],
                "summary": "Fetches finished shifts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Shift"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/shifts/{userID}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Fetches finished shifts for a specific user by their ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "shifts"
                ],
                "summary": "Fetches finished shifts by user ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "userID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/store.Shift"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "main.CreateUserTokenPayload": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                },
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 3
                }
            }
        },
        "main.RegisterUserPayload": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 255
                },
                "password": {
                    "type": "string",
                    "maxLength": 72,
                    "minLength": 3
                }
            }
        },
        "main.UserWithToken": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "main.CreateTimestampPayload": {
            "type": "object",
            "required": [
                "stamp_type",
                "stamp_time"
            ],
            "properties": {
                "stamp_type": {
                    "type": "string"
                },
                "stamp_time": {
                    "type": "string"
                }
            }
        },
        "main.UpdateTimestampPayload": {
            "type": "object",
            "required": [
                "stamp_type",
                "stamp_time"
            ],
            "properties": {
                "stamp_type": {
                    "type": "string"
                },
                "stamp_time": {
                    "type": "string"
                }
            }
        },
        "store.User": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                }
            }
        },
        "store.Timestamp": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "stamp_type": {
                    "type": "string"
                },
                "stamp_time": {
                    "type": "string"
                }
            }
        },
        "store.Shift": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "start_time": {
                    "type": "string"
                },
                "end_time": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Thyme Flies API",
	Description:      "API for Thyme Flies, a time-tracking application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
