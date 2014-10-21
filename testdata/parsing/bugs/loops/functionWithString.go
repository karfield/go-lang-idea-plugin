package main

//
// copies and protects "'s in q
//
func chcopy(q string) string {
	s := ""
	i := 0
	j := 0
	for i = 0; i < len(q); i++ {
		if q[i] == '"' {
			s += q[j:i] + "\\"
			j = i
		}
	}
	return s + q[j:i]
}

func usage() {
	fmt.Fprintf(stderr, "usage: yacc [-o output] [-v parsetable] input\n")
	exit(1)
}

func bitset(set Lkset, bit int) int { return set[bit>>5] & (1 << uint(bit&31)) }

func setbit(set Lkset, bit int) { set[bit>>5] |= (1 << uint(bit&31)) }

func mkset() Lkset { return make([]int, tbitset) }

//
// set a to the union of a and b
// return 1 if b is not a subset of a, 0 otherwise
//
func setunion(a, b []int) int {
	sub := 0
	for i := 0; i < tbitset; i++ {
		x := a[i]
		y := x | b[i]
		a[i] = y
		if y != x {
			sub = 1
		}
	}
	return sub
}

func prlook(p Lkset) {
	if p == nil {
		fmt.Fprintf(foutput, "\tNULL")
		return
	}
	fmt.Fprintf(foutput, " { ")
	for j := 0; j <= ntokens; j++ {
		if bitset(p, j) != 0 {
			fmt.Fprintf(foutput, "%v ", symnam(j))
		}
	}
	fmt.Fprintf(foutput, "}")
}

//
// utility routines
//
var peekrune rune

func isdigit(c rune) bool { return c >= '0' && c <= '9' }

