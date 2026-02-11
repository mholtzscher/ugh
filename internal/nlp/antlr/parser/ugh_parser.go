// Code generated from /home/michael/code/ugh/internal/nlp/antlr/UghParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // UghParser
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type UghParser struct {
	*antlr.BaseParser
}

var UghParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func ughparserParserInit() {
	staticData := &UghParserParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "'#'", "'@'", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "'!'", "'+'", "'-'", "':'",
		"','", "'('", "')'", "'&&'", "'||'",
	}
	staticData.SymbolicNames = []string{
		"", "QUOTED", "HASH_NUMBER", "PROJECT_TAG", "CONTEXT_TAG", "PROJECT_TAG_PREFIX",
		"CONTEXT_TAG_PREFIX", "SET_FIELD", "ADD_FIELD", "REMOVE_FIELD", "CLEAR_FIELD",
		"KW_ADD", "KW_CREATE", "KW_NEW", "KW_SET", "KW_EDIT", "KW_UPDATE", "KW_FIND",
		"KW_SHOW", "KW_LIST", "KW_FILTER", "KW_VIEW", "KW_CONTEXT", "KW_AND",
		"KW_OR", "KW_NOT", "CLEAR_OP", "ADD_OP", "REMOVE_OP", "COLON", "COMMA",
		"LPAREN", "RPAREN", "AND_OP", "OR_OP", "IDENT", "WS",
	}
	staticData.RuleNames = []string{
		"root", "command", "createCommand", "createVerb", "createPart", "createOp",
		"createText", "updateCommand", "updateVerb", "targetRef", "operation",
		"filterCommand", "filterVerb", "filterOrExpr", "filterAndExpr", "filterNotExpr",
		"filterAtom", "filterPredicate", "filterFieldPredicate", "filterTagPredicate",
		"filterTextPredicate", "filterValue", "filterValueWord", "orOp", "andOp",
		"notOp", "viewCommand", "viewTarget", "contextCommand", "contextArg",
		"setOp", "addFieldOp", "removeFieldOp", "clearOp", "tagOp", "opValue",
		"opValueWord", "ident",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 36, 249, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 85,
		8, 1, 1, 2, 1, 2, 5, 2, 89, 8, 2, 10, 2, 12, 2, 92, 9, 2, 1, 3, 1, 3, 1,
		4, 1, 4, 3, 4, 98, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 105, 8, 5,
		1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 111, 8, 6, 1, 7, 1, 7, 3, 7, 115, 8, 7, 1,
		7, 5, 7, 118, 8, 7, 10, 7, 12, 7, 121, 9, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1,
		10, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 132, 8, 10, 1, 11, 1, 11, 1, 11,
		1, 12, 1, 12, 1, 13, 1, 13, 1, 13, 1, 13, 5, 13, 143, 8, 13, 10, 13, 12,
		13, 146, 9, 13, 1, 14, 1, 14, 1, 14, 1, 14, 5, 14, 152, 8, 14, 10, 14,
		12, 14, 155, 9, 14, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 161, 8, 15, 1, 16,
		1, 16, 1, 16, 1, 16, 1, 16, 3, 16, 168, 8, 16, 1, 17, 1, 17, 1, 17, 3,
		17, 173, 8, 17, 1, 18, 1, 18, 1, 18, 1, 19, 1, 19, 1, 20, 1, 20, 1, 21,
		1, 21, 4, 21, 184, 8, 21, 11, 21, 12, 21, 185, 3, 21, 188, 8, 21, 1, 22,
		1, 22, 1, 22, 1, 22, 3, 22, 194, 8, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1,
		25, 1, 25, 1, 26, 1, 26, 3, 26, 204, 8, 26, 1, 27, 1, 27, 1, 28, 1, 28,
		3, 28, 210, 8, 28, 1, 29, 1, 29, 1, 29, 3, 29, 215, 8, 29, 1, 30, 1, 30,
		1, 30, 1, 31, 1, 31, 1, 31, 1, 32, 1, 32, 1, 32, 1, 33, 1, 33, 1, 33, 3,
		33, 229, 8, 33, 1, 34, 1, 34, 1, 35, 1, 35, 4, 35, 235, 8, 35, 11, 35,
		12, 35, 236, 3, 35, 239, 8, 35, 1, 36, 1, 36, 1, 36, 1, 36, 3, 36, 245,
		8, 36, 1, 37, 1, 37, 1, 37, 0, 0, 38, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18,
		20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54,
		56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 0, 9, 1, 0, 11, 13, 1, 0, 14, 16,
		2, 0, 2, 2, 35, 35, 1, 0, 17, 20, 1, 0, 3, 4, 2, 0, 24, 24, 34, 34, 2,
		0, 23, 23, 33, 33, 1, 0, 25, 26, 2, 0, 11, 25, 35, 35, 250, 0, 76, 1, 0,
		0, 0, 2, 84, 1, 0, 0, 0, 4, 86, 1, 0, 0, 0, 6, 93, 1, 0, 0, 0, 8, 97, 1,
		0, 0, 0, 10, 104, 1, 0, 0, 0, 12, 110, 1, 0, 0, 0, 14, 112, 1, 0, 0, 0,
		16, 122, 1, 0, 0, 0, 18, 124, 1, 0, 0, 0, 20, 131, 1, 0, 0, 0, 22, 133,
		1, 0, 0, 0, 24, 136, 1, 0, 0, 0, 26, 138, 1, 0, 0, 0, 28, 147, 1, 0, 0,
		0, 30, 160, 1, 0, 0, 0, 32, 167, 1, 0, 0, 0, 34, 172, 1, 0, 0, 0, 36, 174,
		1, 0, 0, 0, 38, 177, 1, 0, 0, 0, 40, 179, 1, 0, 0, 0, 42, 187, 1, 0, 0,
		0, 44, 193, 1, 0, 0, 0, 46, 195, 1, 0, 0, 0, 48, 197, 1, 0, 0, 0, 50, 199,
		1, 0, 0, 0, 52, 201, 1, 0, 0, 0, 54, 205, 1, 0, 0, 0, 56, 207, 1, 0, 0,
		0, 58, 214, 1, 0, 0, 0, 60, 216, 1, 0, 0, 0, 62, 219, 1, 0, 0, 0, 64, 222,
		1, 0, 0, 0, 66, 228, 1, 0, 0, 0, 68, 230, 1, 0, 0, 0, 70, 238, 1, 0, 0,
		0, 72, 244, 1, 0, 0, 0, 74, 246, 1, 0, 0, 0, 76, 77, 3, 2, 1, 0, 77, 78,
		5, 0, 0, 1, 78, 1, 1, 0, 0, 0, 79, 85, 3, 4, 2, 0, 80, 85, 3, 14, 7, 0,
		81, 85, 3, 22, 11, 0, 82, 85, 3, 52, 26, 0, 83, 85, 3, 56, 28, 0, 84, 79,
		1, 0, 0, 0, 84, 80, 1, 0, 0, 0, 84, 81, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0,
		84, 83, 1, 0, 0, 0, 85, 3, 1, 0, 0, 0, 86, 90, 3, 6, 3, 0, 87, 89, 3, 8,
		4, 0, 88, 87, 1, 0, 0, 0, 89, 92, 1, 0, 0, 0, 90, 88, 1, 0, 0, 0, 90, 91,
		1, 0, 0, 0, 91, 5, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 93, 94, 7, 0, 0, 0,
		94, 7, 1, 0, 0, 0, 95, 98, 3, 10, 5, 0, 96, 98, 3, 12, 6, 0, 97, 95, 1,
		0, 0, 0, 97, 96, 1, 0, 0, 0, 98, 9, 1, 0, 0, 0, 99, 105, 3, 60, 30, 0,
		100, 105, 3, 62, 31, 0, 101, 105, 3, 64, 32, 0, 102, 105, 3, 66, 33, 0,
		103, 105, 3, 68, 34, 0, 104, 99, 1, 0, 0, 0, 104, 100, 1, 0, 0, 0, 104,
		101, 1, 0, 0, 0, 104, 102, 1, 0, 0, 0, 104, 103, 1, 0, 0, 0, 105, 11, 1,
		0, 0, 0, 106, 111, 3, 74, 37, 0, 107, 111, 5, 1, 0, 0, 108, 111, 5, 2,
		0, 0, 109, 111, 5, 30, 0, 0, 110, 106, 1, 0, 0, 0, 110, 107, 1, 0, 0, 0,
		110, 108, 1, 0, 0, 0, 110, 109, 1, 0, 0, 0, 111, 13, 1, 0, 0, 0, 112, 114,
		3, 16, 8, 0, 113, 115, 3, 18, 9, 0, 114, 113, 1, 0, 0, 0, 114, 115, 1,
		0, 0, 0, 115, 119, 1, 0, 0, 0, 116, 118, 3, 20, 10, 0, 117, 116, 1, 0,
		0, 0, 118, 121, 1, 0, 0, 0, 119, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0,
		120, 15, 1, 0, 0, 0, 121, 119, 1, 0, 0, 0, 122, 123, 7, 1, 0, 0, 123, 17,
		1, 0, 0, 0, 124, 125, 7, 2, 0, 0, 125, 19, 1, 0, 0, 0, 126, 132, 3, 60,
		30, 0, 127, 132, 3, 62, 31, 0, 128, 132, 3, 64, 32, 0, 129, 132, 3, 66,
		33, 0, 130, 132, 3, 68, 34, 0, 131, 126, 1, 0, 0, 0, 131, 127, 1, 0, 0,
		0, 131, 128, 1, 0, 0, 0, 131, 129, 1, 0, 0, 0, 131, 130, 1, 0, 0, 0, 132,
		21, 1, 0, 0, 0, 133, 134, 3, 24, 12, 0, 134, 135, 3, 26, 13, 0, 135, 23,
		1, 0, 0, 0, 136, 137, 7, 3, 0, 0, 137, 25, 1, 0, 0, 0, 138, 144, 3, 28,
		14, 0, 139, 140, 3, 46, 23, 0, 140, 141, 3, 28, 14, 0, 141, 143, 1, 0,
		0, 0, 142, 139, 1, 0, 0, 0, 143, 146, 1, 0, 0, 0, 144, 142, 1, 0, 0, 0,
		144, 145, 1, 0, 0, 0, 145, 27, 1, 0, 0, 0, 146, 144, 1, 0, 0, 0, 147, 153,
		3, 30, 15, 0, 148, 149, 3, 48, 24, 0, 149, 150, 3, 30, 15, 0, 150, 152,
		1, 0, 0, 0, 151, 148, 1, 0, 0, 0, 152, 155, 1, 0, 0, 0, 153, 151, 1, 0,
		0, 0, 153, 154, 1, 0, 0, 0, 154, 29, 1, 0, 0, 0, 155, 153, 1, 0, 0, 0,
		156, 157, 3, 50, 25, 0, 157, 158, 3, 32, 16, 0, 158, 161, 1, 0, 0, 0, 159,
		161, 3, 32, 16, 0, 160, 156, 1, 0, 0, 0, 160, 159, 1, 0, 0, 0, 161, 31,
		1, 0, 0, 0, 162, 163, 5, 31, 0, 0, 163, 164, 3, 26, 13, 0, 164, 165, 5,
		32, 0, 0, 165, 168, 1, 0, 0, 0, 166, 168, 3, 34, 17, 0, 167, 162, 1, 0,
		0, 0, 167, 166, 1, 0, 0, 0, 168, 33, 1, 0, 0, 0, 169, 173, 3, 36, 18, 0,
		170, 173, 3, 38, 19, 0, 171, 173, 3, 40, 20, 0, 172, 169, 1, 0, 0, 0, 172,
		170, 1, 0, 0, 0, 172, 171, 1, 0, 0, 0, 173, 35, 1, 0, 0, 0, 174, 175, 5,
		7, 0, 0, 175, 176, 3, 42, 21, 0, 176, 37, 1, 0, 0, 0, 177, 178, 7, 4, 0,
		0, 178, 39, 1, 0, 0, 0, 179, 180, 3, 42, 21, 0, 180, 41, 1, 0, 0, 0, 181,
		188, 5, 1, 0, 0, 182, 184, 3, 44, 22, 0, 183, 182, 1, 0, 0, 0, 184, 185,
		1, 0, 0, 0, 185, 183, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 188, 1, 0,
		0, 0, 187, 181, 1, 0, 0, 0, 187, 183, 1, 0, 0, 0, 188, 43, 1, 0, 0, 0,
		189, 194, 3, 74, 37, 0, 190, 194, 5, 2, 0, 0, 191, 194, 5, 29, 0, 0, 192,
		194, 5, 30, 0, 0, 193, 189, 1, 0, 0, 0, 193, 190, 1, 0, 0, 0, 193, 191,
		1, 0, 0, 0, 193, 192, 1, 0, 0, 0, 194, 45, 1, 0, 0, 0, 195, 196, 7, 5,
		0, 0, 196, 47, 1, 0, 0, 0, 197, 198, 7, 6, 0, 0, 198, 49, 1, 0, 0, 0, 199,
		200, 7, 7, 0, 0, 200, 51, 1, 0, 0, 0, 201, 203, 5, 21, 0, 0, 202, 204,
		3, 54, 27, 0, 203, 202, 1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 53, 1, 0,
		0, 0, 205, 206, 3, 74, 37, 0, 206, 55, 1, 0, 0, 0, 207, 209, 5, 22, 0,
		0, 208, 210, 3, 58, 29, 0, 209, 208, 1, 0, 0, 0, 209, 210, 1, 0, 0, 0,
		210, 57, 1, 0, 0, 0, 211, 215, 5, 3, 0, 0, 212, 215, 5, 4, 0, 0, 213, 215,
		3, 74, 37, 0, 214, 211, 1, 0, 0, 0, 214, 212, 1, 0, 0, 0, 214, 213, 1,
		0, 0, 0, 215, 59, 1, 0, 0, 0, 216, 217, 5, 7, 0, 0, 217, 218, 3, 70, 35,
		0, 218, 61, 1, 0, 0, 0, 219, 220, 5, 8, 0, 0, 220, 221, 3, 70, 35, 0, 221,
		63, 1, 0, 0, 0, 222, 223, 5, 9, 0, 0, 223, 224, 3, 70, 35, 0, 224, 65,
		1, 0, 0, 0, 225, 229, 5, 10, 0, 0, 226, 227, 5, 26, 0, 0, 227, 229, 3,
		74, 37, 0, 228, 225, 1, 0, 0, 0, 228, 226, 1, 0, 0, 0, 229, 67, 1, 0, 0,
		0, 230, 231, 7, 4, 0, 0, 231, 69, 1, 0, 0, 0, 232, 239, 5, 1, 0, 0, 233,
		235, 3, 72, 36, 0, 234, 233, 1, 0, 0, 0, 235, 236, 1, 0, 0, 0, 236, 234,
		1, 0, 0, 0, 236, 237, 1, 0, 0, 0, 237, 239, 1, 0, 0, 0, 238, 232, 1, 0,
		0, 0, 238, 234, 1, 0, 0, 0, 239, 71, 1, 0, 0, 0, 240, 245, 3, 74, 37, 0,
		241, 245, 5, 2, 0, 0, 242, 245, 5, 29, 0, 0, 243, 245, 5, 30, 0, 0, 244,
		240, 1, 0, 0, 0, 244, 241, 1, 0, 0, 0, 244, 242, 1, 0, 0, 0, 244, 243,
		1, 0, 0, 0, 245, 73, 1, 0, 0, 0, 246, 247, 7, 8, 0, 0, 247, 75, 1, 0, 0,
		0, 23, 84, 90, 97, 104, 110, 114, 119, 131, 144, 153, 160, 167, 172, 185,
		187, 193, 203, 209, 214, 228, 236, 238, 244,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// UghParserInit initializes any static state used to implement UghParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewUghParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func UghParserInit() {
	staticData := &UghParserParserStaticData
	staticData.once.Do(ughparserParserInit)
}

