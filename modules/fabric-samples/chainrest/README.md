# ChainREST
## RUN
npm install

DEBUG=chainapi:* npm start

## Functionality and connection to the blockchain from REST

We use a third party module called `` fabric-network-simple: 1.1.0 `` which uses the official library of hyperledger fabric for JS to connect to a blockchain.

We need to specify the blockchain topology and configuration(channel name, smart contract name, number of organizations, etc.) [line 6-routes/chain.js](https://github.com/agustin-marin/ChainREST/blob/main/routes/chain.js#L6)

We can use queryChaincode(readonly smartcontract operation) and invokeChaincode(read-write smartcontract operation) specifying method and parameters: [queryChaincode](https://github.com/agustin-marin/ChainREST/blob/main/routes/chain.js#L98) and [invokeChaincode](https://github.com/agustin-marin/ChainREST/blob/main/routes/chain.js#L220), in order to execute smartcontracts in the blockchain.

## Smart Contract Example
On this [github](https://github.com/hyperledger/fabric-samples/tree/main/chaincode) we have several examples of chaincode written in different languages.

We can see in the [FabCar example](https://github.com/hyperledger/fabric-samples/blob/main/chaincode/fabcar/java/src/main/java/org/hyperledger/fabric/samples/fabcar/FabCar.java), on the java project, a simple example of a java class definning a smartcontract using Labels(@Contract, @Default and @Transaction) and the interface ContractInterface. 

In the method [createCar](https://github.com/hyperledger/fabric-samples/blob/main/chaincode/fabcar/java/src/main/java/org/hyperledger/fabric/samples/fabcar/FabCar.java#L118) we can see an example of ledger request using the stub instance. The stub instance is used to make requests to the ledger(get, put, query(using couchDB)). Remembering that Hyperledger fabric uses key-value json form to save data in the ledger, the method createCar uses stub.getStringState to get the value associated to 'key'. If it is not empty, the car already exists and it gives an error. Else, it creates a new instance of car and uses stub.putStringState(key, value) to store the new car on the ledger.
