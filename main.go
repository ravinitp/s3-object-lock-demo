package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sync"
	"time"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sdk   = s3.New(sess.Copy(&aws.Config{}))
	mutex = sync.Mutex{}
)

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("total time taken: %v\n", time.Now().Sub(now))

	}()
	command := flag.String("command", "", "command")
	bucket := flag.String("bucket", "", "s3 bucket")
	path := flag.String("file", "", "s3 object file")
	flag.Parse()
	fmt.Printf("Command == %s\n", *command)
	switch *command {
	case "lock":
		fmt.Println("Inside Lock")
		versionId, err := lock(*bucket, *path)
		if err != nil {
			fmt.Printf("Lock is already acquired %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("Object locked %s\n", versionId)
		}
	case "unlock":
		fmt.Println("Inside Lock")
		fmt.Println("Enter the versionId")
		var versionId string
		fmt.Scanln(&versionId)
		err := unlock(versionId, *bucket, *path)
		if err != nil {
			fmt.Printf("unlock failed %v\n", err)
			os.Exit(1)
		}
	case "test-lock":

	default:
		fmt.Printf("command %s not supported\n", *command)
	}

}
func lock(bucket, path string) (string, error) {
	// Copy state file
	copyResponse, err := sdk.CopyObject(&s3.CopyObjectInput{
		Key:        aws.String(lockFileName(path)),
		Bucket:     aws.String("lock-demo"),
		CopySource: aws.String(copyFileSource(bucket, path)),
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("copy Response   %v", copyResponse)

	// delete state file
	deleteResponse, err := sdk.DeleteObject(
		&s3.DeleteObjectInput{
			Key:    aws.String(path),
			Bucket: aws.String(bucket),
		})
	if err != nil {
		fmt.Printf("[WARNING] It seems that some one is trying to aquire the lock simultaneously.")
	}
	fmt.Printf("delete Response   %v\n", deleteResponse)
	// Read version of state file
	getResponse, err := sdk.GetObject(&s3.GetObjectInput{
		Key:    aws.String(lockFileName(path)),
		Bucket: aws.String(bucket),
	})
	fmt.Printf("read Response   %v\n", getResponse)

	if *copyResponse.VersionId != *getResponse.VersionId {
		return "", fmt.Errorf("Two Object tried to aquire lock concurrently\n  copy version %s \n read Version %s", *copyResponse.VersionId, *getResponse.VersionId)
	}

	return *getResponse.VersionId, err
}

func unlock(versionId string, bucket, path string) error {
	// restore state file (only for demo)
	response2, err := sdk.CopyObject(&s3.CopyObjectInput{
		Key:        aws.String(path),
		Bucket:     aws.String(bucket),
		CopySource: aws.String(copyFileSource(bucket, lockFileName(path))),
	})

	fmt.Printf("copy Response   %v\n", response2)
	// remove copy
	response1, err := sdk.DeleteObject(
		&s3.DeleteObjectInput{
			Key:       aws.String(lockFileName(path)),
			Bucket:    aws.String(bucket),
			VersionId: aws.String(versionId),
		})

	fmt.Printf("delete Response   %v\n", response1)
	return err
}

func lockFileName(path string) string {
	return fmt.Sprintf("%s.lock", path)
}
func copyFileSource(bucket, path string) string {
	return fmt.Sprintf("%s/%s", bucket, path)
}
