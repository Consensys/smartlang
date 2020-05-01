// GENERATED FILE!!!!!!!!!!!!!!!!!!!!!!!!!!

package scanner

import (
	"regexp"
	"smart/token"
)

var reHEXLITERAL *regexp.Regexp = regexp.MustCompile(`^0x([[:digit:]a-fA-F]{2})*`)
var reIDENT *regexp.Regexp = regexp.MustCompile(`^([[:alpha:]]|_)(_|[[:alnum:]])*`)
var reSTRINGLITERAL *regexp.Regexp = regexp.MustCompile(`^^"(\\"|[^"])*"`)
var reDECLITERAL *regexp.Regexp = regexp.MustCompile(`^[[:digit:]]+`)
var rePLUSEQ *regexp.Regexp = regexp.MustCompile(`^\+=`)
var rePLUS *regexp.Regexp = regexp.MustCompile(`^\+`)
var reLCURL *regexp.Regexp = regexp.MustCompile(`^{`)
var reRCURL *regexp.Regexp = regexp.MustCompile(`^}`)
var reLBRACK *regexp.Regexp = regexp.MustCompile(`^\[`)
var reRBRACK *regexp.Regexp = regexp.MustCompile(`^]`)
var reLPAREN *regexp.Regexp = regexp.MustCompile(`^\(`)
var reRPAREN *regexp.Regexp = regexp.MustCompile(`^\)`)
var reEQEQ *regexp.Regexp = regexp.MustCompile(`^==`)
var reEQ *regexp.Regexp = regexp.MustCompile(`^=`)
var reSEMI *regexp.Regexp = regexp.MustCompile(`^;`)
var reBANG *regexp.Regexp = regexp.MustCompile(`^!`)
var reLARR *regexp.Regexp = regexp.MustCompile(`^<-`)
var reCOMMA *regexp.Regexp = regexp.MustCompile(`^,`)
var reGEQ *regexp.Regexp = regexp.MustCompile(`^>=`)
var reLEQ *regexp.Regexp = regexp.MustCompile(`^<=`)
var reLSHIFT *regexp.Regexp = regexp.MustCompile(`^<<`)
var reRSHIFT *regexp.Regexp = regexp.MustCompile(`^>>`)
var reLT *regexp.Regexp = regexp.MustCompile(`^<`)
var reGT *regexp.Regexp = regexp.MustCompile(`^>`)
var reNEQ *regexp.Regexp = regexp.MustCompile(`^!=`)
var reMINUSEQ *regexp.Regexp = regexp.MustCompile(`^-=`)
var reMINUS *regexp.Regexp = regexp.MustCompile(`^-`)
var reEXP *regexp.Regexp = regexp.MustCompile(`^\*\*`)
var reSTAR *regexp.Regexp = regexp.MustCompile(`^\*`)
var reDIV *regexp.Regexp = regexp.MustCompile(`^/`)
var reDOT *regexp.Regexp = regexp.MustCompile(`^\.`)
var reCOLON *regexp.Regexp = regexp.MustCompile(`^:`)
var reLOGAND *regexp.Regexp = regexp.MustCompile(`^&&`)
var reLOGOR *regexp.Regexp = regexp.MustCompile(`^\|\|`)
var reBITAND *regexp.Regexp = regexp.MustCompile(`^&`)
var reBITOR *regexp.Regexp = regexp.MustCompile(`^\|`)
var reBITNEG *regexp.Regexp = regexp.MustCompile(`^\^`)
var reAT *regexp.Regexp = regexp.MustCompile(`^@`)

