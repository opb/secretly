package kvp

import (
	"encoding/json"

	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SMProvider struct{}

func (smp SMProvider)GetPairs(sourceNames []string) (*OrderedKvps, error) {

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	if region := os.Getenv("SECRETLY_REGION"); region != "" {
		cfg.Region = region
	}

	smSvc := secretsmanager.New(cfg)

	ordered := &OrderedKvps{
		pos:       0,
		order:     []string{},
		data:      make(map[string][]Kvp),
	}

	for _, secretName := range sourceNames {

		awskvps := []Kvp{}

		req := smSvc.GetSecretValueRequest(&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(secretName),
		})

		resp, err := req.Send()
		if err != nil {
			return nil, err
		}

		var respData map[string]string

		err = json.Unmarshal([]byte(*resp.SecretString), &respData)
		if err != nil {
			return nil, err
		}

		for k, v := range respData {
			awskvps = append(awskvps, Kvp{Key: k, Value: v, Source: secretName})
		}
		ordered.append(secretName, awskvps)
	}

	return ordered, nil
}
