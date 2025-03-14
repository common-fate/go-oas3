package generator

import (
	"encoding/json"
	"strings"

	"github.com/common-fate/go-oas3/configurator"
	"github.com/dave/jennifer/jen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/imdario/mergo"
)

const (
	goType            = "x-go-type"
	goTypeStringParse = "x-go-type-string-parse"
	goPointer         = "x-go-pointer"
	goRegex           = "x-go-regex"
	goStringTrimmable = "x-go-string-trimmable"
	goOmitempty       = "x-go-omitempty"
)

type Type struct {
	normalizer *Normalizer          `di.inject:"normalizer"`
	config     *configurator.Config `di.inject:"config"`
}

func (typ *Type) fillJsonTag(into *jen.Statement, schemaRef *openapi3.SchemaRef, name string) {
	var tag = strings.ToLower(name[:1]) + name[1:]
	if typ.getXGoOmitempty(schemaRef.Value) {
		tag = tag + ",omitempty"
	}
	into.Tag(map[string]string{"json": tag})
}

func (typ *Type) fillGoType(into *jen.Statement, parentTypeName string, typeName string, schemaRef *openapi3.SchemaRef, asPointer bool, needAliasing bool) {
	if asPointer || typ.getXGoPointer(schemaRef.Value) {
		into.Op("*")
	}

	if pkg, typee, ok := typ.getXGoType(schemaRef.Value); ok && schemaRef.Value.AdditionalProperties == nil {
		if needAliasing {
			into.Op("=")
		}

		if pkg == "" {
			into.Id(typee)
			return
		}
		into.Qual(pkg, typee)
		return
	}

	schema := schemaRef.Value

	if schema.AnyOf != nil || schema.OneOf != nil {
		into.Interface()
		return
	}

	if schemaRef.Ref != "" {
		into.Qual(typ.config.ComponentsPackage, typ.normalizer.extractNameFromRef(schemaRef.Ref))
		return
	}

	if len(schema.AllOf) > 0 {
		allOfSchema := schema.AllOf[0]

		for i := 1; i < len(schema.AllOf); i++ {
			err := mergo.Merge(&allOfSchema, schema.AllOf[i])
			if err != nil {
				panic(err)
			}
		}

		typ.fillGoType(into, parentTypeName, typeName, allOfSchema, false, needAliasing)

		return
	}

	if len(schema.Enum) > 0 {
		into.Qual(typ.config.ComponentsPackage, typeName)
		return
	}

	switch schema.Type {
	case "object":
		if schemaRef.Ref != "" {
			typeName := typ.normalizer.normalize(typ.normalizer.extractNameFromRef(schemaRef.Ref))
			into.Qual(typ.config.ComponentsPackage, typeName)
			return
		}

		if schema.AdditionalProperties != nil {
			keyCode := jen.Null()

			keyPkg, keyValue, ok := typ.getXGoType(schemaRef.Value)
			if ok {
				if keyPkg == "" {
					keyCode.Id(keyValue)
				} else {
					keyCode.Qual(keyPkg, keyValue)
				}
			} else {
				keyCode.String()
			}

			into.Map(keyCode)

			typ.fillGoType(into, parentTypeName, typeName, schema.AdditionalProperties, false, false)

			//TODO: ANONYMOUS MAP ENTRIES
			//if schema.AdditionalProperties.Ref != "" {
			//	typ.fillGoType(into, parentTypeName, typeName, schema.AdditionalProperties, false, needAliasing)
			//	return
			//}

			//into.Qual(typ.config.ComponentsPackage, parentTypeName+typeName+"MapEntry")

			return
		}

		if len(schema.Properties) == 0 {
			into.Interface()
			return
		}
		return
	case "array":
		into.Index()

		//TODO: ANONYMOUS SLICES
		//if schema.Items.Ref != "" {
		//	typ.fillGoType(into, parentTypeName, typeName, schema.Items, false, needAliasing)
		//	return
		//}

		//into.Qual(typ.config.ComponentsPackage, parentTypeName+typeName+"SliceElement")

		typ.fillGoType(into, parentTypeName, typeName, schema.Items, false, needAliasing)
		return
	case "integer":
		into.Int()
		return
	case "number":
		into.Float64()
		return
	case "boolean":
		into.Bool()
		return
	case "string":
		if needAliasing {
			into.Op("=")
		}

		switch schema.Format {
		case "byte":
			into.Index().Byte()
			return
		case "binary":
			into.Index().Byte()
			return
		case "email":
			into.String()
			return
		case "date":
			into.String()
			return
		case "date-time":
			into.String()
			return
		case "iso4217-currency-code":
			into.Qual("github.com/mikekonan/go-types/currency", "Code")
			return
		case "iso3166-alpha-2":
			into.Qual("github.com/mikekonan/go-types/country", "Alpha2Code")
			return
		case "iso3166-alpha-3":
			into.Qual("github.com/mikekonan/go-types/country", "Alpha3Code")
			return
		case "uuid":
			into.Qual("github.com/google/uuid", "UUID")
			return
		case "json":
			into.Qual("encoding/json", "RawMessage")
			return
		default:
			into.String()
			return
		}
	}

	into.Interface()
}

