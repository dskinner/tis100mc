package tis100mc

import "testing"

const example = `# A COMMENT
  MOV UP, DOWN

 NOP

LOOP:
 MOV   LEFT  ,    RIGHT  
 JMP  LOOP

SWP#COMMENT
SAV

ADD 1 # FOR COMMENT
SUB LEFT
NEG

RET: MOV 0, DOWN`

var expects = []Token{
	{typ: TokenComment, val: "# A COMMENT"},
	{typ: TokenInstruction, val: "MOV"},
	{typ: TokenArgument, val: "UP"},
	{typ: TokenComma, val: ""},
	{typ: TokenArgument, val: "DOWN"},
	{typ: TokenInstruction, val: "NOP"},
	{typ: TokenLabel, val: "LOOP"},
	{typ: TokenInstruction, val: "MOV"},
	{typ: TokenArgument, val: "LEFT"},
	{typ: TokenComma, val: ""},
	{typ: TokenArgument, val: "RIGHT"},
	{typ: TokenInstruction, val: "JMP"},
	{typ: TokenArgument, val: "LOOP"},
	{typ: TokenInstruction, val: "SWP"},
	{typ: TokenComment, val: "#COMMENT"},
	{typ: TokenInstruction, val: "SAV"},
	{typ: TokenInstruction, val: "ADD"},
	{typ: TokenArgument, val: "1"},
	{typ: TokenComment, val: "# FOR COMMENT"},
	{typ: TokenInstruction, val: "SUB"},
	{typ: TokenArgument, val: "LEFT"},
	{typ: TokenInstruction, val: "NEG"},
	{typ: TokenLabel, val: "RET"},
	{typ: TokenInstruction, val: "MOV"},
	{typ: TokenArgument, val: "0"},
	{typ: TokenComma, val: ""},
	{typ: TokenArgument, val: "DOWN"},
}

type LogTokenReceiver struct {
	t        *testing.T
	received []Token
}

func (r *LogTokenReceiver) Receive(tkn Token) {
	r.t.Logf("%s - %q\n", tkn.typ, tkn.val)
	r.received = append(r.received, tkn)
}

func TestLex(t *testing.T) {
	r := &LogTokenReceiver{t: t}
	l := NewLexer(r)
	l.bytes = []byte(example)
	l.Run()
	if len(expects) != len(r.received) {
		t.Fatal("Length of received tokens doesn't match expected.")
	}
	for i := 0; i < len(expects); i++ {
		a := expects[i]
		b := r.received[i]
		if a.typ != b.typ {
			t.Errorf("expected.typ(%s) != received.typ(%s)\n", a.typ, b.typ)
		}
	}
}

// TODO
func testEOF(t *testing.T) {

}
