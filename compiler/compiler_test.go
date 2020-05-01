package compiler

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	asm "smart/assembler"
	"smart/ast"
	"smart/parser"
	"smart/scanner"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
)

type CompilerTestSuite struct {
	suite.Suite
	auth    *bind.TransactOpts
	address common.Address
	gAlloc  core.GenesisAlloc
	sim     *backends.SimulatedBackend
	key     *ecdsa.PrivateKey
	// compiler *Compiler
}

func (suite *CompilerTestSuite) SetupTest() {
	suite.key, _ = crypto.GenerateKey()
	suite.auth = bind.NewKeyedTransactor(suite.key)

	suite.address = suite.auth.From
	suite.gAlloc = map[common.Address]core.GenesisAccount{
		suite.address: {Balance: big.NewInt(10000000000)},
	}

	suite.sim = backends.NewSimulatedBackend(suite.gAlloc, 8000000)
	// chainID := suite.sim.Blockchain().Config().ChainID

	// tx := types.NewContractCreation(0, new(big.Int), 1000000, new(big.Int), common.FromHex(`0x6004600d60003960046000f30011223344`))
	// signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), suite.key)
	// if err != nil {
	// 	panic(err)
	// }
	// suite.sim.SendTransaction(nil, signedTx)
	// suite.sim.Commit()
	// receipt, err := suite.sim.TransactionReceipt(nil, signedTx.Hash())
	// fmt.Println(receipt.ContractAddress)
	// code, err := suite.sim.CodeAt(nil, receipt.ContractAddress, suite.sim.Blockchain().CurrentBlock().Number())
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(common.ToHex(code))

	// _, _, hw, e := DeployCompiler(suite.auth, suite.sim)
	// suite.Compiler = hw
	// suite.Nil(e)
	// suite.sim.Commit()
}

func TestSemantics(t *testing.T) {
	file, err := os.Open("test.smrt")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	s := scanner.New(r)
	p := parser.New(s, "test.smrt")

	tree := p.Parse_prog()
	ast.Check(tree)
}

func (suite *CompilerTestSuite) TestCodeGen() {
	file, err := os.Open("codegen.smrt")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	s := scanner.New(r)
	p := parser.New(s, "codegen.smrt")

	tree := p.Parse_prog()
	ast.Check(tree)
	asm := asm.NewAssembler()
	ast.CodeGen(tree, asm)
	common.ToHex(asm.Compile())
}

func (suite *CompilerTestSuite) TestDeploy() {
	file, err := os.Open("deploy.smrt")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	s := scanner.New(r)
	p := parser.New(s, "deploy.smrt")

	tree := p.Parse_prog()
	ast.Check(tree)
	asm := asm.NewAssembler()
	ast.CodeGen(tree, asm)
	// asm.Print()
	deployable := asm.Compile()
	fmt.Println(common.ToHex(deployable))
	nonce, err := suite.sim.NonceAt(nil, suite.auth.From, suite.sim.Blockchain().CurrentBlock().Number())
	if err != nil {
		panic(err)
	}
	tx := types.NewContractCreation(nonce, new(big.Int), 1000000, new(big.Int), deployable)
	chainID := suite.sim.Blockchain().Config().ChainID
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), suite.key)
	if err != nil {
		panic(err)
	}
	suite.sim.SendTransaction(nil, signedTx)
	suite.sim.Commit()
	receipt, err := suite.sim.TransactionReceipt(nil, signedTx.Hash())
	fmt.Printf("Status: %d\n", receipt.Status)
}

func TestCompilerTestSuite(t *testing.T) {
	suite.Run(t, new(CompilerTestSuite))
}
