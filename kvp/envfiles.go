package kvp

import "github.com/joho/godotenv"

type FileProvider struct{}

func (smp FileProvider) GetPairs(sourceNames []string) (*OrderedKvps, error) {

	ordered := &OrderedKvps{
		pos:   0,
		order: []string{},
		data:  make(map[string][]Kvp),
	}

	for _, sourceName := range sourceNames {

		kvps := []Kvp{}

		envData, err := godotenv.Read(sourceName)
		if err != nil {
			return nil, err
		}

		for k, v := range envData {
			kvps = append(kvps, Kvp{Key: k, Value: v, Source: sourceName})
		}
		ordered.append(sourceName, kvps)
	}

	return ordered, nil
}

func (smp FileProvider) GetAll(sourceNames []string) (map[string]map[string]string, error) {
	output := map[string]map[string]string{}

	for _, sourceName := range sourceNames {
		envData, err := godotenv.Read(sourceName)
		if err != nil {
			return nil, err
		}

		output[sourceName] = envData
	}

	return output, nil
}
