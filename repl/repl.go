package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nayyara-airlangga/basedlang/lexer"
	"github.com/nayyara-airlangga/basedlang/token"
)

const prompt string = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			fmt.Fprintf(out, "%v\n", t)
		}
	}
}
