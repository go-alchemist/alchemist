package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/orochaa/go-clack/prompts"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitConfig() {
	Config = viper.New()
	Config.SetConfigName("alchemist")
	Config.SetConfigType("yaml")
	Config.AddConfigPath(".")
	Config.AddConfigPath("$HOME")
	err := Config.ReadInConfig()
	if err != nil {
		red := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		prompts.Error(red.Render("Error performing operation"))
		fmt.Println(red.Render("\n  Please check if 'alchemist.yaml' exists and is valid."))
		os.Exit(1)
	}
}
