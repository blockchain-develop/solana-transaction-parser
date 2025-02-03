package raydium_amm

import (
	"errors"
	"github.com/blockchain-develop/solana-parser/log"
	"github.com/blockchain-develop/solana-parser/program"
	"github.com/blockchain-develop/solana-parser/types"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/raydium_amm"
)

var (
	Parsers = make(map[uint64]Parser, 0)
)

type Parser func(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta)

func RegisterParser(id uint64, p Parser) {
	Parsers[id] = p
}

func init() {
	program.RegisterParser(raydium_amm.ProgramID, ProgramParser)
	RegisterParser(uint64(raydium_amm.Instruction_Initialize), ParseInitialize)
	RegisterParser(uint64(raydium_amm.Instruction_Initialize2), ParseInitialize2)
	RegisterParser(uint64(raydium_amm.Instruction_MonitorStep), ParseMonitorStep)
	RegisterParser(uint64(raydium_amm.Instruction_Deposit), ParseDeposit)
	RegisterParser(uint64(raydium_amm.Instruction_Withdraw), ParseWithdraw)
	RegisterParser(uint64(raydium_amm.Instruction_MigrateToOpenBook), ParseMigrateToOpenBook)
	RegisterParser(uint64(raydium_amm.Instruction_SetParams), ParseSetParams)
	RegisterParser(uint64(raydium_amm.Instruction_WithdrawPnl), ParseWithdrawPnl)
	RegisterParser(uint64(raydium_amm.Instruction_WithdrawSrm), ParseWithdrawSrm)
	RegisterParser(uint64(raydium_amm.Instruction_SwapBaseIn), ParseSwapBaseIn)
	RegisterParser(uint64(raydium_amm.Instruction_PreInitialize), ParsePreInitialize)
	RegisterParser(uint64(raydium_amm.Instruction_SwapBaseOut), ParseSwapBaseOut)
	RegisterParser(uint64(raydium_amm.Instruction_SimulateInfo), ParseSimulateInfo)
	RegisterParser(uint64(raydium_amm.Instruction_AdminCancelOrders), ParseAdminCancelOrders)
	RegisterParser(uint64(raydium_amm.Instruction_CreateConfigAccount), ParseCreateConfigAccount)
	RegisterParser(uint64(raydium_amm.Instruction_UpdateConfigAccount), ParseUpdateConfigAccount)
}

func ProgramParser(in *types.Instruction, meta *types.Meta) error {
	inst, err := raydium_amm.DecodeInstruction(in.AccountMetas(), in.Instruction.Data)
	if err != nil {
		return err
	}
	id := uint64(inst.TypeID.Uint32())
	parser, ok := Parsers[id]
	if !ok {
		return errors.New("parser not found")
	}
	parser(inst, in, meta)
	return nil
}

func insertAccount(accounts []*solana.AccountMeta, index int) []*solana.AccountMeta {
	s := make([]*solana.AccountMeta, 0)
	s = append(s, accounts[0:index]...)
	s = append(s, &solana.AccountMeta{
		PublicKey:  solana.MustPublicKeyFromBase58("11111111111111111111111111111111"),
		IsWritable: true,
		IsSigner:   false,
	})
	s = append(s, accounts[index:]...)
	return s
}

func ParseInitialize(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse initialize", "program", raydium_amm.ProgramName)
}
func ParseInitialize2(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Initialize2)
	// the latest three transfer
	index := len(in.Children)
	t1 := in.Children[index-3].Event[0].(*types.Transfer)
	t2 := in.Children[index-2].Event[0].(*types.Transfer)
	t3 := in.Children[index-1].Event[0].(*types.MintTo)
	createPool := &types.CreatePool{
		Pool:    inst1.GetAmmAccount().PublicKey,
		User:    inst1.GetAmmAuthorityAccount().PublicKey,
		TokenA:  inst1.GetPcMintAccount().PublicKey,
		TokenB:  inst1.GetCoinMintAccount().PublicKey,
		TokenLP: inst1.GetLpMintAccount().PublicKey,
		VaultA:  inst1.GetPoolPcTokenAccountAccount().PublicKey,
		VaultB:  inst1.GetPoolCoinTokenAccountAccount().PublicKey,
		VaultLP: inst1.GetPoolTempLpAccount().PublicKey,
	}
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetAmmAuthorityAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		TokenLpMint:    t3,
	}
	pool := &types.Pool{
		Hash:     inst1.GetAmmAccount().PublicKey,
		MintA:    inst1.GetCoinMintAccount().PublicKey,
		MintB:    inst1.GetPcMintAccount().PublicKey,
		MintLp:   inst1.GetLpMintAccount().PublicKey,
		VaultA:   inst1.GetPoolCoinTokenAccountAccount().PublicKey,
		VaultB:   inst1.GetPoolPcTokenAccountAccount().PublicKey,
		ReserveA: 0,
		ReserveB: 0,
	}
	in.Event = []interface{}{createPool, addLiquidity}
	in.Receipt = []interface{}{pool}
}
func ParseMonitorStep(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseDeposit(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Deposit)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	addLiquidity := &types.AddLiquidity{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetUserOwnerAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
	}
	in.Event = []interface{}{addLiquidity, addLiquidity}
}
func ParseWithdraw(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.Withdraw)
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	t3 := in.Children[2].Event[0].(*types.Burn)
	removeLiquidity := &types.RemoveLiquidity{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetUserOwnerAccount().PublicKey,
		TokenATransfer: t1,
		TokenBTransfer: t2,
		TokenLpBurn:    t3,
	}
	in.Event = []interface{}{removeLiquidity}
}
func ParseMigrateToOpenBook(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSetParams(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdrawPnl(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseWithdrawSrm(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseSwapBaseIn(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.SwapBaseIn)
	if len(inst1.GetAccounts()) == 17 {
		inst1.SetAccounts(insertAccount(inst1.GetAccounts(), 4))
	}
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetUserSourceOwnerAccount().PublicKey,
		InputTransfer:  t1,
		OutputTransfer: t2,
	}
	in.Event = []interface{}{swap}
}
func ParsePreInitialize(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	log.Logger.Info("ignore parse pre-initialize", "program", raydium_amm.ProgramName)
}
func ParseSwapBaseOut(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	inst1 := inst.Impl.(*raydium_amm.SwapBaseOut)
	if len(inst1.GetAccounts()) == 17 {
		inst1.SetAccounts(insertAccount(inst1.GetAccounts(), 4))
	}
	t1 := in.Children[0].Event[0].(*types.Transfer)
	t2 := in.Children[1].Event[0].(*types.Transfer)
	swap := &types.Swap{
		Pool:           inst1.GetAmmAccount().PublicKey,
		User:           inst1.GetUserSourceOwnerAccount().PublicKey,
		InputTransfer:  t1,
		OutputTransfer: t2,
	}
	in.Event = []interface{}{swap}
}
func ParseSimulateInfo(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseAdminCancelOrders(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseCreateConfigAccount(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}
func ParseUpdateConfigAccount(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
}

// Default
func ParseDefault(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	return
}

// Fault
func ParseFault(inst *raydium_amm.Instruction, in *types.Instruction, meta *types.Meta) {
	panic("not supported")
}
