package program

import (
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
)

var (
	Parsers = make(map[solana.PublicKey]Parser, 0)
)

type Parser func(in *types.Instruction, meta *types.Meta) error

func RegisterParser(program solana.PublicKey, p Parser) {
	Parsers[program] = p
}

func Parse(in *types.Instruction, meta *types.Meta) error {
	programId := in.Instruction.ProgramId
	parser, ok := Parsers[programId]
	if !ok {
		return nil
	}
	return parser(in, meta)
}
