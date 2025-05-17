package lib

import (
	"bytes"
	"log"
	"os"
	"path"
	"strings"

	"github.com/rodpadev/aws-profiles/utils"
)

func LoadAWSProfileData() ([][]byte, error) {
	homeDir, err := utils.GetHomePath()

	if err != nil {
		log.Fatal(err)
	}

	dirPath := path.Join(homeDir, ".aws")

	config, err := os.ReadFile(path.Join(dirPath, "config"))
	if err != nil {
		return nil, err
	}

	credentials, err := os.ReadFile(path.Join(dirPath, "credentials"))
	if err != nil {
		return nil, err
	}

	return [][]byte{config, credentials}, nil
}

func ParseAWSProfileData(data [][]byte) map[string]map[string]string {
	result := map[string]map[string]string{}

	combined := append(data[0], '\n')
	combined = append(combined, data[1]...)

	lines := bytes.Split(combined, []byte{'\n'})

	var currentSection string
	for i, line := range lines {

		// remove comments and ignore empty liens
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}

		// get current section
		if line[0] == '[' && line[len(line)-1] == ']' {
			currentSection = string(line[1 : len(line)-1])
			if _, exists := result[currentSection]; !exists {
				result[currentSection] = make(map[string]string)
			}
			continue
		}

		// key value
		assignIndex := strings.IndexRune(string(line), '=')

		if assignIndex == -1 {
			continue
		}

		key := strings.TrimSpace(string(line[:assignIndex]))
		value := strings.TrimSpace(string(line[assignIndex+1:]))

		if currentSection == "" {
			currentSection = "no_category_" + string(i)
			if _, exists := result[currentSection]; !exists {
				result[currentSection] = map[string]string{}
			}
		}

		if key == "" {
			key = "no_key"
		}

		if value != "" {
			result[currentSection][key] = value
		}

	}

	return result

}
