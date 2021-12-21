package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chwetion/ezsave-for-docker-image/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Client struct {
	*client.Client
	ctx   context.Context
	Cfg   *model.Configuration
	auths map[string]string
	Path  string
}

func NewClient(configuration *model.Configuration, path string) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	c := &Client{
		Client: cli,
		ctx:    context.Background(),
		Cfg:    configuration,
		auths:  make(map[string]string),
		Path:   path,
	}
	c.InitAuth()
	return c, nil
}

func (c *Client) InitAuth() {
	for _, auth := range c.Cfg.Auths {
		authCfg := types.AuthConfig{
			Username:      auth.Username,
			Password:      auth.Password,
			ServerAddress: auth.Address,
		}
		encodedAuth, _ := json.Marshal(authCfg)
		c.auths[auth.Address] = base64.StdEncoding.EncodeToString(encodedAuth)
	}
}

func (c *Client) getFileFullPath(relativePath string) string {
	return c.Path + "/" + relativePath
}

func (c *Client) Forward() {
	for _, p := range c.Cfg.Packages {
		for _, content := range p.Content {
			if err := c.pullAndTag(content.From, content.Name); err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			if err := c.tagAndPush(content.Name, content.Target); err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
		}
	}
}

func (c *Client) Save() {
	for _, p := range c.Cfg.Packages {
		imageNames := make([]string, 0, 5)
		for _, content := range p.Content {
			if err := c.pullAndTag(content.From, content.Name); err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			imageNames = append(imageNames, content.Name)
		}
		if err := c.saveFile(imageNames, c.getFileFullPath(p.File)); err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
	}
}

func (c *Client) Load() {
	for _, p := range c.Cfg.Packages {
		if err := c.loadFile(c.getFileFullPath(p.File)); err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		for _, content := range p.Content {
			if err := c.tagAndPush(content.Name, content.Target); err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
		}
	}
}

func getImageName(image model.Image, def model.Image, defName string) string {
	result := image.DeepCopy()
	result.Registry = getImageRegistry(image, def)
	if image.Project == "" {
		result.Project = def.Project
	}
	if image.Name == "" {
		result.Name = defName
	}
	arr := []string{result.Registry, result.Project, result.Name}
	return strings.Join(arr, "/")
}

func getImageRegistry(image model.Image, def model.Image) string {
	if image.Registry == "" {
		return def.Registry
	}
	return image.Registry
}

func (c *Client) pullAndTag(image model.Image, targetName string) error {
	imageName := getImageName(image, c.Cfg.DefaultFrom, targetName)
	registry := getImageRegistry(image, c.Cfg.DefaultFrom)
	// docker pull
	rc, err := c.ImagePull(c.ctx, imageName, types.ImagePullOptions{
		RegistryAuth: c.auths[registry],
	})
	if err != nil {
		return fmt.Errorf("skip to pull image %s: %s\n", imageName, err)
	}
	defer rc.Close()
	_, _ = io.Copy(os.Stdout, rc)
	// docker tag
	err = c.ImageTag(c.ctx, imageName, targetName)
	if err != nil {
		return fmt.Errorf("skip to tag image %s -> %s: %s\n", imageName, targetName, err)
	}
	return nil
}

func (c *Client) tagAndPush(rawImage string, targetImage model.Image) error {
	targetImageName := getImageName(targetImage, c.Cfg.DefaultTarget, rawImage)
	targetImageRegistry := getImageRegistry(targetImage, c.Cfg.DefaultTarget)
	err := c.ImageTag(c.ctx, rawImage, targetImageName)
	if err != nil {
		return fmt.Errorf("skip to tag image %s -> %s: %s\n", rawImage, targetImageName, err)
	}
	reader, err := c.ImagePush(c.ctx, targetImageName, types.ImagePushOptions{
		RegistryAuth: c.auths[targetImageRegistry],
	})
	if err != nil {
		return fmt.Errorf("failed push image %s: %s\n", targetImageName, err)
	}
	defer reader.Close()
	_, _ = io.Copy(os.Stdout, reader)
	return nil
}

func (c *Client) saveFile(images []string, targetFileName string) error {
	rc, err := c.Client.ImageSave(c.ctx, images)
	if err != nil {
		return fmt.Errorf("save image error: %s\n", err)
	}
	defer rc.Close()
	of, err := os.Create(targetFileName)
	if err != nil {
		return fmt.Errorf("create file error: %s\n", err)
	}
	defer of.Close()
	_, err = io.Copy(of, rc)
	if err != nil {
		return fmt.Errorf("write file error: %s\n", err)
	}
	return nil
}

func (c *Client) loadFile(fileName string) error {
	reader, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("read file %s error: %s\n", fileName, err)
	}
	defer reader.Close()
	res, err := c.ImageLoad(c.ctx, reader, false)
	if err != nil {
		return fmt.Errorf("failed load compressed package: %s\n", err)
	}
	_, _ = io.Copy(os.Stdout, res.Body)
	return nil
}
