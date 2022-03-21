package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
)

const (
	SwaggerVersion2        = "2.0"
	HostPattern            = "^[^{}/ :\\\\]+(?::\\d+)?$"
	BasePathPattern        = "^/"
	HeaderTypePattern      = "[string|number|integer|boolean|array]"
	VendorExtensionPattern = "^x-"
)

// Regex
const (
	refNameRegex = `.*\/(.*)$`
)

type Swagger struct {
	Version             *string                `json:"swagger,omitempty" yaml:"swagger,omitempty"`
	Info                *Info                  `json:"info,omitempty" yaml:"info,omitempty"`
	Paths               *map[string]PathItem   `json:"paths,omitempty" yaml:"paths,omitempty"`
	Host                *string                `json:"host,omitempty" yaml:"host,omitempty"`
	BasePath            *string                `json:"basePath,omitempty" yaml:"basePath,omitempty"`
	Schemes             *[]string              `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Consumes            *[]string              `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces            *[]string              `json:"produces,omitempty" yaml:"produces,omitempty"`
	Definitions         *map[string]Schema     `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	Parameters          *map[string]Parameter  `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Responses           *map[string]Response   `json:"responses,omitempty" yaml:"responses,omitempty"`
	Security            *[]SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
	SecurityDefinitions *map[string]Security   `json:"securityDefinitions,omitempty" yaml:"securityDefinitions,omitempty"`
	Tags                *[]Tag                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs        *ExternalDocs          `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type Info struct {
	Version        *string  `json:"version,omitempty" yaml:"version,omitempty"`
	Title          *string  `json:"title,omitempty" yaml:"title,omitempty"`
	Description    *string  `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService *string  `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
}

type Contact struct {
	Name  *string `json:"name,omitempty" yaml:"name,omitempty"`
	Url   *string `json:"url,omitempty" yaml:"url,omitempty"`
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`
}

type License struct {
	Name *string `json:"name,omitempty" yaml:"name,omitempty"`
	Url  *string `json:"url,omitempty" yaml:"url,omitempty"`
}

