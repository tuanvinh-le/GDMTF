
// Chaincode

type SimpleChaincode struct {
}

func pow(b string) int {
	var hashInt big.Int
	var hash [32]byte

	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	nonce := 0

	// t1 := time.Now() // get current time
	for nonce < maxNonce {
		data := prepareData(nonce, b)

		hash = sha256.Sum256(data)
		// fmt.Printf("#%d = %x\r", nonce, hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			// elapsed := time.Since(t1)
			// fmt.Print("\nApp elapsed: ", elapsed)
			break
		} else {
			nonce++
		}
	}

	return nonce
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("cgublock Init")
	return shim.Success(nil)
}

func (t *SimpleChaincode) lsssEn(m string) string {
	zero := big.NewInt(0)

	loadParams, _ := ioutil.ReadFile("./AttributeParams.load")
	pairing, _ := pbc.NewPairingFromString(string(loadParams[:]))

	//g Y b
	attrPara, _ := ioutil.ReadFile("./AttributeServer.load")

	// 0 g 1 pubkey 2 B 3...N Attri
	pbcPara := strings.Split(string(attrPara), "\n")

	g, _ := pairing.NewG1().SetString(pbcPara[0], 10)
	pubKey, _ := pairing.NewGT().SetString(pbcPara[1], 10)
	b, _ := pairing.NewG1().SetString(pbcPara[2], 10)

	q := make(map[string]*pbc.Element)
	for i := 3; i < len(pbcPara); i++ {
		attris := strings.Split(pbcPara[i], ":")
		q[attris[0]], _ = pairing.NewG1().SetString(attris[1], 10)
	}

	prefix := lsss.InfixToPrefix(serverPbConfig.ServerPolicy)
	// fmt.Println("Access Policy:", prefix)
	attrField := lsss.AccessTree(prefix)

	secrets := make(map[int]*big.Int)
	for _, key := range attrField {
		i := 0
		for range key {
			secrets[i], _ = rand.Prime(rand.Reader, 64)
			i++
		}
		break
	}

	vectors := make(map[string]*big.Int)
	for id, key := range attrField {
		vectors[id] = big.NewInt(0)
		i := 0
		for _, value := range key {
			temp := big.NewInt(int64(value))
			temp.Mul(temp, secrets[i])
			vectors[id] = vectors[id].Add(vectors[id], temp)
			i++
		}
	}

	qr := make(map[string]*pbc.Element)
	nqr := make(map[string]*pbc.Element)
	d := make(map[string]*pbc.Element)
	r := make(map[string]*big.Int)
	c := make(map[string]*pbc.Element)
	convertC := make(map[string]string)
	convertD := make(map[string]string)
	for id, _ := range attrField {
		r[id], _ = rand.Prime(rand.Reader, 36)
		qr[id] = pairing.NewG1().PowBig(q[id], r[id])
		nqr[id] = pairing.NewG1().Invert(qr[id])

		d[id] = pairing.NewG1().Set0()
		d[id].Add(d[id], g)
		d[id].PowBig(d[id], r[id])
		convertD[id] = d[id].String()

		c[id] = pairing.NewG1().Set0()
		c[id].Add(c[id], b)
		c[id].MulBig(c[id], vectors[id])
		if zero.Cmp(vectors[id]) == 1 {
			c[id].Invert(c[id])
		}
		c[id].Mul(c[id], nqr[id])
		convertC[id] = c[id].String()
	}

	privKey := pairing.NewGT().PowBig(pubKey, secrets[0])

	data := m2n(m)

	// hash
	h := pairing.NewG1().SetFromStringHash(m, sha256.New())

	upSig := new(big.Int).Mul(pubKey.X(), secrets[0])
	sig := pairing.NewG1().MulBig(h, upSig)

	// ABSE encrypt
	ciper := data.Mul(data, privKey.X())

	gs := pairing.NewG1().Set0()
	gs.Add(gs, g)
	gs.PowBig(gs, secrets[0])

	convertCjson, _ := json.Marshal(convertC)
	convertDjson, _ := json.Marshal(convertD)

	tm := time.Unix(time.Now().Unix(), 0)

	ct := tm.Format("2006-01-02 15:04:05") + "\n" + serverPbConfig.ServerIP + "\n" + ciper.String() + "\n" + sig.String() + "\n" + gs.String() + "\n" + string(convertCjson) + "\n" + string(convertDjson) + "\n" + serverPbConfig.ServerPolicy

	return base64.StdEncoding.EncodeToString([]byte(ct))
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("cgublock Invoke")

	function, args := stub.GetFunctionAndParameters()
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "add" {
		// the old "Query" is now implemtned in invoke
		return t.add(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"delete\" \"query\" \"add\"")
}

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("cgublock delete")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(ID)

	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) add(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("cgublock add")

	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Write the state back to the ledger
	err = stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("cgublock query")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ID := args[0]

	// Delete the key from the state in ledger
	fileblock, err := stub.GetState(ID)
	if err != nil {
		return shim.Error("Could not fetch application with id")
	}

	if fileblock == nil {
		jsonResp := "{\"Error\":\"Nil data for " + ID + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Printf("Query Response:\n")

	return shim.Success(fileblock)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


// Public blocks (batch signatures) generation

sigall := pairing.NewG1().Set1()

 for i:=0 ; i<=3 ; i++{ 
 signature[i] := pairing.NewG1().SetBytes(ct.signature[i])

  h[i] := pairing.NewG1().SetFromStringHash(message[i], sha256.New())

sigall .Mul(sigall, signature[i])
}


// Smart contract

client, err := ethclient.Dial("https://ropsten.infura.io/v3/a565d0dc884c476fa2e25636ea19fa82")
  if err != nil {
    log.Fatal(err)
  }

  privateKey, err := crypto.HexToECDSA("5F03F06E2B524F4D8FF6135967899992B6B609F8A37B4D0015A1C0154E1A4FDB")
  if err != nil {
    log.Fatal(err)
  }

  publicKey := privateKey.Public()
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
    log.Fatal("error casting public key to ECDSA")
  }

  fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
  nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
  if err != nil {
    log.Fatal(err)
  }
 value := new(big.Int)
  value.SetString("30000000000000000", 10) // in wei (0.3 eth)
  gasLimit := uint64(40000)                 // in units

  gasPrice, err := client.SuggestGasPrice(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  toAddress := common.HexToAddress("0xBa7adA49BffDc8c641D1cB8f3f9aF29F9BD9C66e")
  data := []byte(PB)
  tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
  signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
  if err != nil {
    log.Fatal(err)
  }
  err = client.SendTransaction(context.Background(), signedTx)
  if err != nil {
    log.Fatal(err)
  }
