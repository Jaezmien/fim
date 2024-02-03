﻿using fim.spike;
using fim.spike.Nodes;

namespace fim.celestia
{
    public class Interpreter
    {
        private Report ReportNode;
        internal Interpreter(Report reportNode)
        {
            this.ReportNode = reportNode;
            this.Variables = new VariableManager();
            this.Paragraphs = new Stack<Paragraph>();
            this.ReportName = reportNode.Name;
            this.ReportAuthor = reportNode.Author;

            foreach (Node n in ReportNode.Body)
            {
                if (Utilities.IsSameClass(n.GetType(), typeof(VariableDeclarationNode)))
                {
                    var node = (VariableDeclarationNode)n;
                    var value = EvaluateValueNode(node.Value, out VarType evaluatedType, false);

                    if( node.Type != evaluatedType ) throw new Exception("Expected type " + node.Type + ", got " + evaluatedType);
                    if( Variables.Get(node.Identifier) != null ) throw new Exception("Variable " + node.Identifier + " already exists.");

                    Variables.Push(new Variable(node.Identifier, value, node.Type, node.isConstant), true);
                }
                else if (Utilities.IsSameClass(n.GetType(), typeof(FunctionNode)))
                {
                    var node = (FunctionNode)n;
                    var paragraph = new Paragraph(this, node);

                    if( Paragraphs.FirstOrDefault(p => p.Name == paragraph.Name) != null ) throw new Exception("Paragraph " + paragraph.Name + " already exists.");

                    Paragraphs.Push(paragraph);
                }
                else
                {
                    throw new NotImplementedException("Execution of node " + n.GetType().Name + " is not implemented in report body.");
                }
            }
        }

        public readonly string ReportName;
        public readonly string ReportAuthor;
        public readonly VariableManager Variables;
        public readonly Stack<Paragraph> Paragraphs;
        public Paragraph? MainParagraph
        {
            get
            {
                return Paragraphs.FirstOrDefault(p => p.Main);
            }
        }

        internal object? EvaluateValueNode(ValueNode? node, out VarType resultType, bool local = false)
        {
            resultType = VarType.UNKNOWN;
            if (node == null) { return null; }

            if (Utilities.IsSameClass(node.GetType(), typeof(LiteralNode)))
            {
                var lNode = (LiteralNode)node;
                resultType = lNode.Type;
                return lNode.Value;
            }
            else if (Utilities.IsSameClass(node.GetType(), typeof(IdentifierNode)))
            {
                var iNode = (IdentifierNode)node;
                var variable = Variables.Get(iNode.Identifier, local);
                if (variable != null)
                {
                    resultType = variable.Type;
                    return variable.Value;
                }
            }

            return null;
        }

        internal void EvalauateStatementsNode(StatementsNode node)
        {
            foreach (var statement in node.Statements)
            {
                if (Utilities.IsSameClass(statement.GetType(), typeof(PrintNode)))
                {
                    var pNode = (PrintNode)statement;
                    var value = EvaluateValueNode(pNode.Value, out _, true);
                    Console.Write(value);
                    if (pNode.NewLine) { Console.Write("\n"); }
                }
                if( Utilities.IsSameClass(statement.GetType(), typeof(VariableDeclarationNode)))
                {
                    var vdNode = (VariableDeclarationNode)statement;
                    Variable var = new Variable(vdNode.Identifier, EvaluateValueNode(vdNode.Value, out _, true), vdNode.Type, vdNode.isConstant);
                    Variables.Push(var, false);
                }
                if( Utilities.IsSameClass(statement.GetType(), typeof(VariableModifyNode)))
                {
                    var vmNode = (VariableModifyNode)statement;
                    Variable? var = Variables.Get(vmNode.Identifier, true);

                    if( var == null ) throw new Exception("Variable " + vmNode.Identifier + " not found.");
                    if( var.IsConstant ) throw new Exception("Tried to modify variable " + vmNode.Identifier + ", which is a constant.");

                    var value = EvaluateValueNode(vmNode.Value, out VarType valueType, true);

                    if( var.Type != valueType ) throw new Exception("Expected type " + var.Type + ", got " + valueType);

                    var.Value = value;
                }
            }
        }
    }
}
