const DataStore = require('@google-cloud/datastore');
const request = require('request');
const ethereum = require('./blockchain/ethereum');

const projectId = 'crypto-assistant-dev';

const dataStore = new DataStore({
  projectId: projectId,
});

const calcAssets = (event, callback) => {
  const message = event.data;
  const addressIds = Buffer.from(message.data, 'base64').toString().split(",");

  Promise.all(addressIds.map((addressId) => {
    return calcAsset(addressId);
  }))
  .then(() => {
    return calcPortfolioWebhook();
  })
  .then(() => {
    callback();
  })
  .catch(err => {
    console.error('ERROR:', err);
    callback();
  });
};

const calcPortfolioWebhook = () => {
  return new Promise((resolve, reject) => {
    request.get('https://batch-dot-crypto-assistant-dev.appspot.com/job/broadcast_portfolio')
      .on('response', (response) => {
        if (response.status < 300) {
          resolve();
        } else {
          reject();
        }
      });
  });
};

const calcAsset = (addressId) => {
  const key = dataStore.key(['Address', addressId]);

  return dataStore.get(key)
  .then((results) => {
    if (results.length === 0) {
      return;
    }

    const address = results[0];

    // TODO: 通貨によってアクセスするBlockchainを決める
    return ethereum.calcEther(address).then((result) => {
      return Promise.all([
        updateAsset(address.UserId, addressId, result.amount),
        updateTransaction(address.UserId, addressId, result.transaction)
      ]);
    });
  });
};

const updateAsset = (userId, addressId, amount) => {
  console.log("update asset");
  console.log(amount);
  const key = dataStore.key(['Asset', userId + "-" + addressId]);
  const Asset = {
    key: key,
    data: {
      UserId: userId,
      AddressId: addressId,
      Amount: parseFloat(amount)
    },
  };
  return dataStore.save(Asset);
};

const updateTransaction = (userId, addressId, text) => {
  console.log("update transaction");
  console.log(text);
  const key = dataStore.key(['Transaction', userId + "-" + addressId]);
  // TODO: 暗号化
  const Transaction = {
    key: key,
    data: {
      UserId: userId,
      AddressId: addressId,
      Text: text
    },
  };
  return dataStore.save(Transaction);
};

module.exports = calcAssets;