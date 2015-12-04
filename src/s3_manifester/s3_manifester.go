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
	"os"
	"log"
	"path/filepath"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"bufio"
	"strings"
)

var (
	bucket = flag.String("bucket", "", "Bucket Name to list objects from. REQUIRED")
	region = flag.String("region", "us-east-1", "Region to connect to.")
	creds = flag.String("creds", "default", "Credentials Profile to use")
	search = flag.String("search", "", "Search string to find in object paths")
	t = time.Now()
	dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	name = dir + "/" + strconv(*bucket) + "_" + strconv(*search) + strconv.FormatInt(t.Unix(), 10) + ".log"
)


func listObjectsWorker(objCh chan <- *s3.Object, prefix string, bucket *string, svc s3iface.S3API) {
	params := &s3.ListObjectsInput{
		Bucket: bucket,
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

func caseInsesitiveContains (s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func main() {
	flag.Parse()
	svc := s3.New(session.New(&aws.Config{
		Region:      region,
		Credentials: credentials.NewSharedCredentials("", *creds),
	}))

	if *bucket == "" {
		fmt.Printf("\n%s\n\n", "You Need to specify name of the Bucket to scan")
		return
	}

	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	topLevel, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: bucket, Delimiter: aws.String("/")})
	if err != nil {
		log.Println("Failed to list Top Level objects", err)
		return
	}
	for _, contentKeys := range topLevel.Contents {
		fmt.Println(*contentKeys.Key)
	}

	objCh := make(chan *s3.Object, 10)
	var wg sync.WaitGroup

	for _, commonPrefix := range topLevel.CommonPrefixes {
//		fmt.Println(commonPrefix.Prefix)
		wg.Add(1)
		go func() {
			defer wg.Done()
			listObjectsWorker(objCh, *commonPrefix.Prefix, bucket, svc)
		}()
		go func() {
			wg.Wait()
			close(objCh)
		}()
	}
//	for obj := range objCh {
//		fmt.Println(*obj.Key)
//	}
	w := bufio.NewWriter(f)
	for obj := range objCh {
		switch  {
		case *search == "" :
//			fmt.Println(*obj.Key)
			w.WriteString(*obj.Key + "\n")
		case *search != "" :
			if caseInsesitiveContains(*obj.Key, *search) == true {
				fmt.Println(*obj.Key)
				w.WriteString(*obj.Key + "\n")
			} else {
				continue
			}
		}
	}
	w.Flush()
}