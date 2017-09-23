# Go Deploy

Just a small tool for deploying my site to AWS.

![godeploy](https://raw.githubusercontent.com/JohnBehnke/godeploy/master/godeploy.png)

## Why Go?
 
This could totaly be done in a bash script in about 10 lines. Thats actaully what my inital version was, but I wanted to learn Go, so after ~10~ 17 dependencies and over 200 lines of code, I have a working script. You might ask, why not just use the `aws cli`? Well thats what I _should_ be using, but I have a bit of a unique use case. When I want to deploy changes to my site, I might need to update files that aren't store in my repository for my site, such as a new PDF for my resume, set up scripts, etc. Initally these files were hard coded in the script, but I decided that was a bit lazy, so I introduced a `config.yaml` file for defining these files to copy and upload. Neat 📸

## config.yaml

The script is going to expect a `config.yaml` file in the root directory where the script is executed from. There is an example file provided. 

`bucket` defines the name of the bucket to upload to.
`region` defines the region to upload to.
`files` is an map, where the key is the path in the S3 bucket for the file to be uploaded to, and the value is the path in the OS to the file to upload. You can use `.` as a key to indicate not not nest the file in a folder in S3

`ignore` is an array of files and directories that should not be uploaded. 

## Building 

`go build`

## Running

`./godeploy FLAGS`

## Usage

`-copy` Copies files from outside the site directory into the site directory

`-upload` Uploads the contents of the site directory folder to an AWS Bucket

`-commit` Does the git dance (add, commit, push). You can add a commit message here too

## Example Usage

 `./godeploy -copy -upload -commit "Adding new easter eggs to site"`
 

