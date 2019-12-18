package vsys

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testAccount    = "ATsSvqRGQeTpeQSGt3eAfNphmJwvnGU9dAw"
	testContract   = "CEyGzeTCVcGTGjySkq5thjaYYR1URSzH7vJ"
	testToken      = "TWtB9DcHiUswNG75siMHgVAbpz1TarXee6xhJHs6b"
	testTokenSplit = "TWt31ztsEevZtHtnoKMqmiQtRkeqoFPebddQiJFb1"
	testPrivateKey = "DvwNVbhTdn7XoCZW3W6YhkJrVk8Rq7NopuYRcz13tCzK"
	testSeed       = "save cancel drastic apart shaft mean session quick gap twist minimum borrow vessel art book"
)

func TestSendPaymentTx(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromSeed(testSeed, 0)
	tx := acc.BuildPayment("AU8TRrRkwmrssbCLfD9r8k5nBiLrAuVJEWP", 1e7, "vsys test send payment")
	resp, err := SendPaymentTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}

func TestSendLeasingTx(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromPrivateKey(testPrivateKey)
	tx := acc.BuildLeasing("AU8TRrRkwmrssbCLfD9r8k5nBiLrAuVJEWP", 1e7)
	resp, err := SendLeasingTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
	time.Sleep(5 * time.Second)
	tx = acc.BuildCancelLeasing(resp.Id)
	resp, err = SendCancelLeasingTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}

func TestRegisterContract(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromSeed(testSeed, 0)
	tx := acc.BuildRegisterContract(
		TokenContract,
		100000,
		1e8,
		"vsys test token desc",
		"vsys test contract desc")
	resp, err := SendRegisterContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}

func TestRegisterContract_Split(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromSeed(testSeed, 0)
	tx := acc.BuildRegisterContract(
		TokenContractWithSplit,
		100,
		1e8,
		"vsys test split token desc",
		"vsys test split contract desc")
	resp, err := SendRegisterContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}

func TestExecuteContract(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromSeed(testSeed, 0)
	// test issue
	a := &Contract{
		Amount: 3 * 1e8,
	}
	funcData := a.BuildIssueData()
	tx := acc.BuildExecuteContract(
		TokenId2ContractId(testToken),
		FuncidxIssue,
		funcData,
		"test execute contract - issue")
	resp, err := SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
	time.Sleep(3 * time.Second)

	// test destroy
	a = &Contract{
		Amount: 2 * 1e8,
	}
	funcData = a.BuildDestroyData()
	tx = acc.BuildExecuteContract(TokenId2ContractId(testToken), FuncidxDestroy, funcData, "test execute contract - destroy")
	resp, err = SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
	time.Sleep(3 * time.Second)

	// test send
	a = &Contract{
		Amount:    3 * 1e7, // need mul unity
		Recipient: "AUDRgBJjXM5zFMERzMML7pLPWikajTf8AKh",
	}
	funcData = a.BuildSendData()
	tx = acc.BuildExecuteContract(TokenId2ContractId(testToken), FuncidxSend, funcData, "test execute contract - send")
	resp, err = SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}

func TestExecuteContract_Split(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromSeed(testSeed, 0)

	// test issue
	a := &Contract{
		Amount: 4 * 1e8,
	}
	funcData := a.BuildIssueData()
	tx := acc.BuildExecuteContract(TokenId2ContractId(testTokenSplit), FuncidxIssue, funcData, "test execute contract - issue")
	resp, err := SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
	time.Sleep(5 * time.Second)

	// test destroy
	a = &Contract{
		Amount: 3 * 1e8,
	}
	funcData = a.BuildDestroyData()
	tx = acc.BuildExecuteContract(TokenId2ContractId(testTokenSplit), FuncidxDestroy, funcData, "test execute contract - destroy")
	resp, err = SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
	time.Sleep(5 * time.Second)

	// test send
	a = &Contract{
		ContractId: TokenId2ContractId(testToken),
		Amount:     3 * 1e7, // need mul unity
		Recipient:  "AUDRgBJjXM5zFMERzMML7pLPWikajTf8AKh",
	}
	funcData = a.BuildSendData()
	tx = acc.BuildExecuteContract(TokenId2ContractId(testTokenSplit), FuncidxSendSplit, funcData, "test execute contract - send")
	resp, err = SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
	time.Sleep(5 * time.Second)

	// test split
	a = &Contract{
		NewUnity: 1e4,
	}
	funcData = a.BuildSplitData()
	tx = acc.BuildExecuteContract(TokenId2ContractId(testTokenSplit), FuncidxSplit, funcData, "test execute contract - split")
	resp, err = SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}

func TestSendTokenTransaction(t *testing.T) {
	InitApi("http://test.v.systems:9922", Testnet)
	acc := InitAccount(Testnet)
	acc.BuildFromSeed(testSeed, 0)

	tx := acc.BuildSendTokenTransaction(
		testToken,
		"AUDRgBJjXM5zFMERzMML7pLPWikajTf8AKh",
		3*1e5,
		false,
		"test send token transaction")
	resp, err := SendExecuteContractTx(tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, resp.Error, 0)
}
