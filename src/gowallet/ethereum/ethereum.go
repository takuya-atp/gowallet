package ethereum

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	goethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	GethURL = "http://localhost:8545/"
)

type EthClient struct {
	*ethclient.Client
}

func (e *EthClient) ConfirmBalance(pass string, keyPath string) {
	key, err := e.unlockAccount(pass, keyPath)
	if err != nil {
		return
	}
	balance, err := e.Client.BalanceAt(context.TODO(), key.Address, nil)
	if err != nil {
		fmt.Println("Failed to confirm balance")
		return
	}
	fmt.Printf("Your wallet balance: %d", balance)
}

func (e *EthClient) unlockAccount(pass string, keyPath string) (*keystore.Key, error) {
	if keyPath == "" {
		fmt.Println("Enter keypath of your wallet")
		return nil, errors.New("Failed")
	}
	f, err := os.Open(keyPath)
	if err != nil {
		fmt.Println("Failed to load key file")
		return nil, errors.New("Failed")
	}
	defer f.Close()
	json, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("Failed to read file")
		return nil, errors.New("Failed")
	}
	key, err := keystore.DecryptKey(json, pass)
	if err != nil {
		fmt.Println("Incorrect password")
		return nil, errors.New("Failed")
	}
	return key, nil
}

func (e *EthClient) TransferWei(pass string, wei int, keyPath string, addr string) {
	weib := big.NewInt(int64(wei))
	from, err := e.unlockAccount(pass, keyPath)
	if err != nil {
		return
	}
	if common.IsHexAddress(addr) == false {
		fmt.Println("Addr is not valid address")
		return
	}
	destAccount := common.HexToAddress(addr)
	nonce, err := e.Client.PendingNonceAt(context.TODO(), from.Address)
	if err != nil {
		return
	}

	gasPrice, err := e.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		return
	}
	msg := goethereum.CallMsg{From: from.Address, Value: weib}
	gasLimit, err := e.Client.EstimateGas(context.TODO(), msg)
	if err != nil {
		return
	}

	rawTx := types.NewTransaction(nonce, destAccount, weib, gasLimit, gasPrice, nil)
	signer := types.HomesteadSigner{}
	signature, err := crypto.Sign(signer.Hash(rawTx).Bytes(), from.PrivateKey)
	if err != nil {
		return
	}
	signedTx, err := rawTx.WithSignature(signer, signature)
	if err != nil {
		return
	}
	if err := e.Client.SendTransaction(context.TODO(), signedTx); err != nil {
		return
	}
	return
}

func (e *EthClient) GetAccount(pass string) {
	if pass == "" {
		fmt.Println("Enter your wallet password")
		return
	}
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("error: %#v", errors.Wrap(err, "Failed to find $HOME"))
		return
	}

	ks := keystore.NewKeyStore(home, keystore.StandardScryptN, keystore.StandardScryptP)
	ac, err := ks.NewAccount(pass)
	if err != nil {
		fmt.Printf("error: %#v", errors.Wrap(err, "Failed to create account"))
		return
	}
	fmt.Println("Successed create account, Save your private key")
	fmt.Printf("address: %s, private_key_path: %s", ac.Address.String(), ac.URL.Path)
}

func NewEthClient() (*EthClient, error) {
	cl, err := ethclient.Dial(GethURL)
	if err != nil {
		return nil, err
	}
	return &EthClient{
		Client: cl,
	}, nil
}
