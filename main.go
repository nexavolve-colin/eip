// Copyright (c) 2018, Mike Kaperys <mike@kaperys.io>
// See LICENSE for licensing information
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// version stores the eip version. This is overwritten at compile-time using ldflags.
var version = "0.1.1"

// version stores the eip build commit hash. This is overwritten at compile-time using ldflags.
var build = "dev"

const (
	// The exitInvalidFilter exit code is returned when a provided `--filter` flag is malformed or invalid.
	exitInvalidFilter = 1 << iota

	// The exitBadAWSSession exit code is returned when an AWS SDK session cannot be initialised.
	exitBadAWSSession

	// The exitCannotDescribeInstances exit code is returned when there was an error describing instances.
	// This is usually caused by missing environment variables (like AWS_REGION) or permissions.
	exitCannotDescribeInstances

	// The exitNoInstancesFound exit code is returned when no EC2 instances are found using the given filters.
	exitNoInstancesFound
)

func main() {
	flag.Usage = func() {
		fmt.Printf(`usage: eip [flags]

If no flags are provided, nothing will be returned. Either the --public
or --private flag must be provided. eip makes no assumptions about output.

	--version  show version information and exit
	--filter   filters results using the provided key-value pair (FilterName=FilterValue)
	--public   show the instance public IP address
	--private  show the instance private IP address

Exit Codes:
	%d   an invalid --filter flag was provided
	%d   an AWS session could not be initialised
	%d   EC2 instances could not be described
	%d   no EC2 instances were found\n
`, exitInvalidFilter, exitBadAWSSession, exitCannotDescribeInstances, exitNoInstancesFound)
	}

	var (
		versionFlag bool
		filterFlags filterFlags
		publicFlag  bool
		privateFlag bool
	)

	flag.BoolVar(&versionFlag, "version", false, "show the eip version information")
	flag.Var(&filterFlags, "filter", "filter used to retrieve addresses")
	flag.BoolVar(&publicFlag, "public", false, "show instance public ip address")
	flag.BoolVar(&privateFlag, "private", false, "show the instact private ip address")

	flag.Parse()

	if versionFlag {
		fmt.Printf("eip version %s %s/%s (build %s)\n", version, runtime.GOOS, runtime.GOARCH, build)
		return
	}

	resolve(parseFilters(filterFlags), publicFlag, privateFlag)
}

// resolve parses the given flag data and uses the AWS SDK to describe EC2 instances.
func resolve(ff []*ec2.Filter, pub, priv bool) {
	insIn := &ec2.DescribeInstancesInput{Filters: ff}

	s, err := session.NewSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create the AWS session: %v\n", err)
		os.Exit(exitBadAWSSession)
	}

	result, err := ec2.New(s).DescribeInstances(insIn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not describe EC2 instances: %v\n", err)
		os.Exit(exitCannotDescribeInstances)
	}

	if len(result.Reservations) == 0 {
		fmt.Fprintln(os.Stderr, "no EC2 instances found")
		os.Exit(exitNoInstancesFound)
	}

	for _, r := range result.Reservations {
		for _, inst := range r.Instances {
			if pub && inst.PublicIpAddress != nil {
				fmt.Println(*inst.PublicIpAddress)
			}

			if priv && inst.PrivateIpAddress != nil {
				fmt.Println(*inst.PrivateIpAddress)
			}
		}
	}
}

// parseFilters attempt to parse the given --filter flags. Filter flags are provided
// using the following syntax: `--filter tag:SystemGroup=api,app --filter tag:Name=my-ec2-instance`.
// `--filter` supports flags provided by the AWS SDK:
// https://github.com/datacratic/aws-sdk-go/blob/master/service/ec2/api.go#L9532-L9754
func parseFilters(ff filterFlags) []*ec2.Filter {
	var filters []*ec2.Filter

	if ff != nil {
		for _, f := range ff {
			if !strings.Contains(f, "=") {
				fmt.Fprintf(os.Stderr, "filter %q is invalid\n", f)
				os.Exit(exitInvalidFilter)
			}

			fs := strings.Split(f, "=")
			filters = append(filters, &ec2.Filter{
				Name:   aws.String(fs[0]),
				Values: aws.StringSlice(strings.Split(fs[1], ",")),
			})
		}
	}

	return filters
}

// filterFlags is a type used to represent the `--filter` flag, since there can be many.
type filterFlags []string

// String returns the string representation of the provided filter flags.
func (i *filterFlags) String() string {
	return strings.Join(*i, ",")
}

// Set sets a new filter flag in the slice.
func (i *filterFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
