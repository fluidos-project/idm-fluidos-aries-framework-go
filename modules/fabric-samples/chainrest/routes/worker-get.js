const { parentPort } = require("worker_threads");
const conf = require("./config.json");
const { initGatewayOptions, initGateway, queryChaincode } = require("./fabric-utils");

let contract;

parentPort.once("message", async (reqString) => {
      console.log("execution");
    let req = JSON.parse(reqString);
    let queryParams = req.query;
  let dateFrom = queryParams.dateFrom;
      if (!dateFrom) {
        throw new Error('Se necesita al menos dateFrom');
    }
    let dateTo = queryParams.dateTo || '';
  const gatewayOptions = await initGatewayOptions(conf);
  contract = await initGateway(conf, gatewayOptions);
  const queryResult = await queryChaincode(contract, "get", [dateFrom, dateTo]); //
  console.log("El tipo de objeto de queryResult es:", typeof queryResult); // Verificar el tipo de objeto de queryResult
  if (typeof queryResult === "string") {
    console.log("El contenido de queryResult es:", queryResult); // Imprimir el contenido de queryResult si es una cadena de texto
  }
  parentPort.postMessage(JSON.stringify(queryResult));
});

