package codegen

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/ast"
)

type ArgSet struct {
	Args     []*FieldArgument
	FuncDecl string
}

type FieldArgument struct {
	*ast.ArgumentDefinition
	TypeReference *config.TypeReference
	VarName       string      // The name of the var in go
	Object        *Object     // A link back to the parent object
	Default       interface{} // The default value
	Directives    []*Directive
	Value         interface{} // value set in Data
}

func (f *FieldArgument) Stream() bool {
	return f.Object != nil && f.Object.Stream
}

func (b *builder) buildArg(obj *Object, arg *ast.ArgumentDefinition) (*FieldArgument, error) {
	def := b.Schema.Types[arg.Type.Name()]
	if !def.IsInputType() {
		return nil, errors.Errorf(
			"cannot use %s as argument %s because %s is not a valid input type",
			arg.Type.String(),
			arg.Name,
			def.Kind,
		)
	}

	tr, err := b.Binder.TypeReference(arg.Type)
	if err != nil {
		return nil, err
	}

	argDirs, err := b.getDirectives(arg.Directives)
	if err != nil {
		return nil, err
	}
	newArg := FieldArgument{
		ArgumentDefinition: arg,
		TypeReference:      tr,
		Object:             obj,
		VarName:            templates.ToGoPrivate(arg.Name),
		Directives:         argDirs,
	}

	if arg.DefaultValue != nil {
		newArg.Default, err = arg.DefaultValue.Value(nil)
		if err != nil {
			return nil, errors.Errorf("default value is not valid: %s", err.Error())
		}
	}

	return &newArg, nil
}

func (b *builder) bindArgs(field *Field, params *types.Tuple) error {
	var newArgs []*FieldArgument

nextArg:
	for j := 0; j < params.Len(); j++ {
		param := params.At(j)
		for _, oldArg := range field.Args {
			if strings.EqualFold(oldArg.Name, param.Name()) {
				oldArg.TypeReference.GO = param.Type()
				newArgs = append(newArgs, oldArg)
				continue nextArg
			}
		}

		// no matching arg found, abort
		return fmt.Errorf("arg %s not found on method", param.Name())
	}

	field.Args = newArgs
	return nil
}

func (a *Data) Args() map[string][]*FieldArgument {
	ret := map[string][]*FieldArgument{}
	for _, o := range a.Objects {
		for _, f := range o.Fields {
			if len(f.Args) > 0 {
				ret[f.ArgsFunc()] = f.Args
			}
		}
	}

	for _, d := range a.Directives {
		if len(d.Args) > 0 {
			ret[d.ArgsFunc()] = d.Args
		}
	}
	return ret
}