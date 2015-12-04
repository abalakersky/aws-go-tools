package main

import (
	"time"
	"strconv"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
	"bufio"
	"os"
	"log"
	"path/filepath"
	"strings"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var (
	bucket = flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	region = flag.String("region", "us-east-1", "Region to connect to.")
	creds = flag.String("creds", "default", "Credentials Profile to use")
	search = flag.String("search", "", "Search string to find in object paths")
	t = time.Now()
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	name = dir + "/" + *bucket + "_" + *search + strconv.FormatInt(t.Unix(), 10) + ".log"
)

func prefixes (bucket string , svc s3iface.S3API) {
	topLevel, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucket, Delimiter: aws.String("/")})
	if err != nil {
		log.Println("Failed to list Top Level objects", err)
		return
	}
	var P []string
	var K []string
	for _, commonPrefix := range topLevel.CommonPrefixes {
		P = append(P, *commonPrefix.Prefix)
//		fmt.Println("Directories:", *commonPrefix.Prefix)
	}
	for _, contentKeys := range topLevel.Contents {
		K = append(K, *contentKeys.Key)
	}

	return P
	return K
}

func listObjectsWorker(objCh chan<- *s3.Object, prefix string, bucket string, svc s3iface.S3API) {
		params := &s3.ListObjectsInput{
			Bucket: &bucket,
			Prefix: &prefix,
		}
		err := svc.ListObjectsPages(params,
			func(page *s3.ListObjectsOutput, last bool) bool {
				for _, object := range page.Contents {
					objCh <- object
				}
				return true
			},
		)

		if err != nil {
			fmt.Println("failed to list objects by prefix", prefix, err)
		}
	}

func main() {
	flag.Parse()
	svc := s3.New(session.New(&aws.Config{
		Region:      *region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	if *bucket == "" {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}

	prefixes(*bucket, svc)

	objCh := make(chan *s3.Object, 100)
	var wg sync.WaitGroup

	for _, prefix := range prefixes(*bucket, svc) {
		fmt.Println("Directories:", *commonPrefix.Prefix)
	}
}