package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	awsEndpoint string // Used by customResolver for LocalStack
	awsRegion   string // Used by customResolver for LocalStack
	sdkCfg      aws.Config
)

// Client Data
type CData struct {
	region string
	bucket string
	key    string
}

//	 Constructor Function for CData.
//		Returns a pointer to a struct literal that has been populated with passed-in data.
func NewCData(r, b, k string) *CData {
	return &CData{region: r, bucket: b, key: k}
}

// ListOI interface. Perform a search in S3 based on values provided.
type ListOI interface {
	SearchO() (string, error)
}

// SearchO. Method on type CData, satisfies ListOI
// If an object matches CData.key, return the object that matched.
func (cd *CData) SearchO() (string, error) {
	var found string
	s3client := s3.NewFromConfig(sdkCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.Region = cd.region
	})

	// Retrieve an object list in CData.bucket
	bucketlist, err := s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(cd.bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, object := range bucketlist.Contents {
		if aws.ToString(object.Key) == cd.key {
			log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
			found = aws.ToString(object.Key)
			break
		}
	}
	return found, err
}

// Find. Takes an interface with SearchO() and initiates the search.
func Find(i ListOI) (string, error) {
	found, err := i.SearchO()
	if err != nil {
		fmt.Println("Error: ", err)
		log.Fatal(err)
	}
	return found, nil
}

func main() {
	var localProfile string

	// Default the AWS Profile to "test", but allow the selection of "localstack" with -p
	flag.StringVar(&localProfile, "profile", "test", "local profile name (default: test)")
	flag.Parse()
	fmt.Println("Using: ", localProfile)

	// This will be used if an endpoint is defined for LocalStack,
	// e.g.: a profile from ~/.aws/config that includes `endpoint_url = http://localhost:4566`
	// These variables are assigned according to the loaded profile, they are not to be assigned here.
	//
	//  TODO: Update to V2 of aws.Endpoint, V1 is deprecated.
	//
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// Load Config Profile (~/.aws/config)
	var err error
	sdkCfg, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithSharedConfigProfile(localProfile),
	)
	if err != nil {
		log.Fatalf("Cannot load AWS Profile: %s", err)
	}

	// Build the Data Struct for the ListO function to find the desired object.
	kFind := NewCData("us-west-2", "sre-matttest", "userneeds.png")

	// Perform the search
	finditem, err := Find(kFind)

	// Alternative: Directly with the method
	// bitem, err := kFind.SearchO()

	fmt.Printf("Found bucket item: %q\n", finditem)
}
