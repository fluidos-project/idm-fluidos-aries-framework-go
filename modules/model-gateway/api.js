const express = require('express');
const fs = require('fs');
const path = require('path');
const { enrollAdmin, connectToNetwork, getAdminIdentity } = require('./gateway');

const app = express();
app.use(express.json());

const connectionProfilePath = path.resolve(__dirname, 'connection-profile.json');
const connectionProfileContent = JSON.parse(fs.readFileSync(connectionProfilePath, 'utf8'));
const connectionProfile = connectionProfileContent['connection-profile'];

const channelName = 'mychannel';
const chaincodeName = 'model-treatment';

app.post('/enroll-admin', async (req, res) => {
    try {
        await enrollAdmin(connectionProfile);
        res.json({ success: true, message: 'Admin enrolled successfully' });
    } catch (error) {
        console.error('Failed to enroll admin:', error);
        res.status(500).json({ success: false, error: error.message });
    }
});

app.post('/write-dht', async (req, res) => {
    try {
        const { gateway, contract } = await connectToNetwork(connectionProfile, channelName, chaincodeName);
        const { key, value } = req.body;
        await contract.submitTransaction('writeDHT', key, value);
        await gateway.disconnect();
        res.json({ success: true, message: 'Data written to DHT' });
    } catch (error) {
        console.error(`Failed to write to DHT: ${error}`);
        res.status(500).json({ success: false, error: error.message });
    }
});

app.get('/read-dht/:key', async (req, res) => {
    try {
        const { gateway, contract } = await connectToNetwork(connectionProfile, channelName, chaincodeName);
        const result = await contract.evaluateTransaction('readDHT', req.params.key);
        await gateway.disconnect();
        res.json({ success: true, value: result.toString() });
    } catch (error) {
        console.error(`Failed to read from DHT: ${error}`);
        res.status(500).json({ success: false, error: error.message });
    }
});

app.get('/admin-identity', async (req, res) => {
    try {
        const adminIdentity = await getAdminIdentity();
        res.json({ success: true, adminIdentity });
    } catch (error) {
        console.error(`Error retrieving admin identity: ${error}`);
        res.status(500).json({ success: false, error: error.message });
    }
});

const port = 3000;
app.listen(port, () => {
    console.log(`Model gateway API listening at http://localhost:${port}`);
});

module.exports = {
    app
};
