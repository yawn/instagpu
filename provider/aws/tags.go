package aws

import (
	cf "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	ec2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Tags map[string]string

func (t Tags) ToCF() (tags []cf.Tag) {

	for k, v := range t {
		tags = append(tags, cf.Tag{
			Key:   &k,
			Value: &v,
		})
	}

	return

}

func (t Tags) ToEC2() (tags []ec2.Tag) {

	for k, v := range t {
		tags = append(tags, ec2.Tag{
			Key:   &k,
			Value: &v,
		})
	}

	return

}
