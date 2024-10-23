const express = require("express");
const { Worker } = require("worker_threads");
const router = express.Router();
require("json-circular-stringify");
const winston = require('winston');
const log = winston.createLogger({
    level: 'error',
    format: winston.format.json(),
    defaultMeta: { service: 'user-service' },
    transports: [
        new winston.transports.File({ filename: 'error.log', level: 'error' })
    ]
});

router.get("/json", async (req, res) => {
  try {
    const worker = new Worker("./routes/worker-get.js");
    worker.on("message", (result) => {
      if (result === "[]") {
        res.status(404).send({ message: "No data found for the requested key." });
      } else {
        res.status(200).send(JSON.parse(result));
      }
    });
    worker.postMessage(JSON.stringify(req));
  } catch (error) {
    log.error("Error occurred in /json GET request:", error);
    sendError(res, error);
  }
});

router.post("/json", async (req, res) => {
  try {
    const worker = new Worker("./routes/worker-put.js");
    worker.on("message", (result) => {
      if (Object.keys(result).length === 0) {
        res.status(201).send({ message: "Data successfully added to the ledger." });
      } else {
        res.status(400).send({ message: "Failed to add data to the ledger.", error: result });
      }
    });
    worker.postMessage(JSON.stringify(req));
  } catch (error) {
    log.error("Error occurred in /json POST request:", error);
    sendError(res, error);
  }
});

function sendError(res, error) {
  log.error("An internal server error occurred:", error);
  res.status(500).send({ message: "An internal server error occurred.", error: error.toString() });
}

module.exports = router;
