// Package ast provides access to the native TypeScript AST, parser, and printer.
package ast

import (
	"fmt"

	compilerast "github.com/microsoft/TypeScript/tsc/internal/ast"
	"github.com/microsoft/TypeScript/tsc/internal/core"
	"github.com/microsoft/TypeScript/tsc/internal/parser"
	"github.com/microsoft/TypeScript/tsc/internal/printer"
	"github.com/microsoft/TypeScript/tsc/internal/scanner"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
)

type (
	Node         = compilerast.Node
	NodeFactory  = compilerast.NodeFactory
	ModifierList = compilerast.ModifierList
)

const (
	KindArrayBindingPattern          = compilerast.KindArrayBindingPattern
	KindBooleanKeyword               = compilerast.KindBooleanKeyword
	KindCaseClause                   = compilerast.KindCaseClause
	KindDefaultClause                = compilerast.KindDefaultClause
	KindExclamationEqualsEqualsToken = compilerast.KindExclamationEqualsEqualsToken
	KindExportKeyword                = compilerast.KindExportKeyword
	KindFalseKeyword                 = compilerast.KindFalseKeyword
	KindIdentifier                   = compilerast.KindIdentifier
	KindNeverKeyword                 = compilerast.KindNeverKeyword
	KindQuestionDotToken             = compilerast.KindQuestionDotToken
	KindTrueKeyword                  = compilerast.KindTrueKeyword
	KindTypeKeyword                  = compilerast.KindTypeKeyword
	KindUnknown                      = compilerast.KindUnknown
	KindUnknownKeyword               = compilerast.KindUnknownKeyword
	KindVoidKeyword                  = compilerast.KindVoidKeyword
	NodeFlagsConst                   = compilerast.NodeFlagsConst
	NodeFlagsOptionalChain           = compilerast.NodeFlagsOptionalChain
)

func IsValidIdentifier(name string) bool {
	return scanner.IsValidIdentifier(name)
}

func IsBindingIdentifier(name string, typeName bool) bool {
	kind := scanner.GetIdentifierToken(name)
	return IsValidIdentifier(name) &&
		(kind < compilerast.KindFirstReservedWord || kind > compilerast.KindLastFutureReservedWord) &&
		(!typeName || !isTypeKeyword(kind))
}

func isTypeKeyword(kind compilerast.Kind) bool {
	switch kind {
	case compilerast.KindAnyKeyword, compilerast.KindBigIntKeyword, compilerast.KindBooleanKeyword,
		compilerast.KindIntrinsicKeyword, compilerast.KindNeverKeyword, compilerast.KindNumberKeyword,
		compilerast.KindObjectKeyword, compilerast.KindStringKeyword, compilerast.KindSymbolKeyword,
		compilerast.KindUnknownKeyword:
		return true
	}
	return false
}

// ParseTypeExpression parses a single type and clones it for use in a synthetic AST.
func ParseTypeExpression(factory *NodeFactory, expression string) (*Node, error) {
	input := "type TypeExpression = (\n" + expression + "\n);"
	file := parser.ParseSourceFile(fileOptions("/type-expression.ts"), input, core.ScriptKindTS)
	if diagnostics := file.Diagnostics(); len(diagnostics) > 0 {
		return nil, fmt.Errorf("invalid TypeScript type %q: %s", expression, diagnostics[0].MessageText())
	}
	statements := file.Statements.Nodes
	if len(statements) != 1 || statements[0].Kind != compilerast.KindTypeAliasDeclaration {
		return nil, fmt.Errorf("expected one TypeScript type expression in %q", expression)
	}
	parsed := statements[0].AsTypeAliasDeclaration().Type
	if parsed.Kind != compilerast.KindParenthesizedType {
		return nil, fmt.Errorf("expected a TypeScript type expression in %q", expression)
	}
	return factory.DeepCloneNode(parsed.AsParenthesizedTypeNode().Type), nil
}

// PrintSections emits each group of statements as a separate synthetic source file.
func PrintSections(factory *NodeFactory, sections [][]*Node) []string {
	p := printer.NewPrinter(printer.PrinterOptions{NewLine: core.NewLineKindLF}, printer.PrintHandlers{}, nil)
	output := make([]string, 0, len(sections))
	for _, section := range sections {
		file := factory.NewSourceFile(fileOptions("/generated.ts"), "", factory.NewNodeList(section),
			factory.NewToken(compilerast.KindEndOfFile))
		output = append(output, p.EmitSourceFile(file.AsSourceFile()))
	}
	return output
}

func fileOptions(name string) compilerast.SourceFileParseOptions {
	return compilerast.SourceFileParseOptions{
		FileName: tspath.RootedFilePathFromAbsolute(name),
		PathKey:  tspath.PathKeyFromCanonical(name),
	}
}