// NewUghParser produces a new parser instance for the optional input antlr.TokenStream.
func NewUghParser(input antlr.TokenStream) *UghParser {
	UghParserInit()
	this := new(UghParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &UghParserParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "UghParser.g4"

	return this
}

// UghParser tokens.
const (
	UghParserEOF                = antlr.TokenEOF
	UghParserQUOTED             = 1
	UghParserHASH_NUMBER        = 2
	UghParserPROJECT_TAG        = 3
	UghParserCONTEXT_TAG        = 4
	UghParserPROJECT_TAG_PREFIX = 5
	UghParserCONTEXT_TAG_PREFIX = 6
	UghParserSET_FIELD          = 7
	UghParserADD_FIELD          = 8
	UghParserREMOVE_FIELD       = 9
	UghParserCLEAR_FIELD        = 10
	UghParserKW_ADD             = 11
	UghParserKW_CREATE          = 12
	UghParserKW_NEW             = 13
	UghParserKW_SET             = 14
	UghParserKW_EDIT            = 15
	UghParserKW_UPDATE          = 16
	UghParserKW_FIND            = 17
	UghParserKW_SHOW            = 18
	UghParserKW_LIST            = 19
	UghParserKW_FILTER          = 20
	UghParserKW_VIEW            = 21
	UghParserKW_CONTEXT         = 22
	UghParserKW_AND             = 23
	UghParserKW_OR              = 24
	UghParserKW_NOT             = 25
	UghParserCLEAR_OP           = 26
	UghParserADD_OP             = 27
	UghParserREMOVE_OP          = 28
	UghParserCOLON              = 29
	UghParserCOMMA              = 30
	UghParserLPAREN             = 31
	UghParserRPAREN             = 32
	UghParserAND_OP             = 33
	UghParserOR_OP              = 34
	UghParserIDENT              = 35
	UghParserWS                 = 36
)

// UghParser rules.
const (
	UghParserRULE_root                 = 0
	UghParserRULE_command              = 1
	UghParserRULE_createCommand        = 2
	UghParserRULE_createVerb           = 3
	UghParserRULE_createPart           = 4
	UghParserRULE_createOp             = 5
	UghParserRULE_createText           = 6
	UghParserRULE_updateCommand        = 7
	UghParserRULE_updateVerb           = 8
	UghParserRULE_targetRef            = 9
	UghParserRULE_operation            = 10
	UghParserRULE_filterCommand        = 11
	UghParserRULE_filterVerb           = 12
	UghParserRULE_filterOrExpr         = 13
	UghParserRULE_filterAndExpr        = 14
	UghParserRULE_filterNotExpr        = 15
	UghParserRULE_filterAtom           = 16
	UghParserRULE_filterPredicate      = 17
	UghParserRULE_filterFieldPredicate = 18
	UghParserRULE_filterTagPredicate   = 19
	UghParserRULE_filterTextPredicate  = 20
	UghParserRULE_filterValue          = 21
	UghParserRULE_filterValueWord      = 22
	UghParserRULE_orOp                 = 23
	UghParserRULE_andOp                = 24
	UghParserRULE_notOp                = 25
	UghParserRULE_viewCommand          = 26
	UghParserRULE_viewTarget           = 27
	UghParserRULE_contextCommand       = 28
	UghParserRULE_contextArg           = 29
	UghParserRULE_setOp                = 30
	UghParserRULE_addFieldOp           = 31
	UghParserRULE_removeFieldOp        = 32
	UghParserRULE_clearOp              = 33
	UghParserRULE_tagOp                = 34
	UghParserRULE_opValue              = 35
	UghParserRULE_opValueWord          = 36
	UghParserRULE_ident                = 37
)

// IRootContext is an interface to support dynamic dispatch.
type IRootContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Command() ICommandContext
	EOF() antlr.TerminalNode

	// IsRootContext differentiates from other interfaces.
	IsRootContext()
}

type RootContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRootContext() *RootContext {
	var p = new(RootContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_root
	return p
}

func InitEmptyRootContext(p *RootContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_root
}

func (*RootContext) IsRootContext() {}

func NewRootContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RootContext {
	var p = new(RootContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_root

	return p
}

func (s *RootContext) GetParser() antlr.Parser { return s.parser }

func (s *RootContext) Command() ICommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommandContext)
}

func (s *RootContext) EOF() antlr.TerminalNode {
	return s.GetToken(UghParserEOF, 0)
}

func (s *RootContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RootContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RootContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterRoot(s)
	}
}

func (s *RootContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitRoot(s)
	}
}

