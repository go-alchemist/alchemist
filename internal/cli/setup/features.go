package setup

import "github.com/orochaa/go-clack/prompts"

type Feature struct {
	Name    string
	Enabled bool
}

func (m *model) AddFeatures() {
	featuresOrder := []string{
		"microservice_architecture",
		"monolithic_architecture",
		"api_gateway",
		"graphql_support",
		"database_support",
		"authentication",
		"logging_monitoring",
		"testing_framework",
		"documentation",
		"ci_cd_integration",
		"containerization",
		"deployment_scripts",
		"configuration_management",
		"error_handling",
		"rate_limiting",
		"caching",
		"security_best_practices",
		"performance_optimization",
	}

	featuresAvailable := map[string]string{
		"microservice_architecture": "Microservice Architecture",
		"monolithic_architecture":   "Monolithic Architecture",
		"api_gateway":               "API Gateway",
		"graphql_support":           "GraphQL Support",
		"database_support":          "Database Support",
		"authentication":            "Authentication",
		"logging_monitoring":        "Logging and Monitoring",
		"testing_framework":         "Testing Framework",
		"documentation":             "Documentation",
		"ci_cd_integration":         "CI/CD Integration",
		"containerization":          "Containerization",
		"deployment_scripts":        "Deployment Scripts",
		"configuration_management":  "Configuration Management",
		"error_handling":            "Error Handling",
		"rate_limiting":             "Rate Limiting",
		"caching":                   "Caching",
		"security_best_practices":   "Security Best Practices",
		"performance_optimization":  "Performance Optimization",
	}

	options := make([]*prompts.MultiSelectOption[string], 0, len(featuresOrder))
	for _, key := range featuresOrder {
		label := featuresAvailable[key]
		options = append(options, &prompts.MultiSelectOption[string]{
			Value: key,
			Label: label,
		})
	}
	selectedFeatures, err := prompts.MultiSelect(prompts.MultiSelectParams[string]{
		Message: "Select the features you want to enable:",
		Options: options,
	})
	prompts.ExitOnError(err)

	for key := range featuresAvailable {
		enabled := contains(selectedFeatures, key)
		m.Features = append(m.Features, Feature{
			Name:    key,
			Enabled: enabled,
		})
	}
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
