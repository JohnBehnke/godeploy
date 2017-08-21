# SiteDeploy

Just a small tool for deploying my site to AWS.


## Why Go?
 
This could totaly be done in a bash script in about 10 lines. Thats actaully what my inital version was, but I wanted to learn Go, so after 10 dependencies and almost 200 lines of code, I have a working script. Neat 📸

## Building 

`go build`

## Running

`./SiteDeploy FLAGS`

## Usage

`-copy` Copies files from outside the site directory into the site directory
`-upload` Uploads the contents of the site directory folder to an AWS Bucket
`-commit` Does the git dance (add, commit, push). You can add a commit message here too

## Example Usage

 `./SiteDeploy -copy -upload -commit "Adding new easter eggs to site"
 

