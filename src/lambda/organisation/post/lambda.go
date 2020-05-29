package main

import (
	// "strings"
	"encoding/json"
	// "bytes"
	"net/http"
	// "github.com/google/jsonapi"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/emisgroup/PatientFlowAPI/src/internal/constant"
	"github.com/emisgroup/PatientFlowAPI/src/internal/models/organisation"
	"github.com/emisgroup/PatientFlowAPI/src/internal/utils"
)

// type organisation struct {
// 	SystemType string `jsonapi:"attr,SystemType"`
// 	OrganisationName string `jsonapi:"attr,OrganisationName"`
// 	PK string `jsonapi:"primary,PK"`
// 	SK string `jsonapi:"attr,SK"`
// }

// type organisation struct {
// 	SystemType string `json:"SystemType"`
// 	OrganisationName string `json:"OrganisationName"`
// 	PK string `json:"PK"`
// 	SK string `json:",SK"`
// }

func main() {
	lambda.Start(lambdaHandler)

	// apiGatewayInput := events.APIGatewayProxyRequest{
	// 	HTTPMethod: "Post",
	// 	Headers: map[string]string {"Content-Type" : "application/vnd.api+json"},
	// 	Body: `{
	// 		"PK": "emis:org",
	// 		"SK": "emis:org:60006",
	// 		"OrganisationName": "test 60006 or to put",
	// 		"SystemType": "vision"
	// 	}`,
	// }

	// lambdaHandler(apiGatewayInput)
}

// func getNewSession() (*session.Session, error) {
// 	// creds := credentials.NewSharedCredentials("C:/credentials", "default")
// 	creds := credentials.NewEnvCredentials()
// 	session, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2"), Credentials: creds})

// 	return session, err
// }

func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := utils.CreateSession()
	if err != nil {
		return utils.ServerError(err, constant.SessionCreationError), err
	}
	dbClient := dynamodb.New(sess)
	// Decode post data into a struct

	org := organisation.Organisation{}
	if len(request.Body) == 0 {
		return utils.ClientError(http.StatusBadRequest), err
	}

	// if err := jsonapi.UnmarshalPayload(bytes.NewBufferString(request.Body), &org); err != nil {
	if err := json.Unmarshal([]byte(request.Body), &org); err != nil {
		fmt.Println(err)
		return utils.ServerError(err, constant.UnmarshallError), err
	}

	org.PK = utils.GeneratePK(constant.OrgModuleERN)
	org.SK = utils.GenerateID(constant.OrgModuleERN)

	resp, err := dbClient.PutItem(
		&dynamodb.PutItemInput{
			TableName: aws.String(constant.TableName),
			Item: map[string]*dynamodb.AttributeValue{
				"PK":               {S: aws.String(org.PK)},
				"SK":               {S: aws.String(org.SK)},
				"OrganisationName": {S: aws.String(org.OrganisationName)},
				"SystemType":       {S: aws.String(org.SystemType)},
			},
		},
	)

	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
		return utils.ServerError(err, constant.ErrorPostingData), err
	}

	fmt.Println("Put response", resp)
	return utils.GetAPIGateWayResponse(http.StatusCreated, constant.OrgCreationSuccess), nil
}

// // ClientError process' a given error status code and returns the corret APIGatewayProxyRequest response.
// func clientError(err error) events.APIGatewayProxyResponse {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusInternalServerError,
// 		Body:       err.Error(),
// 	}
// }

// func getAPIGateWayResponse(httpStatus int) events.APIGatewayProxyResponse {
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: httpStatus,
// 		Headers: map[string]string{
// 			"Content-Type":                "application/vnd.api+json",
// 			"Access-Control-Allow-Origin": "*",
// 		},
// 		Body: "Successfully created org",
// 	}
// }
