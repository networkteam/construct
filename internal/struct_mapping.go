package internal

import (
	"go/types"

	"github.com/fatih/structtag"
	"github.com/friendsofgo/errors"
	"golang.org/x/tools/go/packages"
)

type StructMapping struct {
	// TargetName is the exported record name (e.g. "MyRecord")
	TargetName         string
	MappingTypePackage string
	MappingTypeName    string
	FieldMappings      []FieldMapping
}

type FieldMapping struct {
	Name        string
	ReadColDef  *ReadColDef
	WriteColDef *WriteColDef
	FieldType   types.Type
}

type ReadColDef struct {
	Col      string
	Sortable bool
}

type WriteColDef struct {
	Col string
}

func BuildStructMapping(mappingTypePackage string, mappingTypeName string, targetTypeName string) (*StructMapping, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedImports}
	pkgs, err := packages.Load(cfg, mappingTypePackage)
	if err != nil {
		return nil, errors.Wrap(err, "loading package for type info")
	}

	if len(pkgs) != 1 {
		return nil, errors.Errorf("expected single package after load, got %d", len(pkgs))
	}
	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(mappingTypeName)
	if obj == nil {
		return nil, errors.Errorf("%s not found in lookup", mappingTypeName)
	}

	if _, ok := obj.(*types.TypeName); !ok {
		return nil, errors.Errorf("%v is not a named type", obj)
	}
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, errors.Errorf("type %v is a %T, not a struct", obj, obj.Type().Underlying())
	}

	var m StructMapping

	m.TargetName = targetTypeName
	m.MappingTypePackage = mappingTypePackage
	m.MappingTypeName = mappingTypeName

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)

		fieldName := field.Name()
		fieldTagValue := structType.Tag(i)
		tags, err := structtag.Parse(fieldTagValue)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing tags of field %s", fieldName)
		}

		fm := FieldMapping{
			Name:      fieldName,
			FieldType: field.Type(),
		}

		for _, tag := range tags.Tags() {
			if tag.Key == "read_col" {
				fm.ReadColDef = &ReadColDef{
					Col:      tag.Name,
					Sortable: tag.HasOption("sortable"),
				}
			}
			if tag.Key == "write_col" {
				fm.WriteColDef = &WriteColDef{
					Col: tag.Name,
				}
			}
		}

		m.FieldMappings = append(m.FieldMappings, fm)
	}

	return &m, nil
}
