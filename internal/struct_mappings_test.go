package internal_test

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/networkteam/construct/internal"
)

func myTypeStructMapping() *internal.StructMapping {
	return &internal.StructMapping{
		TargetName:         "MyTargetType",
		MappingTypePackage: "github.com/networkteam/construct/internal/fixtures",
		MappingTypeName:    "MyType",
		FieldMappings: []internal.FieldMapping{
			{
				Name: "ID",
				ReadColDef: &internal.ReadColDef{
					Col:      "my_type.id",
					Sortable: false,
				},
				WriteColDef: &internal.WriteColDef{
					Col: "id",
				},
				FieldType: types.NewNamed(types.NewTypeName(token.NoPos, types.NewPackage("github.com/gofrs/uuid", "uuid"), "UUID", nil), nil, nil),
			},
			{
				Name: "Foo",
				ReadColDef: &internal.ReadColDef{
					Col:      "my_type.foo",
					Sortable: true,
				},
				WriteColDef: &internal.WriteColDef{
					Col: "foo",
				},
				FieldType: types.Typ[types.String],
			},
			{
				Name: "Bar",
				ReadColDef: &internal.ReadColDef{
					Col:      "my_type.the_bar",
					Sortable: false,
				},
				WriteColDef: &internal.WriteColDef{
					Col: "the_bar",
				},
				FieldType: types.NewSlice(types.Universe.Lookup("byte").Type()),
			},
			{
				Name: "LastTime",
				ReadColDef: &internal.ReadColDef{
					Col:      "my_type.last_time",
					Sortable: true,
				},
				WriteColDef: &internal.WriteColDef{
					Col: "last_time",
				},
				FieldType: types.NewPointer(types.NewNamed(types.NewTypeName(token.NoPos, types.NewPackage("time", "time"), "Time", nil), nil, nil)),
			},
		},
	}
}

func TestBuildStructMapping(t *testing.T) {
	m, err := internal.BuildStructMapping("github.com/networkteam/construct/internal/fixtures", "MyType", "MyTargetType")
	if err != nil {
		t.Fatalf("error building struct mapping: %v", err)
	}

	expectedStructMapping := myTypeStructMapping()

	if m.MappingTypeName != expectedStructMapping.MappingTypeName {
		t.Errorf("expected mapping type name to be %s, but got %s", expectedStructMapping.MappingTypeName, m.MappingTypeName)
	}
	if m.MappingTypePackage != expectedStructMapping.MappingTypePackage {
		t.Errorf("expected mapping type package to be %s, but got %s", expectedStructMapping.MappingTypePackage, m.MappingTypePackage)
	}
	if m.TargetName != expectedStructMapping.TargetName {
		t.Errorf("expected target name to be %s, but got %s", expectedStructMapping.TargetName, m.TargetName)
	}
	if len(m.FieldMappings) != len(expectedStructMapping.FieldMappings) {
		t.Fatalf("expected %d field mappings, but got %d", len(expectedStructMapping.FieldMappings), len(m.FieldMappings))
	}
	for i, expectedFieldMapping := range expectedStructMapping.FieldMappings {
		actualFieldMapping := m.FieldMappings[i]

		if actualFieldMapping.Name != expectedFieldMapping.Name {
			t.Errorf("expected field mapping %d name to be %s, but got %s", i, expectedFieldMapping.Name, actualFieldMapping.Name)
		}
		if expectedFieldMapping.ReadColDef == nil {
			if actualFieldMapping.ReadColDef != nil {
				t.Errorf("expected field mapping %d read col def to be nil, but it was not", i)
			}
		} else {
			if actualFieldMapping.ReadColDef == nil {
				t.Errorf("expected field mapping %d read col def to be not nil, but it was", i)
			}
			if actualFieldMapping.ReadColDef.Col != expectedFieldMapping.ReadColDef.Col {
				t.Errorf("expected field mapping %d read col def col to be %s, but got %s", i, expectedFieldMapping.ReadColDef.Col, actualFieldMapping.ReadColDef.Col)
			}
			if actualFieldMapping.ReadColDef.Sortable != expectedFieldMapping.ReadColDef.Sortable {
				t.Errorf("expected field mapping %d read col def sortable to be %v, but got %v", i, expectedFieldMapping.ReadColDef.Sortable, actualFieldMapping.ReadColDef.Sortable)
			}
		}
		if expectedFieldMapping.WriteColDef == nil {
			if actualFieldMapping.WriteColDef != nil {
				t.Errorf("expected field mapping %d write col def to be nil, but it was not", i)
			}
		} else {
			if actualFieldMapping.WriteColDef == nil {
				t.Errorf("expected field mapping %d write col def to be not nil, but it was", i)
			}
			if actualFieldMapping.WriteColDef.Col != expectedFieldMapping.WriteColDef.Col {
				t.Errorf("expected field mapping %d write col def col to be %s, but got %s", i, expectedFieldMapping.WriteColDef.Col, actualFieldMapping.WriteColDef.Col)
			}
		}
		if actualFieldMapping.Name != expectedFieldMapping.Name {
			t.Errorf("expected field mapping %d name to be %s, but got %s", i, expectedFieldMapping.Name, actualFieldMapping.Name)
		}
		if actualFieldMapping.FieldType.String() != expectedFieldMapping.FieldType.String() {
			t.Errorf("expected field mapping %d field type to be %s, but got %s", i, expectedFieldMapping.FieldType.String(), actualFieldMapping.FieldType.String())
		}
	}
}
