/*
	A "hello world" plugin in Go,
	which reads a request header and sets a response header.
*/

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"kong-plugin/service"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/hooklift/gowsdl/soap"
)

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.1"
var Priority = 1

type Config struct {
}

func New() interface{} {
	return &Config{}
}

func (config *Config) Access(kong *pdk.PDK) {
	body, err := kong.Request.GetRawBody()
	if err != nil {
		kong.Log.Err(err.Error())
	}

	request := service.NumberToWords{
		XMLName: xml.Name{},
	}

	err = json.Unmarshal([]byte(body), &request)
	if err != nil {
		kong.Log.Err(err.Error())
	}

	envelope := soap.SOAPEnvelope{
		XmlNS: soap.XmlNsSoapEnv,
		Body: soap.SOAPBody{
			Content: &request,
		},
	}

	soap_buffer := new(bytes.Buffer)
	err = xml.NewEncoder(soap_buffer).Encode(envelope)
	if err != nil {
		kong.Log.Err(err.Error())
	}

	kong.ServiceRequest.ClearHeader("Content-Type")
	if err != nil {
		kong.Log.Err(err.Error())
	}
	kong.ServiceRequest.SetHeader("Content-Type", "text/xml; charset=utf-8")
	kong.ServiceRequest.SetRawBody(soap_buffer.String())
}

func (conf Config) Response(kong *pdk.PDK) {
	var (
		code int
		err  error
	)

	if code, err = kong.Response.GetStatus(); err != nil {
		kong.Log.Err(err.Error())
	}

	if code == 200 {
		body, err := kong.ServiceResponse.GetRawBody()
		if err != nil {
			fmt.Print("Error reading Body")
			kong.Log.Err(err.Error())
			return
		}

		response := service.NumberToWordsResponse{}

		respEnvelope := new(soap.SOAPEnvelopeResponse)
		respEnvelope.Body = soap.SOAPBodyResponse{
			Content: &response,
		}

		err = xml.NewDecoder(strings.NewReader(body)).Decode(&respEnvelope)
		if err != nil {
			kong.Log.Err(err.Error())
			return
		}

		buffer := bytes.Buffer{}

		err = json.NewEncoder(&buffer).Encode(&response)
		if err != nil {
			kong.Log.Err(err.Error())
			return
		}
		fmt.Print(buffer.String())
		kong.Response.SetHeader("Content-Type", "application/json")
		kong.Response.Exit(code, buffer.String(), nil)
	}
}
