package main

import (
	"fmt"
	"log"
//	"os"
	"time"
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	t := time.Now()
	bucket := flag.String("bucket", "test", "Bucket Name to list objects from")
	flag.Parse()
//	fmt.Printf("Your bucket name is %q", *bucket)

	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1")}))
	result,  err := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucket})
	if err != nil {
		log.Println("Failed to list objects", err)
		return
	}

	fmt.Printf("Here are the objects in %q bucket on %s\n\n", *bucket, t.Format(time.RFC1123))
//	fmt.Println(result)
	for _, object := range result.Contents {
		fmt.Printf("%s\n", *object.Key)
	}
}
