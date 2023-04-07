package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

var (
	functionArn        = flag.String("function-arn", "", "function arn to update")
	functionNamePrefix = flag.String("function-name-prefix", "", "function name prefix to update")
	region             = flag.String("region", "", "aws region to use")
	architecture       = flag.String("architecture", "x86_64", "CPU architecture")
	dryRun             = flag.Bool("dry-run", false, "validate without applying change")
	publish            = flag.Bool("publish", false, "publish a new version of the function after updating")
	zipFile            = flag.String("zip-file", "", "zip file containing the function code")

	logger     = log.New(os.Stderr, "update-function-code", log.LstdFlags)
)

func getArchitecture() []types.Architecture {
	out := make([]types.Architecture, 1)

	if *architecture == "x86_64" {
		out[0] = types.ArchitectureX8664
	} else if *architecture == "arm64" {
		out[0] = types.ArchitectureArm64
	} else {
		logger.Fatalf("Invalid architecture %s", *architecture)
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

func getFunctionArn(ctx context.Context, client *lambda.Client) (string, error) {
	if len(*functionArn) > 0 {
		return *functionArn, nil
	}

	maxItems := new(int32)
	*maxItems = 50

	params := &lambda.ListFunctionsInput{
		MaxItems: maxItems,
	}
	for {
		res, err := client.ListFunctions(ctx, params)
		if err != nil {
			return "", err
		}
		for _, f := range(res.Functions) {
			if strings.Index(*f.FunctionName, *functionNamePrefix) == 0 {
				logger.Printf("Found match %s -> %s", *f.FunctionName, *f.FunctionArn)
				return *f.FunctionArn, nil
			}

		}
		if res.NextMarker == nil || len(*res.NextMarker) == 0 {
			return "", nil
		}
		params.Marker = res.NextMarker
	}

}

func main() {
	flag.Parse()

	if len(*functionArn) == 0 && len(*functionNamePrefix) == 0{
		logger.Fatal("one of -function-arn or -function-name-prefix needs to be passed")
		os.Exit(1)
	}

	if len(*functionArn) > 0 && len(*functionNamePrefix) > 0{
		logger.Fatal("only one of -function-arn or -function-name-prefix can be passed")
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

	arn, err := getFunctionArn(ctx, cli)

	if err != nil {
		logger.Fatalf("Failed to resolve lambda name: %s", err.Error())
		os.Exit(1)
	}

	if len(arn) == 0 {
		logger.Fatalf("Failed to find a lambda with prefix %s", *functionNamePrefix)
		os.Exit(1)
	}

	getFunctionParams := &lambda.GetFunctionInput {
		FunctionName: aws.String(arn),
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
		FunctionName: aws.String(arn),
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
