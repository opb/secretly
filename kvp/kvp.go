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
	pos   int
	order []string
	data  map[string][]Kvp
}

type Provider interface {
	GetPairs(sourceNames []string) (*OrderedKvps, error)
	GetAll(sourceNames []string) (map[string]map[string]string, error)
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

func Compare(secretNames []string, p Provider) {
	if len(secretNames) < 2 {
		log.Fatalln("Please supply two or more secrets or files in order to compare them")
	}

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

// flipNamesWithKeys takes the map of secretNames, which have the key-value pairs as the contents
// and flips it so that the presences of the keys are recorded across different secretNames
func flipNamesWithKeys(names map[string]map[string]string) (map[string]map[string]bool, bool) {

	output := map[string]map[string]bool{}

	var nameList []string
	var keyList []string

	successFlag := true

	for name, contents := range names {
		nameList = append(nameList, name)
		for k := range contents {
			if _, ok := output[k]; !ok {
				keyList = append(keyList, k)
				output[k] = map[string]bool{name: true}
			} else {
				output[k][name] = true
			}
		}
	}

	// now add false to ones that are missing
	for _, k := range keyList {
		for _, n := range nameList {
			if _, ok := output[k][n]; !ok {
				output[k][n] = false
				successFlag = false
			}
		}
	}

	return output, successFlag
}
