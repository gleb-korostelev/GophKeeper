package swagger

const (
	// SwaggerJSON raw
	SwaggerJSON = `{
		   "openapi":"3.0.3",
		   "info":{
			  "title": "%s",
			  "version":"version not set"
		   },
           "servers": [
				{
				  "url": "%s"
				}
			  ],
		   "paths":{
			 %s
		   }
		}`

	// HandlerRaw raw
	HandlerRaw = `
		"%s":{
			 "%s":{
				"summary": "%s",
				%s
				"responses":{
				   "200":{
					  "description":"A successful response.",
						 %s
				   },
				   "default":{
					  "description":"An unexpected error response.",
						 %s
				   }
				},
				%s
				"tags":[
				   "%s"
				]
			 }
      	}`

	// HandlerBrakesRaw raw
	HandlerBrakesRaw = `
		"%s":{
			%s
      	}`

	HandlerRaw1 = `
		 "%s":{
				"summary": "%s",
				%s
				"responses":{
				   "200":{
					  "description":"A successful response.",
						 %s
				   },
				   "default":{
					  "description":"An unexpected error response.",
						 "content": {
						  "application/json": {
							"schema": %s
						  }
						}
				   }
				},
				%s
				"tags":[
				   "%s"
				]
			 }
	`

	HandlerRawJSON = `
		 "%s":{
				"summary": "%s",
				%s
				"responses":{
				   "200":{
					  "description":"A successful response.",
						 "content": {
						  "application/json": {
							"schema": %s
						  }
						}
				   },
				   "default":{
					  "description":"An unexpected error response.",
						"content": {
						  "application/json": {
							"schema": %s
						  }
						}
				   }
				},
				%s
				"tags":[
				   "%s"
				]
			 }
	`

	// URL parameters raws
	parametersRaw = `
		"parameters": [
          %s
        ],`
	parameter = `
		{
			"name": "%s",
			"in": "%s",
			"required": %t,
			"description": "%s",
			"schema": {
				"type": "%s"
			}
			%s
		}`
	parameterArray = `
		,"items": {
              "type": "%s"
            },
            "collectionFormat": "multi"
		`
	parameterBody = `
		"%s": {
			"type": "%s"
		}`
	parameterFormData = `
		{
			"name": "%s",
			"in": "formData",
			"required": %t,
			"description": "%s",
			"type": "%s"
			%s
		}`

	// EmptyObject case no response
	EmptyObject  = `"schema": { "type":"object" }`
	bodySchema   = `"schema": { "$ref": "#/definitions/%s" }`
	fileResponse = `"content": {
              "application/json": {
                "schema": { "type": "string", "format": "binary" }
              }
            }`

	responseProduce = `
		"produces": [
          %s
        ]
	`
	requestConsumes = `
		"consumes": [
          "%s"
        ]
	`
)
