package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/emisgroup/PatientFlowAPI/src/internal/utils"
)

func main() {
	lambdaHandler()
}

func lambdaHandler() events.APIGatewayProxyResponse {
	sess, err := utils.CreateSession()
	if err != nil {
		return utils.ServerError(err, err.Error())
	}
	dbClient := dynamodb.New(sess)

	resp, err := dbClient.DeleteItem(
		&dynamodb.DeleteItemInput{
			TableName: aws.String("PatientFlow"),
			Key: map[string]*dynamodb.AttributeValue{
				"PK": {S: aws.String("emis:org")},
				"SK": {S: aws.String("emis:org:50005")},
			},
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		return utils.ServerError(err, err.Error())
	}

	fmt.Println(resp)
	return utils.GetAPIGateWayResponse(http.StatusCreated, "Organisation deleted successfully")
}
