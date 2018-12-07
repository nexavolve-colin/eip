package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// version stores the eip build version.
// This is overwritten at compile-time using ldflags.
var version = "dev"

const (
	exitInvalidFlag = iota
	exitBadSession
	exitDescribeInstances
)

func main() {
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
		fmt.Println("eip version", version)
		return
	}

	insIn := &ec2.DescribeInstancesInput{}

	if filterFlags != nil {
		var filters []*ec2.Filter
		for _, f := range filterFlags {
			if !strings.Contains(f, "=") {
				fmt.Fprintf(os.Stderr, "flag %s is invalid\n", f)
				os.Exit(exitInvalidFlag)
			}

			fs := strings.Split(f, "=")
			filters = append(filters, &ec2.Filter{
				Name:   aws.String(fs[0]),
				Values: aws.StringSlice(strings.Split(fs[1], ",")),
			})
		}

		insIn.Filters = filters
	}

	s, err := session.NewSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create the AWS session: %v\n", err)
		os.Exit(exitBadSession)
	}

	result, err := ec2.New(s).DescribeInstances(insIn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not describe EC2 instances: %v\n", err)
		os.Exit(exitDescribeInstances)
	}

	if len(result.Reservations) == 0 {
		return
	}

	for _, r := range result.Reservations {
		for _, inst := range r.Instances {
			if publicFlag && inst.PublicIpAddress != nil {
				fmt.Println(*inst.PublicIpAddress)
			}

			if privateFlag && inst.PrivateIpAddress != nil {
				fmt.Println(*inst.PrivateIpAddress)
			}
		}
	}
}
