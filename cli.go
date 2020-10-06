package construct

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/urfave/cli/v2"

	"github.com/networkteam/construct/internal"
)

// NewCliApp returns a new app that can be executed in a main function.
//
// Example:
//   app := construct.NewCliApp()
//   err := app.Run(os.Args)
//   if err != nil {
//     _, _ = fmt.Fprintf(os.Stderr, "Error: %v", err)
//   }
func NewCliApp() *cli.App {
	return &cli.App{
		Name:      "construct",
		Usage:     "Generate struct mappings and helper functions for SQL",
		ArgsUsage: "[struct type] ([target type name])",
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
		},
		Action: func(c *cli.Context) error {
			// fully qualified type (my/pkg.MyType) where the record mapping is stored in tags
			mappingType := c.Args().Get(0)

			i := strings.LastIndexByte(mappingType, '.')
			if i == -1 {
				return errors.Errorf("invalid mapping type: %q, expected fully qualified type with package and name (e.g. example.com/my/pkg.MyType)", mappingType)
			}
			mappingTypePackage := mappingType[0:i]
			mappingTypeName := mappingType[i+1:]

			var targetTypeName string
			// Either an explicit target type name is given or it is derived from the name part of the mapping type
			if c.NArg() > 1 {
				targetTypeName = c.Args().Get(1)
			} else {
				targetTypeName = mappingTypeName
			}

			goPackage := c.String("go-package")
			goFile := c.String("go-file")

			m, err := internal.BuildStructMapping(mappingTypePackage, mappingTypeName, targetTypeName)
			if err != nil {
				return errors.Wrap(err, "building struct mapping")
			}

			var buf bytes.Buffer
			outputFilename, err := internal.Generate(m, goPackage, goFile, &buf)
			if err != nil {
				return errors.Wrap(err, "generating code")
			}
			if err := ioutil.WriteFile(outputFilename, buf.Bytes(), 0644); err != nil {
				return errors.Wrap(err, "writing output file")
			}

			return nil
		},
	}
}
