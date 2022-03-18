package pom

import (
	"encoding/xml"
	"os"
)

type Pom struct {
	XMLName xml.Name `xml:"project"`

	ArtifactId string `xml:"artifactId"`
	Version string `xml:"version"`
	Description string `xml:"description"`
}

func (p Pom) SuffixedJar(suffix string) string {
	return p.ArtifactId + "-" + p.Version + "-" + suffix + ".jar"
}

func LoadDefault() (*Pom, error) {
	p := Pom{}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pomBytes, err := os.ReadFile(pwd + "/../app/pom.xml")
	if err != nil {
		return nil, err
	}

	xml.Unmarshal(pomBytes, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
