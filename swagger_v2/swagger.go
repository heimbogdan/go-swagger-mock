package swagger_v2

import (
	_ "embed"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema.json
var swaggerSchemaJson string

var swaggerSchemaLoader = gojsonschema.NewStringLoader(swaggerSchemaJson)

var SwaggerV2Schema, _ = gojsonschema.NewSchema(swaggerSchemaLoader)