func isword(c rune) bool {
	return c >= 0xa0 || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func mktemp(t string) string { return t }

//
// return 1 if 2 arrays are equal
// return 0 if not equal
//
func aryeq(a []int, b []int) int {
	n := len(a)
	if len(b) != n {
		return 0
	}
	for ll := 0; ll < n; ll++ {
		if a[ll] != b[ll] {
			return 0
		}
	}
	return 1
}

func putrune(f *bufio.Writer, c int) {
	s := string(c)
	for i := 0; i < len(s); i++ {
		f.WriteByte(s[i])
	}
}

func getrune(f *bufio.Reader) rune {
	var r rune

	if peekrune != 0 {
		if peekrune == EOF {
			return EOF
		}
		r = peekrune
		peekrune = 0
		return r
	}

	c, n, err := f.ReadRune()
	if n == 0 {
		return EOF
	}
	if err != nil {
		errorf("read error: %v", err)
	}
	//fmt.Printf("rune = %v n=%v\n", string(c), n);
	return c
}

func ungetrune(f *bufio.Reader, c rune) {
	if f != finput {
		panic("ungetc - not finput")
	}
	if peekrune != 0 {
		panic("ungetc - 2nd unget")
	}
	peekrune = c
}

func write(f *bufio.Writer, b []byte, n int) int {
	panic("write")
	return 0
}

func open(s string) *bufio.Reader {
	fi, err := os.Open(s)
	if err != nil {
		errorf("error opening %v: %v", s, err)
	}
	//fmt.Printf("open %v\n", s);
	return bufio.NewReader(fi)
}

func create(s string) *bufio.Writer {
	fo, err := os.Create(s)
	if err != nil {
		errorf("error creating %v: %v", s, err)
	}
	//fmt.Printf("create %v mode %v\n", s);
	return bufio.NewWriter(fo)
}

//
// write out error comment
//
func errorf(s string, v ...interface{}) {
	nerrors++
	fmt.Fprintf(stderr, s, v...)
	fmt.Fprintf(stderr, ": %v:%v\n", infile, lineno)
	if fatfl != 0 {
		summary()
		exit(1)
	}
}

func exit(status int) {
	if ftable != nil {
		ftable.Flush()
		ftable = nil
	}
	if foutput != nil {
		foutput.Flush()
		foutput = nil
	}
	if stderr != nil {
		stderr.Flush()
		stderr = nil
	}
	os.Exit(status)
}

var yaccpar string // will be processed version of yaccpartext: s/$$/prefix/g
var yaccpartext = `
/*	parser for yacc output	*/

var $$Debug = 0

type $$Lexer interface {
	Lex(lval *$$SymType) int
	Error(s string)
}

const $$Flag = -1000

func $$Tokname(c int) string {
	if c > 0 && c <= len($$Toknames) {
		if $$Toknames[c-1] != "" {
			return $$Toknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func $$Statname(s int) string {
	if s >= 0 && s < len($$Statenames) {
		if $$Statenames[s] != "" {
			return $$Statenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func $$lex1(lex $$Lexer, lval *$$SymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = $$Tok1[0]
		goto out
	}
	if char < len($$Tok1) {
		c = $$Tok1[char]
		goto out
	}
	if char >= $$Private {
		if char < $$Private+len($$Tok2) {
			c = $$Tok2[char-$$Private]
			goto out
		}
	}
	for i := 0; i < len($$Tok3); i += 2 {
		c = $$Tok3[i+0]
		if c == char {
			c = $$Tok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = $$Tok2[1] /* unknown char */
	}
	if $$Debug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), $$Tokname(c))
	}
	return c
}

func $$Parse($$lex $$Lexer) int {
	var $$n int
	var $$lval $$SymType
	var $$VAL $$SymType
	$$S := make([]$$SymType, $$MaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	$$state := 0
	$$char := -1
	$$p := -1
	goto $$stack

ret0:
	return 0

ret1:
	return 1

$$stack:
	/* put a state and value onto the stack */
	if $$Debug >= 4 {
		fmt.Printf("char %v in %v\n", $$Tokname($$char), $$Statname($$state))
	}

	$$p++
	if $$p >= len($$S) {
		nyys := make([]$$SymType, len($$S)*2)
		copy(nyys, $$S)
		$$S = nyys
	}
	$$S[$$p] = $$VAL
	$$S[$$p].yys = $$state

$$newstate:
	$$n = $$Pact[$$state]
	if $$n <= $$Flag {
		goto $$default /* simple state */
	}
	if $$char < 0 {
		$$char = $$lex1($$lex, &$$lval)
	}
	$$n += $$char
	if $$n < 0 || $$n >= $$Last {
		goto $$default
	}
	$$n = $$Act[$$n]
	if $$Chk[$$n] == $$char { /* valid shift */
		$$char = -1
		$$VAL = $$lval
		$$state = $$n
		if Errflag > 0 {
			Errflag--
		}
		goto $$stack
	}

$$default:
	/* default state action */
	$$n = $$Def[$$state]
	if $$n == -2 {
		if $$char < 0 {
			$$char = $$lex1($$lex, &$$lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if $$Exca[xi+0] == -1 && $$Exca[xi+1] == $$state {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			$$n = $$Exca[xi+0]
			if $$n < 0 || $$n == $$char {
				break
			}
		}
		$$n = $$Exca[xi+1]
		if $$n < 0 {
			goto ret0
		}
	}
	if $$n == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			$$lex.Error("syntax error")
			Nerrs++
			if $$Debug >= 1 {
				fmt.Printf("%s", $$Statname($$state))
				fmt.Printf("saw %s\n", $$Tokname($$char))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for $$p >= 0 {
				$$n = $$Pact[$$S[$$p].yys] + $$ErrCode
				if $$n >= 0 && $$n < $$Last {
					$$state = $$Act[$$n] /* simulate a shift of "error" */
					if $$Chk[$$state] == $$ErrCode {
						goto $$stack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if $$Debug >= 2 {
					fmt.Printf("error recovery pops state %d\n", $$S[$$p].yys)
				}
				$$p--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if $$Debug >= 2 {
				fmt.Printf("error recovery discards %s\n", $$Tokname($$char))
			}
			if $$char == $$EofCode {
				goto ret1
			}
			$$char = -1
			goto $$newstate /* try again in the same state */
		}
	}

	/* reduction by production $$n */
	if $$Debug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", $$n, $$Statname($$state))
	}

	$$nt := $$n
	$$pt := $$p
	_ = $$pt // guard against "declared and not used"

	$$p -= $$R2[$$n]
	$$VAL = $$S[$$p+1]

	/* consult goto table to find next state */
	$$n = $$R1[$$n]
	$$g := $$Pgo[$$n]
	$$j := $$g + $$S[$$p].yys + 1

	if $$j >= $$Last {
		$$state = $$Act[$$g]
	} else {
		$$state = $$Act[$$j]
		if $$Chk[$$state] != -$$n {
			$$state = $$Act[$$g]
		}
	}
	// dummy call; replaced with literal code
	$$run()
	goto $$stack /* stack new state and value */
}
`
/**-----
Go file
  PackageDeclaration(main)
    PsiElement(KEYWORD_PACKAGE)('package')
    PsiWhiteSpace(' ')
    PsiElement(IDENTIFIER)('main')
  PsiWhiteSpace('\n\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// copies and protects "'s in q')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  FunctionDeclaration(chcopy)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('chcopy')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('q')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('string')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('string')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('s')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralStringImpl
            PsiElement(LITERAL_STRING)('""')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('i')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIntegerImpl
            PsiElement(LITERAL_INT)('0')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('j')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIntegerImpl
            PsiElement(LITERAL_INT)('0')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ForWithClausesStmtImpl
        PsiElement(KEYWORD_FOR)('for')
        PsiWhiteSpace(' ')
        AssignStmtImpl
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('i')
          PsiWhiteSpace(' ')
          PsiElement(=)('=')
          PsiWhiteSpace(' ')
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIntegerImpl
                PsiElement(LITERAL_INT)('0')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('i')
          PsiWhiteSpace(' ')
          PsiElement(<)('<')
          PsiWhiteSpace(' ')
          BuiltInCallExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('len')
            PsiElement(()('(')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('q')
            PsiElement())(')')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        IncDecStmt
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('i')
          PsiElement(++)('++')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          IfStmtImpl
            PsiElement(KEYWORD_IF)('if')
            PsiWhiteSpace(' ')
            RelationalExpressionImpl
              IndexExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('q')
                PsiElement([)('[')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('i')
                PsiElement(])(']')
              PsiWhiteSpace(' ')
              PsiElement(==)('==')
              PsiWhiteSpace(' ')
              LiteralExpressionImpl
                LiteralCharImpl
                  PsiElement(LITERAL_CHAR)(''"'')
            PsiWhiteSpace(' ')
            BlockStmtImpl
              PsiElement({)('{')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t\t')
              AssignStmtImpl
                ExpressionListImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('s')
                PsiWhiteSpace(' ')
                PsiElement(+=)('+=')
                PsiWhiteSpace(' ')
                ExpressionListImpl
                  AdditiveExpressionImpl
                    SliceExpressionImpl
                      LiteralExpressionImpl
                        LiteralIdentifierImpl
                          PsiElement(IDENTIFIER)('q')
                      PsiElement([)('[')
                      LiteralExpressionImpl
                        LiteralIdentifierImpl
                          PsiElement(IDENTIFIER)('j')
                      PsiElement(:)(':')
                      LiteralExpressionImpl
                        LiteralIdentifierImpl
                          PsiElement(IDENTIFIER)('i')
                      PsiElement(])(']')
                    PsiWhiteSpace(' ')
                    PsiElement(+)('+')
                    PsiWhiteSpace(' ')
                    LiteralExpressionImpl
                      LiteralStringImpl
                        PsiElement(LITERAL_STRING)('"\\"')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t\t')
              AssignStmtImpl
                ExpressionListImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('j')
                PsiWhiteSpace(' ')
                PsiElement(=)('=')
                PsiWhiteSpace(' ')
                ExpressionListImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('i')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t')
              PsiElement(})('}')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        AdditiveExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('s')
          PsiWhiteSpace(' ')
          PsiElement(+)('+')
          PsiWhiteSpace(' ')
          SliceExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('q')
            PsiElement([)('[')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('j')
            PsiElement(:)(':')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('i')
            PsiElement(])(']')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(usage)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('usage')
    PsiElement(()('(')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('fmt')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Fprintf')
          PsiElement(()('(')
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('stderr')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralStringImpl
                PsiElement(LITERAL_STRING)('"usage: yacc [-o output] [-v parsetable] input\n"')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('exit')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('1')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(bitset)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('bitset')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('set')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('Lkset')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('bit')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace(' ')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        MultiplicativeExpressionImpl
          IndexExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('set')
            PsiElement([)('[')
            MultiplicativeExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('bit')
              PsiElement(>>)('>>')
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('5')
            PsiElement(])(']')
          PsiWhiteSpace(' ')
          PsiElement(&)('&')
          PsiWhiteSpace(' ')
          ParenthesisedExpressionImpl
            PsiElement(()('(')
            MultiplicativeExpressionImpl
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('1')
              PsiWhiteSpace(' ')
              PsiElement(<<)('<<')
              PsiWhiteSpace(' ')
              BuiltInCallExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('uint')
                PsiElement(()('(')
                MultiplicativeExpressionImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('bit')
                  PsiElement(&)('&')
                  LiteralExpressionImpl
                    LiteralIntegerImpl
                      PsiElement(LITERAL_INT)('31')
                PsiElement())(')')
            PsiElement())(')')
      PsiWhiteSpace(' ')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(setbit)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('setbit')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('set')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('Lkset')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('bit')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace(' ')
      AssignStmtImpl
        ExpressionListImpl
          IndexExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('set')
            PsiElement([)('[')
            MultiplicativeExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('bit')
              PsiElement(>>)('>>')
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('5')
            PsiElement(])(']')
        PsiWhiteSpace(' ')
        PsiElement(|=)('|=')
        PsiWhiteSpace(' ')
        ExpressionListImpl
          ParenthesisedExpressionImpl
            PsiElement(()('(')
            MultiplicativeExpressionImpl
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('1')
              PsiWhiteSpace(' ')
              PsiElement(<<)('<<')
              PsiWhiteSpace(' ')
              BuiltInCallExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('uint')
                PsiElement(()('(')
                MultiplicativeExpressionImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('bit')
                  PsiElement(&)('&')
                  LiteralExpressionImpl
                    LiteralIntegerImpl
                      PsiElement(LITERAL_INT)('31')
                PsiElement())(')')
            PsiElement())(')')
      PsiWhiteSpace(' ')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(mkset)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('mkset')
    PsiElement(()('(')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Lkset')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace(' ')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        BuiltInCallExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('make')
          PsiElement(()('(')
          TypeSliceImpl
            PsiElement([)('[')
            PsiElement(])(']')
            TypeNameImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('int')
          PsiElement(,)(',')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('tbitset')
          PsiElement())(')')
      PsiWhiteSpace(' ')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// set a to the union of a and b')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// return 1 if b is not a subset of a, 0 otherwise')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  FunctionDeclaration(setunion)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('setunion')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('a')
        PsiElement(,)(',')
        PsiWhiteSpace(' ')
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('b')
        PsiWhiteSpace(' ')
        TypeSliceImpl
          PsiElement([)('[')
          PsiElement(])(']')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('sub')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIntegerImpl
            PsiElement(LITERAL_INT)('0')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ForWithClausesStmtImpl
        PsiElement(KEYWORD_FOR)('for')
        PsiWhiteSpace(' ')
        ShortVarStmtImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('i')
          PsiWhiteSpace(' ')
          PsiElement(:=)(':=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('i')
          PsiWhiteSpace(' ')
          PsiElement(<)('<')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('tbitset')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        IncDecStmt
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('i')
          PsiElement(++)('++')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ShortVarStmtImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('x')
            PsiWhiteSpace(' ')
            PsiElement(:=)(':=')
            PsiWhiteSpace(' ')
            IndexExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('a')
              PsiElement([)('[')
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('i')
              PsiElement(])(']')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ShortVarStmtImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('y')
            PsiWhiteSpace(' ')
            PsiElement(:=)(':=')
            PsiWhiteSpace(' ')
            AdditiveExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('x')
              PsiWhiteSpace(' ')
              PsiElement(|)('|')
              PsiWhiteSpace(' ')
              IndexExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('b')
                PsiElement([)('[')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('i')
                PsiElement(])(']')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          AssignStmtImpl
            ExpressionListImpl
              IndexExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('a')
                PsiElement([)('[')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('i')
                PsiElement(])(']')
            PsiWhiteSpace(' ')
            PsiElement(=)('=')
            PsiWhiteSpace(' ')
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('y')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          IfStmtImpl
            PsiElement(KEYWORD_IF)('if')
            PsiWhiteSpace(' ')
            RelationalExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('y')
              PsiWhiteSpace(' ')
              PsiElement(!=)('!=')
              PsiWhiteSpace(' ')
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('x')
            PsiWhiteSpace(' ')
            BlockStmtImpl
              PsiElement({)('{')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t\t')
              AssignStmtImpl
                ExpressionListImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('sub')
                PsiWhiteSpace(' ')
                PsiElement(=)('=')
                PsiWhiteSpace(' ')
                ExpressionListImpl
                  LiteralExpressionImpl
                    LiteralIntegerImpl
                      PsiElement(LITERAL_INT)('1')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t')
              PsiElement(})('}')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('sub')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(prlook)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('prlook')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('p')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('Lkset')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('p')
          PsiWhiteSpace(' ')
          PsiElement(==)('==')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              SelectorExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('fmt')
                PsiElement(.)('.')
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('Fprintf')
              PsiElement(()('(')
              ExpressionListImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('foutput')
                PsiElement(,)(',')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralStringImpl
                    PsiElement(LITERAL_STRING)('"\tNULL"')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ReturnStmtImpl
            PsiElement(KEYWORD_RETURN)('return')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('fmt')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Fprintf')
          PsiElement(()('(')
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('foutput')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralStringImpl
                PsiElement(LITERAL_STRING)('" { "')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ForWithClausesStmtImpl
        PsiElement(KEYWORD_FOR)('for')
        PsiWhiteSpace(' ')
        ShortVarStmtImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('j')
          PsiWhiteSpace(' ')
          PsiElement(:=)(':=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('j')
          PsiWhiteSpace(' ')
          PsiElement(<=)('<=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('ntokens')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        IncDecStmt
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('j')
          PsiElement(++)('++')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          IfStmtImpl
            PsiElement(KEYWORD_IF)('if')
            PsiWhiteSpace(' ')
            RelationalExpressionImpl
              CallOrConversionExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('bitset')
                PsiElement(()('(')
                ExpressionListImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('p')
                  PsiElement(,)(',')
                  PsiWhiteSpace(' ')
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('j')
                PsiElement())(')')
              PsiWhiteSpace(' ')
              PsiElement(!=)('!=')
              PsiWhiteSpace(' ')
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('0')
            PsiWhiteSpace(' ')
            BlockStmtImpl
              PsiElement({)('{')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t\t')
              ExpressionStmtImpl
                CallOrConversionExpressionImpl
                  SelectorExpressionImpl
                    LiteralExpressionImpl
                      LiteralIdentifierImpl
                        PsiElement(IDENTIFIER)('fmt')
                    PsiElement(.)('.')
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('Fprintf')
                  PsiElement(()('(')
                  ExpressionListImpl
                    LiteralExpressionImpl
                      LiteralIdentifierImpl
                        PsiElement(IDENTIFIER)('foutput')
                    PsiElement(,)(',')
                    PsiWhiteSpace(' ')
                    LiteralExpressionImpl
                      LiteralStringImpl
                        PsiElement(LITERAL_STRING)('"%v "')
                    PsiElement(,)(',')
                    PsiWhiteSpace(' ')
                    CallOrConversionExpressionImpl
                      LiteralExpressionImpl
                        LiteralIdentifierImpl
                          PsiElement(IDENTIFIER)('symnam')
                      PsiElement(()('(')
                      LiteralExpressionImpl
                        LiteralIdentifierImpl
                          PsiElement(IDENTIFIER)('j')
                      PsiElement())(')')
                  PsiElement())(')')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t')
              PsiElement(})('}')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('fmt')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Fprintf')
          PsiElement(()('(')
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('foutput')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralStringImpl
                PsiElement(LITERAL_STRING)('"}"')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// utility routines')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  VarDeclarationsImpl
    PsiElement(KEYWORD_VAR)('var')
    PsiWhiteSpace(' ')
    VarDeclarationImpl
      LiteralIdentifierImpl
        PsiElement(IDENTIFIER)('peekrune')
      PsiWhiteSpace(' ')
      TypeNameImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('rune')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(isdigit)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('isdigit')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('c')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('rune')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('bool')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace(' ')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LogicalAndExpressionImpl
          RelationalExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('c')
            PsiWhiteSpace(' ')
            PsiElement(>=)('>=')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralCharImpl
                PsiElement(LITERAL_CHAR)(''0'')
          PsiWhiteSpace(' ')
          PsiElement(&&)('&&')
          PsiWhiteSpace(' ')
          RelationalExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('c')
            PsiWhiteSpace(' ')
            PsiElement(<=)('<=')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralCharImpl
                PsiElement(LITERAL_CHAR)(''9'')
      PsiWhiteSpace(' ')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(isword)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('isword')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('c')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('rune')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('bool')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LogicalOrExpressionImpl
          LogicalOrExpressionImpl
            RelationalExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('c')
              PsiWhiteSpace(' ')
              PsiElement(>=)('>=')
              PsiWhiteSpace(' ')
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_HEX)('0xa0')
            PsiWhiteSpace(' ')
            PsiElement(||)('||')
            PsiWhiteSpace(' ')
            ParenthesisedExpressionImpl
              PsiElement(()('(')
              LogicalAndExpressionImpl
                RelationalExpressionImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('c')
                  PsiWhiteSpace(' ')
                  PsiElement(>=)('>=')
                  PsiWhiteSpace(' ')
                  LiteralExpressionImpl
                    LiteralCharImpl
                      PsiElement(LITERAL_CHAR)(''a'')
                PsiWhiteSpace(' ')
                PsiElement(&&)('&&')
                PsiWhiteSpace(' ')
                RelationalExpressionImpl
                  LiteralExpressionImpl
                    LiteralIdentifierImpl
                      PsiElement(IDENTIFIER)('c')
                  PsiWhiteSpace(' ')
                  PsiElement(<=)('<=')
                  PsiWhiteSpace(' ')
                  LiteralExpressionImpl
                    LiteralCharImpl
                      PsiElement(LITERAL_CHAR)(''z'')
              PsiElement())(')')
          PsiWhiteSpace(' ')
          PsiElement(||)('||')
          PsiWhiteSpace(' ')
          ParenthesisedExpressionImpl
            PsiElement(()('(')
            LogicalAndExpressionImpl
              RelationalExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('c')
                PsiWhiteSpace(' ')
                PsiElement(>=)('>=')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralCharImpl
                    PsiElement(LITERAL_CHAR)(''A'')
              PsiWhiteSpace(' ')
              PsiElement(&&)('&&')
              PsiWhiteSpace(' ')
              RelationalExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('c')
                PsiWhiteSpace(' ')
                PsiElement(<=)('<=')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralCharImpl
                    PsiElement(LITERAL_CHAR)(''Z'')
            PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(mktemp)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('mktemp')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('t')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('string')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('string')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace(' ')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('t')
      PsiWhiteSpace(' ')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// return 1 if 2 arrays are equal')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// return 0 if not equal')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  FunctionDeclaration(aryeq)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('aryeq')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('a')
        PsiWhiteSpace(' ')
        TypeSliceImpl
          PsiElement([)('[')
          PsiElement(])(']')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('b')
        PsiWhiteSpace(' ')
        TypeSliceImpl
          PsiElement([)('[')
          PsiElement(])(']')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('n')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        BuiltInCallExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('len')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('a')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          BuiltInCallExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('len')
            PsiElement(()('(')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('b')
            PsiElement())(')')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('n')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ReturnStmtImpl
            PsiElement(KEYWORD_RETURN)('return')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIntegerImpl
                PsiElement(LITERAL_INT)('0')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ForWithClausesStmtImpl
        PsiElement(KEYWORD_FOR)('for')
        PsiWhiteSpace(' ')
        ShortVarStmtImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('ll')
          PsiWhiteSpace(' ')
          PsiElement(:=)(':=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('ll')
          PsiWhiteSpace(' ')
          PsiElement(<)('<')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('n')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        IncDecStmt
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('ll')
          PsiElement(++)('++')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          IfStmtImpl
            PsiElement(KEYWORD_IF)('if')
            PsiWhiteSpace(' ')
            RelationalExpressionImpl
              IndexExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('a')
                PsiElement([)('[')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('ll')
                PsiElement(])(']')
              PsiWhiteSpace(' ')
              PsiElement(!=)('!=')
              PsiWhiteSpace(' ')
              IndexExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('b')
                PsiElement([)('[')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('ll')
                PsiElement(])(']')
            PsiWhiteSpace(' ')
            BlockStmtImpl
              PsiElement({)('{')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t\t')
              ReturnStmtImpl
                PsiElement(KEYWORD_RETURN)('return')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIntegerImpl
                    PsiElement(LITERAL_INT)('0')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t')
              PsiElement(})('}')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIntegerImpl
            PsiElement(LITERAL_INT)('1')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(putrune)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('putrune')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('f')
        PsiWhiteSpace(' ')
        TypePointerImpl
          PsiElement(*)('*')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('bufio')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Writer')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('c')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('s')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        BuiltInCallExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('string')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('c')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ForWithClausesStmtImpl
        PsiElement(KEYWORD_FOR)('for')
        PsiWhiteSpace(' ')
        ShortVarStmtImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('i')
          PsiWhiteSpace(' ')
          PsiElement(:=)(':=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('i')
          PsiWhiteSpace(' ')
          PsiElement(<)('<')
          PsiWhiteSpace(' ')
          BuiltInCallExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('len')
            PsiElement(()('(')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('s')
            PsiElement())(')')
        PsiElement(;)(';')
        PsiWhiteSpace(' ')
        IncDecStmt
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('i')
          PsiElement(++)('++')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              SelectorExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('f')
                PsiElement(.)('.')
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('WriteByte')
              PsiElement(()('(')
              IndexExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('s')
                PsiElement([)('[')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('i')
                PsiElement(])(']')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(getrune)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('getrune')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('f')
        PsiWhiteSpace(' ')
        TypePointerImpl
          PsiElement(*)('*')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('bufio')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Reader')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('rune')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      VarDeclarationsImpl
        PsiElement(KEYWORD_VAR)('var')
        PsiWhiteSpace(' ')
        VarDeclarationImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('r')
          PsiWhiteSpace(' ')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('rune')
      PsiWhiteSpace('\n\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('peekrune')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          IfStmtImpl
            PsiElement(KEYWORD_IF)('if')
            PsiWhiteSpace(' ')
            RelationalExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('peekrune')
              PsiWhiteSpace(' ')
              PsiElement(==)('==')
              PsiWhiteSpace(' ')
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('EOF')
            PsiWhiteSpace(' ')
            BlockStmtImpl
              PsiElement({)('{')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t\t')
              ReturnStmtImpl
                PsiElement(KEYWORD_RETURN)('return')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('EOF')
              PsiWhiteSpace('\n')
              PsiWhiteSpace('\t\t')
              PsiElement(})('}')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          AssignStmtImpl
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('r')
            PsiWhiteSpace(' ')
            PsiElement(=)('=')
            PsiWhiteSpace(' ')
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('peekrune')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          AssignStmtImpl
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('peekrune')
            PsiWhiteSpace(' ')
            PsiElement(=)('=')
            PsiWhiteSpace(' ')
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('0')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ReturnStmtImpl
            PsiElement(KEYWORD_RETURN)('return')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('r')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('c')
        PsiElement(,)(',')
        PsiWhiteSpace(' ')
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('n')
        PsiElement(,)(',')
        PsiWhiteSpace(' ')
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('err')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('f')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('ReadRune')
          PsiElement(()('(')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('n')
          PsiWhiteSpace(' ')
          PsiElement(==)('==')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ReturnStmtImpl
            PsiElement(KEYWORD_RETURN)('return')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('EOF')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('err')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('errorf')
              PsiElement(()('(')
              ExpressionListImpl
                LiteralExpressionImpl
                  LiteralStringImpl
                    PsiElement(LITERAL_STRING)('"read error: %v"')
                PsiElement(,)(',')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('err')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      PsiComment(SL_COMMENT)('//fmt.Printf("rune = %v n=%v\n", string(c), n);')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('c')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(ungetrune)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('ungetrune')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('f')
        PsiWhiteSpace(' ')
        TypePointerImpl
          PsiElement(*)('*')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('bufio')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Reader')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('c')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('rune')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('f')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('finput')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            BuiltInCallExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('panic')
              PsiElement(()('(')
              LiteralExpressionImpl
                LiteralStringImpl
                  PsiElement(LITERAL_STRING)('"ungetc - not finput"')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('peekrune')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            BuiltInCallExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('panic')
              PsiElement(()('(')
              LiteralExpressionImpl
                LiteralStringImpl
                  PsiElement(LITERAL_STRING)('"ungetc - 2nd unget"')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      AssignStmtImpl
        ExpressionListImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('peekrune')
        PsiWhiteSpace(' ')
        PsiElement(=)('=')
        PsiWhiteSpace(' ')
        ExpressionListImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('c')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(write)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('write')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('f')
        PsiWhiteSpace(' ')
        TypePointerImpl
          PsiElement(*)('*')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('bufio')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Writer')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('b')
        PsiWhiteSpace(' ')
        TypeSliceImpl
          PsiElement([)('[')
          PsiElement(])(']')
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('byte')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('n')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypeNameImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('int')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        BuiltInCallExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('panic')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralStringImpl
              PsiElement(LITERAL_STRING)('"write"')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        LiteralExpressionImpl
          LiteralIntegerImpl
            PsiElement(LITERAL_INT)('0')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(open)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('open')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('s')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('string')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypePointerImpl
            PsiElement(*)('*')
            TypeNameImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('bufio')
              PsiElement(.)('.')
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('Reader')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('fi')
        PsiElement(,)(',')
        PsiWhiteSpace(' ')
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('err')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('os')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Open')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('s')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('err')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('errorf')
              PsiElement(()('(')
              ExpressionListImpl
                LiteralExpressionImpl
                  LiteralStringImpl
                    PsiElement(LITERAL_STRING)('"error opening %v: %v"')
                PsiElement(,)(',')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('s')
                PsiElement(,)(',')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('err')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      PsiComment(SL_COMMENT)('//fmt.Printf("open %v\n", s);')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('bufio')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('NewReader')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('fi')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(create)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('create')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('s')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('string')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    FunctionResult
      FunctionParameterListImpl
        FunctionParameterImpl
          TypePointerImpl
            PsiElement(*)('*')
            TypeNameImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('bufio')
              PsiElement(.)('.')
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('Writer')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ShortVarStmtImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('fo')
        PsiElement(,)(',')
        PsiWhiteSpace(' ')
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('err')
        PsiWhiteSpace(' ')
        PsiElement(:=)(':=')
        PsiWhiteSpace(' ')
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('os')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Create')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('s')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('err')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('errorf')
              PsiElement(()('(')
              ExpressionListImpl
                LiteralExpressionImpl
                  LiteralStringImpl
                    PsiElement(LITERAL_STRING)('"error creating %v: %v"')
                PsiElement(,)(',')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('s')
                PsiElement(,)(',')
                PsiWhiteSpace(' ')
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('err')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      PsiComment(SL_COMMENT)('//fmt.Printf("create %v mode %v\n", s);')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ReturnStmtImpl
        PsiElement(KEYWORD_RETURN)('return')
        PsiWhiteSpace(' ')
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('bufio')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('NewWriter')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('fo')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('// write out error comment')
  PsiWhiteSpace('\n')
  PsiComment(SL_COMMENT)('//')
  PsiWhiteSpace('\n')
  FunctionDeclaration(errorf)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('errorf')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('s')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('string')
      PsiElement(,)(',')
      PsiWhiteSpace(' ')
      FunctionParameterVariadicImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('v')
        PsiWhiteSpace(' ')
        PsiElement(...)('...')
        TypeInterfaceImpl
          PsiElement(KEYWORD_INTERFACE)('interface')
          PsiElement({)('{')
          PsiElement(})('}')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IncDecStmt
        LiteralExpressionImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('nerrors')
        PsiElement(++)('++')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('fmt')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Fprintf')
          PsiElement(()('(')
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('stderr')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('s')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('v')
            PsiElement(...)('...')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('fmt')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Fprintf')
          PsiElement(()('(')
          ExpressionListImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('stderr')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralStringImpl
                PsiElement(LITERAL_STRING)('": %v:%v\n"')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('infile')
            PsiElement(,)(',')
            PsiWhiteSpace(' ')
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('lineno')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('fatfl')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIntegerImpl
              PsiElement(LITERAL_INT)('0')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('summary')
              PsiElement(()('(')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('exit')
              PsiElement(()('(')
              LiteralExpressionImpl
                LiteralIntegerImpl
                  PsiElement(LITERAL_INT)('1')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  FunctionDeclaration(exit)
    PsiElement(KEYWORD_FUNC)('func')
    PsiWhiteSpace(' ')
    LiteralIdentifierImpl
      PsiElement(IDENTIFIER)('exit')
    PsiElement(()('(')
    FunctionParameterListImpl
      FunctionParameterImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('status')
        PsiWhiteSpace(' ')
        TypeNameImpl
          LiteralIdentifierImpl
            PsiElement(IDENTIFIER)('int')
    PsiElement())(')')
    PsiWhiteSpace(' ')
    BlockStmtImpl
      PsiElement({)('{')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('ftable')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              SelectorExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('ftable')
                PsiElement(.)('.')
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('Flush')
              PsiElement(()('(')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          AssignStmtImpl
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('ftable')
            PsiWhiteSpace(' ')
            PsiElement(=)('=')
            PsiWhiteSpace(' ')
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('nil')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('foutput')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              SelectorExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('foutput')
                PsiElement(.)('.')
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('Flush')
              PsiElement(()('(')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          AssignStmtImpl
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('foutput')
            PsiWhiteSpace(' ')
            PsiElement(=)('=')
            PsiWhiteSpace(' ')
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('nil')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      IfStmtImpl
        PsiElement(KEYWORD_IF)('if')
        PsiWhiteSpace(' ')
        RelationalExpressionImpl
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('stderr')
          PsiWhiteSpace(' ')
          PsiElement(!=)('!=')
          PsiWhiteSpace(' ')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('nil')
        PsiWhiteSpace(' ')
        BlockStmtImpl
          PsiElement({)('{')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          ExpressionStmtImpl
            CallOrConversionExpressionImpl
              SelectorExpressionImpl
                LiteralExpressionImpl
                  LiteralIdentifierImpl
                    PsiElement(IDENTIFIER)('stderr')
                PsiElement(.)('.')
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('Flush')
              PsiElement(()('(')
              PsiElement())(')')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t\t')
          AssignStmtImpl
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('stderr')
            PsiWhiteSpace(' ')
            PsiElement(=)('=')
            PsiWhiteSpace(' ')
            ExpressionListImpl
              LiteralExpressionImpl
                LiteralIdentifierImpl
                  PsiElement(IDENTIFIER)('nil')
          PsiWhiteSpace('\n')
          PsiWhiteSpace('\t')
          PsiElement(})('}')
      PsiWhiteSpace('\n')
      PsiWhiteSpace('\t')
      ExpressionStmtImpl
        CallOrConversionExpressionImpl
          SelectorExpressionImpl
            LiteralExpressionImpl
              LiteralIdentifierImpl
                PsiElement(IDENTIFIER)('os')
            PsiElement(.)('.')
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('Exit')
          PsiElement(()('(')
          LiteralExpressionImpl
            LiteralIdentifierImpl
              PsiElement(IDENTIFIER)('status')
          PsiElement())(')')
      PsiWhiteSpace('\n')
      PsiElement(})('}')
  PsiWhiteSpace('\n\n')
  VarDeclarationsImpl
    PsiElement(KEYWORD_VAR)('var')
    PsiWhiteSpace(' ')
    VarDeclarationImpl
      LiteralIdentifierImpl
        PsiElement(IDENTIFIER)('yaccpar')
      PsiWhiteSpace(' ')
      TypeNameImpl
        LiteralIdentifierImpl
          PsiElement(IDENTIFIER)('string')
  PsiWhiteSpace(' ')
  PsiComment(SL_COMMENT)('// will be processed version of yaccpartext: s/$$/prefix/g')
  PsiWhiteSpace('\n')
  VarDeclarationsImpl
    PsiElement(KEYWORD_VAR)('var')
    PsiWhiteSpace(' ')
    VarDeclarationImpl
      LiteralIdentifierImpl
        PsiElement(IDENTIFIER)('yaccpartext')
      PsiWhiteSpace(' ')
      PsiElement(=)('=')
      PsiWhiteSpace(' ')
      LiteralExpressionImpl
        LiteralStringImpl
          PsiElement(LITERAL_STRING)('`\n/*\tparser for yacc output\t*/\n\nvar $$Debug = 0\n\ntype $$Lexer interface {\n\tLex(lval *$$SymType) int\n\tError(s string)\n}\n\nconst $$Flag = -1000\n\nfunc $$Tokname(c int) string {\n\tif c > 0 && c <= len($$Toknames) {\n\t\tif $$Toknames[c-1] != "" {\n\t\t\treturn $$Toknames[c-1]\n\t\t}\n\t}\n\treturn fmt.Sprintf("tok-%v", c)\n}\n\nfunc $$Statname(s int) string {\n\tif s >= 0 && s < len($$Statenames) {\n\t\tif $$Statenames[s] != "" {\n\t\t\treturn $$Statenames[s]\n\t\t}\n\t}\n\treturn fmt.Sprintf("state-%v", s)\n}\n\nfunc $$lex1(lex $$Lexer, lval *$$SymType) int {\n\tc := 0\n\tchar := lex.Lex(lval)\n\tif char <= 0 {\n\t\tc = $$Tok1[0]\n\t\tgoto out\n\t}\n\tif char < len($$Tok1) {\n\t\tc = $$Tok1[char]\n\t\tgoto out\n\t}\n\tif char >= $$Private {\n\t\tif char < $$Private+len($$Tok2) {\n\t\t\tc = $$Tok2[char-$$Private]\n\t\t\tgoto out\n\t\t}\n\t}\n\tfor i := 0; i < len($$Tok3); i += 2 {\n\t\tc = $$Tok3[i+0]\n\t\tif c == char {\n\t\t\tc = $$Tok3[i+1]\n\t\t\tgoto out\n\t\t}\n\t}\n\nout:\n\tif c == 0 {\n\t\tc = $$Tok2[1] /* unknown char */\n\t}\n\tif $$Debug >= 3 {\n\t\tfmt.Printf("lex %U %s\n", uint(char), $$Tokname(c))\n\t}\n\treturn c\n}\n\nfunc $$Parse($$lex $$Lexer) int {\n\tvar $$n int\n\tvar $$lval $$SymType\n\tvar $$VAL $$SymType\n\t$$S := make([]$$SymType, $$MaxDepth)\n\n\tNerrs := 0   /* number of errors */\n\tErrflag := 0 /* error recovery flag */\n\t$$state := 0\n\t$$char := -1\n\t$$p := -1\n\tgoto $$stack\n\nret0:\n\treturn 0\n\nret1:\n\treturn 1\n\n$$stack:\n\t/* put a state and value onto the stack */\n\tif $$Debug >= 4 {\n\t\tfmt.Printf("char %v in %v\n", $$Tokname($$char), $$Statname($$state))\n\t}\n\n\t$$p++\n\tif $$p >= len($$S) {\n\t\tnyys := make([]$$SymType, len($$S)*2)\n\t\tcopy(nyys, $$S)\n\t\t$$S = nyys\n\t}\n\t$$S[$$p] = $$VAL\n\t$$S[$$p].yys = $$state\n\n$$newstate:\n\t$$n = $$Pact[$$state]\n\tif $$n <= $$Flag {\n\t\tgoto $$default /* simple state */\n\t}\n\tif $$char < 0 {\n\t\t$$char = $$lex1($$lex, &$$lval)\n\t}\n\t$$n += $$char\n\tif $$n < 0 || $$n >= $$Last {\n\t\tgoto $$default\n\t}\n\t$$n = $$Act[$$n]\n\tif $$Chk[$$n] == $$char { /* valid shift */\n\t\t$$char = -1\n\t\t$$VAL = $$lval\n\t\t$$state = $$n\n\t\tif Errflag > 0 {\n\t\t\tErrflag--\n\t\t}\n\t\tgoto $$stack\n\t}\n\n$$default:\n\t/* default state action */\n\t$$n = $$Def[$$state]\n\tif $$n == -2 {\n\t\tif $$char < 0 {\n\t\t\t$$char = $$lex1($$lex, &$$lval)\n\t\t}\n\n\t\t/* look through exception table */\n\t\txi := 0\n\t\tfor {\n\t\t\tif $$Exca[xi+0] == -1 && $$Exca[xi+1] == $$state {\n\t\t\t\tbreak\n\t\t\t}\n\t\t\txi += 2\n\t\t}\n\t\tfor xi += 2; ; xi += 2 {\n\t\t\t$$n = $$Exca[xi+0]\n\t\t\tif $$n < 0 || $$n == $$char {\n\t\t\t\tbreak\n\t\t\t}\n\t\t}\n\t\t$$n = $$Exca[xi+1]\n\t\tif $$n < 0 {\n\t\t\tgoto ret0\n\t\t}\n\t}\n\tif $$n == 0 {\n\t\t/* error ... attempt to resume parsing */\n\t\tswitch Errflag {\n\t\tcase 0: /* brand new error */\n\t\t\t$$lex.Error("syntax error")\n\t\t\tNerrs++\n\t\t\tif $$Debug >= 1 {\n\t\t\t\tfmt.Printf("%s", $$Statname($$state))\n\t\t\t\tfmt.Printf("saw %s\n", $$Tokname($$char))\n\t\t\t}\n\t\t\tfallthrough\n\n\t\tcase 1, 2: /* incompletely recovered error ... try again */\n\t\t\tErrflag = 3\n\n\t\t\t/* find a state where "error" is a legal shift action */\n\t\t\tfor $$p >= 0 {\n\t\t\t\t$$n = $$Pact[$$S[$$p].yys] + $$ErrCode\n\t\t\t\tif $$n >= 0 && $$n < $$Last {\n\t\t\t\t\t$$state = $$Act[$$n] /* simulate a shift of "error" */\n\t\t\t\t\tif $$Chk[$$state] == $$ErrCode {\n\t\t\t\t\t\tgoto $$stack\n\t\t\t\t\t}\n\t\t\t\t}\n\n\t\t\t\t/* the current p has no shift on "error", pop stack */\n\t\t\t\tif $$Debug >= 2 {\n\t\t\t\t\tfmt.Printf("error recovery pops state %d\n", $$S[$$p].yys)\n\t\t\t\t}\n\t\t\t\t$$p--\n\t\t\t}\n\t\t\t/* there is no state on the stack with an error shift ... abort */\n\t\t\tgoto ret1\n\n\t\tcase 3: /* no shift yet; clobber input char */\n\t\t\tif $$Debug >= 2 {\n\t\t\t\tfmt.Printf("error recovery discards %s\n", $$Tokname($$char))\n\t\t\t}\n\t\t\tif $$char == $$EofCode {\n\t\t\t\tgoto ret1\n\t\t\t}\n\t\t\t$$char = -1\n\t\t\tgoto $$newstate /* try again in the same state */\n\t\t}\n\t}\n\n\t/* reduction by production $$n */\n\tif $$Debug >= 2 {\n\t\tfmt.Printf("reduce %v in:\n\t%v\n", $$n, $$Statname($$state))\n\t}\n\n\t$$nt := $$n\n\t$$pt := $$p\n\t_ = $$pt // guard against "declared and not used"\n\n\t$$p -= $$R2[$$n]\n\t$$VAL = $$S[$$p+1]\n\n\t/* consult goto table to find next state */\n\t$$n = $$R1[$$n]\n\t$$g := $$Pgo[$$n]\n\t$$j := $$g + $$S[$$p].yys + 1\n\n\tif $$j >= $$Last {\n\t\t$$state = $$Act[$$g]\n\t} else {\n\t\t$$state = $$Act[$$j]\n\t\tif $$Chk[$$state] != -$$n {\n\t\t\t$$state = $$Act[$$g]\n\t\t}\n\t}\n\t// dummy call; replaced with literal code\n\t$$run()\n\tgoto $$stack /* stack new state and value */\n}\n`')
