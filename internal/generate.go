package internal

import (
	"fmt"
	"go/types"
	"io"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
)

// Generate Go code for the struct mapping
func Generate(m *StructMapping, goPackage string, goFile string, w io.Writer) (outputFilename string, err error) {
	pkgName := m.MappingTypePackage[strings.LastIndex(m.MappingTypePackage, "/")+1:]

	f := NewFile(goPackage)

	f.PackageComment("Code generated by construct, DO NOT EDIT.")

	generateSchemaVar(f, m)

	generateSortFields(f, m)

	// ChangeSet struct

	changeSetName, err := generateChangeSetStruct(f, m, goPackage)
	if err != nil {
		return "", fmt.Errorf("generating ChangeSet struct: %w", err)
	}

	// toMap() method for ChangeSet

	var toMapBlock []Code
	toMapBlock = append(toMapBlock, Id("m").Op(":=").Make(Map(String()).Interface()))

	for _, fm := range m.FieldMappings {
		if fm.WriteColDef != nil {
			fieldName := fm.Name

			var prepareStmt *Statement
			if fm.WriteColDef.ToJSON {
				prepareStmt = Id("data").Op(",").Id("_").Op(":=").Qual("encoding/json", "Marshal").Call(Id("c").Dot(fieldName))
			}

			mapAssign := Id("m").Index(Lit(fm.WriteColDef.Col)).Op("=")
			if fm.WriteColDef.ToJSON {
				mapAssign.Id("data")
			} else if _, ok := fm.FieldType.(*types.Slice); ok {
				// Do not indirect slice values
				mapAssign.Id("c").Dot(fieldName)
			} else {
				mapAssign.Op("*").Id("c").Dot(fieldName)
			}
			code := If(Id("c").Dot(fieldName).Op("!=").Nil()).Block(prepareStmt, mapAssign)
			toMapBlock = append(toMapBlock, code)
		}
	}

	toMapBlock = append(toMapBlock, Return(Id("m")))

	f.Func().Params(
		Id("c").Id(changeSetName),
	).Id("toMap").Params().Map(String()).Interface().Block(
		toMapBlock...,
	).Line()

	// myRecordToChangeSet() function

	var toChangeSetBlock []Code

	for _, fm := range m.FieldMappings {
		if fm.WriteColDef != nil {
			// This is the default of taking the address of the field value
			code := Id("c").Dot(firstToUpper(fm.Name)).Op("=").
				Op("&").Id("r").Dot(firstToUpper(fm.Name))

			switch v := fm.FieldType.(type) {
			case *types.Slice:
				// Do not take address of slice values
				code = Id("c").Dot(firstToUpper(fm.Name)).Op("=").Id("r").Dot(firstToUpper(fm.Name))
			case *types.Named:
				path := v.Obj().Pkg().Path()
				name := v.Obj().Name()
				fqdnType := path + "." + name
				// Special handling to check for empty values (e.g. UUID and Time) before adding them to the changeset
				switch fqdnType {
				case "github.com/gofrs/uuid.UUID":
					code = If(Id("r").Dot(firstToUpper(fm.Name)).Op("!=").Qual(path, "Nil")).Block(
						code.Clone(),
					)
				case "time.Time":
					code = If(Op("!").Id("r").Dot(firstToUpper(fm.Name)).Dot("IsZero").Call()).Block(
						code.Clone(),
					)
				}
			}

			toChangeSetBlock = append(toChangeSetBlock, code)
		}
	}

	toChangeSetBlock = append(toChangeSetBlock, Return())

	mtp := m.MappingTypePackage
	if goPackage == pkgName {
		mtp = ""
	}
	f.Func().Id(firstToUpper(m.TargetName) + "ToChangeSet").Params(
		Id("r").Qual(mtp, m.MappingTypeName),
	).Params(Id("c").Id(changeSetName)).Block(
		toChangeSetBlock...,
	).Line()

	generateDefaultSelectJsonObject(f, m)

	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]

	outputFilename = baseFilename + "_" + strings.ToLower(m.MappingTypeName) + "_gen.go"
	return outputFilename, f.Render(w)
}

