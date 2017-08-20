///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	// "flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"os"
)

// func isDirectory(path string) bool {

// }

func isDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file.Mode().IsDir()
}

func uploadDirectory(session *session.Session, bucketPrefix string, dirPath string) {

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range files {
		var pathToCheck = dirPath + d.Name()
		if d.Name() != ".git" && d.Name() != ".DS_Store" {

			if isDir(pathToCheck) {
				// fmt.Println(strings.Repeat("\t", num), d.Name())
				// printDir(pathToCheck, num+1)
				uploadDirectory(session, d.Name()+"/", pathToCheck+"/")
			} else {
				// fmt.Println(strings.Repeat("\t", num), d.Name()
				uploadFile(session, "testjlb", bucketPrefix, dirPath, d.Name())
			}
		}

	}

}

func uploadFile(session *session.Session, bucket string, bucketPrefix string, filePath string, fileName string) {
	serviceClient := s3.New(session)
	fmt.Println(filePath + fileName)
	file, err := os.Open(filePath + fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = serviceClient.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketPrefix + fileName),
		Body:   file,
	})
	fmt.Print(file)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
}

func main() {
	const scriptName string = "setup.sh"
	const resumeName string = "JohnBehnkeResume.pdf"
	const tempPath string = "/tmp/deployStaging"
	const resumePath string = "/Users/john/Documents/Work/Resumes/Latest/"
	const installScriptPath string = "/Users/john/Documents/Code/Personal/EasyDevTools"
	const personalSitePath string = "/Users/john/Documents/Code/Personal/johnbehnke.github.io/"

	if _, err := os.Stat(personalSitePath + "assets/files/" + resumeName); !os.IsNotExist(err) {
		fmt.Println("It is there!")
		os.Remove(personalSitePath + "/assets/files/" + resumeName)
	}

	files, err := ioutil.ReadDir(resumePath)
	if err != nil {
		fmt.Println(err)
	}
	var resume = resumePath + (files[len(files)-1]).Name()

	sourcePDF, err := os.Open(resume + "/" + resumeName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sourcePDF.Close()

	destinationPDF, err := os.Create(personalSitePath + "/assets/files/" + resumeName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer destinationPDF.Close()

	_, err = io.Copy(destinationPDF, sourcePDF)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = destinationPDF.Sync()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := os.Stat(personalSitePath + scriptName); !os.IsNotExist(err) {
		fmt.Println("It is there!")
		os.Remove(personalSitePath + scriptName + resumeName)
	}

	sourceScript, err := os.Open(installScriptPath + "/" + scriptName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sourceScript.Close()

	destinationScript, err := os.Create(personalSitePath + scriptName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer destinationScript.Close()

	_, err = io.Copy(destinationScript, sourceScript)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = destinationScript.Sync()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// siteFiles, err := ioutil.ReadDir(personalSitePath)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// printDir(personalSitePath, 0)

	// var bucket string
	// var key string

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		fmt.Print(err)
	}

	uploadDirectory(sess, "", personalSitePath)

	// uploadFile(sess, "testjlb", "/", filePath, fileName)
	// serviceClient := s3.New(sess)

	// flag.StringVar(&bucket, "b", "", "Bucket Name.")
	// flag.StringVar(&key, "k", "", "Object key name.")
	// flag.Parse()

	// t, err := serviceClient.PutObject(&s3.PutObjectInput{
	// 	Bucket: aws.String(bucket),
	// 	Key:    aws.String(key),
	// 	Body:   os.Stdin,
	// })
	// fmt.Print(t)
	// if err != nil {
	// 	fmt.Fprint(os.Stderr, "Error %v\n", err)
	// 	os.Exit(1)
	// }
}
