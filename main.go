// Copyright (c) 2018, Mike Kaperys <mike@kaperys.io>.
// See LICENSE for licensing information.
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

var (
	// version stores the eip version. This value is overwritten at compile-time using ldflags.
	version = "0.2.0"

	// version stores the eip build commit hash. This value is overwritten at compile-time using ldflags.
	build = "dev"
)

const (
	_ = iota

	// The exitInvalidFilter exit code is returned when a provided `--filter` flag is malformed or invalid.
	exitInvalidFilter

	// The exitBadAWSSession exit code is returned when an AWS SDK session cannot be initialised.
	exitBadAWSSession

	// The exitCannotDescribeInstances exit code is returned when there was an error describing instances.
	// This is usually caused by a missing required environment variable (like AWS_REGION) or bad permissions.
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
	--all      show private IP addresses associated with all instance network interfaces (requires the --private flag)

Exit codes:
	%d   an invalid --filter flag was provided
	%d   an AWS session could not be initialised
	%d   EC2 instances could not be described
	%d   no EC2 instances were found
`, exitInvalidFilter, exitBadAWSSession, exitCannotDescribeInstances, exitNoInstancesFound)
	}

	var (
		versionFlag bool
		filterFlags filterFlags
		publicFlag  bool
		privateFlag bool
		allFlag     bool
	)

	flag.BoolVar(&versionFlag, "version", false, "show the eip version information")
	flag.Var(&filterFlags, "filter", "filter used to retrieve addresses")
	flag.BoolVar(&publicFlag, "public", false, "show instance public ip address")
	flag.BoolVar(&privateFlag, "private", false, "show the instact private ip address")
	flag.BoolVar(&allFlag, "all", false, "return all private IP addresses (from all network interfaces)")

	flag.Parse()

	if versionFlag {
		fmt.Printf("eip version %s %s/%s (build %s)\n", version, runtime.GOOS, runtime.GOARCH, build)
		return
	}

	resolve(parseFilters(filterFlags), publicFlag, privateFlag, allFlag)
}

// resolve creates an AWS SDK session and attempts to describe EC2 instances using the given
// flag data.
//
// Results are filtered using the provided slice of *ec2.Filters. For the matching results the
// following logic is applied:
//  If `pub` is true the associated public IP address is returned.
//  If `pri` is true the associated private IP address is returned.
//  If `pri` and `all` are both true, the private IP addresses associated with all instance network interfaces are returned.
func resolve(ff []*ec2.Filter, pub, pri, all bool) {
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
			if all && pri && len(inst.NetworkInterfaces) != 0 {
				for _, nic := range inst.NetworkInterfaces {
					fmt.Println(*nic.PrivateIpAddress)
				}

				continue
			}

			if pri && inst.PrivateIpAddress != nil {
				fmt.Println(*inst.PrivateIpAddress)
				continue
			}

			if pub && inst.PublicIpAddress != nil {
				fmt.Println(*inst.PublicIpAddress)
				continue
			}
		}
	}
}

// parseFilters attempts to parse the given --filter flags. Filter flags are provided
// using the following syntax: `--filter tag:SystemGroup=api,app --filter tag:Name=my-ec2-instance`.
//
// `--filter` supports all filters provided by the AWS Go SDK. See the following link:
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

// filterFlags is a type used to represent the `--filter` flag, since many are supported.
type filterFlags []string

// String returns the string representation of the provided filter flags.
func (i *filterFlags) String() string {
	return strings.Join(*i, ",")
}

// Set adds a new filter flag to the slice.
func (i *filterFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
