package internal_test

import (
	"go/token"
	"go/types"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/construct/v2/internal"
)

func myTypeStructMapping() *internal.StructMapping {
	return &internal.StructMapping{
		TargetName:         "MyTargetType",
		MappingTypePackage: "github.com/networkteam/construct/v2/internal/fixtures",
		MappingTypeName:    "MyType",
		TableName:          "my_type",
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
				Name: "Baz",
				ReadColDef: &internal.ReadColDef{
					Col:      "my_type.baz",
					Sortable: false,
				},
				WriteColDef: &internal.WriteColDef{
					Col:    "baz",
					ToJSON: true,
				},
				FieldType: types.NewNamed(types.NewTypeName(token.NoPos, types.NewPackage("github.com/networkteam/construct/v2/internal/fixtures", "MyEmbeddedType"), "MyEmbeddedType", nil), nil, nil),
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
	m, err := internal.BuildStructMapping("github.com/networkteam/construct/v2/internal/fixtures", "MyType", "MyTargetType")
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
				t.Fatalf("expected field mapping %d read col def to be not nil, but it was", i)
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
				t.Fatalf("expected field mapping %d write col def to be not nil, but it was", i)
			}
			if actualFieldMapping.WriteColDef.Col != expectedFieldMapping.WriteColDef.Col {
				t.Errorf("expected field mapping %d write col def col to be %s, but got %s", i, expectedFieldMapping.WriteColDef.Col, actualFieldMapping.WriteColDef.Col)
			}
			if actualFieldMapping.WriteColDef.ToJSON != expectedFieldMapping.WriteColDef.ToJSON {
				t.Errorf("expected field mapping %d write col def json to be %v, but got %v", i, expectedFieldMapping.WriteColDef.ToJSON, actualFieldMapping.WriteColDef.ToJSON)
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

func TestDiscoverStructMappings(t *testing.T) {
	structMappings, err := internal.DiscoverStructMappings("github.com/networkteam/construct/v2/internal/fixtures")
	require.NoError(t, err)

	require.Len(t, structMappings, 2)

	// order by MappingTypeName
	sort.Slice(structMappings, func(i, j int) bool {
		return structMappings[i].MappingTypeName < structMappings[j].MappingTypeName
	})

	assert.Equal(t, "MyOtherType", structMappings[0].MappingTypeName)
	assert.Equal(t, "MyOtherType", structMappings[0].TargetName)
	assert.Equal(t, "MyType", structMappings[1].MappingTypeName)
	assert.Equal(t, "MyType", structMappings[1].TargetName)
}
