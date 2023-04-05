package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

var (
	functionName = flag.String("function-name", "", "function name to update")
	region       = flag.String("region", "", "aws region to use")
	architecture = flag.String("architecture", "x86_64", "CPU architecture")
	dryRun       = flag.Bool("dry-run", false, "validate without applying change")
	publish      = flag.Bool("publish", false, "publish a new version of the function after updating")
	zipFile      = flag.String("zip-file", "", "zip file containing the function code")

	logger     = log.New(os.Stderr, "update-function-code", log.LstdFlags)
)

func getArchitecture() []types.Architecture {
	out := make([]types.Architecture, 1)

	if *architecture == "x86_64" {
		out[0] = types.ArchitectureX8664
	} else if *architecture == "arm64" {
		out[0] = types.ArchitectureArm64
	} else {
		logger.Fatal("Invalid architecture %s", *architecture)
		os.Exit(1)
	}
	return out
}

func getEncodedZipFile() (*bytes.Buffer, error) {
	reader, err := os.Open(*zipFile)

	if err != nil {
		logger.Printf("Failed to open zip file: %s", err.Error())
		return nil, err
	}

	defer reader.Close()

	zipBuffer := new(bytes.Buffer)

	encoder := base64.NewEncoder(base64.StdEncoding, zipBuffer)

	defer encoder.Close()

	ibuf := make([]byte, 4096)

	for {
		n, err := reader.Read(ibuf)
		if n > 0 {
			encoder.Write(ibuf[:n])
			continue
		}
		if err == io.EOF {
			break
		}
		return nil, err
	}

	return zipBuffer, nil
}

func main() {
	flag.Parse()

	if len(*functionName) == 0 {
		logger.Fatal("-function-name missing")
		os.Exit(1)
	}

	if len(*region) == 0 {
		logger.Fatal("-region missing")
		os.Exit(1)
	}

	if len(*zipFile) == 0 {
		logger.Fatal("-zip-file missing")
		os.Exit(1)
	}

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))

	if err != nil {
		logger.Fatalf("Failed to get aws config: %s", err.Error())
		os.Exit(1)
	}

	cli := lambda.NewFromConfig(cfg)

	getFunctionParams := &lambda.GetFunctionInput {
		FunctionName: aws.String(*functionName),
	}
	_, err = cli.GetFunction(ctx, getFunctionParams) // verify we can retrieve the function details

	if err != nil {
		logger.Fatalf("Failed to get function details: %s", err.Error())
		os.Exit(1)
	}

	zipBuffer, err := getEncodedZipFile()

	if err != nil {
		logger.Fatalf("Failed to retrieve zip file: %s", err.Error())
		os.Exit(1)
	}

	updateParams := &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(*functionName),
		Architectures: getArchitecture(),
		DryRun: *dryRun,
		Publish: *publish,
		ZipFile: zipBuffer.Bytes(),
	}

	_, err = cli.UpdateFunctionCode(ctx, updateParams)

	if err != nil {
		logger.Fatalf("Failed to update function code: %s", err.Error())
		os.Exit(1)
	}
}
