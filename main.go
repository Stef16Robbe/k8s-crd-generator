package main

import (
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"

	. "github.com/dave/jennifer/jen"
)

func capitalize(input string) string {
	return strings.ToUpper(input[:1]) + input[1:]
}

func createStructs(f *File, parent string, props map[string]v1.JSONSchemaProps, needParent bool) {
	var fields []Code

	for k, props := range props {
		k = capitalize(k)
		switch props.Type {
		case "object":
			if needParent {
				k = fmt.Sprintf("%v%v", parent, k)
			}
			if props.Properties == nil {
				fields = append(fields, Id(k).Map(String()).Interface())
			} else {
				fields = append(fields, Id(k).Id(k))
			}
		case "string":
			fields = append(fields, Id(k).String())
		case "integer":
			fields = append(fields, Id(k).Int())
		case "array":
			// TODO:
			// almost done...? I hope? It just duplicates some structs for some reason.
			// but we're almost there I think :)

			// fmt.Printf("%+v: %+v\n", k, props.Items.Schema.Properties)
			if len(props.Items.Schema.Properties) > 1 {
				fields = append(fields, Id(fmt.Sprintf("%v%v", parent, k)).Index().Id(fmt.Sprintf("%v%v", parent, k)))
				createStructs(f, k, props.Items.Schema.Properties, true)
			} else {
				createStructs(f, k, props.Properties, false)
			}
			// os.Exit(0)
		}
		if props.Properties != nil {
			createStructs(f, k, props.Properties, false)
		} //else if props.Items.Schema.Properties != nil {

		//}
	}

	f.Type().Id(parent).Struct(fields...)
}

func main() {
	data, err := os.ReadFile("crds/subscriptions.operators.coreos.com.yaml")
	if err != nil {
		panic(err.Error())
	}

	var obj v1.CustomResourceDefinition
	err = yaml.Unmarshal(data, &obj)
	if err != nil {
		panic(err.Error())
	}

	// fmt.Printf("%+v", obj.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"].Properties["config"].Properties["selector"].Properties["matchExpressions"].Items.Schema)

	kind := obj.Spec.Names.Kind
	f := NewFile("generated")

	createStructs(f, kind, obj.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties, false)

	err = f.Save(fmt.Sprintf("generated/%v.go", kind))
	if err != nil {
		panic(err.Error())
	}
}
