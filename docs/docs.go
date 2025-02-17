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
        "/jrp": {
            "get": {
                "description": "returns a randomly generated Japanese phrase.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "jrp"
                ],
                "summary": "get a random Japanese phrase.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_yanosea_jrp_v2_app_presentation_api_jrp-server_formatter.JrpJsonOutputDto"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_yanosea_jrp_v2_app_presentation_api_jrp-server_formatter.JrpJsonOutputDto": {
            "description": "response format for jrp",
            "type": "object",
            "properties": {
                "phrase": {
                    "description": "@Description Generated Japanese phrase",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "JRP API",
	Description:      "jrp api server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
