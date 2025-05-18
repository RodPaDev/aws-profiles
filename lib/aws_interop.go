package lib

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"github.com/rodpadev/aws-profiles/utils"
)

const DEFAULT_SECTION = "[default]"

func BackupAWSProfileData(homeDir string) error {

	dirPath := path.Join(homeDir, ".aws")
	configPath := path.Join(dirPath, "config")
	credentialsPath := path.Join(dirPath, "credentials")

	backupConfigPath := fmt.Sprintf("%s.backup", configPath)
	backupCredentialsPathPath := fmt.Sprintf("%s.backup", credentialsPath)

	if err := utils.CopyFile(configPath, backupConfigPath); err != nil {
		return err
	}
	if err := utils.CopyFile(credentialsPath, backupCredentialsPathPath); err != nil {
		return err
	}

	return nil
}

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
		line = bytes.TrimSpace(line)

		// remove comments and ignore empty liens
		if len(line) == 0 || line[0] == '#' || line[0] == ';' {
			continue
		}

		// get current section
		if line[0] == '[' && line[len(line)-1] == ']' {
			currentSection = strings.TrimSpace(strings.Replace(string(line[1:len(line)-1]), "profile", "", 1))

			if currentSection == "default" {
				currentSection = "⚠️ RENAME [awscli default]"
			}

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
			currentSection = "no_category_" + strconv.Itoa(i)
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

func EncodeAWSProfileData(selectedProfile Profile) ([2]string, error) {
	configFile := DEFAULT_SECTION + "\n"
	credentialsFile := DEFAULT_SECTION + "\n"

	credVal := reflect.ValueOf(selectedProfile.Credential)
	credType := reflect.TypeOf(selectedProfile.Credential)

	for i := 0; i < credVal.NumField(); i++ {
		tag := credType.Field(i).Tag.Get("json")
		value := credVal.Field(i).String()
		credentialsFile += fmt.Sprintf("%s = %s\n", tag, value)
	}

	configVal := reflect.ValueOf(selectedProfile.Config)
	configType := reflect.TypeOf(selectedProfile.Config)

	for i := 0; i < configVal.NumField(); i++ {
		tag := configType.Field(i).Tag.Get("json")
		val := configVal.Field(i).String()
		if val != "" {
			configFile += fmt.Sprintf("%s = %s\n", tag, val)
		}
	}

	return [2]string{configFile, credentialsFile}, nil
}

func SaveAWSProfileData(homeDir string, awsData [2]string) error {

	configFile := awsData[0]
	credentialsFile := awsData[1]

	awsDir := path.Join(homeDir, ".aws")

	if err := os.WriteFile(path.Join(awsDir, "config"), []byte(configFile), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(awsDir, "credentials"), []byte(credentialsFile), 0600); err != nil {
		return err
	}

	return nil
}
