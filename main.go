package main

import (
	"net/url"
	"strings"
	"net/http"
	"io/ioutil"
	"fmt"

	"github.com/tkanos/gonfig"
	"github.com/gin-gonic/gin"
)

var conf Configuration

type Configuration struct {
	API_URL string
	RELATIVE_PATH string
	PORT string
}

type Context struct {
	MethodApi   string
	QueryParams url.Values
}

func (context *Context) SetMethodApi(value string) {
	context.MethodApi = value
}

func (context *Context) SetQueryParams(query url.Values) {
	context.QueryParams = query
}

func (context *Context) GetRawQueryParams() string {
	var arrRawQuery []string
	for key, value := range context.QueryParams {
		arrRawQuery = append(arrRawQuery, key+"="+strings.Join(value, ","))
	}
	return strings.Join(arrRawQuery, "&")
}

func init() {
	err := gonfig.GetConf(".conf", &conf)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	router := gin.Default()
	router.GET(conf.RELATIVE_PATH+"/:methodApi", Api)
	router.Run(":"+conf.PORT)
}

func Api(ginContext *gin.Context) {
	var context Context
	context.SetMethodApi(ginContext.Param("methodApi"))
	context.SetQueryParams(ginContext.Request.URL.Query())

	response, err := http.Get(conf.API_URL + context.MethodApi + "?" + context.GetRawQueryParams()) // onError
	if err != nil {
		ginContext.JSON(400, err.Error())
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	ginContext.JSON(200, string(body))
}
