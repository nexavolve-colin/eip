package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestParseFlags(t *testing.T) {
	tt := []struct {
		Name        string
		FilterFlags filterFlags
		Expected    []*ec2.Filter
	}{
		{
			Name:        "2 valid flags, 1 multiple value",
			FilterFlags: []string{"tag:SystemGroup=api,app", "tag:Name=my-ec2-instance"},
			Expected: []*ec2.Filter{
				{Name: aws.String("tag:SystemGroup"), Values: aws.StringSlice([]string{"api", "app"})},
				{Name: aws.String("tag:Name"), Values: aws.StringSlice([]string{"my-ec2-instance"})},
			},
		},
		{
			Name:        "1 invalid flag",
			FilterFlags: []string{"tag:SystemGroup"},
			Expected:    nil,
		},
	}

	for _, tc := range tt {
		if f := parseFilters(tc.FilterFlags); !reflect.DeepEqual(f, tc.Expected) {
			t.Fatalf("TestParseFlags: %q failed: have %+v, want %+v\n", tc.Name, f, tc.Expected)
		}
	}
}
