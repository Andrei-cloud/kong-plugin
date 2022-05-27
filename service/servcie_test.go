package service_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"kong-plugin/service"
	"testing"

	"github.com/hooklift/gowsdl/soap"
)

func TestNewType(t *testing.T) {

	request := service.NumberToWords{
		XMLName: xml.Name{},
		UbiNum:  500,
	}

	envelope := soap.SOAPEnvelope{
		XmlNS: soap.XmlNsSoapEnv,
	}

	envelope.Body.Content = request
	soap_buffer := new(bytes.Buffer)
	soap_encoder := xml.NewEncoder(soap_buffer)
	soap_encoder.Indent("", " ")
	err := soap_encoder.Encode(envelope)
	if err != nil {
		t.Error(err)
	}

	err = soap_encoder.Flush()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(soap_buffer.String())

	json_buffer := new(bytes.Buffer)
	json_encoder := json.NewEncoder(json_buffer)
	json_encoder.SetIndent("", " ")
	err = json_encoder.Encode(request)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(json_buffer.String())

}
