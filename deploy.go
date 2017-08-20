///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"io/ioutil"
	"os"
)

func isDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file.Mode().IsDir()
}
func getFileCount(path string) int {
	count := 0
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range files {
		var pathToCheck = path + d.Name()
		if d.Name() != ".git" && d.Name() != ".DS_Store" {
			if isDir(pathToCheck) {
				count = count + getFileCount(pathToCheck+"/")
			} else {
				count = count + 1
			}
		}
	}
	return count

}

func uploadDirectory(session *session.Session, progessBar *pb.ProgressBar, bucketPrefix string, dirPath string) {

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range files {
		var pathToCheck = dirPath + d.Name()
		if d.Name() != ".git" && d.Name() != ".DS_Store" {
			if isDir(pathToCheck) {
				uploadDirectory(session, progessBar, d.Name()+"/", pathToCheck+"/")
			} else {
				uploadFile(session, progessBar, "testjlb", bucketPrefix, dirPath, d.Name())
			}
		}
	}
}

func uploadFile(session *session.Session, progressBar *pb.ProgressBar, bucket string, bucketPrefix string, filePath string, fileName string) {
	serviceClient := s3.New(session)
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
	if err != nil {
		fmt.Fprint(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
	// fmt.Println(t)
	progressBar.Increment()
}

func main() {
	const scriptName string = "setup.sh"
	const resumeName string = "JohnBehnkeResume.pdf"
	const tempPath string = "/tmp/deployStaging"
	const resumePath string = "/Users/john/Documents/Work/Resumes/Latest/"
	const installScriptPath string = "/Users/john/Documents/Code/Personal/EasyDevTools"
	const personalSitePath string = "/Users/john/Documents/Code/Personal/johnbehnke.github.io/"

	// uploadFile(sess, "testjlb", "/", filePath, fileName)
	// serviceClient := s3.New(sess)

	// flag.StringVar(&bucket, "b", "", "Bucket Name.")
	// flag.StringVar(&key, "k", "", "Object key name.")
	// flag.Parse()

	// t, err := serviceClient.PutObject(&s3.PutObjectInput{
	//  Bucket: aws.String(bucket),
	//  Key:    aws.String(key),
	//  Body:   os.Stdin,
	// })
	// fmt.Print(t)
	// if err != nil {
	//  fmt.Fprint(os.Stderr, "Error %v\n", err)
	//  os.Exit(1)
	// }

	shouldCopy := flag.Bool("c", false, "Copy new files to site dir")
	shouldUpload := flag.Bool("u", false, "Upload site to AWS Bucket")
	flag.Parse()

	if *shouldCopy {
		fmt.Println("Copying over files...")
		progressBar := pb.StartNew(2)
		progressBar.Format("<=👉 >")
		if _, err := os.Stat(personalSitePath + "assets/files/" + resumeName); !os.IsNotExist(err) {
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
		progressBar.Increment()

		if _, err := os.Stat(personalSitePath + scriptName); !os.IsNotExist(err) {
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
		progressBar.Increment()
		progressBar.FinishPrint("Files successfully copied over 🍻")
	}
	if *shouldUpload {

		sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println("Uploading Files")
		progressBar := pb.StartNew(getFileCount(personalSitePath))
		progressBar.Format("<=👉 >")

		uploadDirectory(sess, progressBar, "", personalSitePath)
		progressBar.FinishPrint("All Files uploaded! 🍻")

	}
}
