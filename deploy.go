///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/cheggaaa/pb.v1"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	Bucket string            `yaml:"bucket"`
	Region string            `yaml:"region"`
	Files  map[string]string `yaml:"files"`
	Ignore []string          `yaml:"ignore"`
}

//Determine if an item at a given path is a directory or not
func isDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return file.Mode().IsDir()
}

//Count the number of files in a directory. Ignores anything in .git and the .DS_Store file
func getFileCount(path string, ignore []string) int {
	count := 0
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, d := range files {
		var pathToCheck = strings.Join([]string{path, d.Name()}, "/")
		if !exists(ignore, d.Name()) {
			if isDir(pathToCheck) {
				count = count + getFileCount(pathToCheck+"/", ignore)
			} else {
				count += 1
			}
		}
	}
	return count
}

func exists(testingArray []string, target string) bool {
	var returnValue = false

	for x := range testingArray {
		if testingArray[x] == target {
			returnValue = true
			break
		}

	}
	return returnValue
}

//Upload the contents of a directory to AWS S3
func uploadDirectory(session *session.Session, progessBar *pb.ProgressBar, bucketPrefix string, dirPath string, ignore []string, bucket string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, d := range files {
		var pathToCheck = strings.Join([]string{dirPath, d.Name()}, "/")
		if !exists(ignore, d.Name()) {
			if isDir(pathToCheck) {
				uploadDirectory(session, progessBar, bucketPrefix+d.Name()+"/", pathToCheck+"/", ignore, bucket)
			} else {
				uploadFile(session, progessBar, bucket, bucketPrefix, dirPath, d.Name())
			}
		}
	}
}

//Upload a specific file to AWS S3 Bucket
func uploadFile(session *session.Session, progressBar *pb.ProgressBar, bucket string, bucketPrefix string, filePath string, fileName string) {
	serviceClient := s3.New(session)
	file, err := os.Open(filePath + fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	if path.Ext(filePath+fileName) == ".css" {
		fileType = "text/css"
	}

	_, err = serviceClient.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(bucketPrefix + fileName),
		Body:          fileBytes,
		ContentType:   aws.String(fileType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	progressBar.Increment()
}

func deleteFiles(region string, bucket string) {
	fmt.Println("Deleting Files from s3...")

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	serviceClient := s3.New(sess)

	resp, err := serviceClient.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var count = 0

	for _, item := range resp.Contents {
		if !strings.Contains(*item.Key, "cdn") && !strings.Contains(*item.Key, "s3") {
			count += 1
		}
	}
	progressBar := pb.StartNew(count)
	progressBar.Format("<=👉 >")

	for _, item := range resp.Contents {

		if !strings.Contains(*item.Key, "cdn") && !strings.Contains(*item.Key, "s3") {
			input := &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(*item.Key),
			}

			_, err := serviceClient.DeleteObject(input)
			if err != nil {
				fmt.Println(err)
			}
			progressBar.Increment()
		}

	}

	progressBar.FinishPrint("Files deleted! 🍻")
}

func copyFiles(localDir string, filesToCopy map[string]string) {

	fmt.Println("Copying over files...")
	progressBar := pb.StartNew(len(filesToCopy))
	progressBar.Format("<=👉 >")
	for targetPath, targetPayload := range filesToCopy {
		var targetFile string = filepath.Base(targetPayload)
		if _, err := os.Stat(strings.Join([]string{localDir, targetPath, targetFile}, "/")); !os.IsNotExist(err) {
			os.Remove(strings.Join([]string{localDir, targetPath, targetFile}, "/"))
		}
		source, err := os.Open(targetPayload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer source.Close()

		destination, err := os.Create(strings.Join([]string{localDir, targetPath, targetFile}, "/"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = destination.Sync()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		progressBar.Increment()

	}
	progressBar.FinishPrint("Files successfully copied over 🍻")
}

func main() {

	var shouldPush string
	shouldCopy := flag.Bool("copy", false, "Copy new files to site dir")
	shouldUpload := flag.Bool("upload", false, "Upload site to AWS Bucket")
	shouldDelete := flag.Bool("delete", false, "Delete the contents of the AWS buckets")
	flag.StringVar(&shouldPush, "commit", "", "Commit message to be used in a push")
	flag.Parse()

	var config Config
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var localPath string = dir

	if *shouldCopy {
		copyFiles(localPath, config.Files)
	}

	if *shouldDelete {
		deleteFiles(config.Region, config.Bucket)
	}
	if *shouldUpload {

		sess, err := session.NewSession(&aws.Config{Region: aws.String(config.Region)})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Uploading Files")
		progressBar := pb.StartNew(getFileCount(localPath, config.Ignore))
		progressBar.Format("<=👉 >")
		uploadDirectory(sess, progressBar, "", localPath+"/", config.Ignore, config.Bucket)
		progressBar.FinishPrint("All Files uploaded! 🍻")

	}

	if shouldPush != "" {

		cmd := "git"
		args := []string{"-C", localPath, "add", "."}
		args2 := []string{"-C", localPath, "commit", "-m", shouldPush}
		args3 := []string{"-C", localPath, "push"}

		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := exec.Command(cmd, args2...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := exec.Command(cmd, args3...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	}
}
