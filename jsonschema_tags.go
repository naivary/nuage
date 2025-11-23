package nuage

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

type jsonSchemaTagOpts struct {
	// metdata
	// default cannot be used because its a reserved go keyword
	dflt       json.RawMessage
	deprecated bool
	readOnly   bool
	writeOnly  bool

	// validation
	enum []any

	multipleOf       *float64
	minimum          *float64
	maximum          *float64
	exclusiveMinimum *float64
	exclusiveMaximum *float64
	minLength        *int
	maxLength        *int
	pattern          string

	// arrays
	minItems    *int
	maxItems    *int
	uniqueItems bool
	minContains *int
	maxContains *int

	// objects
	minProperties     *int
	maxProperties     *int
	dependentRequired map[string][]string
}

func parseJSONSchemaTagOpts(field reflect.StructField) (*jsonSchemaTagOpts, error) {
	opts := jsonSchemaTagOpts{}
	dflt, found := field.Tag.Lookup("default")
	if found {
		opts.dflt = json.RawMessage(fmt.Sprintf(`"%s"`, dflt))
	}
	deprecated, found := field.Tag.Lookup("deprecated")
	if found {
		deprecated, err := strconv.ParseBool(deprecated)
		if err != nil {
			return nil, err
		}
		opts.deprecated = deprecated
	}
	readOnly, found := field.Tag.Lookup("readOnly")
	if found {
		readOnly, err := strconv.ParseBool(readOnly)
		if err != nil {
			return nil, err
		}
		opts.readOnly = readOnly
	}
	writeOnly, found := field.Tag.Lookup("writeOnly")
	if found {
		writeOnly, err := strconv.ParseBool(writeOnly)
		if err != nil {
			return nil, err
		}
		opts.writeOnly = writeOnly
	}
	enum, found := field.Tag.Lookup("enum")
	if found {
		for el := range strings.SplitSeq(enum, ",") {
			opts.enum = append(opts.enum, any(el))
		}
	}
	multipleOf, found := field.Tag.Lookup("multipleOf")
	if found {
		f, err := strconv.ParseFloat(multipleOf, 64)
		if err != nil {
			return nil, err
		}
		opts.multipleOf = &f
	}
	minimum, found := field.Tag.Lookup("minimum")
	if found {
		f, err := strconv.ParseFloat(minimum, 64)
		if err != nil {
			return nil, err
		}
		opts.minimum = &f
	}
	maximum, found := field.Tag.Lookup("maximum")
	if found {
		f, err := strconv.ParseFloat(maximum, 64)
		if err != nil {
			return nil, err
		}
		opts.maximum = &f
	}
	exclusiveMinimum, found := field.Tag.Lookup("exclusiveMinimum")
	if found {
		f, err := strconv.ParseFloat(exclusiveMinimum, 64)
		if err != nil {
			return nil, err
		}
		opts.exclusiveMinimum = &f
	}
	exclusiveMaximum, found := field.Tag.Lookup("exclusiveMaximum")
	if found {
		f, err := strconv.ParseFloat(exclusiveMaximum, 64)
		if err != nil {
			return nil, err
		}
		opts.exclusiveMaximum = &f
	}
	minLength, found := field.Tag.Lookup("minLength")
	if found {
		i, err := strconv.Atoi(minLength)
		if err != nil {
			return nil, err
		}
		opts.minLength = &i
	}
	maxLength, found := field.Tag.Lookup("maxLength")
	if found {
		i, err := strconv.Atoi(maxLength)
		if err != nil {
			return nil, err
		}
		opts.maxLength = &i
	}
	pattern, found := field.Tag.Lookup("pattern")
	if found {
		opts.pattern = pattern
	}
	minItems, found := field.Tag.Lookup("minItems")
	if found {
		i, err := strconv.Atoi(minItems)
		if err != nil {
			return nil, err
		}
		opts.minItems = &i
	}
	maxItems, found := field.Tag.Lookup("maxItems")
	if found {
		i, err := strconv.Atoi(maxItems)
		if err != nil {
			return nil, err
		}
		opts.maxItems = &i
	}
	uniqueItems, found := field.Tag.Lookup("uniqueItems")
	if found {
		uniqueItems, err := strconv.ParseBool(uniqueItems)
		if err != nil {
			return nil, err
		}
		opts.uniqueItems = uniqueItems
	}
	minContains, found := field.Tag.Lookup("minContains")
	if found {
		i, err := strconv.Atoi(minContains)
		if err != nil {
			return nil, err
		}
		opts.minContains = &i
	}
	maxContains, found := field.Tag.Lookup("maxContains")
	if found {
		i, err := strconv.Atoi(maxContains)
		if err != nil {
			return nil, err
		}
		opts.maxContains = &i
	}
	minProperties, found := field.Tag.Lookup("minProperties")
	if found {
		i, err := strconv.Atoi(minProperties)
		if err != nil {
			return nil, err
		}
		opts.minProperties = &i
	}
	maxProperties, found := field.Tag.Lookup("maxProperties")
	if found {
		i, err := strconv.Atoi(maxProperties)
		if err != nil {
			return nil, err
		}
		opts.maxProperties = &i
	}
	dependentRequired, found := field.Tag.Lookup("dependentRequired")
	if found {
		jsonName := jsonNameOf(field)
		if opts.dependentRequired == nil {
			opts.dependentRequired = make(map[string][]string)
		}
		opts.dependentRequired[jsonName] = strings.Split(dependentRequired, ",")
	}
	return &opts, nil
}