func generateSchemaVar(f *File, m *StructMapping) {
	f.Var().Id(firstToLower(m.MappingTypeName)).Op("=").StructFunc(func(g *Group) {
		if m.TableName != "" {
			g.Qual("github.com/networkteam/qrb/builder", "Identer")
		}
		for _, fm := range m.FieldMappings {
			if fm.ReadColDef == nil {
				continue
			}
			g.Id(fm.Name).Qual("github.com/networkteam/qrb/builder", "IdentExp")
		}
	}).Values(DictFunc(func(d Dict) {
		if m.TableName != "" {
			d[Id("Identer")] = Qual("github.com/networkteam/qrb", "N").Call(Lit(m.TableName))
		}
		for _, fm := range m.FieldMappings {
			if fm.ReadColDef != nil {
				d[Id(fm.Name)] = Qual("github.com/networkteam/qrb", "N").Call(Lit(fm.ReadColDef.Col))
			}
		}
	})).Line()
}

func readColVarName(m *StructMapping, fm FieldMapping) string {
	return firstToLower(m.MappingTypeName) + "." + fm.Name
}

func generateDefaultSelectJsonObject(f *File, m *StructMapping) {
	varName := firstToLower(m.TargetName + "DefaultJson")

	code := Qual("github.com/networkteam/qrb/fn", "JsonBuildObject").Call()
	for _, fm := range m.FieldMappings {
		if fm.ReadColDef != nil {
			propValueStmt := Id(readColVarName(m, fm))

			// Special handling for []byte (needs encoding to Base64 (used by encoding/json for unmarshal)
			if v, ok := fm.FieldType.(*types.Slice); ok {
				if v, ok := v.Elem().(*types.Basic); ok && v.Kind() == types.Byte {
					propValueStmt = Qual("github.com/networkteam/qrb", "Func").Call(Lit("ENCODE"), Id(readColVarName(m, fm)), Qual("github.com/networkteam/qrb", "String").Call(Lit("BASE64")))
				}
			}
			code.Op(".").Line().Id("Prop").Call(Lit(fm.Name), propValueStmt)
		}
	}

	f.Var().Id(varName).Op("=").Add(code).Line()
}

func generateChangeSetStruct(f *File, m *StructMapping, pkgDest string) (changeSetName string, err error) {
	var structFields []Code

	for _, fm := range m.FieldMappings {
		if fm.WriteColDef != nil {
			code := Id(firstToUpper(fm.Name))
			switch v := fm.FieldType.(type) {
			case *types.Basic:
				code.Op("*").Id(v.String())
			case *types.Named:
				typeName := v.Obj()
				pkgPath := typeName.Pkg().Path()
				pkgName := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
				if pkgDest == pkgName {
					pkgPath = ""
				}
				code.Op("*").Qual(
					pkgPath,
					typeName.Name(),
				)
			case *types.Pointer:
				code.Op("*")
				elemType := v.Elem()
				switch v := elemType.(type) {
				case *types.Basic:
					code.Op("*").Id(v.String())
				case *types.Named:
					typeName := v.Obj()
					pkgPath := typeName.Pkg().Path()
					pkgName := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
					if pkgDest == pkgName {
						pkgPath = ""
					}
					code.Op("*").Qual(
						pkgPath,
						typeName.Name(),
					)
				case *types.Slice:
					return "", fmt.Errorf("pointer to slice is not supported")
				default:
					return "", fmt.Errorf("pointer type not handled: %T", v)
				}
			case *types.Slice:
				code.Op("[]")
				elemType := v.Elem()
				switch v := elemType.(type) {
				case *types.Basic:
					code.Id(v.String())
				case *types.Named:
					typeName := v.Obj()
					pkgPath := typeName.Pkg().Path()
					pkgName := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
					if pkgDest == pkgName {
						pkgPath = ""
					}
					code.Qual(
						pkgPath,
						typeName.Name(),
					)
				default:
					return "", fmt.Errorf("slice type not handled: %T", v)
				}
			default:
				return "", fmt.Errorf("struct field type not handled: %T (for %s)", v, fm.FieldType)
			}
			structFields = append(structFields, code)
		}
	}

	changeSetName = firstToUpper(m.TargetName + "ChangeSet")
	f.Type().Id(changeSetName).Struct(structFields...)

	return changeSetName, nil
}

func generateSortFields(f *File, m *StructMapping) {
	f.Var().Id(firstToLower(m.TargetName+"SortFields")).Op("=").Map(String()).Qual("github.com/networkteam/qrb/builder", "IdentExp").
		Values(DictFunc(func(d Dict) {
			for _, fm := range m.FieldMappings {
				if fm.ReadColDef != nil && fm.ReadColDef.Sortable {
					d[Lit(strings.ToLower(fm.Name))] = Id(readColVarName(m, fm))
				}
			}
		}))
}
