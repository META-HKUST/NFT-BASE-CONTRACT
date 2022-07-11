package chaincode

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


func (c *TokenERC721Contract) SetBaseURI(ctx contractapi.TransactionContextInterface,uri string) string{
	err := ctx.GetClientIdentity().AssertAttributeValue("unifit.creator", "true")
	if err != nil {
		return "submitting client not authorized to set BaseURI, does not have unifit.creator role"
	}
	
	baseURI = uri

	err = ctx.GetStub().PutState(baseURIPrefix, []byte(baseURI))
	if err != nil {
		return "failed to PutState baseURI"
	}
	value :="BaseURI set to "+baseURI
	return value
}






func (c *TokenERC721Contract) PublicMint(ctx contractapi.TransactionContextInterface) (string, error) {
	baseURI = c.BaseURI(ctx)
	if len(baseURI) == 0 {
		return "baseURI is nil", fmt.Errorf("baseURI is nil, failed to mint NFT token")
	}
	nft, err := _mint(ctx)
	if err != nil {
		return "something wrong", fmt.Errorf("failed to mint NFT token: %v", err)
	}
	return nft.TokenId, err
}


func (c *TokenERC721Contract) PublicNtimesMint(ctx contractapi.TransactionContextInterface, ntimes int) (string, error) {
	baseURI = c.BaseURI(ctx)
	if len(baseURI) == 0 {
		return "baseURI is nil", fmt.Errorf("baseURI is nil, failed to mint NFT token")
	}
	
	minter64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "user is nil", fmt.Errorf("failed to get minter id: %v", err)
	}

	minterBytes, err := base64.StdEncoding.DecodeString(minter64)
	if err != nil {
		return "user is nil", fmt.Errorf("failed to DecodeString minter64: %v", err)
	}
	minter := string(minterBytes)
	
	i := _countOf(ctx, minter)
	n := i+ntimes
	for i<n {
		id:= strconv.Itoa(i)
		_, err := _ntimesmint(ctx,id)
		if err != nil {
			p := "Number"+id+"is failed to mint"
			return p, fmt.Errorf("failed to mint NFT token: %v", err)
		}
		i++
	}
	
	return "finished", nil
}



func (c *TokenERC721Contract) PrivateBurn(ctx contractapi.TransactionContextInterface, tokenId string) (bool, error) {	
	err := ctx.GetClientIdentity().AssertAttributeValue("unifit.creator", "true")
	if err != nil {
		return false, fmt.Errorf("submitting client not authorized to burn NFT, does not have unifit.creator role")
	}
	

	burnResult, err := _burn(ctx, tokenId)
	if err != nil {
		return false, fmt.Errorf("failed to get Burn: %v", err)
	}
	if !burnResult {
		return false, fmt.Errorf("failed to burn NFT token")
	}
	return true, nil
}


// ClientAccountBalance returns the balance of the requesting client's account.
// returns {Number} Returns the account balance
func (c *TokenERC721Contract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {
	// Get ID of submitting client identity
	clientAccountID64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, fmt.Errorf("failed to GetClientIdentity minter: %v", err)
	}

	clientAccountIDBytes, err := base64.StdEncoding.DecodeString(clientAccountID64)
	if err != nil {
		return 0, fmt.Errorf("failed to DecodeString sender: %v", err)
	}

	clientAccountID := string(clientAccountIDBytes)

	return c.BalanceOf(ctx, clientAccountID), nil
}



// ClientAccountID returns the id of the requesting client's account.
// In this implementation, the client account ID is the clientId itself.
// Users can use this function to get their own account id, which they can then give to others as the payment address

func (c *TokenERC721Contract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get ID of submitting client identity
	clientAccountID64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to GetClientIdentity minter: %v", err)
	}

	clientAccountBytes, err := base64.StdEncoding.DecodeString(clientAccountID64)
	if err != nil {
		return "", fmt.Errorf("failed to DecodeString clientAccount64: %v", err)
	}
	clientAccount := string(clientAccountBytes)

	return clientAccount, nil
}

// TotalSupply counts non-fungible tokens tracked by this contract.
//
// @param {Context} ctx the transaction context
// @returns {Number} Returns a count of valid non-fungible tokens tracked by this contract,
// where each one of them has an assigned and queryable owner.

func (c *TokenERC721Contract) TotalSupply(ctx contractapi.TransactionContextInterface) int {
	// There is a key record for every non-fungible token in the format of nftPrefix.tokenId.
	// TotalSupply() queries for and counts all records matching nftPrefix.*

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(nftPrefix, []string{})
	if err != nil {
		panic("Error creating GetStateByPartialCompositeKey:" + err.Error())
	}
	// Count the number of returned composite keys

	totalSupply := 0
	for iterator.HasNext() {
		_, err := iterator.Next()
		if err != nil {
			return 0
		}
		totalSupply++

	}
	return totalSupply

}
