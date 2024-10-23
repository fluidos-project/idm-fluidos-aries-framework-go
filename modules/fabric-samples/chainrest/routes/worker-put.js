const { parentPort } = require("worker_threads");
const conf = require("./config.json");
const { initGatewayOptions, initGateway, queryChaincode } = require("./fabric-utils");

let contract;

parentPort.once("message", async (reqString) => {
  let req = JSON.parse(reqString);
  const gatewayOptions = await initGatewayOptions(conf);
  contract = await initGateway(conf, gatewayOptions);
  const body = req.body;
  await queryChaincode(contract, "put", [JSON.stringify(body)]);
  parentPort.postMessage({});
});
