# What is Injeolmi?

Injeolmi is AWS Lambda function for terraform distribution through Gitlab. 

It is very inspired by [atlantis](https://github.com/runatlantis/atlantis). The difference is that the terraform-cli is performed on the Gitlab-Pipeline.

## Architecture

![arch](.assets/arch.png)


# AWS Infrastructure

A series of AWS resources for deploying **injeolmi** can be build with terraform in `/infra`.
