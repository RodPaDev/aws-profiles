package lib

import (
	"encoding/json"
	"os"
)

func LoadStore(filePath string) (ProfileMap, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	profileMap := ProfileMap{}

	err = json.Unmarshal(data, &profileMap)

	if err != nil {
		return nil, err
	}

	return profileMap, err

}

func SaveStore(profileMap ProfileMap, filePath string) error {

	data, err := json.Marshal(profileMap)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
