package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main(){
	lambdaHandler();
}

func getNewSession()(*session.Session, error){
	creds := credentials.NewSharedCredentials("C:/credentials", "default")
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-2"),Credentials: creds})

	return sess, err
}

func lambdaHandler(){

	sess, err := getNewSession()
	dbClient := dynamodb.New(sess)

	resp, err := dbClient.UpdateItem(
		&dynamodb.UpdateItemInput{
			TableName: aws.String("PatientFlow"),
			Key: map[string]*dynamodb.AttributeValue{
				"PK": {S: aws.String("emis:org")},
				"SK": {S: aws.String("emis:org:60001")},
			},
			ExpressionAttributeNames: map[string]*string{
				"#ORGNAME": aws.String("OrganisationName"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":orgname": {S: aws.String("Org updated second time")},
			},
			UpdateExpression: aws.String("SET #ORGNAME = :orgname"),
			// UpdateExpression: aws.String("SET #ORGNAME = :orgname, #SITENAME = :sitename"), //For updating multiple values
			// ReturnValues:     aws.String("ALL_NEW"),  // Not sure why it is needed
		},
	)

	if err != nil{
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp)
}


