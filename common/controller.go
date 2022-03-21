package common

import (
	"github.com/gin-gonic/gin"
	"github.com/heimbogdan/go-swagger-mock/swagger_v2/models"
	"github.com/sirupsen/logrus"
	"strings"
)

type BaseController struct {
	Path    string
	Methods []Method
}

type Method struct {
	Type    MethodType
	Handler gin.HandlerFunc
}

type MethodType int

const (
	UNKNOWN MethodType = iota
	GET
	PUT
	POST
	DELETE
	PATCH
	OPTIONS
	HEAD
)

func (m MethodType) toString() string {
	switch m {
	case GET:
		return "get"
	case PUT:
		return "put"
	case POST:
		return "post"
	case PATCH:
		return "patch"
	case DELETE:
		return "delete"
	case OPTIONS:
		return "options"
	case HEAD:
		return "head"
	default:
		logrus.Warn("unknown method type")
		return "unknown"
	}
}

func MethodTypeFromString(m string) MethodType {
	switch strings.ToLower(m) {
	case "get":
		return GET
	case "put":
		return PUT
	case "post":
		return POST
	case "patch":
		return PATCH
	case "delete":
		return DELETE
	case "options":
		return OPTIONS
	case "head":
		return HEAD
	default:
		logrus.Warn("unknown method type")
		return UNKNOWN
	}
}

func CreateController(path string, item models.PathItem) BaseController {
	return BaseController{
		Path:    ToGinPath(path),
		Methods: CreateMethods(item),
	}
}

func CreateMethods(item models.PathItem) []Method {
	globalParams := item.Parameters
	methods := make([]Method, 0)
	if nil != item.Get {
		get := item.Get
		m := Method{
			Type:    GET,
			Handler: CreateHandler(get, globalParams),
		}
		methods = append(methods, m)
	}
	return methods
}

func CreateHandler(op *models.Operation, gParams *[]models.Parameter) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//TODO implement
	}
}
