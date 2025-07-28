package setup

import (
	"os"
	"path/filepath"
	"time"

	"github.com/orochaa/go-clack/prompts"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type model struct {
	ProjectName           string
	Default               bool
	FolderStructure       string
	ConfigFile            string
	CreateStructureFolder bool
	CustomFolders         []string
	Features              []Feature
	FeaturesOptions       []string
}

func Init() *model {
	return &model{
		ProjectName:           "",
		Default:               false,
		FolderStructure:       "",
		ConfigFile:            "",
		CreateStructureFolder: false,
		Features:              []Feature{},
		FeaturesOptions:       []string{},
	}
}

func RunSetup(ctx *cli.Context) error {
	m := Init()
	m.InitSetup()
	return nil
}

func (m *model) InitSetup() {
	m.SelectProjectName()

	m.DefaultSettings()

	m.SelectFolderStructure()

	m.SelectDefaultConfig()

	if !m.Default {
		m.AddFeatures()
	}

	if m.FolderStructure != "custom" {
		m.CreateStructure()
	}

	m.Setup()
}

type Config struct {
	Version       int                               `yaml:"version"`
	ProjectName   string                            `yaml:"project_name"`
	Debug         bool                              `yaml:"debug"`
	PathStructure string                            `yaml:"path_structure"`
	CustomPaths   map[string]string                 `yaml:"custom_path,omitempty"`
	Config        ConfigDetails                     `yaml:"config"`
	Features      map[string]map[string]interface{} `yaml:"features,omitempty"`
}

type ConfigDetails struct {
	File     string         `yaml:"file"`
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Schema   string `yaml:"schema"`
}

func (m *model) Setup() {

	cfg := &Config{
		Version:       1,
		ProjectName:   m.ProjectName,
		Debug:         false,
		PathStructure: m.FolderStructure,
	}

	var defaultFeatures = map[string]map[string]interface{}{
		"microservice_architecture": {
			"regex":        ".*service.*",
			"directory":    "./modules",
			"multiple_dbs": false,
		},
		"rest_api": {
			"prefix": "/api/v1",
		},
	}
	cfg.Features = make(map[string]map[string]interface{})
	featureEnabled := make(map[string]bool)
	prompts.Tasks(
		[]prompts.Task{
			{
				Title: "Initializing project....",
				Task: func(message func(msg string)) (string, error) {
					start := time.Now()
					_ = os.MkdirAll(m.ProjectName, 0755)

					if m.FolderStructure == "custom" {
						cfg.CustomPaths = make(map[string]string)
						for _, folder := range m.CustomFolders {
							cfg.CustomPaths[folder] = "./" + folder
						}
					}

					cfg.Config = ConfigDetails{
						File: m.ConfigFile,
						Database: DatabaseConfig{
							Driver:   "DB_DRIVER",
							Host:     "DB_HOST",
							Port:     "DB_PORT",
							User:     "DB_USER",
							Password: "DB_PASSWORD",
							Name:     "DB_NAME",
							Schema:   "DB_SCHEMA",
						},
					}

					for _, f := range m.Features {
						featureMap := map[string]interface{}{}
						if def, ok := defaultFeatures[f.Name]; ok {
							for k, v := range def {
								featureMap[k] = v
							}
						}
						featureMap["enabled"] = f.Enabled
						cfg.Features[f.Name] = featureMap
						featureEnabled[f.Name] = f.Enabled
					}

					if m.CreateStructureFolder {
						if featureEnabled["microservice_architecture"] {
							serviceBase := filepath.Join(m.ProjectName, "modules", "user-service")
							createMicroserviceStructure(serviceBase, m.FolderStructure, m.CustomFolders)
						} else {
							createProjectStructure(m.ProjectName, m.FolderStructure, m.CustomFolders)
						}

						createGoFilesWithFeatures(m.ProjectName, featureEnabled)
					}

					_ = SaveConfigFile(cfg, filepath.Join(m.ProjectName, "alchemist.yaml"))

					elapsed := time.Since(start)
					if elapsed < 2*time.Second {
						time.Sleep(2*time.Second - elapsed)
					}
					return "Project initialized successfully", nil
				},
			},
		},
		prompts.SpinnerOptions{},
	)

	m.FinalSteps(featureEnabled)
}

func SaveConfigFile(cfg *Config, filename string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

var layeredSubfolders = []string{"handler", "repository", "model", "service"}
var cleanArchSubfolders = map[string][]string{
	"domain":         {"model", "repository"},
	"application":    {"service"},
	"infrastructure": {"repository"},
	"interfaces":     {"http", "grpc"},
}
var dddSubfolders = []string{"handler", "repository", "model", "service"}
var defaultModules = []string{"user", "order"}

func createProjectStructure(base, structure string, customFolders []string) {
	var folders []string
	switch structure {
	case "standard_layout":
		folders = []string{
			filepath.Join(base, "cmd"),
			filepath.Join(base, "pkg"),
			filepath.Join(base, "internal"),
		}
		for _, sub := range layeredSubfolders {
			folders = append(folders, filepath.Join(base, "internal", sub))
		}
	case "domain_driven":
		domains := []string{"user", "order"}
		for _, domain := range domains {
			for _, sub := range dddSubfolders {
				folders = append(folders, filepath.Join(base, domain, sub))
			}
		}
	case "layered":
		for _, sub := range layeredSubfolders {
			folders = append(folders, filepath.Join(base, sub))
		}
	case "clean_architecture":
		for mainFolder, subs := range cleanArchSubfolders {
			for _, sub := range subs {
				folders = append(folders, filepath.Join(base, mainFolder, sub))
			}
		}
	case "modular":
		for _, mod := range defaultModules {
			for _, sub := range layeredSubfolders {
				folders = append(folders, filepath.Join(base, "modules", mod, sub))
			}
		}
	case "custom":
		for _, f := range customFolders {
			folders = append(folders, filepath.Join(base, f))
		}
	}
	for _, f := range folders {
		_ = os.MkdirAll(f, 0755)
	}
}

func createMicroserviceStructure(base, structure string, customFolders []string) {
	createProjectStructure(base, structure, customFolders)
}
