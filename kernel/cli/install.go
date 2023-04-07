package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/krishpranav/gomicro/register"
)

type ServiceCreate struct {
	register.Command
}

func (c *ServiceCreate) Register() {
	c.Signature = "service:create [service-name]"
	c.Description = "Create new Go-Web service"
}

func (c *ServiceCreate) Run() {
	if len(c.Args) == 0 {
		c.Help()
		return
	}

	fmt.Printf("Creating service %s...\n", c.Args)
	if err := c.clone(c.Args); err != nil {
		log.Fatalf("Error: %s", err)
	}

	if err := c.reset_git(); err != nil {
		log.Fatalf("Error: %s", err)
	}

	if err := c.update(); err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Println("Service created successfully!")
}

func (c *ServiceCreate) Help() {
	log.Println("Usage: create-service [service-name]")
}

func (c *ServiceCreate) clone(destination string) error {
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL:      "https://github.com/krishpranav/gomicro.git",
		Progress: nil,
	})

	return err
}

func (c *ServiceCreate) reset_git() error {
	path := fmt.Sprintf("%s/.git", c.Args)
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = c.Args
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (c *ServiceCreate) update() error {
	cmd := exec.Command("go", "get", "-u", "github.com/krishpranav/gomicro")
	cmd.Dir = c.Args

	return cmd.Run()
}
