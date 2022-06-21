/*
Copyright 2021 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"path"
	"time"
)

const (
	mspID         = "Org2MSP"
	cryptoPath    = "../../test-network/organizations/peerOrganizations/org2.example.com"
	certPath      = cryptoPath + "/users/User1@org2.example.com/msp/signcerts/cert.pem"
	keyPath       = cryptoPath + "/users/User1@org2.example.com/msp/keystore/"
	tlsCertPath   = cryptoPath + "/peers/peer0.org2.example.com/tls/ca.crt"
	peerEndpoint  = "localhost:9051"
	gatewayPeer   = "peer0.org2.example.com"
	channelName   = "mychannel"
	chaincodeName = "basic"
)

var now = time.Now()
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

func main() {
	log.Println("============ application-golang starts ============")

	// The gRPC client connection should be shared by all Gateway connections to this endpoint
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	defer gateway.Close()

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)


	fmt.Println("mint:")
	mint(contract)
	
	//fmt.Println("owner:")
	//owner(contract)
	
	//fmt.Println("transfer:")
	//transfer(contract)
	
	
	//fmt.Println("owner:")
	//owner(contract)
	
	//fmt.Println("burn:")
	//burn(contract)
	
	//fmt.Println("owner:")
	//owner(contract)


	log.Println("============ application-golang ends ============")
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func newSign() identity.Sign {
	files, err := ioutil.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := ioutil.ReadFile(path.Join(keyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}


func mint(contract *client.Contract) {
	fmt.Printf("Submit Transaction: mint, creates new erc721 with tokenID, tokenURI \n")
	contract.SubmitTransaction("SetBaseURI","http://www.try1try.com")
	ID, err := contract.SubmitTransaction("PublicMint")
	if err != nil {
		panic(fmt.Errorf("failed to SubmitTransaction: %w", err))
	}
	fmt.Printf(string(ID))
}


func owner(contract *client.Contract) {
	fmt.Printf("evaluate Transaction: ownerOf \n")

	ownerer, err := contract.EvaluateTransaction("OwnerOf", "103")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := ownerer
	fmt.Printf("*** Result:%s\n", result)
	
}

func transfer(contract *client.Contract) {
	fmt.Printf("Submit Transaction: transfer \n")

	_, err := contract.SubmitTransaction("TransferFrom", "x509::CN=user1,OU=client,O=Hyperledger,ST=North Carolina,C=US::CN=ca.org2.example.com,O=org2.example.com,L=Hursley,ST=Hampshire,C=UK", "x509::CN=user1,OU=client,O=Hyperledger,ST=North Carolina,C=US::CN=ca.org1.example.com,O=org1.example.com,L=Durham,ST=North Carolina,C=US","102")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}

func GetChaincodeOwner(contract *client.Contract) {
	fmt.Printf("evaluate Transaction: GetChaincodeOwner \n")

	chainCodeOwner,err := contract.EvaluateTransaction("GetChainCodeOwner")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	result := chainCodeOwner
	fmt.Printf("*** Result:%s\n", result)
	fmt.Printf("*** Transaction evaluated successfully\n")
}

func SetChaincodeOwner(contract *client.Contract) {
	fmt.Printf("Submit Transaction: SetChaincodeOwner \n")
	chainCodeOwner, err := contract.SubmitTransaction("SetChainCodeOwner","x509")
	if err != nil {
		panic(fmt.Errorf("failed to SubmitTransaction: %w", err))
	}
	fmt.Printf(string(chainCodeOwner))
}


func burn(contract *client.Contract) {
	fmt.Printf("Submit Transaction: burn \n")

	_, err := contract.SubmitTransaction("Burn", "102")
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Printf("*** Transaction committed successfully\n")
}


//Format JSON data
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, " ", ""); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return prettyJSON.String()
}
