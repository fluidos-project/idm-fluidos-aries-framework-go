const { Wallets, Gateway } = require("fabric-network");

async function initGatewayOptions(config) {
  const wallet = await Wallets.newInMemoryWallet();
  const x509Identity = {
    credentials: {
      certificate: config.identity.certificate,
      privateKey: config.identity.privateKey,
    },
    mspId: config.identity.mspid,
    type: "X.509",
  };
  await wallet.put(config.identity.mspid, x509Identity);
  const gatewayOptions = {
    identity: config.identity.mspid,
    wallet,
    discovery: {
      enabled: config.settings.enableDiscovery,
      asLocalhost: config.settings.asLocalhost,
    },
  };
  return gatewayOptions;
}

async function initGateway(config, gatewayOptions) {
  const gateway = new Gateway();
  const currentDate = new Date();
  const timestamp = currentDate.getTime();
  config.connectionProfile["name"] = "umu.fabric." + timestamp;
  config.connectionProfile["version"] = "1.0.0" + timestamp;
  await gateway.connect(config.connectionProfile, gatewayOptions);
  const network = await gateway.getNetwork(config.channelName);
  const contract = network.getContract(config.contractName);
  return contract;
}

async function queryChaincode(contract, transaction, args) {
  try {
    const queryResult = await contract.submitTransaction(transaction, ...args);
    let result = "[]";
    if (queryResult) {
      result = queryResult.toString();
    }
    return result;
  } catch (error) {
    console.error(
      `Failed to query transaction: "${transaction}" with arguments: "${args}", error: "${error}"` + error.toString()
    );
    throw error;
  }
}

module.exports = { initGatewayOptions, initGateway, queryChaincode };
