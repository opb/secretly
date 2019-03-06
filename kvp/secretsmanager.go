package kvp

import (
	"encoding/json"

	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SMProvider struct{}

func (smp SMProvider) GetPairs(sourceNames []string) (*OrderedKvps, error) {

	smSvc, err := getSmService()
	if err != nil {
		return nil, err
	}

	ordered := &OrderedKvps{
		pos:   0,
		order: []string{},
		data:  make(map[string][]Kvp),
	}

	for _, secretName := range sourceNames {

		awskvps := []Kvp{}

		respData, err := getSecretsMapForSecret(smSvc, secretName)
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

func (smp SMProvider) GetAll(sourceNames []string) (map[string]map[string]string, error) {
	output := map[string]map[string]string{}

	smSvc, err := getSmService()
	if err != nil {
		return nil, err
	}

	for _, secretName := range sourceNames {
		respData, err := getSecretsMapForSecret(smSvc, secretName)
		if err != nil {
			return nil, err
		}

		output[secretName] = respData
	}

	return output, nil
}

func getSmService() (*secretsmanager.SecretsManager, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	if region := os.Getenv("SECRETLY_REGION"); region != "" {
		cfg.Region = region
	}

	return secretsmanager.New(cfg), nil
}

func getSecretsMapForSecret(svc *secretsmanager.SecretsManager, name string) (map[string]string, error) {
	req := svc.GetSecretValueRequest(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
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

	return respData, nil
}
