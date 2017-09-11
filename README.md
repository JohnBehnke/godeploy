# Go Deploy

Just a small tool for deploying my site to AWS.


## Why Go?
 
This could totaly be done in a bash script in about 10 lines. Thats actaully what my inital version was, but I wanted to learn Go, so after ~10~ 17 dependencies and over 200 lines of code, I have a working script. You might ask, why not just use the `aws cli`? Well thats what I _should_ be using, but I have a bit of a unique use case. When I want to deploy changes to my site, I might need to update files that aren't store in my repository for my site, such as a new PDF for my resume, set up scripts, etc. Initally these files were hard coded in the script, but I decided that was a bit lazy, so I introduced a `config.yaml` file for defining these files to copy and upload. Neat 📸

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
 

