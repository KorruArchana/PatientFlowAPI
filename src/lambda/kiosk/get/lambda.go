package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/emisgroup/PatientFlowAPI/src/internal/constant"
	"github.com/emisgroup/PatientFlowAPI/src/internal/models/kiosk"
	"github.com/emisgroup/PatientFlowAPI/src/internal/utils"
)

func main() {
	// lambdaHandler(events.APIGatewayProxyRequest{Body: `
	// 	{
	// 		"kioskGUID": "emis:kiosk:27f6ad41-38ab-4ff2-a7ce-ee4f5928dd46",
	// 		"machineID": "test1234"
	// 	}
	// `})

	lambda.Start(lambdaHandler)
}

func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	kioskGUID := request.QueryStringParameters["kioskGUID"]
	machineID := request.QueryStringParameters["machineID"]

	// if len(request.QueryStringParameters) < 1 {
	// 	response, err := GetAllKiosks(request)
	// 	if err != nil {
	// 		return getGatewayResponse(response), nil
	// 	}
	// }
	// kioskGUID := "emis:kiosk:27f6ad41-38ab-4ff2-a7ce-ee4f5928dd46"
	// machineID := "test1234"

	sess, err := utils.CreateSession()
	if err != nil {
		return utils.ServerError(err, constant.SessionCreationError), err
	}

	dbClient := dynamodb.New(sess)
	kioskPK := utils.GeneratePK(constant.KioskModuleERN)
	// response, err := dbClient.GetItem(
	// 	&dynamodb.GetItemInput{
	// 		ConsistentRead: aws.Bool(true), //By default It'll be eventually consistent read, We neeed to set this to have strongly consistent read. This eventually consistent read will be slow when compared with strongly consistent read
	// 		TableName:      aws.String(constant.TableName),
	// 		Key: map[string]*dynamodb.AttributeValue{
	// 			"PK": {S: aws.String(kioskPK)},
	// 			"SK": {S: aws.String(kioskGUID)},
	// 			// "MachineID": {S: aws.String(machineID)},
	// 		},
	// 	},
	// )

	response, err := dbClient.Query(
		&dynamodb.QueryInput{
			ConsistentRead:         aws.Bool(true), //By default It'll be eventually consistent read, We neeed to set this to have strongly consistent read. This eventually consistent read will be slow when compared with strongly consistent read
			TableName:              aws.String(constant.TableName),
			KeyConditionExpression: aws.String("#pk=:pk and #sk=:sk"),
			ExpressionAttributeNames: map[string]*string{
				"#mi": aws.String("MachineID"),
				"#sk": aws.String("SK"),
				"#pk": aws.String("PK"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":mi": {S: aws.String(machineID)},
				":sk": {S: aws.String(kioskGUID)},
				":pk": {S: aws.String(kioskPK)},
			},
			FilterExpression: aws.String("#mi = :mi"),
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		return utils.ServerError(err, constant.ErrorGettingData), err
	}

	// We need to unmarshall the return values
	item := kiosk.Kiosk{}
	for _, resp := range response.Items {
		err = dynamodbattribute.UnmarshalMap(resp, &item)

		if err != nil {
			fmt.Println("Failed to marshall the record")
			return utils.ServerError(err, constant.UnmarshallError), err
		}
	}

	fmt.Println("Kiosk Response:", response)
	fmt.Println("Item:", item)

	return getGatewayResponse(item), nil
}

// GetAllKiosks to get all the kiosks data
func GetAllKiosks(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	item := kiosk.Kiosk{}
	return getGatewayResponse(item), nil
}

func getGatewayResponse(item kiosk.Kiosk) events.APIGatewayProxyResponse {

	kioskData, _ := json.Marshal(item)
	gatewayResponse := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/vnd.api+json",
			"Access-control-Allow-Origin": "*",
		},
		Body: string(kioskData),
	}

	return gatewayResponse
}

func getGatewayErrorResponse(err error) events.APIGatewayProxyResponse {

	gatewayErrorResponse := events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError) + ": " + err.Error(),
	}

	return gatewayErrorResponse
}
