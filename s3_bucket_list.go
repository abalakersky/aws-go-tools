package main

import (
	"fmt"
	"log"
	"time"

	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"my_tools/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	t := time.Now()
	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println("Failed to list buckets", err)
		return
	}

	fmt.Println("Here are your buckets on", t.Format(time.RFC1123), "\n")
	//	fmt.Println(result)
	for _, bucket := range result.Buckets {
		fmt.Printf("%s\n", *bucket.Name)
	}
}
