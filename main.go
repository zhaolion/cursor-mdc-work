package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

//go:embed mapping.json
var embeddedFiles embed.FS

// MappingConfig - Nested structure to support categorization
type MappingConfig struct {
	Mappings map[string]map[string][]string `json:"mappings"`
}

var (
	// Command line parameters
	configFile   string
	targetDir    string
	ruleTypes    []string
	ruleCategory string
)

func main() {
	// Create root command
	rootCmd := &cobra.Command{
		Use:   "cursor-mdc-work",
		Short: "MDC Rules Copy Tool",
		Long:  "Copy rule files from submodules (like cursor/rules-mdc) to a specified target directory.",
	}

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available rule types",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadConfig(configFile)
			if err != nil {
				fmt.Printf("Failed to load configuration file: %v\n", err)
				os.Exit(1)
			}

			if ruleCategory != "" {
				// List only types under a specific category
				if category, exists := config.Mappings[ruleCategory]; exists {
					fmt.Printf("Rule types under category '%s':\n", ruleCategory)
					
					// Create a sorted list of types
					var types []string
					for t := range category {
						types = append(types, t)
					}
					sort.Strings(types)
					
					for _, t := range types {
						fmt.Printf("- %s\n", t)
					}
				} else {
					fmt.Printf("Error: Category '%s' not found\n", ruleCategory)
					// List all available categories
					fmt.Println("\nAvailable categories:")
					
					// Create a sorted list of categories
					var categories []string
					for c := range config.Mappings {
						categories = append(categories, c)
					}
					sort.Strings(categories)
					
					for _, c := range categories {
						fmt.Printf("- %s\n", c)
					}
					os.Exit(1)
				}
			} else {
				// List all categories and their types
				fmt.Println("All rule categories and types:")
				
				// Create a sorted list of categories
				var categories []string
				for category := range config.Mappings {
					categories = append(categories, category)
				}
				sort.Strings(categories)
				
				for _, category := range categories {
					fmt.Printf("\n%s:\n", category)
					
					// Create a sorted list of types within each category
					var types []string
					for t := range config.Mappings[category] {
						types = append(types, t)
					}
					sort.Strings(types)
					
					for _, t := range types {
						fmt.Printf("  - %s\n", t)
					}
				}
			}
		},
	}

	// Copy command
	copyCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy rule files to target directory",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if targetDir == "" {
				return fmt.Errorf("error: Please specify a target directory (--target)")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := loadConfig(configFile)
			if err != nil {
				fmt.Printf("Failed to load configuration file: %v\n", err)
				os.Exit(1)
			}

			// If a category is specified
			if ruleCategory != "" {
				categoryMap, exists := config.Mappings[ruleCategory]
				if !exists {
					fmt.Printf("Error: Category '%s' not found\n", ruleCategory)
					fmt.Println("\nAvailable categories:")
					
					// Create a sorted list of categories
					var categories []string
					for c := range config.Mappings {
						categories = append(categories, c)
					}
					sort.Strings(categories)
					
					for _, c := range categories {
						fmt.Printf("- %s\n", c)
					}
					os.Exit(1)
				}

				// If specific types are specified
				if len(ruleTypes) > 0 {
					typesFound := false
					for _, ruleType := range ruleTypes {
						filePaths, exists := categoryMap[ruleType]
						if !exists {
							fmt.Printf("Warning: Type '%s' not found in category '%s'\n", ruleType, ruleCategory)
							continue
						}

						typesFound = true
						// Copy files of the specified type
						for _, path := range filePaths {
							copyFile(path, targetDir)
						}
						fmt.Printf("Successfully copied rules for %s / %s to %s\n", ruleCategory, ruleType, targetDir)
					}

					if !typesFound {
						fmt.Printf("Error: No specified types found in category '%s'\n", ruleCategory)
						os.Exit(1)
					}
				} else {
					// Copy all types in the category
					for ruleType, filePaths := range categoryMap {
						for _, path := range filePaths {
							copyFile(path, targetDir)
						}
						fmt.Printf("Successfully copied rules for %s / %s to %s\n", ruleCategory, ruleType, targetDir)
					}
				}
			} else if len(ruleTypes) > 0 {
				// Search for specified types in all categories
				typesFound := false
				for category, categoryMap := range config.Mappings {
					for _, requestedType := range ruleTypes {
						if filePaths, exists := categoryMap[requestedType]; exists {
							typesFound = true
							// Copy files of the specified type
							for _, path := range filePaths {
								copyFile(path, targetDir)
							}
							fmt.Printf("Successfully copied rules for %s / %s to %s\n", category, requestedType, targetDir)
						}
					}
				}

				if !typesFound {
					fmt.Println("Error: No specified types found in any category")
					os.Exit(1)
				}
			} else {
				// Copy all types in all categories
				for category, categoryMap := range config.Mappings {
					for ruleType, filePaths := range categoryMap {
						for _, path := range filePaths {
							copyFile(path, targetDir)
						}
						fmt.Printf("Successfully copied rules for %s / %s to %s\n", category, ruleType, targetDir)
					}
				}
				fmt.Printf("Successfully copied all rule files to %s\n", targetDir)
			}
		},
	}

	// Set global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to mapping configuration file (if empty, uses embedded default)")

	// Set list command flags
	listCmd.Flags().StringVar(&ruleCategory, "category", "", "Specify rule category to list")

	// Set copy command flags
	copyCmd.Flags().StringVar(&targetDir, "target", "", "Target directory path")
	copyCmd.Flags().StringSliceVar(&ruleTypes, "types", []string{}, "Rule types to copy, can specify multiple (e.g.: --types=python,javascript)")
	copyCmd.Flags().StringVar(&ruleCategory, "category", "", "Specify rule category to copy")

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(copyCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Load configuration file
func loadConfig(filePath string) (*MappingConfig, error) {
	var data []byte
	var err error

	// If no config file is specified, use the embedded one
	if filePath == "" {
		fmt.Println("Using embedded default configuration")
		data, err = embeddedFiles.ReadFile("mapping.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded configuration: %v", err)
		}
	} else {
		// Check if the specified configuration file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// If the user specified a file but it doesn't exist, create an example
			fmt.Printf("Config file %s not found. Creating example configuration...\n", filePath)
			
			// Create example configuration file
			exampleConfig := MappingConfig{
				Mappings: map[string]map[string][]string{
					"frontend-frameworks": {
						"react": {
							"cursor/rules-mdc/react.mdc",
							"cursor/rules-mdc/javascript.mdc",
							"cursor/rules-mdc/typescript.mdc",
							"cursor/rules-mdc/html.mdc",
							"cursor/rules-mdc/css.mdc",
						},
					},
					"backend-languages": {
						"python": {
							"cursor/rules-mdc/python.mdc",
							"cursor/rules-mdc/fastapi.mdc",
						},
					},
				},
			}
			
			configData, err := json.MarshalIndent(exampleConfig, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to generate example configuration: %v", err)
			}
			
			err = os.WriteFile(filePath, configData, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to write example configuration file: %v", err)
			}
			
			fmt.Printf("Created example configuration file: %s\n", filePath)
			fmt.Println("Please edit this file to suit your needs, then run the program again.")
			os.Exit(0)
		}
		
		// Read the specified configuration file
		data, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read configuration file: %v", err)
		}
	}
	
	// Parse configuration file
	var config MappingConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %v", err)
	}
	
	return &config, nil
}

// Copy file to target directory
func copyFile(sourcePath, targetDir string) {
	// Get source file name
	fileName := filepath.Base(sourcePath)
	
	// Build target file path
	targetPath := filepath.Join(targetDir, fileName)
	
	// Create target directory (if it doesn't exist)
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create target directory: %v\n", err)
		return
	}
	
	// Open source file
	source, err := os.Open(sourcePath)
	if err != nil {
		fmt.Printf("Failed to open source file %s: %v\n", sourcePath, err)
		return
	}
	defer source.Close()
	
	// Create target file
	target, err := os.Create(targetPath)
	if err != nil {
		fmt.Printf("Failed to create target file %s: %v\n", targetPath, err)
		return
	}
	defer target.Close()
	
	// Copy file content
	_, err = io.Copy(target, source)
	if err != nil {
		fmt.Printf("Failed to copy file %s to %s: %v\n", sourcePath, targetPath, err)
		return
	}
	
	fmt.Printf("Copied: %s -> %s\n", sourcePath, targetPath)
} 