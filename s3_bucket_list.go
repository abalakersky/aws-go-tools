package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result,  err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println("Failed to list buckets", err)
		return
	}

	log.Println("Buckets:")
	fmt.Println(result)
//	for _, bucket := range result.Buckets {
//		log.Printf("%s : %s\n", *bucket.Name, bucket.CreationDate)
//	}
}
