package lib

type ProfileCredential struct {
	AccessKey string `json:"aws_access_key_id"`
	SecretKey string `json:"aws_secret_access_key"`
}

type ProfileConfig struct {
	Region string `json:"region,omitempty"`
	Output string `json:"output,omitempty"`
}

type Profile struct {
	Name       string            `json:"name"`
	Credential ProfileCredential `json:"credential"`
	Config     ProfileConfig     `json:"config"`
}

type ProfileMap map[string]Profile

func BuildProfileMap(data map[string]map[string]string) (ProfileMap, error) {
	profiles := ProfileMap{}

	for section, kv := range data {
		profiles[section] = Profile{
			Name: section,
			Credential: ProfileCredential{
				AccessKey: kv["aws_access_key_id"],
				SecretKey: kv["aws_secret_access_key"],
			},
			Config: ProfileConfig{
				Region: kv["region"],
				Output: kv["output"],
			},
		}

	}

	return profiles, nil
}

func (p *Profile) SetField(field, value string) {
	switch field {
	case "name":
		p.Name = value
	case "aws_access_key_id":
		p.Credential.AccessKey = value
	case "aws_secret_access_key":
		p.Credential.SecretKey = value
	case "region":
		p.Config.Region = value
	case "output":
		p.Config.Output = value
	}
}

func (p Profile) GetField(field string) string {
	switch field {
	case "name":
		return p.Name
	case "aws_access_key_id":
		return p.Credential.AccessKey
	case "aws_secret_access_key":
		return p.Credential.SecretKey
	case "region":
		return p.Config.Region
	case "output":
		return p.Config.Output
	default:
		return ""
	}
}
