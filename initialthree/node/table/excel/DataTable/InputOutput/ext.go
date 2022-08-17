package InputOutput

import (
	"initialthree/node/common/inoutput"
)

func (this *InputOutput) InOut(count int32) ([]inoutput.ResDesc, []inoutput.ResDesc) {
	input := make([]inoutput.ResDesc, 0, len(this.Input))
	output := make([]inoutput.ResDesc, 0, len(this.Output))

	for _, v := range this.Input {
		input = append(input, inoutput.ResDesc{ID: v.ID, Type: int(v.Type), Count: v.Count * count})
	}

	for _, v := range this.Output {
		output = append(output, inoutput.ResDesc{ID: v.ID, Type: int(v.Type), Count: v.Count * count})
	}
	return input, output
}