type ExternalDocs struct {
	Url         *string `json:"url,omitempty" yaml:"url,omitempty"`
	Description *string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Operation struct {
	Responses    *map[string]Response   `json:"responses,omitempty" yaml:"responses,omitempty"`
	Tags         *[]string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      *string                `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  *string                `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs          `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationId  *string                `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Consumes     *[]string              `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces     *[]string              `json:"produces,omitempty" yaml:"produces,omitempty"`
	Parameters   *[]Parameter           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Schemes      *[]Schema              `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Deprecated   *bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     *[]SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty"`
}

func (op *Operation) GetQueryParameters() []Parameter {
	return op.getParametersByLocation("query")
}

func (op *Operation) GetPathParameters() []Parameter {
	return op.getParametersByLocation("path")
}

func (op *Operation) GetBodyParameters() []Parameter {
	return op.getParametersByLocation("body")
}

func (op *Operation) GetCookieParameters() []Parameter {
	return op.getParametersByLocation("cookie")
}

func (op *Operation) GetHeaderParameters() []Parameter {
	return op.getParametersByLocation("header")
}

func (op *Operation) GetFormDataParameters() []Parameter {
	return op.getParametersByLocation("formData")
}

func (op *Operation) getParametersByLocation(location string) []Parameter {
	query := make([]Parameter, 0)
	if op.Parameters != nil {
		for _, p := range *op.Parameters {
			if p.In != nil && *p.In == location {
				query = append(query, p)
			}
		}
	}
	return query
}

type PathItem struct {
	Ref        *string      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Get        *Operation   `json:"get,omitempty" yaml:"get,omitempty"`
	Put        *Operation   `json:"put,omitempty" yaml:"put,omitempty"`
	Post       *Operation   `json:"post,omitempty" yaml:"post,omitempty"`
	Delete     *Operation   `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options    *Operation   `json:"options,omitempty" yaml:"options,omitempty"`
	Head       *Operation   `json:"head,omitempty" yaml:"head,omitempty"`
	Patch      *Operation   `json:"patch,omitempty" yaml:"patch,omitempty"`
	Parameters *[]Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func (pi *PathItem) ToShortString() string {
	if pi != nil {
		buffer := bytes.NewBufferString("PathItem:[")
		if pi.Parameters != nil {
			buffer.WriteString("Global:")
			buffer.WriteString(ToShortString(pi.Parameters))
		}
		buffer.WriteString(ToShortStringPrefixed("Get:", GetOperationParameters(pi.Get)))
		buffer.WriteString(ToShortStringPrefixed("Put:", GetOperationParameters(pi.Put)))
		buffer.WriteString(ToShortStringPrefixed("Post:", GetOperationParameters(pi.Post)))
		buffer.WriteString(ToShortStringPrefixed("Delete:", GetOperationParameters(pi.Delete)))
		buffer.WriteString(ToShortStringPrefixed("Patch:", GetOperationParameters(pi.Patch)))
		buffer.WriteString(ToShortStringPrefixed("Head:", GetOperationParameters(pi.Head)))
		buffer.WriteString(ToShortStringPrefixed("Options:", GetOperationParameters(pi.Options)))
		buffer.WriteString("]")
		return buffer.String()
	}
	return ""
}

func GetOperationParameters(op *Operation) *[]Parameter {
	if op != nil {
		return op.Parameters
	}
	return nil
}

type Response struct {
	Description *string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Ref         *string                 `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Schema      *Schema                 `json:"schema,omitempty" yaml:"schema,omitempty"`
	Headers     *map[string]Header      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Examples    *map[string]interface{} `json:"examples,omitempty" yaml:"examples,omitempty"`
}

type Header struct {
	TypeStruct
	Restrictions
	Items            *PrimitivesItems `json:"items,omitempty" yaml:"items,omitempty"`
	CollectionFormat *string          `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	Description      *string          `json:"description,omitempty" yaml:"description,omitempty"`
}

type Parameter struct {
	TypeStruct
	Restrictions
	Ref              *string          `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Name             *string          `json:"name,omitempty" yaml:"name,omitempty"`
	In               *string          `json:"in,omitempty" yaml:"in,omitempty"`
	Description      *string          `json:"description,omitempty" yaml:"description,omitempty"`
	Schema           *Schema          `json:"schema,omitempty" yaml:"schema,omitempty"`
	Required         *bool            `json:"required,omitempty" yaml:"required,omitempty"`
	Items            *PrimitivesItems `json:"items,omitempty" yaml:"items,omitempty"`
	CollectionFormat *string          `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	AllowEmptyValue  *bool            `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
}

func (p *Parameter) GetLocationAndName() string {
	buffer := bytes.NewBufferString("")
	if p.In != nil {
		buffer.WriteString(*p.In)
	}
	buffer.WriteString(":")
	if p.Name != nil {
		buffer.WriteString(*p.Name)
	}
	return buffer.String()
}

func (p *Parameter) ToShortString() string {
	buffer := bytes.NewBufferString("[")
	buffer.WriteString(p.GetLocationAndName())
	buffer.WriteString("(")
	if p.Ref != nil {
		buffer.WriteString(p.GetRefNamePrefixed())
	} else if p.Schema != nil {
		s := *p.Schema
		if s.Ref != nil && len(*s.Ref) > 0 {
			buffer.WriteString(s.GetRefNamePrefixed())
		} else {
			buffer.WriteString(s.GetSchemaTypeAndFormat())
		}
	} else {
		buffer.WriteString(p.GetTypeAndFormat())
	}
	buffer.WriteString(")]")
	return buffer.String()
}

func (p *Parameter) GetRefName() string {
	if p.Ref != nil {
		matcher := regexp.MustCompile(refNameRegex)
		return string(matcher.ReplaceAll([]byte(*p.Ref), []byte(`$1`)))
	}
	return ""
}

func (p *Parameter) GetRefNamePrefixed() string {
	if p.Ref != nil {
		ref := p.GetRefName()
		return fmt.Sprintf("ref:%s", ref)
	}
	return ""
}

func (s *Schema) GetRefName() string {
	if s.Ref != nil {
		matcher := regexp.MustCompile(refNameRegex)
		return string(matcher.ReplaceAll([]byte(*s.Ref), []byte(`$1`)))
	}
	return ""
}

func (s *Schema) GetRefNamePrefixed() string {
	if s.Ref != nil {
		ref := s.GetRefName()
		return fmt.Sprintf("ref:%s", ref)
	}
	return ""
}

func (s *Schema) GetSchemaTypeAndFormat() string {
	buffer := bytes.NewBufferString("")
	if s.Ref != nil && len(*s.Ref) > 0 {
		buffer.WriteString(s.GetRefNamePrefixed())
	} else {
		if s.Type != nil {
			buffer.WriteString(*s.Type)
		}
		buffer.WriteString(":")
		if s.Type != nil && *s.Type == "array" {
			items := s.GetItems()
			if items != nil && len(items) > 0 {
				buffer.WriteString("[")
				for _, item := range items {
					buffer.WriteString("(")
					buffer.WriteString(item.GetSchemaTypeAndFormat())
					buffer.WriteString(")")
				}
				buffer.WriteString("]")
			}
		} else if s.Format != nil {
			buffer.WriteString(*s.Format)
		}
	}
	return buffer.String()
}

func ToShortString(arr *[]Parameter) string {
	if arr != nil {
		buffer := bytes.NewBufferString("Parameters:[")
		for i, p := range *arr {
			buffer.WriteString(p.ToShortString())
			if i < len(*arr)-1 {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString("]")
		return buffer.String()
	}
	return ""
}

func ToShortStringPrefixed(prefix string, arr *[]Parameter) string {
	if arr != nil && len(*arr) > 0 {
		buffer := bytes.NewBufferString(prefix)
		buffer.WriteString(ToShortString(arr))
		return buffer.String()
	}
	return ""
}

type Schema struct {
	Restrictions
	FileSchema
	Ref                  *string            `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	MaxProperties        *int               `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        *int               `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	AdditionalProperties interface{}        `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"` //TODO check if is Schema or boolean
	Items                interface{}        `json:"items,omitempty" yaml:"items,omitempty"`                               // TODO will check if array of Schema
	AllOf                *[]Schema          `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	Properties           *map[string]Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Discriminator        *string            `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	Xml                  *Xml               `json:"xml,omitempty" yaml:"xml,omitempty"`
}

func (s *Schema) GetItems() (arr []Schema) {
	if s.Items != nil {
		data, err := ToString(s.Items)
		if err != nil {
			logrus.Error(err)
		} else {
			switch s.Items.(type) {
			case map[string]interface{}:
				sObj := Schema{}
				err = json.Unmarshal([]byte(data), &sObj)
				if err != nil {
					logrus.Error(err)
				} else {
					return append(arr, sObj)
				}
			case []map[string]interface{}:
				var sObj []Schema
				err = json.Unmarshal([]byte(data), &sObj)
				if err != nil {
					logrus.Error(err)
				} else {
					return append(arr, sObj...)
				}
			default:
				return nil
			}

		}
	}
	return nil
}

func ToString[T any](i T) (string, error) {
	data, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type FileSchema struct {
	TypeStruct
	Title        *string       `json:"title,omitempty" yaml:"title,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	Required     *[]string     `json:"required,omitempty" yaml:"required,omitempty"`
	ReadOnly     *bool         `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Enum         *[]string     `json:"enum,omitempty" yaml:"enum,omitempty"`
}

type PrimitivesItems struct {
	TypeStruct
	Restrictions
	Items            *PrimitivesItems `json:"items,omitempty" yaml:"items,omitempty"`
	CollectionFormat *string          `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
}

type SecurityRequirement map[string][]string

type Xml struct {
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace *string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Prefix    *string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Attribute *string `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Wrapped   *bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

type Tag struct {
	Name         *string       `json:"name,omitempty" yaml:"name,omitempty"`
	Description  *string       `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

type Security struct {
	TypeStruct
	Name             *string            `json:"name,omitempty" yaml:"name,omitempty"`
	In               *string            `json:"in,omitempty" yaml:"in,omitempty"`
	Flow             *string            `json:"flow,omitempty" yaml:"flow,omitempty"`
	Scopes           *map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Description      *string            `json:"description,omitempty" yaml:"description,omitempty"`
	AuthorizationUrl *string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenUrl         *string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
}

type TypeStruct struct {
	Type    *string     `json:"type,omitempty" yaml:"type,omitempty"`
	Format  *string     `json:"format,omitempty" yaml:"format,omitempty"`
	Default interface{} `json:"default,omitempty" yaml:"default,omitempty"`
}

func (t *TypeStruct) GetTypeAndFormat() string {
	buffer := bytes.NewBufferString("")
	if t.Type != nil {
		buffer.WriteString(*t.Type)
	}
	buffer.WriteString(":")
	if t.Format != nil {
		buffer.WriteString(*t.Format)
	}
	return buffer.String()
}

type Restrictions struct {
	Maximum          *int      `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum *bool     `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum          *int      `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum *bool     `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	MaxLength        *int      `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength        *int      `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern          *string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MaxItems         *int      `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems         *int      `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems      *bool     `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Enum             *[]string `json:"enum,omitempty" yaml:"enum,omitempty"`
	MultipleOf       *int      `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
}