func (s *RootContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitRoot(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) Root() (localctx IRootContext) {
	localctx = NewRootContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, UghParserRULE_root)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(76)
		p.Command()
	}
	{
		p.SetState(77)
		p.Match(UghParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommandContext is an interface to support dynamic dispatch.
type ICommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CreateCommand() ICreateCommandContext
	UpdateCommand() IUpdateCommandContext
	FilterCommand() IFilterCommandContext
	ViewCommand() IViewCommandContext
	ContextCommand() IContextCommandContext

	// IsCommandContext differentiates from other interfaces.
	IsCommandContext()
}

type CommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandContext() *CommandContext {
	var p = new(CommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_command
	return p
}

func InitEmptyCommandContext(p *CommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_command
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_command

	return p
}

func (s *CommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandContext) CreateCommand() ICreateCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateCommandContext)
}

func (s *CommandContext) UpdateCommand() IUpdateCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUpdateCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUpdateCommandContext)
}

func (s *CommandContext) FilterCommand() IFilterCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterCommandContext)
}

func (s *CommandContext) ViewCommand() IViewCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IViewCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IViewCommandContext)
}

func (s *CommandContext) ContextCommand() IContextCommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContextCommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContextCommandContext)
}

func (s *CommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterCommand(s)
	}
}

func (s *CommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitCommand(s)
	}
}

func (s *CommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) Command() (localctx ICommandContext) {
	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, UghParserRULE_command)
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(79)
			p.CreateCommand()
		}

	case UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(80)
			p.UpdateCommand()
		}

	case UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(81)
			p.FilterCommand()
		}

	case UghParserKW_VIEW:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(82)
			p.ViewCommand()
		}

	case UghParserKW_CONTEXT:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(83)
			p.ContextCommand()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateCommandContext is an interface to support dynamic dispatch.
type ICreateCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CreateVerb() ICreateVerbContext
	AllCreatePart() []ICreatePartContext
	CreatePart(i int) ICreatePartContext

	// IsCreateCommandContext differentiates from other interfaces.
	IsCreateCommandContext()
}

type CreateCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateCommandContext() *CreateCommandContext {
	var p = new(CreateCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createCommand
	return p
}

func InitEmptyCreateCommandContext(p *CreateCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createCommand
}

func (*CreateCommandContext) IsCreateCommandContext() {}

func NewCreateCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateCommandContext {
	var p = new(CreateCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_createCommand

	return p
}

func (s *CreateCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateCommandContext) CreateVerb() ICreateVerbContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateVerbContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateVerbContext)
}

func (s *CreateCommandContext) AllCreatePart() []ICreatePartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICreatePartContext); ok {
			len++
		}
	}

	tst := make([]ICreatePartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICreatePartContext); ok {
			tst[i] = t.(ICreatePartContext)
			i++
		}
	}

	return tst
}

func (s *CreateCommandContext) CreatePart(i int) ICreatePartContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreatePartContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreatePartContext)
}

func (s *CreateCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterCreateCommand(s)
	}
}

func (s *CreateCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitCreateCommand(s)
	}
}

func (s *CreateCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitCreateCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) CreateCommand() (localctx ICreateCommandContext) {
	localctx = NewCreateCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, UghParserRULE_createCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(86)
		p.CreateVerb()
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&35567697822) != 0 {
		{
			p.SetState(87)
			p.CreatePart()
		}

		p.SetState(92)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateVerbContext is an interface to support dynamic dispatch.
type ICreateVerbContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_ADD() antlr.TerminalNode
	KW_CREATE() antlr.TerminalNode
	KW_NEW() antlr.TerminalNode

	// IsCreateVerbContext differentiates from other interfaces.
	IsCreateVerbContext()
}

type CreateVerbContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateVerbContext() *CreateVerbContext {
	var p = new(CreateVerbContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createVerb
	return p
}

func InitEmptyCreateVerbContext(p *CreateVerbContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createVerb
}

func (*CreateVerbContext) IsCreateVerbContext() {}

func NewCreateVerbContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateVerbContext {
	var p = new(CreateVerbContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_createVerb

	return p
}

func (s *CreateVerbContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateVerbContext) KW_ADD() antlr.TerminalNode {
	return s.GetToken(UghParserKW_ADD, 0)
}

func (s *CreateVerbContext) KW_CREATE() antlr.TerminalNode {
	return s.GetToken(UghParserKW_CREATE, 0)
}

func (s *CreateVerbContext) KW_NEW() antlr.TerminalNode {
	return s.GetToken(UghParserKW_NEW, 0)
}

func (s *CreateVerbContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateVerbContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateVerbContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterCreateVerb(s)
	}
}

func (s *CreateVerbContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitCreateVerb(s)
	}
}

func (s *CreateVerbContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitCreateVerb(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) CreateVerb() (localctx ICreateVerbContext) {
	localctx = NewCreateVerbContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, UghParserRULE_createVerb)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(93)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&14336) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreatePartContext is an interface to support dynamic dispatch.
type ICreatePartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CreateOp() ICreateOpContext
	CreateText() ICreateTextContext

	// IsCreatePartContext differentiates from other interfaces.
	IsCreatePartContext()
}

type CreatePartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreatePartContext() *CreatePartContext {
	var p = new(CreatePartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createPart
	return p
}

func InitEmptyCreatePartContext(p *CreatePartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createPart
}

func (*CreatePartContext) IsCreatePartContext() {}

func NewCreatePartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreatePartContext {
	var p = new(CreatePartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_createPart

	return p
}

func (s *CreatePartContext) GetParser() antlr.Parser { return s.parser }

func (s *CreatePartContext) CreateOp() ICreateOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateOpContext)
}

func (s *CreatePartContext) CreateText() ICreateTextContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateTextContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateTextContext)
}

func (s *CreatePartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreatePartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreatePartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterCreatePart(s)
	}
}

func (s *CreatePartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitCreatePart(s)
	}
}

func (s *CreatePartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitCreatePart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) CreatePart() (localctx ICreatePartContext) {
	localctx = NewCreatePartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, UghParserRULE_createPart)
	p.SetState(97)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserPROJECT_TAG, UghParserCONTEXT_TAG, UghParserSET_FIELD, UghParserADD_FIELD, UghParserREMOVE_FIELD, UghParserCLEAR_FIELD, UghParserCLEAR_OP:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(95)
			p.CreateOp()
		}

	case UghParserQUOTED, UghParserHASH_NUMBER, UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserCOMMA, UghParserIDENT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(96)
			p.CreateText()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateOpContext is an interface to support dynamic dispatch.
type ICreateOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SetOp() ISetOpContext
	AddFieldOp() IAddFieldOpContext
	RemoveFieldOp() IRemoveFieldOpContext
	ClearOp() IClearOpContext
	TagOp() ITagOpContext

	// IsCreateOpContext differentiates from other interfaces.
	IsCreateOpContext()
}

type CreateOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateOpContext() *CreateOpContext {
	var p = new(CreateOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createOp
	return p
}

func InitEmptyCreateOpContext(p *CreateOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createOp
}

func (*CreateOpContext) IsCreateOpContext() {}

func NewCreateOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateOpContext {
	var p = new(CreateOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_createOp

	return p
}

func (s *CreateOpContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateOpContext) SetOp() ISetOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetOpContext)
}

func (s *CreateOpContext) AddFieldOp() IAddFieldOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAddFieldOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAddFieldOpContext)
}

func (s *CreateOpContext) RemoveFieldOp() IRemoveFieldOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRemoveFieldOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRemoveFieldOpContext)
}

func (s *CreateOpContext) ClearOp() IClearOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClearOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClearOpContext)
}

func (s *CreateOpContext) TagOp() ITagOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagOpContext)
}

func (s *CreateOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterCreateOp(s)
	}
}

func (s *CreateOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitCreateOp(s)
	}
}

func (s *CreateOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitCreateOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) CreateOp() (localctx ICreateOpContext) {
	localctx = NewCreateOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, UghParserRULE_createOp)
	p.SetState(104)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserSET_FIELD:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(99)
			p.SetOp()
		}

	case UghParserADD_FIELD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(100)
			p.AddFieldOp()
		}

	case UghParserREMOVE_FIELD:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(101)
			p.RemoveFieldOp()
		}

	case UghParserCLEAR_FIELD, UghParserCLEAR_OP:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(102)
			p.ClearOp()
		}

	case UghParserPROJECT_TAG, UghParserCONTEXT_TAG:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(103)
			p.TagOp()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateTextContext is an interface to support dynamic dispatch.
type ICreateTextContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ident() IIdentContext
	QUOTED() antlr.TerminalNode
	HASH_NUMBER() antlr.TerminalNode
	COMMA() antlr.TerminalNode

	// IsCreateTextContext differentiates from other interfaces.
	IsCreateTextContext()
}

type CreateTextContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateTextContext() *CreateTextContext {
	var p = new(CreateTextContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createText
	return p
}

func InitEmptyCreateTextContext(p *CreateTextContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_createText
}

func (*CreateTextContext) IsCreateTextContext() {}

func NewCreateTextContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateTextContext {
	var p = new(CreateTextContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_createText

	return p
}

func (s *CreateTextContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateTextContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *CreateTextContext) QUOTED() antlr.TerminalNode {
	return s.GetToken(UghParserQUOTED, 0)
}

func (s *CreateTextContext) HASH_NUMBER() antlr.TerminalNode {
	return s.GetToken(UghParserHASH_NUMBER, 0)
}

func (s *CreateTextContext) COMMA() antlr.TerminalNode {
	return s.GetToken(UghParserCOMMA, 0)
}

func (s *CreateTextContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateTextContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateTextContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterCreateText(s)
	}
}

func (s *CreateTextContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitCreateText(s)
	}
}

