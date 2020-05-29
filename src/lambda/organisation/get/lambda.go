package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/emisgroup/PatientFlowAPI/src/internal/constant"
	"github.com/emisgroup/PatientFlowAPI/src/internal/models/organisation"
	"github.com/emisgroup/PatientFlowAPI/src/internal/utils"

	// "github.com/KorruArchana/PatientFlowAPI/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	// "github.com/google/jsonapi"
)

func main() {
	lambda.Start(lambdaHandler)

	// lambdaHandler()
}

// type organisation struct {
// 	// ID         string `jsonapi:"attr,SK"`
// 	SystemType       string `jsonapi:"attr,SystemType"`
// 	OrganisationName string `jsonapi:"attr,OrganisationName"`
// 	PK string `jsonapi:"primary,PK"`
// 	SK string `jsonapi:"attr,SK"`
// }

// func getNewSession() (*session.Session, error) {
// 	creds := credentials.NewEnvCredentials() // For deployment.. It'll get the credentials from the environment
// 	// creds := credentials.NewSharedCredentials("C:/credentials", "default") //For debugging in local, Give local path
// 	session, err := session.NewSession(&aws.Config{
// 		Region: aws.String("eu-west-2"), Credentials: creds})

// 	return session, err
// }

func lambdaHandler() (events.APIGatewayProxyResponse, error) {
	//create new session
	sess, err := utils.CreateSession()
	if err != nil {
		return utils.ServerError(err, constant.SessionCreationError), err
	}
	//create dynamodb client
	dbClient := dynamodb.New(sess)

	// getOrganisations();

	//All the lines in getorganisations can be replaced by below statement
	// response, err := dbClient.GetItem(
	// 	&dynamodb.GetItemInput{
	// 		ConsistentRead: aws.Bool(true), //By default It'll be eventually consistent read, We neeed to set this to have strongly consistent read. This eventually consistent read will be slow when compared with strongly consistent read
	// 		TableName: aws.String("PatientFlow"),
	// 		Key: map[string]*dynamodb.AttributeValue{
	// 			"PK": {S: aws.String("emis:org")},
	// 			"SK": {S: aws.String("emis:org:50002")},
	// 			},
	// 	},
	// )

	// getItemInput := dynamodb.BatchGetItemInput{
	// 	RequestItems: map[string]*dynamodb.KeysAndAttributes{
	// 		"PatientFlow": &dynamodb.KeysAndAttributes{
	// 			Keys: []map[string]*dynamodb.AttributeValue{
	// 				map[string]*dynamodb.AttributeValue{
	// 					"PK": &dynamodb.AttributeValue{S: aws.String("emis:org")},
	// 					"SK": &dynamodb.AttributeValue{S: aws.String("emis:org:50002")},
	// 				},
	// 			},
	// 			ProjectionExpression: aws.String("OrganisationName, SystemType"),
	// 		},
	// 	},
	// }

	// //Check what are expression attributes and projection attributes(Output elements to get)... Then how to get multiple values and how to get
	// response, err := dbClient.BatchGetItem( &getItemInput)

	orgPK := utils.GeneratePK(constant.OrgModuleERN)

	scanInput := dynamodb.ScanInput{
		TableName:        aws.String(constant.TableName),
		FilterExpression: aws.String("PK = :PK AND begins_with(SK, :PK)"),
		// ExpressionAttributeNames: map[string]*string{
		// 	"pk" : aws.String("PK"),
		// 	"sk" : aws.String("SK"),
		// },
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":PK": {S: aws.String(orgPK)},
			// "SK" : {S: aws.String("emis:org:")}
		},
	}

	response, err := dbClient.Scan(&scanInput)

	if err != nil {
		fmt.Println(err.Error())
		return utils.ServerError(err, constant.ErrorGettingData), err
	}

	// We need to unmarshall the return values
	item := []organisation.Organisation{}
	unMarshallerr := dynamodbattribute.UnmarshalListOfMaps(response.Items, &item)

	if unMarshallerr != nil {
		fmt.Println("Failed to marshall the record")
	}

	fmt.Println("Organisation record:", response)
	fmt.Println("Organisation records:", item)

	return getGatewayResponse(item), nil
}

func getGatewayResponse(item []organisation.Organisation) events.APIGatewayProxyResponse {

	orgData, _ := json.Marshal(item)
	gatewayResponse := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                "application/vnd.api+json",
			"Access-control-Allow-Origin": "*",
		},
		Body: string(orgData),
	}

	return gatewayResponse
}

// // Gets organisation data using dbClient's GetItem in detailed manner
// func getOrganisations(){
// 	tableName := "PatientFlow"
// 	orgPK := "emis:org"
// 	orgSk := "emis:org:50002"

// 	attributeValue1 := dynamodb.AttributeValue{
// 		S : &orgPK,
// 	}

// 	attributeValue2 := dynamodb.AttributeValue{
// 		S : &orgSk,
// 	}

// 	//Prepare Query input
// 	getItemInput := dynamodb.GetItemInput{
// 		// ConsistentRead: aws.Bool(true),
// 		TableName: aws.String(tableName),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			 "PK": &attributeValue1,
// 			 "SK": &attributeValue2,
// 			 },
// 	}

// 	//Get the records
// 	response, err := dbClient.GetItem(&getItemInput)

// 	return response, err
// }

// // Gets organisation data using dbClient's GetItem
// func getOrganisationsSimply()(response, error){

// 	response, err := dbClient.GetItem(
// 		&dynamodb.GetItemInput{
// 			ConsistentRead: aws.Bool(true), //By default It'll be eventually consistent read, We neeed to set this to have strongly consistent read. This eventually consistent read will be slow when compared with strongly consistent read
// 			TableName: aws.String("PatientFlow"),
// 			Key: map[string]*dynamodb.AttributeValue{
// 				"PK": {S: aws.String("emis:org")},
// 				"SK": {S: aws.String("emis:org:50002")},
// 				},
// 		},
// 	)

// 	return {response, err}
// }
