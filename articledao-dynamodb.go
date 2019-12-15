package articledao

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
)

//ArticleDAODynamoDB is a DynamoDB implementation of ArticleDAO
type ArticleDAODynamoDB struct {
}

var sess *session.Session
var tableName string
var svc *dynamodb.DynamoDB

func init() {
	tableName = "nine-test-dev"
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = dynamodb.New(sess)
}

//GetArticle finds one article by id
func (dao *ArticleDAODynamoDB) GetArticle(articleID string) (Article, error) {

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(articleID),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return Article{}, nil
	}

	article := Article{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &article)

	return article, err
}

//SaveArticle adds the article to an in-memory storage, returns article id
func (dao *ArticleDAODynamoDB) SaveArticle(article Article) (string, error) {
	id, _ := uuid.NewUUID()
	uuidStr, _ := id.MarshalText()
	dateStr := strings.ReplaceAll(article.Date, "-", "")
	article.ID = dateStr + "-" + string(uuidStr)
	av, err := dynamodbattribute.MarshalMap(article)
	if err != nil {
		fmt.Println(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return "", err
	}

	return article.ID, nil
}

//FindByTagAndDate finds articles by tag and date
func (dao *ArticleDAODynamoDB) FindByTagAndDate(tag string, date string) (FindResponse, error) {

	if !IsDateParamValid(date) {
		return FindResponse{}, errors.New("date path variable is invalid")
	}

	filt := expression.Name("id").BeginsWith(date)
	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
	}

	response := FindResponse{}
	response.Tag = tag
	response.Count = 0
	response.Articles = make([]string, 0, 100)
	response.RelatedTags = make([]string, 0, 100)

	for _, i := range result.Items {
		a := Article{}
		dynamodbattribute.UnmarshalMap(i, &a)
		dateStr := strings.ReplaceAll(a.Date, "-", "")
		if dateStr == date {
			response.RelatedTags = MergeUnique(response.RelatedTags, a.Tags, tag)
			if Contains(a.Tags, tag) {
				response.Count = response.Count + 1
				response.Articles = append(response.Articles, a.ID)
			}
		}
	}

	return response, nil
}
