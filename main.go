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

func createStructs(f *File, parent string, props map[string]v1.JSONSchemaProps) {
	var fields []Code

	for k, props := range props {
		k = capitalize(k)
		switch props.Type {
		case "object":
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
			// TODO: the array type is the type of the child
			// array type has inside it items: and that one has the type object or whatever
			// need to make this recursive somehow - we need to have props.Items.Schema.Type, but since we use jennifer,
			// we need to map these to code like we do in this switch statement
			// the problem is ofcourse that an array can have arrays in arrays, etc ...
			// should be somewhat easily fixable by providing the type in this function...?
			// i hope.
			fmt.Printf("%+v: %+v\n", k, props.Items.Schema.Properties)
			fields = append(fields, Id(k).Index().String())
			os.Exit(0)
		}
		if props.Properties != nil {
			createStructs(f, k, props.Properties)
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

	createStructs(f, kind, obj.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties)

	f.Save(fmt.Sprintf("generated/%v.go", kind))
}