func (s *Scanner) Scan0(src []byte) (tok token.Token, lit string) {
	s.skipWhitespace()

	if s.offset == len(s.src) {
		return token.EOF, ""
	}

	var loc []int

	loc = reHEXLITERAL.FindIndex(src)
	if loc != nil {
		tok = token.HEXLITERAL
		goto found
	}
	loc = reIDENT.FindIndex(src)
	if loc != nil {
		tok = token.IDENT
		goto found
	}
	loc = reSTRINGLITERAL.FindIndex(src)
	if loc != nil {
		tok = token.STRINGLITERAL
		goto found
	}
	loc = reDECLITERAL.FindIndex(src)
	if loc != nil {
		tok = token.DECLITERAL
		goto found
	}
	loc = rePLUSEQ.FindIndex(src)
	if loc != nil {
		tok = token.PLUSEQ
		goto found
	}
	loc = rePLUS.FindIndex(src)
	if loc != nil {
		tok = token.PLUS
		goto found
	}
	loc = reLCURL.FindIndex(src)
	if loc != nil {
		tok = token.LCURL
		goto found
	}
	loc = reRCURL.FindIndex(src)
	if loc != nil {
		tok = token.RCURL
		goto found
	}
	loc = reLBRACK.FindIndex(src)
	if loc != nil {
		tok = token.LBRACK
		goto found
	}
	loc = reRBRACK.FindIndex(src)
	if loc != nil {
		tok = token.RBRACK
		goto found
	}
	loc = reLPAREN.FindIndex(src)
	if loc != nil {
		tok = token.LPAREN
		goto found
	}
	loc = reRPAREN.FindIndex(src)
	if loc != nil {
		tok = token.RPAREN
		goto found
	}
	loc = reEQEQ.FindIndex(src)
	if loc != nil {
		tok = token.EQEQ
		goto found
	}
	loc = reEQ.FindIndex(src)
	if loc != nil {
		tok = token.EQ
		goto found
	}
	loc = reSEMI.FindIndex(src)
	if loc != nil {
		tok = token.SEMI
		goto found
	}
	loc = reBANG.FindIndex(src)
	if loc != nil {
		tok = token.BANG
		goto found
	}
	loc = reLARR.FindIndex(src)
	if loc != nil {
		tok = token.LARR
		goto found
	}
	loc = reCOMMA.FindIndex(src)
	if loc != nil {
		tok = token.COMMA
		goto found
	}
	loc = reGEQ.FindIndex(src)
	if loc != nil {
		tok = token.GEQ
		goto found
	}
	loc = reLEQ.FindIndex(src)
	if loc != nil {
		tok = token.LEQ
		goto found
	}
	loc = reLSHIFT.FindIndex(src)
	if loc != nil {
		tok = token.LSHIFT
		goto found
	}
	loc = reRSHIFT.FindIndex(src)
	if loc != nil {
		tok = token.RSHIFT
		goto found
	}
	loc = reLT.FindIndex(src)
	if loc != nil {
		tok = token.LT
		goto found
	}
	loc = reGT.FindIndex(src)
	if loc != nil {
		tok = token.GT
		goto found
	}
	loc = reNEQ.FindIndex(src)
	if loc != nil {
		tok = token.NEQ
		goto found
	}
	loc = reMINUSEQ.FindIndex(src)
	if loc != nil {
		tok = token.MINUSEQ
		goto found
	}
	loc = reMINUS.FindIndex(src)
	if loc != nil {
		tok = token.MINUS
		goto found
	}
	loc = reEXP.FindIndex(src)
	if loc != nil {
		tok = token.EXP
		goto found
	}
	loc = reSTAR.FindIndex(src)
	if loc != nil {
		tok = token.STAR
		goto found
	}
	loc = reDIV.FindIndex(src)
	if loc != nil {
		tok = token.DIV
		goto found
	}
	loc = reDOT.FindIndex(src)
	if loc != nil {
		tok = token.DOT
		goto found
	}
	loc = reCOLON.FindIndex(src)
	if loc != nil {
		tok = token.COLON
		goto found
	}
	loc = reLOGAND.FindIndex(src)
	if loc != nil {
		tok = token.LOGAND
		goto found
	}
	loc = reLOGOR.FindIndex(src)
	if loc != nil {
		tok = token.LOGOR
		goto found
	}
	loc = reBITAND.FindIndex(src)
	if loc != nil {
		tok = token.BITAND
		goto found
	}
	loc = reBITOR.FindIndex(src)
	if loc != nil {
		tok = token.BITOR
		goto found
	}
	loc = reBITNEG.FindIndex(src)
	if loc != nil {
		tok = token.BITNEG
		goto found
	}
	loc = reAT.FindIndex(src)
	if loc != nil {
		tok = token.AT
		goto found
	}
	return token.ILLEGAL, "<error>"
found:
	str := string(src[loc[0]:loc[1]])
	return tok, str
}