func (s *CreateTextContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitCreateText(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) CreateText() (localctx ICreateTextContext) {
	localctx = NewCreateTextContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, UghParserRULE_createText)
	p.SetState(110)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserIDENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(106)
			p.Ident()
		}

	case UghParserQUOTED:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(107)
			p.Match(UghParserQUOTED)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserHASH_NUMBER:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(108)
			p.Match(UghParserHASH_NUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCOMMA:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(109)
			p.Match(UghParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUpdateCommandContext is an interface to support dynamic dispatch.
type IUpdateCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UpdateVerb() IUpdateVerbContext
	TargetRef() ITargetRefContext
	AllOperation() []IOperationContext
	Operation(i int) IOperationContext

	// IsUpdateCommandContext differentiates from other interfaces.
	IsUpdateCommandContext()
}

type UpdateCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUpdateCommandContext() *UpdateCommandContext {
	var p = new(UpdateCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_updateCommand
	return p
}

func InitEmptyUpdateCommandContext(p *UpdateCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_updateCommand
}

func (*UpdateCommandContext) IsUpdateCommandContext() {}

func NewUpdateCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UpdateCommandContext {
	var p = new(UpdateCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_updateCommand

	return p
}

func (s *UpdateCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *UpdateCommandContext) UpdateVerb() IUpdateVerbContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUpdateVerbContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUpdateVerbContext)
}

func (s *UpdateCommandContext) TargetRef() ITargetRefContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITargetRefContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITargetRefContext)
}

func (s *UpdateCommandContext) AllOperation() []IOperationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOperationContext); ok {
			len++
		}
	}

	tst := make([]IOperationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOperationContext); ok {
			tst[i] = t.(IOperationContext)
			i++
		}
	}

	return tst
}

func (s *UpdateCommandContext) Operation(i int) IOperationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperationContext)
}

func (s *UpdateCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UpdateCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UpdateCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterUpdateCommand(s)
	}
}

func (s *UpdateCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitUpdateCommand(s)
	}
}

func (s *UpdateCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitUpdateCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) UpdateCommand() (localctx IUpdateCommandContext) {
	localctx = NewUpdateCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, UghParserRULE_updateCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(112)
		p.UpdateVerb()
	}
	p.SetState(114)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == UghParserHASH_NUMBER || _la == UghParserIDENT {
		{
			p.SetState(113)
			p.TargetRef()
		}

	}
	p.SetState(119)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&67110808) != 0 {
		{
			p.SetState(116)
			p.Operation()
		}

		p.SetState(121)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUpdateVerbContext is an interface to support dynamic dispatch.
type IUpdateVerbContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_SET() antlr.TerminalNode
	KW_EDIT() antlr.TerminalNode
	KW_UPDATE() antlr.TerminalNode

	// IsUpdateVerbContext differentiates from other interfaces.
	IsUpdateVerbContext()
}

type UpdateVerbContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUpdateVerbContext() *UpdateVerbContext {
	var p = new(UpdateVerbContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_updateVerb
	return p
}

func InitEmptyUpdateVerbContext(p *UpdateVerbContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_updateVerb
}

func (*UpdateVerbContext) IsUpdateVerbContext() {}

func NewUpdateVerbContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UpdateVerbContext {
	var p = new(UpdateVerbContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_updateVerb

	return p
}

func (s *UpdateVerbContext) GetParser() antlr.Parser { return s.parser }

func (s *UpdateVerbContext) KW_SET() antlr.TerminalNode {
	return s.GetToken(UghParserKW_SET, 0)
}

func (s *UpdateVerbContext) KW_EDIT() antlr.TerminalNode {
	return s.GetToken(UghParserKW_EDIT, 0)
}

func (s *UpdateVerbContext) KW_UPDATE() antlr.TerminalNode {
	return s.GetToken(UghParserKW_UPDATE, 0)
}

func (s *UpdateVerbContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UpdateVerbContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UpdateVerbContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterUpdateVerb(s)
	}
}

func (s *UpdateVerbContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitUpdateVerb(s)
	}
}

func (s *UpdateVerbContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitUpdateVerb(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) UpdateVerb() (localctx IUpdateVerbContext) {
	localctx = NewUpdateVerbContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, UghParserRULE_updateVerb)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(122)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&114688) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITargetRefContext is an interface to support dynamic dispatch.
type ITargetRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HASH_NUMBER() antlr.TerminalNode
	IDENT() antlr.TerminalNode

	// IsTargetRefContext differentiates from other interfaces.
	IsTargetRefContext()
}

type TargetRefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTargetRefContext() *TargetRefContext {
	var p = new(TargetRefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_targetRef
	return p
}

func InitEmptyTargetRefContext(p *TargetRefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_targetRef
}

func (*TargetRefContext) IsTargetRefContext() {}

func NewTargetRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TargetRefContext {
	var p = new(TargetRefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_targetRef

	return p
}

func (s *TargetRefContext) GetParser() antlr.Parser { return s.parser }

func (s *TargetRefContext) HASH_NUMBER() antlr.TerminalNode {
	return s.GetToken(UghParserHASH_NUMBER, 0)
}

func (s *TargetRefContext) IDENT() antlr.TerminalNode {
	return s.GetToken(UghParserIDENT, 0)
}

func (s *TargetRefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TargetRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TargetRefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterTargetRef(s)
	}
}

func (s *TargetRefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitTargetRef(s)
	}
}

func (s *TargetRefContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitTargetRef(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) TargetRef() (localctx ITargetRefContext) {
	localctx = NewTargetRefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, UghParserRULE_targetRef)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(124)
		_la = p.GetTokenStream().LA(1)

		if !(_la == UghParserHASH_NUMBER || _la == UghParserIDENT) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperationContext is an interface to support dynamic dispatch.
type IOperationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SetOp() ISetOpContext
	AddFieldOp() IAddFieldOpContext
	RemoveFieldOp() IRemoveFieldOpContext
	ClearOp() IClearOpContext
	TagOp() ITagOpContext

	// IsOperationContext differentiates from other interfaces.
	IsOperationContext()
}

type OperationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperationContext() *OperationContext {
	var p = new(OperationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_operation
	return p
}

func InitEmptyOperationContext(p *OperationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_operation
}

func (*OperationContext) IsOperationContext() {}

func NewOperationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperationContext {
	var p = new(OperationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_operation

	return p
}

func (s *OperationContext) GetParser() antlr.Parser { return s.parser }

func (s *OperationContext) SetOp() ISetOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetOpContext)
}

func (s *OperationContext) AddFieldOp() IAddFieldOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAddFieldOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAddFieldOpContext)
}

func (s *OperationContext) RemoveFieldOp() IRemoveFieldOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRemoveFieldOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRemoveFieldOpContext)
}

func (s *OperationContext) ClearOp() IClearOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IClearOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IClearOpContext)
}

func (s *OperationContext) TagOp() ITagOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITagOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITagOpContext)
}

func (s *OperationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterOperation(s)
	}
}

func (s *OperationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitOperation(s)
	}
}

func (s *OperationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitOperation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) Operation() (localctx IOperationContext) {
	localctx = NewOperationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, UghParserRULE_operation)
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserSET_FIELD:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(126)
			p.SetOp()
		}

	case UghParserADD_FIELD:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(127)
			p.AddFieldOp()
		}

	case UghParserREMOVE_FIELD:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(128)
			p.RemoveFieldOp()
		}

	case UghParserCLEAR_FIELD, UghParserCLEAR_OP:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(129)
			p.ClearOp()
		}

	case UghParserPROJECT_TAG, UghParserCONTEXT_TAG:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(130)
			p.TagOp()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterCommandContext is an interface to support dynamic dispatch.
type IFilterCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FilterVerb() IFilterVerbContext
	FilterOrExpr() IFilterOrExprContext

	// IsFilterCommandContext differentiates from other interfaces.
	IsFilterCommandContext()
}

type FilterCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterCommandContext() *FilterCommandContext {
	var p = new(FilterCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterCommand
	return p
}

func InitEmptyFilterCommandContext(p *FilterCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterCommand
}

func (*FilterCommandContext) IsFilterCommandContext() {}

func NewFilterCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterCommandContext {
	var p = new(FilterCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterCommand

	return p
}

func (s *FilterCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterCommandContext) FilterVerb() IFilterVerbContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterVerbContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterVerbContext)
}

func (s *FilterCommandContext) FilterOrExpr() IFilterOrExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterOrExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterOrExprContext)
}

func (s *FilterCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterCommand(s)
	}
}

func (s *FilterCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterCommand(s)
	}
}

func (s *FilterCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterCommand() (localctx IFilterCommandContext) {
	localctx = NewFilterCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, UghParserRULE_filterCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(133)
		p.FilterVerb()
	}
	{
		p.SetState(134)
		p.FilterOrExpr()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterVerbContext is an interface to support dynamic dispatch.
type IFilterVerbContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_FIND() antlr.TerminalNode
	KW_SHOW() antlr.TerminalNode
	KW_LIST() antlr.TerminalNode
	KW_FILTER() antlr.TerminalNode

	// IsFilterVerbContext differentiates from other interfaces.
	IsFilterVerbContext()
}

type FilterVerbContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterVerbContext() *FilterVerbContext {
	var p = new(FilterVerbContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterVerb
	return p
}

func InitEmptyFilterVerbContext(p *FilterVerbContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterVerb
}

func (*FilterVerbContext) IsFilterVerbContext() {}

func NewFilterVerbContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterVerbContext {
	var p = new(FilterVerbContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterVerb

	return p
}

func (s *FilterVerbContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterVerbContext) KW_FIND() antlr.TerminalNode {
	return s.GetToken(UghParserKW_FIND, 0)
}

func (s *FilterVerbContext) KW_SHOW() antlr.TerminalNode {
	return s.GetToken(UghParserKW_SHOW, 0)
}

func (s *FilterVerbContext) KW_LIST() antlr.TerminalNode {
	return s.GetToken(UghParserKW_LIST, 0)
}

func (s *FilterVerbContext) KW_FILTER() antlr.TerminalNode {
	return s.GetToken(UghParserKW_FILTER, 0)
}

func (s *FilterVerbContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterVerbContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterVerbContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterVerb(s)
	}
}

func (s *FilterVerbContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterVerb(s)
	}
}

func (s *FilterVerbContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterVerb(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterVerb() (localctx IFilterVerbContext) {
	localctx = NewFilterVerbContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, UghParserRULE_filterVerb)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(136)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1966080) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterOrExprContext is an interface to support dynamic dispatch.
type IFilterOrExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFilterAndExpr() []IFilterAndExprContext
	FilterAndExpr(i int) IFilterAndExprContext
	AllOrOp() []IOrOpContext
	OrOp(i int) IOrOpContext

	// IsFilterOrExprContext differentiates from other interfaces.
	IsFilterOrExprContext()
}

type FilterOrExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterOrExprContext() *FilterOrExprContext {
	var p = new(FilterOrExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterOrExpr
	return p
}

func InitEmptyFilterOrExprContext(p *FilterOrExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterOrExpr
}

func (*FilterOrExprContext) IsFilterOrExprContext() {}

func NewFilterOrExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterOrExprContext {
	var p = new(FilterOrExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterOrExpr

	return p
}

func (s *FilterOrExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterOrExprContext) AllFilterAndExpr() []IFilterAndExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterAndExprContext); ok {
			len++
		}
	}

	tst := make([]IFilterAndExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterAndExprContext); ok {
			tst[i] = t.(IFilterAndExprContext)
			i++
		}
	}

	return tst
}

