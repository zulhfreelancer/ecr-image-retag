package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/spf13/cobra"
)

// Version is ecr-image-retag CLI version
var cliVersion string

// AWS Client
var sess *session.Session

// To store important ECR image metadata
type ecrImage struct {
	imageDigest string
	imageTag    string
}

// ecrImages, sorted by time (oldest first)
var ecrImages = make(map[int]ecrImage)

// Flags for CLI
var tagName string
var destImageDigest string
var profileName string
var regionName string
var ecrRepo string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ecr-image-retag",
	Short:   "A helper CLI to help with retagging process of ECR images i.e. moving 'latest' tag from image A to image B",
	Version: cliVersion,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		initClient()

		newImage, err := getImageByDigest()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newImageManifest := newImage.ImageManifest
		err = removeTagFromImages()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = addTagToImage(newImageManifest)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\nRetag success. The '%s' tag is now applied to an image with %s digest.\n", tagName, destImageDigest)
		fmt.Printf("Please restart your AWS ECS services now. Otherwise, they will continue running on the old '%s' image.\n", tagName)
	},
}

func initClient() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(regionName)},
		Profile: profileName,
	}))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().StringVarP(&tagName, "tag-name", "t", "", "The tag name that will be dropped from current images and to be applied to --new-image-digest")
	rootCmd.Flags().StringVarP(&destImageDigest, "new-image-digest", "d", "", "The new image digest that will receive the --tag-name")
	rootCmd.Flags().StringVarP(&profileName, "profile", "p", "", "The AWS profile name from ~/.aws/credentials file")
	rootCmd.Flags().StringVarP(&regionName, "region", "r", "", "The AWS region where the ECR repo is located")
	rootCmd.Flags().StringVarP(&ecrRepo, "ecr-repo", "e", "", "The AWS ECR repo name containing the images")

	rootCmd.MarkFlagRequired("tag-name")
	rootCmd.MarkFlagRequired("new-image-digest")
	rootCmd.MarkFlagRequired("profile")
	rootCmd.MarkFlagRequired("region")
	rootCmd.MarkFlagRequired("ecr-repo")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getImageByDigest() (*ecr.Image, error) {
	svc := ecr.New(sess)
	input := &ecr.BatchGetImageInput{
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageDigest: aws.String(destImageDigest),
			},
		},
		RepositoryName: aws.String(ecrRepo),
	}

	result, err := svc.BatchGetImage(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	image := result.Images[0]
	return image, nil
}

func removeTagFromImages() error {
	svc := ecr.New(sess)
	input := &ecr.BatchDeleteImageInput{
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageTag: aws.String(tagName),
			},
		},
		RepositoryName: aws.String(ecrRepo),
	}

	_, err := svc.BatchDeleteImage(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	return nil
}

func addTagToImage(newImageManifest *string) error {
	svc := ecr.New(sess)
	input := &ecr.PutImageInput{
		ImageManifest:  newImageManifest,
		ImageTag:       aws.String(tagName),
		RepositoryName: aws.String(ecrRepo),
	}

	_, err := svc.PutImage(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	return nil
}

// -------------------------------------------------------------------
// Everything below is not used (for now)
// -------------------------------------------------------------------

func sortImages() {
	svc := ecr.New(sess)
	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(ecrRepo),
	}

	result, err := svc.DescribeImages(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	for _, image := range result.ImageDetails {
		imageDigest := *image.ImageDigest
		imageTag := "untagged"

		var imageTags []*string = image.ImageTags
		if len(imageTags) > 0 {
			imageTag = *image.ImageTags[0]
		}

		img := ecrImage{
			imageDigest: imageDigest,
			imageTag:    imageTag,
		}

		pushedAt := *image.ImagePushedAt
		mapKey := int(pushedAt.UnixNano())
		ecrImages[mapKey] = img
	}

	keys := make([]int, 0, len(ecrImages))
	for k := range ecrImages {
		keys = append(keys, k)
	}
	sort.Ints(keys)
}
