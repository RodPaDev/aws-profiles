package lib

import (
	"os"
	"path"
	"testing"
)

func TestIntegrationFlow(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "aws-profiles-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create mock AWS directory
	awsDir := path.Join(tempDir, ".aws")
	err = os.Mkdir(awsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create mock .aws directory: %v", err)
	}

	// Create mock config file
	mockConfig :=
		`
		[default]
		region = us-west-2
		output = json

		[dev]
		region = us-east-1
		output = text

		[prod]
		region = us-west-1
		output = json
		`

	err = os.WriteFile(path.Join(awsDir, "config"), []byte(mockConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock config file: %v", err)
	}

	// Create mock credentials file
	mockCredentials :=
		`
		[default]
		aws_access_key_id = default-access-key
		aws_secret_access_key = default-secret-key

		[dev]
		aws_access_key_id = dev-access-key
		aws_secret_access_key = dev-secret-key

		[prod]
		aws_access_key_id = prod-access-key
		aws_secret_access_key = prod-secret-key
		`

	err = os.WriteFile(path.Join(awsDir, "credentials"), []byte(mockCredentials), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock credentials file: %v", err)
	}

	// Create a store file path
	storePath := path.Join(tempDir, "store.json")

	// Test BackupAWSProfileData
	err = BackupAWSProfileData(tempDir)
	if err != nil {
		t.Fatalf("Failed to backup AWS profile data: %v", err)
	}

	// Verify backup files were created
	_, err = os.Stat(path.Join(awsDir, "config.backup"))
	if err != nil {
		t.Errorf("Config backup file was not created: %v", err)
	}
	_, err = os.Stat(path.Join(awsDir, "credentials.backup"))
	if err != nil {
		t.Errorf("Credentials backup file was not created: %v", err)
	}

	// Override LoadAWSProfileData for testing
	originalHomePath := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomePath)

	// Load AWS profile data
	data, err := LoadAWSProfileData()
	if err != nil {
		t.Fatalf("Failed to load AWS profile data: %v", err)
	}

	// Verify data was loaded
	if len(data) != 2 {
		t.Fatalf("Expected 2 data items, got %d", len(data))
	}

	// Parse AWS profile data
	parsedData, _ := ParseAWSProfileData(data)

	// Verify parsed data
	if len(parsedData) != 3 { // default, profile dev, profile prod, dev, prod
		t.Fatalf("Expected 3 profiles, got %d", len(parsedData))
	}

	// Check if default profile exists
	defaultProfile, exists := parsedData["default"]
	if !exists {
		t.Fatalf("Default profile not found in parsed data")
	}

	// Verify default profile values
	if defaultProfile["region"] != "us-west-2" {
		t.Errorf("Expected default region to be us-west-2, got %s", defaultProfile["region"])
	}
	if defaultProfile["aws_access_key_id"] != "default-access-key" {
		t.Errorf("Expected default access key to be default-access-key, got %s", defaultProfile["aws_access_key_id"])
	}

	// Build profile map
	profileMap, err := BuildProfileMap(parsedData)
	if err != nil {
		t.Fatalf("Failed to build profile map: %v", err)
	}

	// Verify profile map
	if len(profileMap) != 3 {
		t.Fatalf("Expected 3 profiles in map, got %d", len(profileMap))
	}

	// Check default profile in map
	defaultMapProfile, exists := profileMap["default"]
	if !exists {
		t.Fatalf("Default profile not found in profile map")
	}

	// Verify default profile in map
	if defaultMapProfile.Config.Region != "us-west-2" {
		t.Errorf("Expected default region in map to be us-west-2, got %s", defaultMapProfile.Config.Region)
	}
	if defaultMapProfile.Credential.AccessKey != "default-access-key" {
		t.Errorf("Expected default access key in map to be default-access-key, got %s", defaultMapProfile.Credential.AccessKey)
	}

	// Test saving to store
	err = SaveStore(profileMap, storePath)
	if err != nil {
		t.Fatalf("Failed to save profile map to store: %v", err)
	}

	// Verify store file was created
	_, err = os.Stat(storePath)
	if err != nil {
		t.Errorf("Store file was not created: %v", err)
	}

	// Test loading from store
	loadedProfileMap, err := LoadStore(storePath)
	if err != nil {
		t.Fatalf("Failed to load profile map from store: %v", err)
	}

	// Verify loaded profile map
	if len(loadedProfileMap) != len(profileMap) {
		t.Fatalf("Expected %d profiles in loaded map, got %d", len(profileMap), len(loadedProfileMap))
	}

	// Test encoding profile data
	encodedData, err := EncodeAWSProfileData(defaultMapProfile)
	if err != nil {
		t.Fatalf("Failed to encode AWS profile data: %v", err)
	}

	// Verify encoded data
	if len(encodedData) != 2 {
		t.Fatalf("Expected 2 encoded data items, got %d", len(encodedData))
	}

	// Test saving AWS profile data
	err = SaveAWSProfileData(tempDir, encodedData)
	if err != nil {
		t.Fatalf("Failed to save AWS profile data: %v", err)
	}

	// Verify saved files
	configContent, err := os.ReadFile(path.Join(awsDir, "config"))
	if err != nil {
		t.Fatalf("Failed to read saved config file: %v", err)
	}
	if string(configContent) != encodedData[0] {
		t.Errorf("Saved config content does not match encoded data")
	}

	credContent, err := os.ReadFile(path.Join(awsDir, "credentials"))
	if err != nil {
		t.Fatalf("Failed to read saved credentials file: %v", err)
	}
	if string(credContent) != encodedData[1] {
		t.Errorf("Saved credentials content does not match encoded data")
	}
}
