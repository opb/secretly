package kvp

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Kvp struct {
	Key    string
	Value  string
	Source string
}

type OrderedKvps struct {
	pos       int
	order     []string
	data      map[string][]Kvp
}

type Provider interface {
	GetPairs(sourceNames []string) (*OrderedKvps, error)
}

func List(secretNames []string, p Provider) {
	ordered, err := p.GetPairs(secretNames)
	if err != nil {
		log.Fatalln("Could not get secrets from Source:", err.Error())
	}

	format := fmt.Sprintf("%%-%vv", 23)

	fmt.Printf(format, "Key")
	fmt.Print("From Source\n")
	fmt.Printf(format, "===")
	fmt.Print("===========\n")

	for {
		secretName, kvps := ordered.next()
		if kvps == nil {
			break
		}
		for _, kvp := range kvps {
			fmt.Printf(format, kvp.Key)
			fmt.Println(secretName)
		}

	}
}

func (o *OrderedKvps) EnvVars() map[string]string {
	out := make(map[string]string)

	for {
		_, kvps := o.next()
		if kvps == nil {
			break
		}
		for _, kvp := range kvps {
			out[kvp.Key] = kvp.Value
		}

	}

	return out
}

func MergedEnvPairs(sourceNames []string, p Provider) ([]string, error) {
	// get map of env vars from the OS
	initialVars := envVarsFromOS(os.Environ())

	// now get ordered keys from the Provider
	orderedVars, err := p.GetPairs(sourceNames)
	if err != nil {
		return nil, err
	}

	kvps := orderedVars.EnvVars()

	for k, v := range kvps {
		initialVars[k] = v
	}

	var out []string
	for k, v := range initialVars {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}

	return out, nil
}

func (o *OrderedKvps) next() (string, []Kvp) {
	cur := o.pos
	if o.pos >= len(o.order) {
		return "", nil
	}
	out := o.data[o.order[cur]]
	o.pos++

	return o.order[cur], out
}

func (o *OrderedKvps) append(secretName string, kvps []Kvp) {
	o.order = append(o.order, secretName)
	o.data[secretName] = kvps
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
