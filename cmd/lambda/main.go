// main.go
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"leetsignal/internal/client"
	"leetsignal/internal/config"
	"leetsignal/internal/ntfy"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func Handler(ctx context.Context) (string, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./cmd/lambda/config.json" // default path used for local testing
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return "", err
	}

	client := client.NewClient()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("Failed to load AWS config: %v", err)
		return "", err
	}
	ssmClient := ssm.NewFromConfig(awsCfg)

	for _, profile := range cfg.Profiles {
		currentSolved, err := client.GetNumSolved(profile)
		if err != nil {
			log.Printf("Error fetching solved problems for %s: %v", profile, err)
			ntfy.SendErrorAlert(cfg, fmt.Sprintf("Error fetching solved problems for %s: %v", profile, err))
		}

		paramName := fmt.Sprintf("/leetcode/%s/solved_count", profile)
		var storedSolved int

		paramResp, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String(paramName),
		})
		if err != nil {
			var nfe *types.ParameterNotFound
			if !errors.As(err, &nfe) {
				log.Printf("Failed to get parameter %s: %v", paramName, err)
			} else {
				storedSolved = 0 // first run, parameter doesn't exist yet
			}
		} else if paramResp.Parameter != nil {
			storedSolved, _ = strconv.Atoi(*paramResp.Parameter.Value)
		}

		if currentSolved > storedSolved {
			submission, err := client.GetMostRecentAcceptedSubmission(profile)
			if err != nil {
				log.Printf("Error fetching latest submission for %s: %v", profile, err)
				ntfy.SendErrorAlert(cfg, fmt.Sprintf("Error fetching latest submission for %s: %v", profile, err))
				continue
			}
			ntfy.SendLeetSignal(cfg, profile, *submission)

			// 4. Update Parameter Store
			_, err = ssmClient.PutParameter(ctx, &ssm.PutParameterInput{
				Name:      aws.String(paramName),
				Value:     aws.String(strconv.Itoa(currentSolved)),
				Type:      "String",
				Overwrite: aws.Bool(true),
			})
			if err != nil {
				log.Printf("Failed to update Parameter Store for %s: %v", profile, err)
			}
		} else {
			log.Printf("%s has no new solves (current: %d, stored: %d)", profile, currentSolved, storedSolved)
		}
	}

	return "done", nil
}

// For local testing, run with `go run cmd/lambda/main.go -local`
func main() {
	local := flag.Bool("local", false, "Run locally without Lambda")
	flag.Parse()

	if *local {
		ctx := context.Background()
		_, err := Handler(ctx)
		if err != nil {
			fmt.Println("Error running handler locally:", err)
			return
		}
	} else {
		lambda.Start(Handler)
	}
}
