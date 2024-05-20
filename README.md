# S3 Object Lock Demo
This Go project is a demonstration of a simple command-line tool to manage object locking in an AWS S3 bucket. It allows users to acquire and release locks on a specified file in an S3 bucket. This project is designed to showcase the concept of object locking in AWS S3 and is a minimalistic implementation for educational purposes.
## Table of Contents
1. [Prerequisites](#Prerequisites)
2. [Building the Project](#Building the Project)
3. [Usage](#Usage)
4. [Running Unit Tests](#Usage)
5. [Design](#Design)
6. [Important Note](#Important Note)
## Prerequisites
Before using this tool, make sure you have the following prerequisites:

Go (Golang) installed on your machine.
AWS credentials configured on your system. You can set these up using the aws configure command.

## Building the Project
To build the project, follow these steps:

Clone the repository to your local machine:

```bash
git clone https://github.com/ravinitp/s3-object-lock-demo.git
```
Change to the project directory:

```bash
cd s3-object-lock-demo
```
Build the project using the following command:

```bash
go build .
```


## Usage
The tool provides two main commands: lock and unlock. These commands allow you to acquire and release locks on an object in an S3 bucket.

### Lock Command
To acquire a lock on an S3 object, use the lock command as follows:

```bash
./s3-object-lock-demo -command=lock -bucket=<S3_BUCKET_NAME> -file=<OBJECT_KEY>
```
<S3_BUCKET_NAME> is the name of the S3 bucket where the object is stored.
<OBJECT_KEY> is the key (path) to the object you want to lock.
### Unlock Command
To release a lock on an S3 object, use the unlock command as follows:

```bash
./s3-object-lock-demo -command=unlock -bucket=<S3_BUCKET_NAME> -file=<OBJECT_KEY>
```
you will need to provide the versionId obtained during the lock operation to release the lock properly.

<S3_BUCKET_NAME> is the name of the S3 bucket where the object is stored.
<OBJECT_KEY> is the key (path) to the object you want to unlock.
## Running Unit Tests
Unit tests are included in this project to validate the locking mechanism. These tests simulate multiple threads simultaneously attempting to acquire a lock on the same object.

To run the unit tests, use the following command:

```
go test .
```

## Design
The locking mechanism in this project follows these steps:

Copy the object from <OBJECT_KEY> to <OBJECT_KEY>.lock to create a locked version of the object.
Take note of the version of the locked object.
Delete the original object <OBJECT_KEY>.
When another process or user wants to acquire the lock, it must check the version of <OBJECT_KEY>.lock.
If the version matches the one acquired earlier, the lock is granted; otherwise, it indicates that someone else has acquired the lock.
This design guarantees the elimination of race conditions when acquiring the lock.

Important: Ensure that your S3 bucket has versioning enabled for this mechanism to work correctly.

## Important Note
This project is intended for educational purposes and as a demonstration of object locking in AWS S3. In a production environment, it is recommended to use AWS's built-in S3 Object Lock feature, which provides a more robust and secure locking mechanism.

Please use this tool responsibly and ensure that your AWS credentials and permissions are correctly configured for S3 bucket operations.




