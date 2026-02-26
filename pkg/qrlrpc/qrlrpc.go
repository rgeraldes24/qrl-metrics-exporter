package qrlrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
)

// QRLError - qrl error
type QRLError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err QRLError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type qrlResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *QRLError       `json:"error"`
}

type qrlRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// QRLRPC - QRL rpc client
type QRLRPC struct {
	url    string
	client httpClient
	log    logger
	Debug  bool
}

// New create new rpc client with given url
func New(url string, options ...func(rpc *QRLRPC)) *QRLRPC {
	rpc := &QRLRPC{
		url:    url,
		client: http.DefaultClient,
		log:    log.New(os.Stderr, "", log.LstdFlags),
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// NewQRLRPC create new rpc client with given url
func NewQRLRPC(url string, options ...func(rpc *QRLRPC)) *QRLRPC {
	return New(url, options...)
}

func (rpc *QRLRPC) call(method string, target interface{}, params ...interface{}) error {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// URL returns client url
func (rpc *QRLRPC) URL() string {
	return rpc.url
}

// Call returns raw response of method call
func (rpc *QRLRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := qrlRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if rpc.Debug {
		rpc.log.Println(fmt.Sprintf("%s\nRequest: %s\nResponse: %s\n", method, body, data))
	}

	resp := new(qrlResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil
}

// RawCall returns raw response of method call (Deprecated)
func (rpc *QRLRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

// Web3ClientVersion returns the current client version.
func (rpc *QRLRPC) Web3ClientVersion() (string, error) {
	var clientVersion string

	err := rpc.call("web3_clientVersion", &clientVersion)
	return clientVersion, err
}

// Web3Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (rpc *QRLRPC) Web3Sha3(data []byte) (string, error) {
	var hash string

	err := rpc.call("web3_sha3", &hash, fmt.Sprintf("0x%x", data))
	return hash, err
}

// NetVersion returns the current network protocol version.
func (rpc *QRLRPC) NetVersion() (string, error) {
	var version string

	err := rpc.call("net_version", &version)
	return version, err
}

// NetListening returns true if client is actively listening for network connections.
func (rpc *QRLRPC) NetListening() (bool, error) {
	var listening bool

	err := rpc.call("net_listening", &listening)
	return listening, err
}

// NetPeerCount returns number of peers currently connected to the client.
func (rpc *QRLRPC) NetPeerCount() (int, error) {
	var response string
	if err := rpc.call("net_peerCount", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// QRLProtocolVersion returns the current qrl protocol version.
func (rpc *QRLRPC) QRLProtocolVersion() (string, error) {
	var protocolVersion string

	err := rpc.call("qrl_protocolVersion", &protocolVersion)
	return protocolVersion, err
}

// QRLSyncing returns an object with data about the sync status or false.
func (rpc *QRLRPC) QRLSyncing() (*Syncing, error) {
	result, err := rpc.RawCall("qrl_syncing")
	if err != nil {
		return nil, err
	}
	syncing := new(Syncing)
	if bytes.Equal(result, []byte("false")) {
		return syncing, nil
	}
	err = json.Unmarshal(result, syncing)
	return syncing, err
}

// QRLCoinbase returns the client coinbase address
func (rpc *QRLRPC) QRLCoinbase() (string, error) {
	var address string

	err := rpc.call("qrl_coinbase", &address)
	return address, err
}

// QRLMining returns true if client is actively mining new blocks.
func (rpc *QRLRPC) QRLMining() (bool, error) {
	var mining bool

	err := rpc.call("qrl_mining", &mining)
	return mining, err
}

// QRLGasPrice returns the current price per gas in planck.
func (rpc *QRLRPC) QRLGasPrice() (big.Int, error) {
	var response string
	if err := rpc.call("qrl_gasPrice", &response); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// QRLAccounts returns a list of addresses owned by client.
func (rpc *QRLRPC) QRLAccounts() ([]string, error) {
	accounts := []string{}

	err := rpc.call("qrl_accounts", &accounts)
	return accounts, err
}

// QRLBlockNumber returns the number of most recent block.
func (rpc *QRLRPC) QRLBlockNumber() (int, error) {
	var response string
	if err := rpc.call("qrl_blockNumber", &response); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// QRLGetBalance returns the balance of the account of given address in planck.
func (rpc *QRLRPC) QRLGetBalance(address, block string) (big.Int, error) {
	var response string
	if err := rpc.call("qrl_getBalance", &response, address, block); err != nil {
		return big.Int{}, err
	}

	return ParseBigInt(response)
}

// QRLGetStorageAt returns the value from a storage position at a given address.
func (rpc *QRLRPC) QRLGetStorageAt(data string, position int, tag string) (string, error) {
	var result string

	err := rpc.call("qrl_getStorageAt", &result, data, IntToHex(position), tag)
	return result, err
}

// QRLGetTransactionCount returns the number of transactions sent from an address.
func (rpc *QRLRPC) QRLGetTransactionCount(address, block string) (int, error) {
	var response string

	if err := rpc.call("qrl_getTransactionCount", &response, address, block); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// QRLGetBlockTransactionCountByHash returns the number of transactions in a block from a block matching the given block hash.
func (rpc *QRLRPC) QRLGetBlockTransactionCountByHash(hash string) (int, error) {
	var response string

	if err := rpc.call("qrl_getBlockTransactionCountByHash", &response, hash); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// QRLGetBlockTransactionCountByNumber returns the number of transactions in a block from a block matching the given block
func (rpc *QRLRPC) QRLGetBlockTransactionCountByNumber(number int) (int, error) {
	var response string

	if err := rpc.call("qrl_getBlockTransactionCountByNumber", &response, IntToHex(number)); err != nil {
		return 0, err
	}

	return ParseInt(response)
}

// QRLGetCode returns code at a given address.
func (rpc *QRLRPC) QRLGetCode(address, block string) (string, error) {
	var code string

	err := rpc.call("qrl_getCode", &code, address, block)
	return code, err
}

// QRLSign signs data with a given address.
// Calculates an QRLereum specific signature with: sign(keccak256("\x19QRLereum Signed Message:\n" + len(message) + message)))
func (rpc *QRLRPC) QRLSign(address, data string) (string, error) {
	var signature string

	err := rpc.call("qrl_sign", &signature, address, data)
	return signature, err
}

// QRLSendTransaction creates new message call transaction or a contract creation, if the data field contains code.
func (rpc *QRLRPC) QRLSendTransaction(transaction T) (string, error) {
	var hash string

	err := rpc.call("qrl_sendTransaction", &hash, transaction)
	return hash, err
}

// QRLSendRawTransaction creates new message call transaction or a contract creation for signed transactions.
func (rpc *QRLRPC) QRLSendRawTransaction(data string) (string, error) {
	var hash string

	err := rpc.call("qrl_sendRawTransaction", &hash, data)
	return hash, err
}

// QRLCall executes a new message call immediately without creating a transaction on the block chain.
func (rpc *QRLRPC) QRLCall(transaction T, tag string) (string, error) {
	var data string

	err := rpc.call("qrl_call", &data, transaction, tag)
	return data, err
}

// QRLEstimateGas makes a call or transaction, which won't be added to the blockchain and returns the used gas, which can be used for estimating the used gas.
func (rpc *QRLRPC) QRLEstimateGas(transaction T) (int, error) {
	var response string

	err := rpc.call("qrl_estimateGas", &response, transaction)
	if err != nil {
		return 0, err
	}

	return ParseInt(response)
}

func (rpc *QRLRPC) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	result, err := rpc.RawCall(method, params...)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(result, []byte("null")) {
		return nil, nil
	}

	var response proxyBlock
	if withTransactions {
		response = new(proxyBlockWithTransactions)
	} else {
		response = new(proxyBlockWithoutTransactions)
	}

	err = json.Unmarshal(result, response)
	if err != nil {
		return nil, err
	}

	block := response.toBlock()
	return &block, nil
}

// QRLGetBlockByHash returns information about a block by hash.
func (rpc *QRLRPC) QRLGetBlockByHash(hash string, withTransactions bool) (*Block, error) {
	return rpc.getBlock("qrl_getBlockByHash", withTransactions, hash, withTransactions)
}

// QRLGetBlockByNumber returns information about a block by block number.
func (rpc *QRLRPC) QRLGetBlockByNumber(number int, withTransactions bool) (*Block, error) {
	return rpc.getBlock("qrl_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
}

func (rpc *QRLRPC) getTransaction(method string, params ...interface{}) (*Transaction, error) {
	transaction := new(Transaction)

	err := rpc.call(method, transaction, params...)
	return transaction, err
}

// QRLGetTransactionByHash returns the information about a transaction requested by transaction hash.
func (rpc *QRLRPC) QRLGetTransactionByHash(hash string) (*Transaction, error) {
	return rpc.getTransaction("qrl_getTransactionByHash", hash)
}

// QRLGetTransactionByBlockHashAndIndex returns information about a transaction by block hash and transaction index position.
func (rpc *QRLRPC) QRLGetTransactionByBlockHashAndIndex(blockHash string, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("qrl_getTransactionByBlockHashAndIndex", blockHash, IntToHex(transactionIndex))
}

// QRLGetTransactionByBlockNumberAndIndex returns information about a transaction by block number and transaction index position.
func (rpc *QRLRPC) QRLGetTransactionByBlockNumberAndIndex(blockNumber, transactionIndex int) (*Transaction, error) {
	return rpc.getTransaction("qrl_getTransactionByBlockNumberAndIndex", IntToHex(blockNumber), IntToHex(transactionIndex))
}

// QRLGetTransactionReceipt returns the receipt of a transaction by transaction hash.
// Note That the receipt is not available for pending transactions.
func (rpc *QRLRPC) QRLGetTransactionReceipt(hash string) (*TransactionReceipt, error) {
	transactionReceipt := new(TransactionReceipt)

	err := rpc.call("qrl_getTransactionReceipt", transactionReceipt, hash)
	if err != nil {
		return nil, err
	}

	return transactionReceipt, nil
}

// QRLGetCompilers returns a list of available compilers in the client.
func (rpc *QRLRPC) QRLGetCompilers() ([]string, error) {
	compilers := []string{}

	err := rpc.call("qrl_getCompilers", &compilers)
	return compilers, err
}

// QRLNewFilter creates a new filter object.
func (rpc *QRLRPC) QRLNewFilter(params FilterParams) (string, error) {
	var filterID string
	err := rpc.call("qrl_newFilter", &filterID, params)
	return filterID, err
}

// QRLNewBlockFilter creates a filter in the node, to notify when a new block arrives.
// To check if the state has changed, call QRLGetFilterChanges.
func (rpc *QRLRPC) QRLNewBlockFilter() (string, error) {
	var filterID string
	err := rpc.call("qrl_newBlockFilter", &filterID)
	return filterID, err
}

// QRLNewPendingTransactionFilter creates a filter in the node, to notify when new pending transactions arrive.
// To check if the state has changed, call QRLGetFilterChanges.
func (rpc *QRLRPC) QRLNewPendingTransactionFilter() (string, error) {
	var filterID string
	err := rpc.call("qrl_newPendingTransactionFilter", &filterID)
	return filterID, err
}

// QRLUninstallFilter uninstalls a filter with given id.
func (rpc *QRLRPC) QRLUninstallFilter(filterID string) (bool, error) {
	var res bool
	err := rpc.call("qrl_uninstallFilter", &res, filterID)
	return res, err
}

// QRLGetFilterChanges polling method for a filter, which returns an array of logs which occurred since last poll.
func (rpc *QRLRPC) QRLGetFilterChanges(filterID string) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("qrl_getFilterChanges", &logs, filterID)
	return logs, err
}

// QRLGetFilterLogs returns an array of all logs matching filter with given id.
func (rpc *QRLRPC) QRLGetFilterLogs(filterID string) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("qrl_getFilterLogs", &logs, filterID)
	return logs, err
}

// QRLGetLogs returns an array of all logs matching a given filter object.
func (rpc *QRLRPC) QRLGetLogs(params FilterParams) ([]Log, error) {
	var logs = []Log{}
	err := rpc.call("qrl_getLogs", &logs, params)
	return logs, err
}

// QRL1 returns 1 quanta value (10^18 planck)
func (rpc *QRLRPC) QRL1() *big.Int {
	return QRL1()
}

// QRL1 returns 1 quanta value (10^18 planck)
func QRL1() *big.Int {
	return big.NewInt(1000000000000000000)
}
