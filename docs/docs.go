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
        "/admin/events/visit": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "event"
                ],
                "summary": "admin visit event",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.AdminVisitEventInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/events/archived": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "event"
                ],
                "summary": "get archived events",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/usecase.ClosedEventsOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/events/upcoming": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "event"
                ],
                "summary": "get upcoming events",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.UpcomingEventInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/usecase.UpcomingEventList"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/events/visit": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "event"
                ],
                "summary": "visit event",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.VisitEventInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/llm": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "llm"
                ],
                "summary": "chat with llm",
                "parameters": [
                    {
                        "description": "input",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.LLMChatInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ai.StreamResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/login": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "login with telegram",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/usecase.LoginOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/users/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "get user me",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/usecase.UserMeOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "ai.StreamResponse": {
            "type": "object",
            "properties": {
                "choices": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "finish_reason": {
                                "type": "string"
                            },
                            "index": {
                                "type": "integer"
                            },
                            "logprobs": {
                                "type": "string"
                            },
                            "stop_reason": {
                                "type": "string"
                            },
                            "text": {
                                "type": "string"
                            }
                        }
                    }
                },
                "created": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "object": {
                    "type": "string"
                },
                "usage": {
                    "type": "string"
                }
            }
        },
        "entities.AchievementType": {
            "type": "object",
            "properties": {
                "achievementTypeID": {
                    "type": "integer"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Event"
                    }
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "entities.Company": {
            "type": "object",
            "properties": {
                "company": {
                    "type": "string"
                },
                "companyID": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Event"
                    }
                },
                "logo": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Product"
                    }
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "entities.Event": {
            "type": "object",
            "properties": {
                "achievementType": {
                    "$ref": "#/definitions/entities.AchievementType"
                },
                "achievementTypeID": {
                    "type": "integer"
                },
                "coin": {
                    "type": "number"
                },
                "company": {
                    "$ref": "#/definitions/entities.Company"
                },
                "companyID": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "endDs": {
                    "type": "string"
                },
                "erpID": {
                    "type": "integer"
                },
                "eventID": {
                    "type": "integer"
                },
                "eventName": {
                    "type": "string"
                },
                "eventType": {
                    "type": "string"
                },
                "maxUsers": {
                    "type": "integer"
                },
                "startDs": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "entities.Product": {
            "type": "object",
            "properties": {
                "avalibility": {
                    "type": "integer"
                },
                "company": {
                    "$ref": "#/definitions/entities.Company"
                },
                "companyID": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "productID": {
                    "type": "integer"
                }
            }
        },
        "usecase.AdminVisitEventInput": {
            "type": "object",
            "properties": {
                "achievement_type_id": {
                    "type": "integer"
                },
                "eventID": {
                    "type": "integer"
                },
                "userID": {
                    "type": "integer"
                }
            }
        },
        "usecase.ClosedEventsOutput": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Event"
                    }
                }
            }
        },
        "usecase.LLMChatInput": {
            "type": "object",
            "properties": {
                "promnt": {
                    "type": "string"
                }
            }
        },
        "usecase.LoginOutput": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "usecase.UpcomingEventInput": {
            "type": "object",
            "properties": {
                "period": {
                    "type": "string",
                    "enum": [
                        "today",
                        "tomorrow",
                        "week",
                        "month"
                    ]
                }
            }
        },
        "usecase.UpcomingEventList": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Event"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "usecase.UserMeOutput": {
            "type": "object",
            "properties": {
                "birth_date": {
                    "type": "string"
                },
                "coin": {
                    "type": "number"
                },
                "email": {
                    "type": "string"
                },
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Event"
                    }
                },
                "l_surname": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "usecase.VisitEventInput": {
            "type": "object",
            "properties": {
                "eventID": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "B8boost API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
