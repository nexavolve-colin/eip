package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestParseFilters(t *testing.T) {
	tt := []struct {
		Name        string
		FilterFlags filterFlags
		Expected    []*ec2.Filter
	}{
		{
			Name:        "1 valid filter",
			FilterFlags: []string{"tag:Name=my-ec2-instance"},
			Expected: []*ec2.Filter{
				{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{"my-ec2-instance"})},
			},
		},
		{
			Name:        "1 valid filter, 1 multiple value",
			FilterFlags: []string{"tag:Name=my-ec2-instance,my-other-ec2-instance"},
			Expected: []*ec2.Filter{
				{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{"my-ec2-instance", "my-other-ec2-instance"})},
			},
		},
		{
			Name:        "2 valid filters, 1 multiple value",
			FilterFlags: []string{"tag:SystemGroup=api,app", "tag:Name=my-ec2-instance"},
			Expected: []*ec2.Filter{
				{Name: aws.String("tag:SystemGroup"), Values: aws.StringSlice([]string{"api", "app"})},
				{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{"my-ec2-instance"})},
			},
		},
		{
			Name:        "2 valid filters, 2 multiple value",
			FilterFlags: []string{"tag:SystemGroup=api,app", "tag:Name=my-ec2-instance,my-other-ec2-instance"},
			Expected: []*ec2.Filter{
				{Name: aws.String("tag:SystemGroup"), Values: aws.StringSlice([]string{"api", "app"})},
				{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{"my-ec2-instance", "my-other-ec2-instance"})},
			},
		},
	}

	for _, tc := range tt {
		if f := parseFilters(tc.FilterFlags); !reflect.DeepEqual(f, tc.Expected) {
			t.Fatalf("TestParseFilters: %q failed: have %+v, want %+v\n", tc.Name, f, tc.Expected)
		}
	}
}
