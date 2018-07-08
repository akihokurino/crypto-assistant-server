const DataStore = require('@google-cloud/datastore');
const etherscan = require('etherscan-api').init('FQFUE4MMQW4BIWC2I31176GCQN5ZMJH319');
const projectId = 'crypto-assistant-dev';
const Units = require('ethereumjs-units');

const dataStore = new DataStore({
  projectId: projectId,
});

const calcAssets = (event, callback) => {
  const message = event.data;
  const addressId = Buffer.from(message.data, 'base64').toString();

  const key = dataStore.key(['Address', addressId]);

  let address;

  dataStore
  .get(key)
  .then((results) => {
    if (results.length === 0) {
      return;
    }

    address = results[0];
    const addressText = results[0].Value;
    const balance = etherscan.account.balance(addressText);
    const txList = etherscan.account.txlist(addressText);

    return Promise.all([balance, txList]);
  })
  .then((results) => {
    const amount = Units.convert(results[0].result, 'wei', 'eth');
    const transaction = results[1].result;

    return Promise.all([
      updateAsset(address.UserId, addressId, amount),
      updateTransaction(address.UserId, addressId, transaction)
    ]);
  })
  .then(() => {
    callback();
  })
  .catch(err => {
    console.error('ERROR:', err);
    callback();
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

const calcAssetsLocal = (address) => {
  const balance = etherscan.account.balance(address);
  const txList = etherscan.account.txlist(address);

  Promise.all([balance, txList])
  .then((results) => {
    const amount = Units.convert(results[0].result, 'wei', 'eth');
    console.log(amount);
    console.log(results[1].result);
  })
  .catch(err => {
    console.error('ERROR:', err);
  });
};

module.exports = calcAssets;