func (s *FilterOrExprContext) FilterAndExpr(i int) IFilterAndExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterAndExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterAndExprContext)
}

func (s *FilterOrExprContext) AllOrOp() []IOrOpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOrOpContext); ok {
			len++
		}
	}

	tst := make([]IOrOpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOrOpContext); ok {
			tst[i] = t.(IOrOpContext)
			i++
		}
	}

	return tst
}

func (s *FilterOrExprContext) OrOp(i int) IOrOpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrOpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrOpContext)
}

func (s *FilterOrExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterOrExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterOrExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterOrExpr(s)
	}
}

func (s *FilterOrExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterOrExpr(s)
	}
}

func (s *FilterOrExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterOrExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterOrExpr() (localctx IFilterOrExprContext) {
	localctx = NewFilterOrExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, UghParserRULE_filterOrExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(138)
		p.FilterAndExpr()
	}
	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == UghParserKW_OR || _la == UghParserOR_OP {
		{
			p.SetState(139)
			p.OrOp()
		}
		{
			p.SetState(140)
			p.FilterAndExpr()
		}

		p.SetState(146)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterAndExprContext is an interface to support dynamic dispatch.
type IFilterAndExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllFilterNotExpr() []IFilterNotExprContext
	FilterNotExpr(i int) IFilterNotExprContext
	AllAndOp() []IAndOpContext
	AndOp(i int) IAndOpContext

	// IsFilterAndExprContext differentiates from other interfaces.
	IsFilterAndExprContext()
}

type FilterAndExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterAndExprContext() *FilterAndExprContext {
	var p = new(FilterAndExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterAndExpr
	return p
}

func InitEmptyFilterAndExprContext(p *FilterAndExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterAndExpr
}

func (*FilterAndExprContext) IsFilterAndExprContext() {}

func NewFilterAndExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterAndExprContext {
	var p = new(FilterAndExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterAndExpr

	return p
}

func (s *FilterAndExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterAndExprContext) AllFilterNotExpr() []IFilterNotExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterNotExprContext); ok {
			len++
		}
	}

	tst := make([]IFilterNotExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterNotExprContext); ok {
			tst[i] = t.(IFilterNotExprContext)
			i++
		}
	}

	return tst
}

func (s *FilterAndExprContext) FilterNotExpr(i int) IFilterNotExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterNotExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterNotExprContext)
}

func (s *FilterAndExprContext) AllAndOp() []IAndOpContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAndOpContext); ok {
			len++
		}
	}

	tst := make([]IAndOpContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAndOpContext); ok {
			tst[i] = t.(IAndOpContext)
			i++
		}
	}

	return tst
}

func (s *FilterAndExprContext) AndOp(i int) IAndOpContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAndOpContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAndOpContext)
}

func (s *FilterAndExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterAndExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterAndExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterAndExpr(s)
	}
}

func (s *FilterAndExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterAndExpr(s)
	}
}

func (s *FilterAndExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterAndExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterAndExpr() (localctx IFilterAndExprContext) {
	localctx = NewFilterAndExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, UghParserRULE_filterAndExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(147)
		p.FilterNotExpr()
	}
	p.SetState(153)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == UghParserKW_AND || _la == UghParserAND_OP {
		{
			p.SetState(148)
			p.AndOp()
		}
		{
			p.SetState(149)
			p.FilterNotExpr()
		}

		p.SetState(155)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterNotExprContext is an interface to support dynamic dispatch.
type IFilterNotExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NotOp() INotOpContext
	FilterAtom() IFilterAtomContext

	// IsFilterNotExprContext differentiates from other interfaces.
	IsFilterNotExprContext()
}

type FilterNotExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterNotExprContext() *FilterNotExprContext {
	var p = new(FilterNotExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterNotExpr
	return p
}

func InitEmptyFilterNotExprContext(p *FilterNotExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterNotExpr
}

func (*FilterNotExprContext) IsFilterNotExprContext() {}

func NewFilterNotExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterNotExprContext {
	var p = new(FilterNotExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterNotExpr

	return p
}

func (s *FilterNotExprContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterNotExprContext) NotOp() INotOpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INotOpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INotOpContext)
}

func (s *FilterNotExprContext) FilterAtom() IFilterAtomContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterAtomContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterAtomContext)
}

func (s *FilterNotExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterNotExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterNotExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterNotExpr(s)
	}
}

func (s *FilterNotExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterNotExpr(s)
	}
}

func (s *FilterNotExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterNotExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterNotExpr() (localctx IFilterNotExprContext) {
	localctx = NewFilterNotExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, UghParserRULE_filterNotExpr)
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(156)
			p.NotOp()
		}
		{
			p.SetState(157)
			p.FilterAtom()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(159)
			p.FilterAtom()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterAtomContext is an interface to support dynamic dispatch.
type IFilterAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	FilterOrExpr() IFilterOrExprContext
	RPAREN() antlr.TerminalNode
	FilterPredicate() IFilterPredicateContext

	// IsFilterAtomContext differentiates from other interfaces.
	IsFilterAtomContext()
}

type FilterAtomContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterAtomContext() *FilterAtomContext {
	var p = new(FilterAtomContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterAtom
	return p
}

func InitEmptyFilterAtomContext(p *FilterAtomContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterAtom
}

func (*FilterAtomContext) IsFilterAtomContext() {}

func NewFilterAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterAtomContext {
	var p = new(FilterAtomContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterAtom

	return p
}

func (s *FilterAtomContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterAtomContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(UghParserLPAREN, 0)
}

func (s *FilterAtomContext) FilterOrExpr() IFilterOrExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterOrExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterOrExprContext)
}

func (s *FilterAtomContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(UghParserRPAREN, 0)
}

func (s *FilterAtomContext) FilterPredicate() IFilterPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterPredicateContext)
}

func (s *FilterAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterAtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterAtom(s)
	}
}

func (s *FilterAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterAtom(s)
	}
}

func (s *FilterAtomContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterAtom(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterAtom() (localctx IFilterAtomContext) {
	localctx = NewFilterAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, UghParserRULE_filterAtom)
	p.SetState(167)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserLPAREN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(162)
			p.Match(UghParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(163)
			p.FilterOrExpr()
		}
		{
			p.SetState(164)
			p.Match(UghParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserQUOTED, UghParserHASH_NUMBER, UghParserPROJECT_TAG, UghParserCONTEXT_TAG, UghParserSET_FIELD, UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserCOLON, UghParserCOMMA, UghParserIDENT:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(166)
			p.FilterPredicate()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterPredicateContext is an interface to support dynamic dispatch.
type IFilterPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FilterFieldPredicate() IFilterFieldPredicateContext
	FilterTagPredicate() IFilterTagPredicateContext
	FilterTextPredicate() IFilterTextPredicateContext

	// IsFilterPredicateContext differentiates from other interfaces.
	IsFilterPredicateContext()
}

type FilterPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterPredicateContext() *FilterPredicateContext {
	var p = new(FilterPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterPredicate
	return p
}

func InitEmptyFilterPredicateContext(p *FilterPredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterPredicate
}

func (*FilterPredicateContext) IsFilterPredicateContext() {}

func NewFilterPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterPredicateContext {
	var p = new(FilterPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterPredicate

	return p
}

func (s *FilterPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterPredicateContext) FilterFieldPredicate() IFilterFieldPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterFieldPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterFieldPredicateContext)
}

func (s *FilterPredicateContext) FilterTagPredicate() IFilterTagPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterTagPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterTagPredicateContext)
}

func (s *FilterPredicateContext) FilterTextPredicate() IFilterTextPredicateContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterTextPredicateContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterTextPredicateContext)
}

func (s *FilterPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterPredicate(s)
	}
}

func (s *FilterPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterPredicate(s)
	}
}

func (s *FilterPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterPredicate() (localctx IFilterPredicateContext) {
	localctx = NewFilterPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, UghParserRULE_filterPredicate)
	p.SetState(172)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserSET_FIELD:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(169)
			p.FilterFieldPredicate()
		}

	case UghParserPROJECT_TAG, UghParserCONTEXT_TAG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(170)
			p.FilterTagPredicate()
		}

	case UghParserQUOTED, UghParserHASH_NUMBER, UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserCOLON, UghParserCOMMA, UghParserIDENT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(171)
			p.FilterTextPredicate()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterFieldPredicateContext is an interface to support dynamic dispatch.
type IFilterFieldPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SET_FIELD() antlr.TerminalNode
	FilterValue() IFilterValueContext

	// IsFilterFieldPredicateContext differentiates from other interfaces.
	IsFilterFieldPredicateContext()
}

type FilterFieldPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterFieldPredicateContext() *FilterFieldPredicateContext {
	var p = new(FilterFieldPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterFieldPredicate
	return p
}

func InitEmptyFilterFieldPredicateContext(p *FilterFieldPredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterFieldPredicate
}

func (*FilterFieldPredicateContext) IsFilterFieldPredicateContext() {}

func NewFilterFieldPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterFieldPredicateContext {
	var p = new(FilterFieldPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterFieldPredicate

	return p
}

func (s *FilterFieldPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterFieldPredicateContext) SET_FIELD() antlr.TerminalNode {
	return s.GetToken(UghParserSET_FIELD, 0)
}

func (s *FilterFieldPredicateContext) FilterValue() IFilterValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterValueContext)
}

func (s *FilterFieldPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterFieldPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterFieldPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterFieldPredicate(s)
	}
}

func (s *FilterFieldPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterFieldPredicate(s)
	}
}

