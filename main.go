package main

import (
	"io/ioutil"

	"errors"

	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type ProjectConfig struct {
	Name string
	IdePath string
	RootPath string
	Lang string
}

type IdeConfig struct {
	Lang string
	Path string
}

type ProjectConfigsList struct {
	Projects []ProjectConfig
	Ides []IdeConfig 
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadProjectConfigs(pathToProjectConfigsFile string) ProjectConfigsList {
	yamlConfig, err := ioutil.ReadFile(pathToProjectConfigsFile)
	handleError(err)

	projectsConfigsList := ProjectConfigsList{}

	err = yaml.Unmarshal(yamlConfig, &projectsConfigsList)
	handleError(err)
	
	return projectsConfigsList
}

func searchProjectConfig(projectName string, configs []ProjectConfig) (ProjectConfig, error) {
	for _, config := range configs {
		if config.Name == projectName {
			return config, nil
		}
	}
	return  ProjectConfig{}, errors.New("project with this name not fount. Check your config file")
}

func searchIdeConfig(project ProjectConfig, ides []IdeConfig) (IdeConfig, error) {
	for _, ide := range ides {
		if project.Lang == ide.Lang {
			return ide, nil
		}
	}
	return  IdeConfig{}, errors.New("ide for this project not fount. Check your config file")
}

func main()  {
	home, _ := os.LookupEnv("HOME")
	pathToProjectConfigsFile :=  home + "/.config/projects.yaml"

	args := os.Args
	if len(args) < 2 {
		panic("Not found project name")	
	}
	projectName := args[1]
	
	projectConfigsList := loadProjectConfigs(pathToProjectConfigsFile)
	projectConfigs := projectConfigsList.Projects
	idesConfigs := projectConfigsList.Ides
	
	projectConfig, err := searchProjectConfig(projectName, projectConfigs)
	handleError(err)
	ideConfig, err := searchIdeConfig(projectConfig, idesConfigs)
	handleError(err)
	
	ideRunCmd := exec.Cmd{
		Path: ideConfig.Path,
		Args: []string{ideConfig.Path, projectConfig.IdePath},
		Stdout: os.Stdout,
		Stdin: os.Stdin,
	}
	
	err = ideRunCmd.Run()
	handleError(err)
}
