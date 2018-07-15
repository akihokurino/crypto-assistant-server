const etherscan = require('etherscan-api').init('FQFUE4MMQW4BIWC2I31176GCQN5ZMJH319');
const Units = require('ethereumjs-units');

const calcEther = (address) => {
  const addressText = address.Value;
  const balance = etherscan.account.balance(addressText);
  const txList = etherscan.account.txlist(addressText);

  return Promise.all([balance, txList])
  .then((results) => {
    const amount = Units.convert(results[0].result, 'wei', 'eth');
    const transaction = results[1].result;

    return {amount: amount, transaction: transaction}
  });
};

module.exports = {
  calcEther: calcEther,
};