func (s *FilterFieldPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterFieldPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterFieldPredicate() (localctx IFilterFieldPredicateContext) {
	localctx = NewFilterFieldPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, UghParserRULE_filterFieldPredicate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(174)
		p.Match(UghParserSET_FIELD)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(175)
		p.FilterValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterTagPredicateContext is an interface to support dynamic dispatch.
type IFilterTagPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PROJECT_TAG() antlr.TerminalNode
	CONTEXT_TAG() antlr.TerminalNode

	// IsFilterTagPredicateContext differentiates from other interfaces.
	IsFilterTagPredicateContext()
}

type FilterTagPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterTagPredicateContext() *FilterTagPredicateContext {
	var p = new(FilterTagPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterTagPredicate
	return p
}

func InitEmptyFilterTagPredicateContext(p *FilterTagPredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterTagPredicate
}

func (*FilterTagPredicateContext) IsFilterTagPredicateContext() {}

func NewFilterTagPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterTagPredicateContext {
	var p = new(FilterTagPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterTagPredicate

	return p
}

func (s *FilterTagPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterTagPredicateContext) PROJECT_TAG() antlr.TerminalNode {
	return s.GetToken(UghParserPROJECT_TAG, 0)
}

func (s *FilterTagPredicateContext) CONTEXT_TAG() antlr.TerminalNode {
	return s.GetToken(UghParserCONTEXT_TAG, 0)
}

func (s *FilterTagPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterTagPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterTagPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterTagPredicate(s)
	}
}

func (s *FilterTagPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterTagPredicate(s)
	}
}

func (s *FilterTagPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterTagPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterTagPredicate() (localctx IFilterTagPredicateContext) {
	localctx = NewFilterTagPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, UghParserRULE_filterTagPredicate)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(177)
		_la = p.GetTokenStream().LA(1)

		if !(_la == UghParserPROJECT_TAG || _la == UghParserCONTEXT_TAG) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterTextPredicateContext is an interface to support dynamic dispatch.
type IFilterTextPredicateContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FilterValue() IFilterValueContext

	// IsFilterTextPredicateContext differentiates from other interfaces.
	IsFilterTextPredicateContext()
}

type FilterTextPredicateContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterTextPredicateContext() *FilterTextPredicateContext {
	var p = new(FilterTextPredicateContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterTextPredicate
	return p
}

func InitEmptyFilterTextPredicateContext(p *FilterTextPredicateContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterTextPredicate
}

func (*FilterTextPredicateContext) IsFilterTextPredicateContext() {}

func NewFilterTextPredicateContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterTextPredicateContext {
	var p = new(FilterTextPredicateContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterTextPredicate

	return p
}

func (s *FilterTextPredicateContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterTextPredicateContext) FilterValue() IFilterValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterValueContext)
}

func (s *FilterTextPredicateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterTextPredicateContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterTextPredicateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterTextPredicate(s)
	}
}

func (s *FilterTextPredicateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterTextPredicate(s)
	}
}

func (s *FilterTextPredicateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterTextPredicate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterTextPredicate() (localctx IFilterTextPredicateContext) {
	localctx = NewFilterTextPredicateContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, UghParserRULE_filterTextPredicate)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(179)
		p.FilterValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterValueContext is an interface to support dynamic dispatch.
type IFilterValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QUOTED() antlr.TerminalNode
	AllFilterValueWord() []IFilterValueWordContext
	FilterValueWord(i int) IFilterValueWordContext

	// IsFilterValueContext differentiates from other interfaces.
	IsFilterValueContext()
}

type FilterValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterValueContext() *FilterValueContext {
	var p = new(FilterValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterValue
	return p
}

func InitEmptyFilterValueContext(p *FilterValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterValue
}

func (*FilterValueContext) IsFilterValueContext() {}

func NewFilterValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterValueContext {
	var p = new(FilterValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterValue

	return p
}

func (s *FilterValueContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterValueContext) QUOTED() antlr.TerminalNode {
	return s.GetToken(UghParserQUOTED, 0)
}

func (s *FilterValueContext) AllFilterValueWord() []IFilterValueWordContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFilterValueWordContext); ok {
			len++
		}
	}

	tst := make([]IFilterValueWordContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFilterValueWordContext); ok {
			tst[i] = t.(IFilterValueWordContext)
			i++
		}
	}

	return tst
}

func (s *FilterValueContext) FilterValueWord(i int) IFilterValueWordContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFilterValueWordContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFilterValueWordContext)
}

func (s *FilterValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterValue(s)
	}
}

func (s *FilterValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterValue(s)
	}
}

func (s *FilterValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterValue() (localctx IFilterValueContext) {
	localctx = NewFilterValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, UghParserRULE_filterValue)
	var _alt int

	p.SetState(187)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserQUOTED:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(181)
			p.Match(UghParserQUOTED)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserHASH_NUMBER, UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserCOLON, UghParserCOMMA, UghParserIDENT:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(183)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(182)
					p.FilterValueWord()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(185)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 13, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFilterValueWordContext is an interface to support dynamic dispatch.
type IFilterValueWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ident() IIdentContext
	HASH_NUMBER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	COMMA() antlr.TerminalNode

	// IsFilterValueWordContext differentiates from other interfaces.
	IsFilterValueWordContext()
}

type FilterValueWordContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFilterValueWordContext() *FilterValueWordContext {
	var p = new(FilterValueWordContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterValueWord
	return p
}

func InitEmptyFilterValueWordContext(p *FilterValueWordContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_filterValueWord
}

func (*FilterValueWordContext) IsFilterValueWordContext() {}

func NewFilterValueWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FilterValueWordContext {
	var p = new(FilterValueWordContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_filterValueWord

	return p
}

func (s *FilterValueWordContext) GetParser() antlr.Parser { return s.parser }

func (s *FilterValueWordContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *FilterValueWordContext) HASH_NUMBER() antlr.TerminalNode {
	return s.GetToken(UghParserHASH_NUMBER, 0)
}

func (s *FilterValueWordContext) COLON() antlr.TerminalNode {
	return s.GetToken(UghParserCOLON, 0)
}

func (s *FilterValueWordContext) COMMA() antlr.TerminalNode {
	return s.GetToken(UghParserCOMMA, 0)
}

func (s *FilterValueWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FilterValueWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FilterValueWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterFilterValueWord(s)
	}
}

func (s *FilterValueWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitFilterValueWord(s)
	}
}

func (s *FilterValueWordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitFilterValueWord(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) FilterValueWord() (localctx IFilterValueWordContext) {
	localctx = NewFilterValueWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, UghParserRULE_filterValueWord)
	p.SetState(193)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserIDENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(189)
			p.Ident()
		}

	case UghParserHASH_NUMBER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(190)
			p.Match(UghParserHASH_NUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCOLON:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(191)
			p.Match(UghParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCOMMA:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(192)
			p.Match(UghParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOrOpContext is an interface to support dynamic dispatch.
type IOrOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OR_OP() antlr.TerminalNode
	KW_OR() antlr.TerminalNode

	// IsOrOpContext differentiates from other interfaces.
	IsOrOpContext()
}

type OrOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrOpContext() *OrOpContext {
	var p = new(OrOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_orOp
	return p
}

func InitEmptyOrOpContext(p *OrOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_orOp
}

func (*OrOpContext) IsOrOpContext() {}

func NewOrOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrOpContext {
	var p = new(OrOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_orOp

	return p
}

func (s *OrOpContext) GetParser() antlr.Parser { return s.parser }

func (s *OrOpContext) OR_OP() antlr.TerminalNode {
	return s.GetToken(UghParserOR_OP, 0)
}

func (s *OrOpContext) KW_OR() antlr.TerminalNode {
	return s.GetToken(UghParserKW_OR, 0)
}

func (s *OrOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterOrOp(s)
	}
}

func (s *OrOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitOrOp(s)
	}
}

func (s *OrOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitOrOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) OrOp() (localctx IOrOpContext) {
	localctx = NewOrOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, UghParserRULE_orOp)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(195)
		_la = p.GetTokenStream().LA(1)

		if !(_la == UghParserKW_OR || _la == UghParserOR_OP) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAndOpContext is an interface to support dynamic dispatch.
type IAndOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AND_OP() antlr.TerminalNode
	KW_AND() antlr.TerminalNode

	// IsAndOpContext differentiates from other interfaces.
	IsAndOpContext()
}

type AndOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAndOpContext() *AndOpContext {
	var p = new(AndOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_andOp
	return p
}

func InitEmptyAndOpContext(p *AndOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_andOp
}

func (*AndOpContext) IsAndOpContext() {}

func NewAndOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AndOpContext {
	var p = new(AndOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_andOp

	return p
}

func (s *AndOpContext) GetParser() antlr.Parser { return s.parser }

func (s *AndOpContext) AND_OP() antlr.TerminalNode {
	return s.GetToken(UghParserAND_OP, 0)
}

func (s *AndOpContext) KW_AND() antlr.TerminalNode {
	return s.GetToken(UghParserKW_AND, 0)
}

func (s *AndOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AndOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterAndOp(s)
	}
}

func (s *AndOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitAndOp(s)
	}
}

func (s *AndOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitAndOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) AndOp() (localctx IAndOpContext) {
	localctx = NewAndOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, UghParserRULE_andOp)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(197)
		_la = p.GetTokenStream().LA(1)

		if !(_la == UghParserKW_AND || _la == UghParserAND_OP) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INotOpContext is an interface to support dynamic dispatch.
type INotOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLEAR_OP() antlr.TerminalNode
	KW_NOT() antlr.TerminalNode

	// IsNotOpContext differentiates from other interfaces.
	IsNotOpContext()
}

type NotOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNotOpContext() *NotOpContext {
	var p = new(NotOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_notOp
	return p
}

func InitEmptyNotOpContext(p *NotOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_notOp
}

func (*NotOpContext) IsNotOpContext() {}

func NewNotOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NotOpContext {
	var p = new(NotOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_notOp

	return p
}

func (s *NotOpContext) GetParser() antlr.Parser { return s.parser }

func (s *NotOpContext) CLEAR_OP() antlr.TerminalNode {
	return s.GetToken(UghParserCLEAR_OP, 0)
}

func (s *NotOpContext) KW_NOT() antlr.TerminalNode {
	return s.GetToken(UghParserKW_NOT, 0)
}

func (s *NotOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NotOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NotOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterNotOp(s)
	}
}

func (s *NotOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitNotOp(s)
	}
}

