package main

import (
	"encoding/json"
	"fmt"
	"github.com/heimbogdan/go-swagger-mock/common"
	v2 "github.com/heimbogdan/go-swagger-mock/swagger_v2"
	v2m "github.com/heimbogdan/go-swagger-mock/swagger_v2/models"
	"github.com/xeipuuv/gojsonschema"
)

func main() {
	petStore := gojsonschema.NewReferenceLoader("https://petstore.swagger.io/v2/swagger.json")
	validate, err := v2.SwaggerV2Schema.Validate(petStore)
	if err != nil {
		panic(err.Error())
	}
	if validate.Valid() {
		fmt.Println("valid")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range validate.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	petJson, _ := petStore.LoadJSON()
	swagg := v2m.Swagger{}
	data, _ := common.ToString(petJson)
	err = json.Unmarshal([]byte(data), &swagg)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(common.ToString(swagg))
	fmt.Println(common.ToString(swagg.Info))

	p1 := common.ToGinPath("/test/{prop1}/ddd-as/{prop2}/231_asd/{prop3}")
	fmt.Println(p1)
	p2 := common.ToSwaggerPath(p1)
	fmt.Println(p2)

	for k, v := range *swagg.Paths {
		fmt.Println(k, v2m.ToShortString(v.Parameters))
		fmt.Println(v.ToShortString())
	}
}
