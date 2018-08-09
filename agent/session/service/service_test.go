// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package service is a wrapper for the message gateway Service
package service

import (
	"encoding/xml"
	"testing"

	"github.com/aws/amazon-ssm-agent/agent/log"
	"github.com/aws/amazon-ssm-agent/agent/rip"
	mgsConfig "github.com/aws/amazon-ssm-agent/agent/session/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/stretchr/testify/assert"
	"github.com/twinj/uuid"
)

var (
	signer = &v4.Signer{
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", "SESSION")}
	region     = "us-east-1"
	instanceId = "i-12345678"
	sessionId  = "s-12345678"
	token      = "abcdefg"
)

func TestGetRegion(t *testing.T) {
	service := getService()

	result := service.GetRegion()

	assert.Equal(t, region, result)
}

func TestGetV4Signer(t *testing.T) {
	service := getService()

	result := service.GetV4Signer()

	assert.Equal(t, result, signer)
}

func TestCreateControlChannel(t *testing.T) {
	service := getService()
	createControlChannelInput := &CreateControlChannelInput{
		MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		RequestId:            aws.String(uuid.NewV4().String()),
	}
	makeRestcall = func(request []byte, methodType string, url string, region string, signer *v4.Signer) ([]byte, error) {
		output := &CreateControlChannelOutput{
			TokenValue:           aws.String(token),
			MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		}
		return xml.Marshal(output)
	}
	output, err := service.CreateControlChannel(log.NewMockLog(), createControlChannelInput, instanceId)

	assert.Nil(t, err)
	assert.Equal(t, token, *output.TokenValue)
}

func TestCreateDataChannel(t *testing.T) {
	service := getService()
	createDataChannelInput := &CreateDataChannelInput{
		MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		RequestId:            aws.String(uuid.NewV4().String()),
		ClientId:             aws.String(uuid.NewV4().String()),
	}
	makeRestcall = func(request []byte, methodType string, url string, region string, signer *v4.Signer) ([]byte, error) {
		output := &CreateDataChannelOutput{
			TokenValue:           aws.String(token),
			MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		}
		return xml.Marshal(output)
	}
	output, err := service.CreateDataChannel(log.NewMockLog(), createDataChannelInput, sessionId)

	assert.Nil(t, err)
	assert.Equal(t, token, *output.TokenValue)
}

func TestDeleteControlChannel(t *testing.T) {
	service := getService()
	deleteControlChannelInput := &DeleteChannelInput{
		MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		RequestId:            aws.String(uuid.NewV4().String()),
	}
	makeRestcall = func(request []byte, methodType string, url string, region string, signer *v4.Signer) ([]byte, error) {
		output := &DeleteChannelOutput{
			ChannelId:            aws.String(instanceId),
			MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		}
		return xml.Marshal(output)
	}
	output, err := service.DeleteControlChannel(log.NewMockLog(), deleteControlChannelInput, instanceId)

	assert.Nil(t, err)
	assert.Equal(t, instanceId, *output.ChannelId)
}

func TestDeleteDataChannel(t *testing.T) {
	service := getService()
	deleteDataChannelInput := &DeleteChannelInput{
		MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		RequestId:            aws.String(uuid.NewV4().String()),
	}
	makeRestcall = func(request []byte, methodType string, url string, region string, signer *v4.Signer) ([]byte, error) {
		output := &DeleteChannelOutput{
			ChannelId:            aws.String(sessionId),
			MessageSchemaVersion: aws.String(mgsConfig.MessageSchemaVersion),
		}
		return xml.Marshal(output)
	}
	output, err := service.DeleteControlChannel(log.NewMockLog(), deleteDataChannelInput, instanceId)

	assert.Nil(t, err)
	assert.Equal(t, sessionId, *output.ChannelId)
}

func TestGetBaseUrl(t *testing.T) {
	// data channel url test
	dataChannelUrlResult, err := getMGSBaseUrl(log.NewMockLog(), mgsConfig.DataChannel, sessionId, region)

	host := rip.GetMgsEndpoint(region)
	expectedDataChannelUrl := "https://" + host + "/v1/data-channel/" + sessionId
	assert.Nil(t, err)
	assert.Equal(t, expectedDataChannelUrl, dataChannelUrlResult)

	// control channel url test
	controlChannelUrlResult, err := getMGSBaseUrl(log.NewMockLog(), mgsConfig.ControlChannel, instanceId, region)

	expectedControlChannelUrl := "https://" + host + "/v1/control-channel/" + instanceId
	assert.Nil(t, err)
	assert.Equal(t, expectedControlChannelUrl, controlChannelUrlResult)
}

func getService() Service {
	return &MessageGatewayService{
		region: "us-east-1",
		signer: signer,
	}
}