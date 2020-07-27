A helper CLI to help with retagging process of ECR images i.e. moving 'latest' tag from image A to image B

## Usage

```
$ ecr-image-retag --help

A helper CLI to help with retagging process of ECR images i.e. moving 'latest' tag from image A to image B

Usage:
  ecr-image-retag [flags]

Flags:
  -e, --ecr-repo string           The AWS ECR repo name containing the images
  -h, --help                      help for ecr-image-retag
  -d, --new-image-digest string   The new image digest that will receive the --tag-name
  -p, --profile string            The AWS profile name from ~/.aws/credentials file
  -r, --region string             The AWS region where the ECR repo is located
  -t, --tag-name string           The tag name that will be dropped from current images and to be applied to --new-image-digest
```

## Example

The following example command will move the `latest` tag from image(s) in `my-repo` ECR repository to an image with digest `sha256:66386bebbfe612be82286dc40e4fbb10f93ab85ad8c13d00dd73dfe822e32a01` in the same ECR repository.

```
$ ecr-image-retag --ecr-repo my-repo --new-image-digest sha256:66386bebbfe612be82286dc40e4fbb10f93ab85ad8c13d00dd73dfe822e32a01 --profile my-aws-profile --region ap-southeast-1 --tag-name latest
```

## Assumptions

This CLI assumes that you have at least one profile under _~/.aws/credentials_ path containing AWS Access Key ID and Secret Access Key like this:

```
[my-aws-profile]
aws_access_key_id = XXX
aws_secret_access_key = XXX
region = ap-southeast-1
```

## IAM Policy

The user for the `--profile` should have the following permissions in its IAM policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "EcrImageRetagCLI",
      "Effect": "Allow",
      "Action": ["ecr:BatchGetImage", "ecr:BatchDeleteImage", "ecr:PutImage"],
      "Resource": "*"
    }
  ]
}
```

To make it more secure, the wildcard for the `Resource` key above can be replaced with a complete ARN i.e. `arn:aws:ecr:ap-southeast-1:AWS_ACCOUNT_ID:repository/ECR_REPO_NAME`.

## Previous Releases

If you need to refer at specific version of this package, it's available [here](https://github.com/zulhfreelancer/ecr-image-retag/releases)

## Contribute

Feel free to fork and submit PRs for this project. I'm more than happy to review and merge it. If you have any questions regarding contributing, feel free to reach out to me on [Twitter](https://twitter.com/zulhhandyplast).
