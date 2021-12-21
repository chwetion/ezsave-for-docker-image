package model

type RegistryAuth struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	Address  string `yaml:"address,omitempty"`
}

type Image struct {
	Registry string `yaml:"registry,omitempty"`
	Project  string `yaml:"project,omitempty"`
	Name     string `yaml:"name,omitempty"`
}

type Configuration struct {
	Auths         []RegistryAuth       `yaml:"auths,omitempty"`
	Packages      []CompressedPackages `yaml:"packages,omitempty"`
	DefaultTarget Image                `yaml:"defaultTarget,omitempty"`
	DefaultFrom   Image                `yaml:"defaultFrom,omitempty"`
}

type Content struct {
	Name   string `yaml:"name,omitempty"`
	From   Image  `yaml:"from,omitempty"`
	Target Image  `yaml:"export,omitempty"`
}

type CompressedPackages struct {
	File    string    `yaml:"file,omitempty"`
	Content []Content `yaml:"content,omitempty"`
}

func (i *Image) DeepCopy() *Image {
	return &Image{
		Registry: i.Registry,
		Project:  i.Project,
		Name:     i.Name,
	}
}
