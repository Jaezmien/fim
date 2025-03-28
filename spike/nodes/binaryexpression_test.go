package nodes

import (
	"testing"

	"git.jaezmien.com/Jaezmien/fim/twilight"
	"github.com/stretchr/testify/assert"

	. "git.jaezmien.com/Jaezmien/fim/spike/node"
)

func TestBinaryExpression(t *testing.T) {
	t.Run("should create an add expression", func(t *testing.T) {
		source := "1 plus 1"

		tokens := twilight.Parse(source)
		tokens = tokens[:len(tokens)-1] // Ignore EOF token

		valueNode, err := CreateValueNode(tokens, CreateValueNodeOptions{})
		if !assert.NoError(t, err) {
			return
		}
		binaryExpression := valueNode.(*BinaryExpressionNode)

		assert.Equal(t, TYPE_BINARYEXPRESSION, binaryExpression.Type(), "Expected BinaryExpressionNode")
		assert.Equal(t, BINARYTYPE_ARITHMETIC, binaryExpression.BinaryType, "Expected operator to be Arithmetic")
		assert.Equal(t, BINARYOPERATOR_ADD, binaryExpression.Operator, "Expected operator to be Add")
		assert.Equal(t, TYPE_LITERAL, binaryExpression.Left.Type(), "Expected left node to be LiteralNode")
		assert.Equal(t, TYPE_LITERAL, binaryExpression.Right.Type(), "Expected right node to be LiteralNode")
	})
	t.Run("should handle multiple expressions", func(t *testing.T) {
		source := "1 plus 2 minus 3"

		tokens := twilight.Parse(source)
		tokens = tokens[:len(tokens)-1] // Ignore EOF token

		valueNode, err := CreateValueNode(tokens, CreateValueNodeOptions{})
		if !assert.NoError(t, err) {
			return
		}
		binaryExpression := valueNode.(*BinaryExpressionNode)

		assert.Equal(t, TYPE_BINARYEXPRESSION, binaryExpression.Type(), "Expected BinaryExpressionNode")
		assert.Equal(t, BINARYTYPE_ARITHMETIC, binaryExpression.BinaryType, "Expected operator to be Arithmetic")
		assert.Equal(t, BINARYOPERATOR_ADD, binaryExpression.Operator, "Expected operator to be Add")
		assert.Equal(t, TYPE_LITERAL, binaryExpression.Left.Type(), "Expected left node to be LiteralNode")
		assert.Equal(t, TYPE_BINARYEXPRESSION, binaryExpression.Right.Type(), "Expected right node to be BinaryExpresionNode")
		assert.Equal(t, 1.0, binaryExpression.Left.(*LiteralNode).GetValueNumber(), "Expected left node value to be 1")

		rightNode := binaryExpression.Right.(*BinaryExpressionNode)

		assert.Equal(t, TYPE_BINARYEXPRESSION, rightNode.Type(), "Expected BinaryExpressionNode")
		assert.Equal(t, BINARYTYPE_ARITHMETIC, rightNode.BinaryType, "Expected operator to be Arithmetic")
		assert.Equal(t, BINARYOPERATOR_SUB, rightNode.Operator, "Expected operator to be Subtract")
		assert.Equal(t, TYPE_LITERAL, rightNode.Left.Type(), "Expected left node to be LiteralNode")
		assert.Equal(t, TYPE_LITERAL, rightNode.Right.Type(), "Expected right node to be LiteralNode")
		assert.Equal(t, 2.0, rightNode.Left.(*LiteralNode).GetValueNumber(), "Expected left node value to be 2")
		assert.Equal(t, 3.0, rightNode.Right.(*LiteralNode).GetValueNumber(), "Expected right node value to be 3")
	})

	t.Run("should handle relationals", func(t *testing.T) {
		source := "correct is equal to true"

		tokens := twilight.Parse(source)
		tokens = tokens[:len(tokens)-1] // Ignore EOF token

		valueNode, err := CreateValueNode(tokens, CreateValueNodeOptions{})
		if !assert.NoError(t, err) {
			return
		}
		binaryExpression := valueNode.(*BinaryExpressionNode)

		assert.Equal(t, TYPE_BINARYEXPRESSION, binaryExpression.Type(), "Expected BinaryExpressionNode")
		assert.Equal(t, BINARYTYPE_RELATIONAL, binaryExpression.BinaryType, "Expected expression type to be Relational")
		assert.Equal(t, BINARYOPERATOR_EQ, binaryExpression.Operator, "Expected operator to be Equal")
		assert.Equal(t, TYPE_LITERAL, binaryExpression.Left.Type(), "Expected left node to be LiteralNode")
		assert.Equal(t, TYPE_LITERAL, binaryExpression.Right.Type(), "Expected right node to be BinaryExpresionNode")
	})
}
