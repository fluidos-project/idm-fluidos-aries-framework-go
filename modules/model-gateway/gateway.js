const { Wallets, Gateway } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const fs = require('fs');
const path = require('path');

const connectionProfilePath = path.resolve(__dirname, 'connection-profile.json');
const connectionProfileContent = JSON.parse(fs.readFileSync(connectionProfilePath, 'utf8'));
const connectionProfile = connectionProfileContent['connection-profile'];

// Check if the peers property exists and has the expected structure
if (connectionProfile && connectionProfile.peers && connectionProfile.peers['peer0.org1.example.com']) {
    connectionProfile.peers['peer0.org1.example.com'].url = 'grpcs://peer0.org1.example.com:7051';
} else {
    console.warn("Warning: Unable to set peer URL. Check the connection profile structure.");
}

const channelName = 'mychannel';
const chaincodeName = 'model-treatment';

async function enrollAdmin(connectionProfile, caName = 'ca.org1.example.com', mspId = 'Org1MSP') {
    try {
        const caInfo = connectionProfile.certificateAuthorities[caName];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        const identity = await wallet.get('admin');
        if (identity) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: mspId,
            type: 'X.509',
        };
        await wallet.put('admin', x509Identity);
        console.log('Successfully enrolled admin user "admin" and imported it into the wallet');
    } catch (error) {
        console.error(`Failed to enroll admin user "admin": ${error}`);
        throw error;
    }
}

async function enrollUser(connectionProfile, username, password, caName = 'ca.org1.example.com', mspId = 'Org1MSP') {
    try {
        const caInfo = connectionProfile.certificateAuthorities[caName];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);

        const userIdentity = await wallet.get(username);
        if (userIdentity) {
            console.log(`An identity for the user ${username} already exists in the wallet`);
            return;
        }

        const enrollment = await ca.enroll({ enrollmentID: username, enrollmentSecret: password });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: mspId,
            type: 'X.509',
        };
        await wallet.put(username, x509Identity);
        console.log(`Successfully enrolled user ${username} and imported it into the wallet`);
    } catch (error) {
        console.error(`Failed to enroll user ${username}: ${error}`);
        throw error;
    }
}

async function connectToNetwork(connectionProfile, channelName, chaincodeName) {
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    const gateway = new Gateway();

    try {
        const identity = await wallet.get('admin');
        if (!identity) {
            throw new Error('Admin identity not found in the wallet');
        }

        await gateway.connect(connectionProfile, {
            wallet,
            identity: 'admin',
            discovery: { enabled: true, asLocalhost: true }
        });

        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        return { gateway, contract };
    } catch (error) {
        console.error(`Failed to connect to the network: ${error}`);
        throw error;
    }
}

async function getAdminIdentity() {
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    
    const adminIdentity = await wallet.get('admin');
    if (!adminIdentity) {
        throw new Error('Admin identity not found in the wallet');
    }
    
    return adminIdentity;
}

module.exports = {
    enrollAdmin,
    enrollUser,
    connectToNetwork,
    getAdminIdentity,
};
