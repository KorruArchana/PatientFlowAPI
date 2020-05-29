package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/emisgroup/PatientFlowAPI/src/internal/constant"
	"github.com/emisgroup/PatientFlowAPI/src/internal/utils"
)

func main() {
	lambda.Start(lambdaHandler)
}

type kiosk struct {
	PK            string `json:"pk"`
	SK            string `json:"sk"`
	KioskName     string `json:"kioskName"`
	AdminPassword string `json:"adminPassword"`
	ScreenTimeout string `json:"screenTimeOut"`
	MachineID     string `json:"machineID"`
}

func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := utils.CreateSession()
	if err != nil {
		return utils.ServerError(err, constant.SessionCreationError), err
	}

	dbClient := dynamodb.New(sess)

	kiosk := kiosk{}

	if len(request.Body) == 0 {
		return utils.ClientError(http.StatusBadRequest), err
	}

	err = json.Unmarshal([]byte(request.Body), &kiosk)
	if err != nil {
		return utils.ServerError(err, constant.UnmarshallError), err
	}

	kiosk.PK = utils.GeneratePK(constant.KioskModuleERN)
	kiosk.SK = utils.GenerateID(constant.KioskModuleERN)

	kioskInput := &dynamodb.PutItemInput{
		TableName: aws.String(constant.TableName),
		Item: map[string]*dynamodb.AttributeValue{
			"PK":            {S: aws.String(kiosk.PK)},
			"SK":            {S: aws.String(kiosk.SK)},
			"KioskName":     {S: aws.String(kiosk.KioskName)},
			"AdminPassword": {S: aws.String(kiosk.AdminPassword)},
			"ScreenTimeOut": {S: aws.String(kiosk.ScreenTimeout)},
			"MachineID":     {S: aws.String(kiosk.MachineID)},
		},
	}

	_, err = dbClient.PutItem(kioskInput)
	if err != nil {
		return utils.ServerError(err, constant.ErrorPostingData), err
	}

	//Once The kiosk is successfully created we need to add it in Cognito and generate client Id and secret key and save it in DB for that kiosk
	return utils.GetAPIGateWayResponse(http.StatusCreated, constant.KioskCreationSuccess), nil
}