func (s *NotOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitNotOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) NotOp() (localctx INotOpContext) {
	localctx = NewNotOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, UghParserRULE_notOp)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(199)
		_la = p.GetTokenStream().LA(1)

		if !(_la == UghParserKW_NOT || _la == UghParserCLEAR_OP) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IViewCommandContext is an interface to support dynamic dispatch.
type IViewCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_VIEW() antlr.TerminalNode
	ViewTarget() IViewTargetContext

	// IsViewCommandContext differentiates from other interfaces.
	IsViewCommandContext()
}

type ViewCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyViewCommandContext() *ViewCommandContext {
	var p = new(ViewCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_viewCommand
	return p
}

func InitEmptyViewCommandContext(p *ViewCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_viewCommand
}

func (*ViewCommandContext) IsViewCommandContext() {}

func NewViewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ViewCommandContext {
	var p = new(ViewCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_viewCommand

	return p
}

func (s *ViewCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ViewCommandContext) KW_VIEW() antlr.TerminalNode {
	return s.GetToken(UghParserKW_VIEW, 0)
}

func (s *ViewCommandContext) ViewTarget() IViewTargetContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IViewTargetContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IViewTargetContext)
}

func (s *ViewCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ViewCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ViewCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterViewCommand(s)
	}
}

func (s *ViewCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitViewCommand(s)
	}
}

func (s *ViewCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitViewCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) ViewCommand() (localctx IViewCommandContext) {
	localctx = NewViewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, UghParserRULE_viewCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(201)
		p.Match(UghParserKW_VIEW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&34426845184) != 0 {
		{
			p.SetState(202)
			p.ViewTarget()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IViewTargetContext is an interface to support dynamic dispatch.
type IViewTargetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ident() IIdentContext

	// IsViewTargetContext differentiates from other interfaces.
	IsViewTargetContext()
}

type ViewTargetContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyViewTargetContext() *ViewTargetContext {
	var p = new(ViewTargetContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_viewTarget
	return p
}

func InitEmptyViewTargetContext(p *ViewTargetContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_viewTarget
}

func (*ViewTargetContext) IsViewTargetContext() {}

func NewViewTargetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ViewTargetContext {
	var p = new(ViewTargetContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_viewTarget

	return p
}

func (s *ViewTargetContext) GetParser() antlr.Parser { return s.parser }

func (s *ViewTargetContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ViewTargetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ViewTargetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ViewTargetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterViewTarget(s)
	}
}

func (s *ViewTargetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitViewTarget(s)
	}
}

func (s *ViewTargetContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitViewTarget(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) ViewTarget() (localctx IViewTargetContext) {
	localctx = NewViewTargetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, UghParserRULE_viewTarget)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(205)
		p.Ident()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IContextCommandContext is an interface to support dynamic dispatch.
type IContextCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	KW_CONTEXT() antlr.TerminalNode
	ContextArg() IContextArgContext

	// IsContextCommandContext differentiates from other interfaces.
	IsContextCommandContext()
}

type ContextCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContextCommandContext() *ContextCommandContext {
	var p = new(ContextCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_contextCommand
	return p
}

func InitEmptyContextCommandContext(p *ContextCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_contextCommand
}

func (*ContextCommandContext) IsContextCommandContext() {}

func NewContextCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContextCommandContext {
	var p = new(ContextCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_contextCommand

	return p
}

func (s *ContextCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *ContextCommandContext) KW_CONTEXT() antlr.TerminalNode {
	return s.GetToken(UghParserKW_CONTEXT, 0)
}

func (s *ContextCommandContext) ContextArg() IContextArgContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IContextArgContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IContextArgContext)
}

func (s *ContextCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContextCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ContextCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterContextCommand(s)
	}
}

func (s *ContextCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitContextCommand(s)
	}
}

func (s *ContextCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitContextCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) ContextCommand() (localctx IContextCommandContext) {
	localctx = NewContextCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, UghParserRULE_contextCommand)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(207)
		p.Match(UghParserKW_CONTEXT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(209)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&34426845208) != 0 {
		{
			p.SetState(208)
			p.ContextArg()
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IContextArgContext is an interface to support dynamic dispatch.
type IContextArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PROJECT_TAG() antlr.TerminalNode
	CONTEXT_TAG() antlr.TerminalNode
	Ident() IIdentContext

	// IsContextArgContext differentiates from other interfaces.
	IsContextArgContext()
}

type ContextArgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyContextArgContext() *ContextArgContext {
	var p = new(ContextArgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_contextArg
	return p
}

func InitEmptyContextArgContext(p *ContextArgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_contextArg
}

func (*ContextArgContext) IsContextArgContext() {}

func NewContextArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContextArgContext {
	var p = new(ContextArgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_contextArg

	return p
}

func (s *ContextArgContext) GetParser() antlr.Parser { return s.parser }

func (s *ContextArgContext) PROJECT_TAG() antlr.TerminalNode {
	return s.GetToken(UghParserPROJECT_TAG, 0)
}

func (s *ContextArgContext) CONTEXT_TAG() antlr.TerminalNode {
	return s.GetToken(UghParserCONTEXT_TAG, 0)
}

func (s *ContextArgContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ContextArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ContextArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ContextArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterContextArg(s)
	}
}

func (s *ContextArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitContextArg(s)
	}
}

func (s *ContextArgContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitContextArg(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) ContextArg() (localctx IContextArgContext) {
	localctx = NewContextArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, UghParserRULE_contextArg)
	p.SetState(214)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserPROJECT_TAG:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(211)
			p.Match(UghParserPROJECT_TAG)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCONTEXT_TAG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(212)
			p.Match(UghParserCONTEXT_TAG)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserIDENT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(213)
			p.Ident()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISetOpContext is an interface to support dynamic dispatch.
type ISetOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SET_FIELD() antlr.TerminalNode
	OpValue() IOpValueContext

	// IsSetOpContext differentiates from other interfaces.
	IsSetOpContext()
}

type SetOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetOpContext() *SetOpContext {
	var p = new(SetOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_setOp
	return p
}

func InitEmptySetOpContext(p *SetOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_setOp
}

func (*SetOpContext) IsSetOpContext() {}

func NewSetOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetOpContext {
	var p = new(SetOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_setOp

	return p
}

func (s *SetOpContext) GetParser() antlr.Parser { return s.parser }

func (s *SetOpContext) SET_FIELD() antlr.TerminalNode {
	return s.GetToken(UghParserSET_FIELD, 0)
}

func (s *SetOpContext) OpValue() IOpValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOpValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOpValueContext)
}

func (s *SetOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterSetOp(s)
	}
}

func (s *SetOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitSetOp(s)
	}
}

func (s *SetOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitSetOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) SetOp() (localctx ISetOpContext) {
	localctx = NewSetOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, UghParserRULE_setOp)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(216)
		p.Match(UghParserSET_FIELD)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(217)
		p.OpValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAddFieldOpContext is an interface to support dynamic dispatch.
type IAddFieldOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ADD_FIELD() antlr.TerminalNode
	OpValue() IOpValueContext

	// IsAddFieldOpContext differentiates from other interfaces.
	IsAddFieldOpContext()
}

type AddFieldOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAddFieldOpContext() *AddFieldOpContext {
	var p = new(AddFieldOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_addFieldOp
	return p
}

func InitEmptyAddFieldOpContext(p *AddFieldOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_addFieldOp
}

func (*AddFieldOpContext) IsAddFieldOpContext() {}

func NewAddFieldOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AddFieldOpContext {
	var p = new(AddFieldOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_addFieldOp

	return p
}

func (s *AddFieldOpContext) GetParser() antlr.Parser { return s.parser }

func (s *AddFieldOpContext) ADD_FIELD() antlr.TerminalNode {
	return s.GetToken(UghParserADD_FIELD, 0)
}

func (s *AddFieldOpContext) OpValue() IOpValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOpValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOpValueContext)
}

func (s *AddFieldOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AddFieldOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AddFieldOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterAddFieldOp(s)
	}
}

func (s *AddFieldOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitAddFieldOp(s)
	}
}

func (s *AddFieldOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitAddFieldOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) AddFieldOp() (localctx IAddFieldOpContext) {
	localctx = NewAddFieldOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, UghParserRULE_addFieldOp)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(219)
		p.Match(UghParserADD_FIELD)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(220)
		p.OpValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IRemoveFieldOpContext is an interface to support dynamic dispatch.
type IRemoveFieldOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REMOVE_FIELD() antlr.TerminalNode
	OpValue() IOpValueContext

	// IsRemoveFieldOpContext differentiates from other interfaces.
	IsRemoveFieldOpContext()
}

type RemoveFieldOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRemoveFieldOpContext() *RemoveFieldOpContext {
	var p = new(RemoveFieldOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_removeFieldOp
	return p
}

func InitEmptyRemoveFieldOpContext(p *RemoveFieldOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_removeFieldOp
}

func (*RemoveFieldOpContext) IsRemoveFieldOpContext() {}

func NewRemoveFieldOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RemoveFieldOpContext {
	var p = new(RemoveFieldOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_removeFieldOp

	return p
}

func (s *RemoveFieldOpContext) GetParser() antlr.Parser { return s.parser }

func (s *RemoveFieldOpContext) REMOVE_FIELD() antlr.TerminalNode {
	return s.GetToken(UghParserREMOVE_FIELD, 0)
}

func (s *RemoveFieldOpContext) OpValue() IOpValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOpValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOpValueContext)
}

func (s *RemoveFieldOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RemoveFieldOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RemoveFieldOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterRemoveFieldOp(s)
	}
}

func (s *RemoveFieldOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitRemoveFieldOp(s)
	}
}

func (s *RemoveFieldOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitRemoveFieldOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) RemoveFieldOp() (localctx IRemoveFieldOpContext) {
	localctx = NewRemoveFieldOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, UghParserRULE_removeFieldOp)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(222)
		p.Match(UghParserREMOVE_FIELD)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(223)
		p.OpValue()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IClearOpContext is an interface to support dynamic dispatch.
type IClearOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLEAR_FIELD() antlr.TerminalNode
	CLEAR_OP() antlr.TerminalNode
	Ident() IIdentContext

	// IsClearOpContext differentiates from other interfaces.
	IsClearOpContext()
}

type ClearOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyClearOpContext() *ClearOpContext {
	var p = new(ClearOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_clearOp
	return p
}

func InitEmptyClearOpContext(p *ClearOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_clearOp
}

func (*ClearOpContext) IsClearOpContext() {}

func NewClearOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ClearOpContext {
	var p = new(ClearOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_clearOp

	return p
}

func (s *ClearOpContext) GetParser() antlr.Parser { return s.parser }

func (s *ClearOpContext) CLEAR_FIELD() antlr.TerminalNode {
	return s.GetToken(UghParserCLEAR_FIELD, 0)
}

func (s *ClearOpContext) CLEAR_OP() antlr.TerminalNode {
	return s.GetToken(UghParserCLEAR_OP, 0)
}

func (s *ClearOpContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ClearOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ClearOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ClearOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterClearOp(s)
	}
}

func (s *ClearOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitClearOp(s)
	}
}

func (s *ClearOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitClearOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) ClearOp() (localctx IClearOpContext) {
	localctx = NewClearOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, UghParserRULE_clearOp)
	p.SetState(228)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserCLEAR_FIELD:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(225)
			p.Match(UghParserCLEAR_FIELD)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCLEAR_OP:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(226)
			p.Match(UghParserCLEAR_OP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(227)
			p.Ident()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITagOpContext is an interface to support dynamic dispatch.
type ITagOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PROJECT_TAG() antlr.TerminalNode
	CONTEXT_TAG() antlr.TerminalNode

	// IsTagOpContext differentiates from other interfaces.
	IsTagOpContext()
}

type TagOpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagOpContext() *TagOpContext {
	var p = new(TagOpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_tagOp
	return p
}

func InitEmptyTagOpContext(p *TagOpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_tagOp
}

func (*TagOpContext) IsTagOpContext() {}

func NewTagOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagOpContext {
	var p = new(TagOpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_tagOp

	return p
}

func (s *TagOpContext) GetParser() antlr.Parser { return s.parser }

func (s *TagOpContext) PROJECT_TAG() antlr.TerminalNode {
	return s.GetToken(UghParserPROJECT_TAG, 0)
}

func (s *TagOpContext) CONTEXT_TAG() antlr.TerminalNode {
	return s.GetToken(UghParserCONTEXT_TAG, 0)
}

func (s *TagOpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagOpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterTagOp(s)
	}
}

func (s *TagOpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitTagOp(s)
	}
}

func (s *TagOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitTagOp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) TagOp() (localctx ITagOpContext) {
	localctx = NewTagOpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, UghParserRULE_tagOp)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(230)
		_la = p.GetTokenStream().LA(1)

		if !(_la == UghParserPROJECT_TAG || _la == UghParserCONTEXT_TAG) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOpValueContext is an interface to support dynamic dispatch.
type IOpValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QUOTED() antlr.TerminalNode
	AllOpValueWord() []IOpValueWordContext
	OpValueWord(i int) IOpValueWordContext

	// IsOpValueContext differentiates from other interfaces.
	IsOpValueContext()
}

type OpValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOpValueContext() *OpValueContext {
	var p = new(OpValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_opValue
	return p
}

func InitEmptyOpValueContext(p *OpValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_opValue
}

func (*OpValueContext) IsOpValueContext() {}

func NewOpValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OpValueContext {
	var p = new(OpValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_opValue

	return p
}

func (s *OpValueContext) GetParser() antlr.Parser { return s.parser }

func (s *OpValueContext) QUOTED() antlr.TerminalNode {
	return s.GetToken(UghParserQUOTED, 0)
}

func (s *OpValueContext) AllOpValueWord() []IOpValueWordContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOpValueWordContext); ok {
			len++
		}
	}

	tst := make([]IOpValueWordContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOpValueWordContext); ok {
			tst[i] = t.(IOpValueWordContext)
			i++
		}
	}

	return tst
}

func (s *OpValueContext) OpValueWord(i int) IOpValueWordContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOpValueWordContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOpValueWordContext)
}

func (s *OpValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OpValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OpValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterOpValue(s)
	}
}

func (s *OpValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitOpValue(s)
	}
}

func (s *OpValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitOpValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) OpValue() (localctx IOpValueContext) {
	localctx = NewOpValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, UghParserRULE_opValue)
	var _alt int

	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserQUOTED:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(232)
			p.Match(UghParserQUOTED)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserHASH_NUMBER, UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserCOLON, UghParserCOMMA, UghParserIDENT:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(234)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(233)
					p.OpValueWord()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(236)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOpValueWordContext is an interface to support dynamic dispatch.
type IOpValueWordContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ident() IIdentContext
	HASH_NUMBER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	COMMA() antlr.TerminalNode

	// IsOpValueWordContext differentiates from other interfaces.
	IsOpValueWordContext()
}

type OpValueWordContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOpValueWordContext() *OpValueWordContext {
	var p = new(OpValueWordContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_opValueWord
	return p
}

func InitEmptyOpValueWordContext(p *OpValueWordContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_opValueWord
}

func (*OpValueWordContext) IsOpValueWordContext() {}

func NewOpValueWordContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OpValueWordContext {
	var p = new(OpValueWordContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_opValueWord

	return p
}

func (s *OpValueWordContext) GetParser() antlr.Parser { return s.parser }

func (s *OpValueWordContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *OpValueWordContext) HASH_NUMBER() antlr.TerminalNode {
	return s.GetToken(UghParserHASH_NUMBER, 0)
}

func (s *OpValueWordContext) COLON() antlr.TerminalNode {
	return s.GetToken(UghParserCOLON, 0)
}

func (s *OpValueWordContext) COMMA() antlr.TerminalNode {
	return s.GetToken(UghParserCOMMA, 0)
}

func (s *OpValueWordContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OpValueWordContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OpValueWordContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterOpValueWord(s)
	}
}

func (s *OpValueWordContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitOpValueWord(s)
	}
}

func (s *OpValueWordContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitOpValueWord(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) OpValueWord() (localctx IOpValueWordContext) {
	localctx = NewOpValueWordContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, UghParserRULE_opValueWord)
	p.SetState(244)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case UghParserKW_ADD, UghParserKW_CREATE, UghParserKW_NEW, UghParserKW_SET, UghParserKW_EDIT, UghParserKW_UPDATE, UghParserKW_FIND, UghParserKW_SHOW, UghParserKW_LIST, UghParserKW_FILTER, UghParserKW_VIEW, UghParserKW_CONTEXT, UghParserKW_AND, UghParserKW_OR, UghParserKW_NOT, UghParserIDENT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(240)
			p.Ident()
		}

	case UghParserHASH_NUMBER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(241)
			p.Match(UghParserHASH_NUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCOLON:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(242)
			p.Match(UghParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case UghParserCOMMA:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(243)
			p.Match(UghParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIdentContext is an interface to support dynamic dispatch.
type IIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENT() antlr.TerminalNode
	KW_ADD() antlr.TerminalNode
	KW_CREATE() antlr.TerminalNode
	KW_NEW() antlr.TerminalNode
	KW_SET() antlr.TerminalNode
	KW_EDIT() antlr.TerminalNode
	KW_UPDATE() antlr.TerminalNode
	KW_FIND() antlr.TerminalNode
	KW_SHOW() antlr.TerminalNode
	KW_LIST() antlr.TerminalNode
	KW_FILTER() antlr.TerminalNode
	KW_VIEW() antlr.TerminalNode
	KW_CONTEXT() antlr.TerminalNode
	KW_AND() antlr.TerminalNode
	KW_OR() antlr.TerminalNode
	KW_NOT() antlr.TerminalNode

	// IsIdentContext differentiates from other interfaces.
	IsIdentContext()
}

type IdentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentContext() *IdentContext {
	var p = new(IdentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_ident
	return p
}

func InitEmptyIdentContext(p *IdentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = UghParserRULE_ident
}

func (*IdentContext) IsIdentContext() {}

func NewIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentContext {
	var p = new(IdentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = UghParserRULE_ident

	return p
}

func (s *IdentContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentContext) IDENT() antlr.TerminalNode {
	return s.GetToken(UghParserIDENT, 0)
}

func (s *IdentContext) KW_ADD() antlr.TerminalNode {
	return s.GetToken(UghParserKW_ADD, 0)
}

func (s *IdentContext) KW_CREATE() antlr.TerminalNode {
	return s.GetToken(UghParserKW_CREATE, 0)
}

func (s *IdentContext) KW_NEW() antlr.TerminalNode {
	return s.GetToken(UghParserKW_NEW, 0)
}

func (s *IdentContext) KW_SET() antlr.TerminalNode {
	return s.GetToken(UghParserKW_SET, 0)
}

func (s *IdentContext) KW_EDIT() antlr.TerminalNode {
	return s.GetToken(UghParserKW_EDIT, 0)
}

func (s *IdentContext) KW_UPDATE() antlr.TerminalNode {
	return s.GetToken(UghParserKW_UPDATE, 0)
}

func (s *IdentContext) KW_FIND() antlr.TerminalNode {
	return s.GetToken(UghParserKW_FIND, 0)
}

func (s *IdentContext) KW_SHOW() antlr.TerminalNode {
	return s.GetToken(UghParserKW_SHOW, 0)
}

func (s *IdentContext) KW_LIST() antlr.TerminalNode {
	return s.GetToken(UghParserKW_LIST, 0)
}

func (s *IdentContext) KW_FILTER() antlr.TerminalNode {
	return s.GetToken(UghParserKW_FILTER, 0)
}

func (s *IdentContext) KW_VIEW() antlr.TerminalNode {
	return s.GetToken(UghParserKW_VIEW, 0)
}

func (s *IdentContext) KW_CONTEXT() antlr.TerminalNode {
	return s.GetToken(UghParserKW_CONTEXT, 0)
}

func (s *IdentContext) KW_AND() antlr.TerminalNode {
	return s.GetToken(UghParserKW_AND, 0)
}

func (s *IdentContext) KW_OR() antlr.TerminalNode {
	return s.GetToken(UghParserKW_OR, 0)
}

func (s *IdentContext) KW_NOT() antlr.TerminalNode {
	return s.GetToken(UghParserKW_NOT, 0)
}

func (s *IdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.EnterIdent(s)
	}
}

func (s *IdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(UghParserListener); ok {
		listenerT.ExitIdent(s)
	}
}

func (s *IdentContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case UghParserVisitor:
		return t.VisitIdent(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *UghParser) Ident() (localctx IIdentContext) {
	localctx = NewIdentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, UghParserRULE_ident)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(246)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&34426845184) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
