// generated by stringer -type=Instruction,Port -output=types_string.go; DO NOT EDIT

package tis100mc

import "fmt"

const _Instruction_name = "NOPMOVSWPSAVADDSUBNEGJMPJEZJNZJGZJLZJRO"

var _Instruction_index = [...]uint8{3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36, 39}

func (i Instruction) String() string {
	i -= 1000
	if i < 0 || i >= Instruction(len(_Instruction_index)) {
		return fmt.Sprintf("Instruction(%d)", i+1000)
	}
	hi := _Instruction_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _Instruction_index[i-1]
	}
	return _Instruction_name[lo:hi]
}

const _Port_name = "ACCBAKNILLEFTRIGHTUPDOWNANYLAST"

var _Port_index = [...]uint8{3, 6, 9, 13, 18, 20, 24, 27, 31}

func (i Port) String() string {
	i -= 1013
	if i < 0 || i >= Port(len(_Port_index)) {
		return fmt.Sprintf("Port(%d)", i+1013)
	}
	hi := _Port_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _Port_index[i-1]
	}
	return _Port_name[lo:hi]
}
