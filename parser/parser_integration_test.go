package parser_test

// import (
// 	"io/ioutil"
// 	"monkey/lexer"
// 	"monkey/parser"
// 	"testing"
// )

// func TestParserIntegration(t *testing.T) {
// 	file := "./tests/let.monkey"
// 	content, err := ioutil.ReadFile(file)

// 	if err != nil {
// 		t.Fatalf("Can not read test file: %q", err)
// 	}

// 	source := string(content)

// 	l := lexer.New(source)
// 	p := parser.New(l)

// 	program := p.ParseProgram()
// }
