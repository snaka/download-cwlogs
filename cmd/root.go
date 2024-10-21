/*
Copyright © 2024 Shinji Nakamatsu <snaka@hey.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/spf13/cobra"
)

var (
	logGroupName   string
	logStreamName  string
	region         string
	outputFileName string
)

var rootCmd = &cobra.Command{
	Use:  "download-cwlogs",
	Long: `CloudWatch Logs から指定された LogGroup の LogStream を一括でダウンロードします`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a session
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			fmt.Println("Failed to load configuration,", err)
			return
		}

		// Create a CloudWatch Logs client
		svc := cloudwatchlogs.NewFromConfig(cfg)

		// Make params for GetLogEvents
		params := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(logGroupName),
			LogStreamName: aws.String(logStreamName),
			StartFromHead: aws.Bool(true),
		}

		// Create a CSV file
		file, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println("Failed to create a CSV file,", err)
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		writer.Write([]string{"timestamp", "message"})

		// Write log events
		for {
			resp, err := svc.GetLogEvents(context.TODO(), params)
			if err != nil {
				fmt.Println("Failed to get log events,", err)
				return
			}

			// Write log events to the CSV file
			for _, event := range resp.Events {
				writer.Write([]string{fmt.Sprintf("%v", *event.Timestamp), *event.Message})
			}

			fmt.Printf("NextToken: %v\n", aws.ToString(params.NextToken))

			// If not change next forward token, break the loop
			if aws.ToString(resp.NextForwardToken) == aws.ToString(params.NextToken) {
				break
			}

			params.NextToken = resp.NextForwardToken
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&logGroupName, "log-group", "g", "", "Log Group Name")
	rootCmd.Flags().StringVarP(&logStreamName, "log-stream", "s", "", "Log Stream Name")
	rootCmd.Flags().StringVarP(&region, "region", "r", "ap-northeast-1", "AWS Region")
	rootCmd.Flags().StringVarP(&outputFileName, "output-file", "o", "output.csv", "Output file name")

	rootCmd.MarkFlagRequired("log-group")
	rootCmd.MarkFlagRequired("log-stream")
}
