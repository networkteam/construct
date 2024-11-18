package construct

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/urfave/cli/v2"

	"github.com/networkteam/construct/v2/internal"
)

// NewCliApp returns a new app that can be executed in a main function.
//
// Example:
//
//	func main() {
//	  app := construct.NewCliApp()
//	  err := app.Run(os.Args)
//	  if err != nil {
//	    _, _ = fmt.Fprintf(os.Stderr, "Error: %v", err)
//	  }
//	}
func NewCliApp() *cli.App {
	return &cli.App{
		Name:      "construct",
		Usage:     "Generate struct mappings and helper functions for SQL",
		ArgsUsage: "[struct type | pkg] ([target type name])",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "go-package",
				Required: true,
				EnvVars:  []string{"GOPACKAGE"},
			},
			&cli.StringFlag{
				Name:     "go-file",
				Required: true,
				EnvVars:  []string{"GOFILE"},
			},
			&cli.BoolFlag{
				Name:  "split-files",
				Value: true,
				Usage: "Split generated code into multiple files",
			},
		},
		Action: func(c *cli.Context) error {
			// fully qualified type (my/pkg.MyType) or package where the record mapping is stored in tags
			mappingType := c.Args().Get(0)
			mappingTypeName, mappingTypePackage, err := getPackageAndTypeName(mappingType)
			if err != nil {
				return err
			}

			goPackage := c.String("go-package")
			goFile := c.String("go-file")

			if mappingTypeName != "" {
				var targetTypeName string
				// Either an explicit target type name is given or it is derived from the name part of the mapping type
				if c.NArg() > 1 {
					targetTypeName = c.Args().Get(1)
				} else {
					targetTypeName = mappingTypeName
				}

				m, err := internal.BuildStructMapping(mappingTypePackage, mappingTypeName, targetTypeName)
				if err != nil {
					return fmt.Errorf("building struct mapping: %w", err)
				}

				var buf bytes.Buffer
				f := internal.StartFile(goPackage)
				err = internal.GenerateMapping(f, m, goPackage)
				if err != nil {
					return fmt.Errorf("generating code: %w", err)
				}
				outputFilename := internal.SplitOutputFilename(m, goFile)
				if err := os.WriteFile(outputFilename, buf.Bytes(), 0644); err != nil {
					return fmt.Errorf("writing output file: %w", err)
				}

				return nil
			}

			mappings, err := internal.DiscoverStructMappings(mappingTypePackage)
			if err != nil {
				return fmt.Errorf("discovering struct mappings: %w", err)
			}

			splitFiles := c.Bool("split-files")

			var f *jen.File
			if !splitFiles {
				f = internal.StartFile(goPackage)
			}

			for _, m := range mappings {
				if splitFiles {
					f = internal.StartFile(goPackage)
				}

				err = internal.GenerateMapping(f, m, goPackage)
				if err != nil {
					return fmt.Errorf("generating code: %w", err)
				}

				if splitFiles {
					outputFilename := internal.SplitOutputFilename(m, goFile)
					err = writeFile(outputFilename, f)
					if err != nil {
						return err
					}
				}
			}

			if !splitFiles {
				outputFilename := internal.CombinedOutputFilename(goFile)
				err = writeFile(outputFilename, f)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func writeFile(filename string, f *jen.File) error {
	outputFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("opening output file: %w", err)
	}
	defer outputFile.Close()
	err = f.Render(outputFile)
	if err != nil {
		return fmt.Errorf("rendering output file: %w", err)
	}
	return nil
}

func getPackageAndTypeName(mappingType string) (string, string, error) {
	i := strings.LastIndexByte(mappingType, '/')
	if i == -1 {
		return "", "", fmt.Errorf("invalid mapping type: %q, expected fully qualified type with package and name (e.g. example.com/my/pkg.MyType)", mappingType)
	}
	lastPackageAndTypeName := mappingType[i+1:]

	// Check if last package part has a ".", if so, a type name is specified
	j := strings.LastIndexByte(lastPackageAndTypeName, '.')
	if j != -1 {
		// Split mappingType by last "." to get the package and the type name
		k := strings.LastIndexByte(mappingType, '.')
		mappingTypePackage := mappingType[:k]
		mappingTypeName := mappingType[k+1:]

		return mappingTypeName, mappingTypePackage, nil
	}

	return "", mappingType, nil
}
