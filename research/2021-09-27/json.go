package main

import (
	"fmt"

	"vimagination.zapto.org/javascript"
	"vimagination.zapto.org/javascript/walk"
	"vimagination.zapto.org/parser"
)

func jsonArray(al *javascript.ArrayLiteral) bool {
	if al.SpreadElement != nil {
		return false
	}
	for _, ae := range al.ElementList {
		if !jsonAE(&ae) {
			return false
		}
	}
	return true
}

func jsonObject(ol *javascript.ObjectLiteral) bool {
	for _, pd := range ol.PropertyDefinitionList {
		if pd.MethodDefinition != nil || pd.PropertyName == nil || pd.PropertyName.LiteralPropertyName == nil || pd.AssignmentExpression == nil {
			return false
		}
		if !jsonAE(pd.AssignmentExpression) {
			return false
		}
	}
	return true
}

func jsonAE(ae *javascript.AssignmentExpression) bool {
	if ae.AssignmentOperator != javascript.AssignmentNone || ae.ConditionalExpression == nil {
		return false
	}
	uc := javascript.UnwrapConditional(ae.ConditionalExpression)
	switch uc := uc.(type) {
	case *javascript.PrimaryExpression:
		if uc.Literal == nil {
			return false
		}
	case *javascript.ArrayLiteral:
		if !jsonArray(uc) {
			return false
		}
	case *javascript.ObjectLiteral:
		if !jsonObject(uc) {
			return false
		}
	default:
		return false
	}
	return true
}

type JSON []string

func (j *JSON) Handle(t javascript.Type) error {
	switch t := t.(type) {
	case *javascript.PrimaryExpression:
		if t.Literal != nil {
			*j = append(*j, t.Literal.Data)
		}
	case *javascript.ArrayLiteral:
		if jsonArray(t) {
			*j = append(*j, fmt.Sprintf("%s", t))
			return nil
		}
	case *javascript.ObjectLiteral:
		if jsonObject(t) {
			*j = append(*j, fmt.Sprintf("%s", t))
			return nil
		}
	}
	walk.Walk(t, j)
	return nil
}

func main() {
	m, _ := javascript.ParseModule(parser.NewStringTokeniser(`
var x = 'lorem ipsum';
var y = 999;
var z = [1, 2, "", null];
var c = {a: 1, b: 2};
var badA = [() => {}];
var badO = {a: () => {}, b: 1};
var badU = [undefined];
	`))
	var j JSON
	walk.Walk(m, &j)
	fmt.Println(j)
}