func (opts *jsonSchemaTagOpts) applyToSchema(schema *jsonschema.Schema, isRoot bool) error {
	// type agnostic options
	schema.Default = opts.dflt
	schema.Deprecated = opts.deprecated
	schema.ReadOnly = opts.readOnly
	schema.WriteOnly = opts.writeOnly
	schema.Enum = opts.enum

	switch schema.Type {
	case "boolean":
		// only the type agnostic options are available
	case "integer", "number":
		schema.MultipleOf = opts.multipleOf
		schema.ExclusiveMinimum = opts.exclusiveMinimum
		schema.ExclusiveMaximum = opts.exclusiveMaximum
		if opts.minimum != nil {
			switch schema.Minimum {
			case nil:
				schema.Minimum = opts.minimum
			default:
				if *opts.minimum > *schema.Minimum {
					schema.Minimum = opts.minimum
				}
				return fmt.Errorf("jsonschema: invalid minimum value %f", *opts.minimum)
			}
		}
		if opts.maximum != nil {
			switch schema.Maximum {
			case nil:
				schema.Maximum = opts.maximum
			default:
				if *opts.maximum <= *schema.Maximum {
					schema.Maximum = opts.maximum
				}
				return fmt.Errorf("jsonschema: invalid maximum value %f", *opts.maximum)
			}
		}
	case "string":
		schema.MinLength = opts.minLength
		schema.MaxLength = opts.maxLength
		schema.Pattern = opts.pattern
	case "array":
		schema.UniqueItems = opts.uniqueItems
		schema.MinContains = opts.minContains
		schema.MaxContains = opts.maxContains
		if schema.MinItems == nil {
			schema.MinItems = opts.minItems
		}
		if schema.MaxItems == nil {
			schema.MaxItems = opts.maxItems
		}
		return opts.applyToSchema(schema.Items, false)
	case "object":
		if !isRoot {
			schema.MinProperties = opts.minProperties
			schema.MaxProperties = opts.maxProperties
		}
		if len(opts.dependentRequired) > 0 && isRoot {
			for jsonName, requiredMembers := range opts.dependentRequired {
				if slices.Contains(schema.Required, jsonName) {
					return fmt.Errorf("jsonschema: dependentRequired cannot be used with a required field %s", jsonName)
				}
				if schema.DependentRequired == nil {
					schema.DependentRequired = make(map[string][]string)
				}
				schema.DependentRequired[jsonName] = requiredMembers
			}
		}
	}
	return nil
}