func (typ *Type) hasXGoType(schema *openapi3.Schema) bool {
	if len(schema.Extensions) > 0 && schema.Extensions[goType] != nil {
		return true
	}

	return false
}

func (typ *Type) hasXGoPointer(schema *openapi3.Schema) bool {
	if len(schema.Extensions) > 0 && schema.Extensions[goPointer] != nil {
		return true
	}

	return false
}

func (typ *Type) hasXGoTypeStringParse(schema *openapi3.Schema) bool {
	if typ.hasXGoType(schema) && schema.Extensions[goTypeStringParse] != nil {
		return true
	}

	return false
}

func (typ *Type) getXGoTypeStringParse(schema *openapi3.Schema) (string, string, bool) {
	if typ.hasXGoTypeStringParse(schema) {
		var customType string

		if err := json.Unmarshal(schema.Extensions[goTypeStringParse].(json.RawMessage), &customType); err != nil {
			panic(err)
		}

		index := strings.LastIndex(customType, ".")

		return customType[:index], customType[index+1:], true
	}

	return "", "", false
}

func (typ *Type) getXGoType(schema *openapi3.Schema) (string, string, bool) {
	if typ.hasXGoType(schema) {
		var customType string

		if err := json.Unmarshal(schema.Extensions[goType].(json.RawMessage), &customType); err != nil {
			panic(err)
		}

		index := strings.LastIndex(customType, ".")
		if index == -1 {
			return "", customType, true
		}

		return customType[:index], customType[index+1:], true
	}

	return "", "", false
}

func (typ *Type) getXGoPointer(schema *openapi3.Schema) bool {
	var value = false

	if typ.hasXGoPointer(schema) {
		if err := json.Unmarshal(schema.Extensions[goPointer].(json.RawMessage), &value); err != nil {
			panic(err)
		}
	}

	return value
}

func (typ *Type) hasXGoOmitempty(schema *openapi3.Schema) bool {
	if len(schema.Extensions) > 0 && schema.Extensions[goOmitempty] != nil {
		return true
	}

	return false
}

func (typ *Type) getXGoOmitempty(schema *openapi3.Schema) bool {
	var value = false

	if typ.hasXGoOmitempty(schema) {
		if err := json.Unmarshal(schema.Extensions[goOmitempty].(json.RawMessage), &value); err != nil {
			panic(err)
		}
	}

	return value
}

func (typ *Type) isCustomType(schema *openapi3.Schema) bool {
	return schema.Type == "string" && (schema.Format != "" || typ.hasXGoTypeStringParse(schema))
}
