package kvp

import (
	"encoding/json"

	"os"
	"strings"

	"fmt"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsKvp struct {
	Key         string
	Value       string
	SecretAlias string
}

type OrderedAwsKvps struct {
	pos       int
	order     []string
	maxKeyLen int
	data      map[string][]AwsKvp
}

func (o *OrderedAwsKvps) next() (string, []AwsKvp) {
	cur := o.pos
	if o.pos >= len(o.order) {
		return "", nil
	}
	out := o.data[o.order[cur]]
	o.pos++

	return o.order[cur], out
}

func (o *OrderedAwsKvps) Append(secretName string, awskvps []AwsKvp) {
	for _, awskvp := range awskvps {
		if len(awskvp.Key) > o.maxKeyLen {
			o.maxKeyLen = len(awskvp.Key)
		}
	}
	o.order = append(o.order, secretName)
	o.data[secretName] = awskvps
}

func (o *OrderedAwsKvps) EnvVars() map[string]string {
	out := make(map[string]string)

	for {
		_, awskvps := o.next()
		if awskvps == nil {
			break
		}
		for _, awskvp := range awskvps {
			out[awskvp.Key] = awskvp.Value
		}

	}

	return out
}

func (o *OrderedAwsKvps) KeysFromSecrets() map[string]string {
	out := make(map[string]string)

	for {
		secretName, awskvps := o.next()
		if awskvps == nil {
			break
		}
		for _, awskvp := range awskvps {
			out[awskvp.Key] = secretName
		}

	}

	return out
}

var cfg aws.Config

func PrintKeyListBySecret(secretNames []string) {
	ordered, err := kvpsFromAWS(secretNames)
	if err != nil {
		log.Fatalln("Could not get secrets from AWS:", err.Error())
	}

	format := fmt.Sprintf("%%-%vv", ordered.maxKeyLen+3)

	fmt.Printf(format, "Key")
	fmt.Print("From Secret\n")
	fmt.Printf(format, "===")
	fmt.Print("===========\n")

	for {
		secretName, awskvps := ordered.next()
		if awskvps == nil {
			break
		}
		for _, awskvp := range awskvps {
			fmt.Printf(format, awskvp.Key)
			fmt.Println(secretName)
		}

	}
}

func EnvPairs(secretsNames []string) ([]string, error) {
	initialVars := envVarsFromOS(os.Environ())

	orderedAwsVars, err := kvpsFromAWS(secretsNames)
	if err != nil {
		return nil, err
	}

	awsKvps := orderedAwsVars.EnvVars()

	for k, v := range awsKvps {
		initialVars[k] = v
	}

	var out []string
	for k, v := range initialVars {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}

	return out, nil
}

func envVarsFromOS(pairs []string) map[string]string {

	e := make(map[string]string)

	for _, v := range pairs {
		kvpair := strings.Split(v, "=")
		if len(kvpair) > 1 {
			e[kvpair[0]] = strings.Join(kvpair[1:], "")
		}
	}

	return e
}

func kvpsFromAWS(secretsNames []string) (*OrderedAwsKvps, error) {
	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	if region := os.Getenv("SECRETLY_REGION"); region != "" {
		cfg.Region = region
	}

	smSvc := secretsmanager.New(cfg)

	ordered := &OrderedAwsKvps{
		pos:       0,
		order:     []string{},
		maxKeyLen: 0,
		data:      make(map[string][]AwsKvp),
	}

	for _, secretName := range secretsNames {

		awskvps := []AwsKvp{}

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
			awskvps = append(awskvps, AwsKvp{Key: k, Value: v, SecretAlias: secretName})
		}
		ordered.Append(secretName, awskvps)
	}

	return ordered, nil
}
