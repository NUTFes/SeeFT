// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/users": {
            "get": {
                tags: ["user"],
                "description": "usersの一覧を取得",
                "responses": {
                    "200": {
                        "description": "usersの一覧の取得",
                    }
                }
            },
        },
        "/users/{grade_id}/{department_id}/{bureau_id}/{role_id}": {
            "post": {
                tags: ["user"],
                "description": "userをcreate",
                "responses": {
                    "200": {
                        "description": "createしたuserが返ってくる",
                    }
                },
                "parameters": [
                    {
                        "name": "name",
                        "in": "query",
                        "description": "name",
                        "type": "string",
                        "required": true
                    },
                    {
                        "name": "mail",
                        "in": "query",
                        "description": "mail",
                        "type": "string",
                        "required": true
                    },
                    {
                        "name": "student_number",
                        "in": "query",
                        "description": "student_number",
                        "type": "integer",
                        "required": true,
                    },
                    {
                        "name": "grade_id",
                        "in": "path",
                        "description": "grade_id",
                        "type": "integer",
                        "required": true
                    },
                    {
                        "name": "bureau_id",
                        "in": "path",
                        "description": "bureau_id",
                        "type": "integer",
                        "required": true
                    },
                    {
                        "name": "department_id",
                        "in": "path",
                        "description": "department_id",
                        "type": "integer",
                        "required": true
                    },
                    {
                        "name": "role_id",
                        "in": "path",
                        "description": "role_id",
                        "type": "integer",
                        "required": true,
                    },
                    {
                        "name": "tel",
                        "in": "query",
                        "description": "tel",
                        "type": "string",
                        "required": true,
                    },
                    {
                        "name": "password",
                        "in": "query",
                        "description": "password",
                        "type": "string",
                        "required": true,
                    }
                ]
            }
        },
        "/grades": {
            "get": {
                tags: ["grade"],
                "description": "gradesの一覧を取得",
                "responses": {
                    "200": {
                        "description": "gradesの一覧の取得",
                    }
                }
            }
        },
        "/bureaus": {
            "get": {
                tags: ["bureau"],
                "description": "bureausの一覧を取得",
                "responses": {
                    "200": {
                        "description": "bureausの一覧を取得",
                    }
                }
            }
        },
        "/shifts": {
            "get": {
                tags: ["shift"],
                "description": "shiftsの一覧を取得",
                "responses": {
                    "200": {
                        "description": "shiftsの一覧を取得",
                    }
                }
            }
        },
        "/tasks": {
            "get": {
                tags: ["task"],
                "description": "tasksの一覧を取得",
                "responses": {
                    "200": {
                        "description": "tasksの一覧を取得",
                    }
                }
            }
        },
        "/times": {
            "get": {
                tags: ["time"],
                "description": "timesの一覧を取得",
                "responses": {
                    "200": {
                        "description": "timesの一覧を取得",
                    }
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
