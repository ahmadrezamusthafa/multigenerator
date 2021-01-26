package structgen

import (
	"bytes"
	"github.com/ahmadrezamusthafa/multigenerator/shared/consts"
	"github.com/ahmadrezamusthafa/multigenerator/shared/enums/valuetype"
	"github.com/ahmadrezamusthafa/multigenerator/shared/types"
	"time"
)

type StructGen struct {
}

var (
	operatorMap = map[string]interface{}{
		consts.OperatorEqual:            nil,
		consts.OperatorLessThan:         nil,
		consts.OperatorGreaterThan:      nil,
		consts.OperatorLessThanEqual:    nil,
		consts.OperatorGreaterThanEqual: nil,
	}

	logicalOperatorMap = map[string]string{
		consts.LogicalOperatorAndSyntax: consts.LogicalOperatorAnd,
		consts.LogicalOperatorOrSyntax:  consts.LogicalOperatorOr,
	}
)

func (s *StructGen) GenerateCondition(query string) (types.Condition, error) {
	tokenAttributes := getTokenAttributes(query)
	if len(tokenAttributes) == 0 {
		return types.Condition{Attribute: &types.Attribute{}}, nil
	}
	_, condition := buildCondition(types.Condition{}, tokenAttributes)
	return condition, nil
}

func buildCondition(condition types.Condition, attrs []*types.TokenAttribute) (int, types.Condition) {
	var (
		conditionItem *types.Condition
		lastPos       int
		operator      string
	)
	for i := 0; i < len(attrs); i++ {
		lastPos = i
		attr := attrs[i]
		if attr.HasCalled {
			continue
		}
		attr.HasCalled = true
		if attr.Value == ")" {
			break
		}
		if attr.Value == "(" {
			newCondition := types.Condition{
				Operator: operator,
			}
			lastPos, resp := buildCondition(newCondition, attrs[i+1:])
			condition.Conditions = append(condition.Conditions, &resp)
			i = i + lastPos + 1
			continue
		}

		if val, ok := logicalOperatorMap[attr.Value]; ok {
			operator = val
			conditionItem = nil
		} else if _, ok := operatorMap[attr.Value]; ok {
			if conditionItem != nil {
				conditionItem.Attribute.Operator = attr.Value
			}
		} else {
			if conditionItem == nil {
				conditionItem = &types.Condition{
					Attribute: &types.Attribute{
						Name: attr.Value,
					},
				}
				conditionItem.Attribute = &types.Attribute{
					Name: attr.Value,
				}
			} else {
				conditionItem.Attribute.Value = attr.Value
				//todo : enable this later
				/*if !attr.IsAlphanumeric {
					conditionItem.Attribute.Type = getValueType(attr.Value)
				}*/
				if condition.Conditions == nil {
					condition.Conditions = []*types.Condition{}
				}
				conditionItem.Operator = operator
				condition.Conditions = append(condition.Conditions, conditionItem)
			}
		}
	}
	return lastPos, condition
}

func getTokenAttributes(query string) []*types.TokenAttribute {
	var tokenAttributes []*types.TokenAttribute
	buffer := &bytes.Buffer{}
	isOpenQuote := false
	isAlphanumeric := false
	for _, char := range query {
		switch char {
		case ' ', '\n', '\'':
			if !isOpenQuote {
				continue
			} else {
				buffer.WriteRune(char)
			}
		case '|', '&', '<', '>':
			if buffer.Len() > 0 {
				bufBytes := buffer.Bytes()
				switch bufBytes[0] {
				case consts.ByteVerticalBar:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, consts.LogicalOperatorOrSyntax, isAlphanumeric)
					isAlphanumeric = false
				case consts.ByteAmpersand:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, consts.LogicalOperatorAndSyntax, isAlphanumeric)
					isAlphanumeric = false
				default:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes), isAlphanumeric)
					isAlphanumeric = false
					buffer.WriteRune(char)
				}
			} else {
				buffer.WriteRune(char)
			}
		case '=', '(', ')':
			if buffer.Len() > 0 {
				bufBytes := buffer.Bytes()
				switch bufBytes[0] {
				case consts.ByteLessThan, consts.ByteGreaterThan:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes)+string(char), isAlphanumeric)
					isAlphanumeric = false
					continue
				default:
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufBytes), isAlphanumeric)
					isAlphanumeric = false
				}
			}
			tokenAttributes = append(tokenAttributes, &types.TokenAttribute{
				Value: string(char),
			})
		case '"':
			isOpenQuote = !isOpenQuote
			if !isOpenQuote {
				isAlphanumeric = true
			}
		default:
			if buffer.Len() > 0 {
				bufByte := buffer.Bytes()[0]
				if bufByte == consts.ByteLessThan || bufByte == consts.ByteGreaterThan {
					tokenAttributes = appendAttribute(tokenAttributes, buffer, string(bufByte), isAlphanumeric)
					isAlphanumeric = false
				}
			}
			buffer.WriteRune(char)
		}
	}
	if buffer.Len() > 0 {
		tokenAttributes = appendAttribute(tokenAttributes, buffer, buffer.String(), isAlphanumeric)
		isAlphanumeric = false
	}
	return tokenAttributes
}

func appendAttribute(tokenAttributes []*types.TokenAttribute, buffer *bytes.Buffer, value string, isAlphanumeric bool) []*types.TokenAttribute {
	tokenAttributes = append(tokenAttributes, &types.TokenAttribute{
		Value:          value,
		IsAlphanumeric: isAlphanumeric,
	})
	buffer.Reset()
	return tokenAttributes
}

func getValueType(value string) valuetype.ValueType {
	varType, indexVal, dotCount := valuetype.Alphanumeric, 0, 0
	for _, char := range value {
		if char == ',' {
			continue
		}
		if '0' <= char && char <= '9' {
			if indexVal == 0 || (indexVal > 0 && dotCount == 1) {
				varType = valuetype.Numeric
			}
		} else if char == '.' {
			if indexVal > 0 && varType == valuetype.Numeric {
				dotCount++
				varType = valuetype.Alphanumeric
			}
			if dotCount > 1 {
				varType = valuetype.Alphanumeric
				break
			}
		} else {
			varType = valuetype.Alphanumeric
			break
		}
		indexVal++
	}
	if varType == valuetype.Alphanumeric {
		if _, err := time.Parse(consts.DateTimeFormat, value); err == nil {
			varType = valuetype.Date
		}
	}
	return varType
}
