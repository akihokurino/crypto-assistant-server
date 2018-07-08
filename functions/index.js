// gcloud beta functions deploy calcAssets --trigger-resource calc-assets --trigger-event google.pubsub.topic.publish
const calcAssets = require('./calc_assets');

exports.calcAssets = calcAssets;

