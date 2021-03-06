package secrets

import (
	"io/ioutil"
	"os/exec"

	mig "github.com/mozilla/mig/pgp"
	"gopkg.in/yaml.v2"
)

func generatePGPKeyPair(email, name, desc string) (KeyPair, error) {
	pubkey, privkey, fp, err := mig.GenerateKeyPair(name, desc, email)
	if err != nil {
		return KeyPair{}, err
	}
	return KeyPair{
		Public:      string(pubkey),
		Private:     string(privkey),
		Fingerprint: string(fp),
	}, nil
}

type SopsYaml struct {
	CreationRules []CreationRule `yaml:"creation_rules"`
}

type CreationRule struct {
	PGP string `yaml:"pgp"`
}

func GeneratePGPSopsYaml(keyPair KeyPair) error {
	sops := SopsYaml{
		CreationRules: []CreationRule{
			{PGP: keyPair.Fingerprint},
		},
	}

	content, err := yaml.Marshal(sops)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(".sops.yaml", content, 0644)
}

func LoadPGPSopsYaml() (SopsYaml, error) {
	contents, err := ioutil.ReadFile(".sops.yaml")
	if err != nil {
		return SopsYaml{}, err
	}

	var sops SopsYaml
	return sops, yaml.Unmarshal(contents, &sops)
}

func ImportPGPKeypair(keyPair KeyPair) error {
	cmd := exec.Command("gpg", "--import")
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write(([]byte)(keyPair.Private)); err != nil {
		return err
	}
	in.Close()
	return cmd.Wait()
}

func RemovePGPKeypairs(fingerprints []string) error {
	for _, fingerprint := range fingerprints {
		if err := exec.Command("gpg", "--batch", "--yes", "--delete-secret-keys", fingerprint).Run(); err != nil {
			return err
		}
		if err := exec.Command("gpg", "--batch", "--yes", "--delete-key", fingerprint).Run(); err != nil {
			return err
		}
	}
	return nil
}
