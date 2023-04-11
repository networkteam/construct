package internal

import (
	"fmt"
	"go/types"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/packages"
)

// StructMapping contains mapping information for a struct type
type StructMapping struct {
	// TargetName is the exported record name (e.g. "MyRecord")
	TargetName string
	// MappingTypePackage is the package of the type with mapping information
	MappingTypePackage string
	// MappingTypePackage is the name of the type with mapping information
	MappingTypeName string
	// TableName contains the table name (if set)
	TableName string
	// FieldMappings contains all field mappings derived from the type
	FieldMappings []FieldMapping
}

// FieldMapping contains mapping information for a field of a mapping
type FieldMapping struct {
	// Name is the Go field name
	Name string
	// ReadColDef is the read column definition (nil if not given)
	ReadColDef *ReadColDef
	// WriteColDef is the write column definition (nil if not given)
	WriteColDef *WriteColDef
	// FieldType is the type in the Go struct
	FieldType types.Type
}

// ReadColDef is the read column definition
type ReadColDef struct {
	// Col is the column name
	Col string
	// Sortable is true, if the column should appear in the list of sortable columns
	Sortable bool
}

// WriteColDef is the write column definition
type WriteColDef struct {
	// Col is the column name
	Col string
	// ToJSON when writing a column value
	ToJSON bool
}

// BuildStructMapping builds a struct mapping for a given mapping type and target type
func BuildStructMapping(mappingTypePackage string, mappingTypeName string, targetTypeName string) (*StructMapping, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedImports}
	pkgs, err := packages.Load(cfg, mappingTypePackage)
	if err != nil {
		return nil, fmt.Errorf("loading package for type info: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected single package after load, got %d", len(pkgs))
	}
	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(mappingTypeName)
	if obj == nil {
		return nil, fmt.Errorf("%s not found in lookup", mappingTypeName)
	}

	if _, ok := obj.(*types.TypeName); !ok {
		return nil, fmt.Errorf("%v is not a named type", obj)
	}
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("type %v is a %T, not a struct", obj, obj.Type().Underlying())
	}

	return buildStructMapping(targetTypeName, mappingTypePackage, mappingTypeName, structType)
}

// DiscoverStructMappings discovers all struct mappings in a given package
func DiscoverStructMappings(mappingTypePackage string) (mappings []*StructMapping, err error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax | packages.NeedImports}
	pkgs, err := packages.Load(cfg, mappingTypePackage)
	if err != nil {
		return nil, fmt.Errorf("loading package for type info: %w", err)
	}

	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected single package after load, got %d", len(pkgs))
	}
	pkg := pkgs[0]

	names := pkg.Types.Scope().Names()

	for _, name := range names {
		obj := pkg.Types.Scope().Lookup(name)
		if _, ok := obj.(*types.TypeName); !ok {
			continue
		}
		structType, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}

		m, err := buildStructMapping(name, mappingTypePackage, name, structType)
		if err != nil {
			return nil, fmt.Errorf("building struct mapping for %s: %w", name, err)
		}

		// Only include mappings with read or write columns defined
		if m.hasColDef() {
			mappings = append(mappings, m)
		}
	}

	return mappings, nil
}

func buildStructMapping(targetTypeName, mappingTypePackage, mappingTypeName string, structType *types.Struct) (*StructMapping, error) {
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
			return nil, fmt.Errorf("parsing tags of field %s: %w", fieldName, err)
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
					Col:    tag.Name,
					ToJSON: tag.HasOption("json"),
				}
			}
			if tag.Key == "table_name" {
				if m.TableName != "" {
					return nil, fmt.Errorf("table_name tag defined multiple times")
				}
				m.TableName = tag.Name
			}
		}

		if fm.ReadColDef != nil || fm.WriteColDef != nil {
			m.FieldMappings = append(m.FieldMappings, fm)
		}
	}

	return &m, nil
}

func (m *StructMapping) hasColDef() bool {
	for _, fm := range m.FieldMappings {
		if fm.ReadColDef != nil || fm.WriteColDef != nil {
			return true
		}
	}
	return false
}
