// A helper script to quickly generate the necessary env variables.
// Copy here the Json configuration from the Firebase console, and run this script in a REPL.
var config = {
    apiKey: "DummyApiKey",
    authDomain: "fireball-development.firebaseapp.com",
    databaseURL: "https://fireball-development.firebaseio.com",
    projectId: "fireball-development",
    storageBucket: "fireball-development.appspot.com",
    messagingSenderId: "445151709427",
    certContent: "DummyCertificateContent",
    keyContent: "DummyKeyContent",
    uploadUrl: "https://us-central1-travix-development.cloudfunctions.net/http"
  };

console.log(
`For PowerShell:
$env:TRAVIX_CERT_CONTENT='${config.certContent}'
$env:TRAVIX_KEY_CONTENT='${config.keyContent}'
$env:TRAVIX_FIREBASE_API_KEY='${config.apiKey}'
$env:TRAVIX_FIREBASE_DATABASE_URL='${config.databaseURL}'
$env:TRAVIX_FIREBASE_STORAGE_BUCKET='${config.storageBucket}'
$env:TRAVIX_FIREBASE_MESSAGING_SENDER_ID='${config.messagingSenderId}'
$env:TRAVIX_FIREBASE_AUTH_DOMAIN='${config.authDomain}'
$env:TRAVIX_FIREBASE_REFRESH_TOKEN_URL='https://securetoken.googleapis.com/v1/token?key='
$env:TRAVIX_DEVELOPER_PROFILE_URL='https://developerprofile.${config.projectId.endsWith("development") ? "development." : config.projectId.endsWith("staging") ? "staging." : ""}travix.com/'
$env:TRAVIX_LOGGER_URL='https://frogger.travix.com/'
$env:TRAVIX_UPLOAD_URL='${config.uploadUrl}'

For bash:
TRAVIX_CERT_CONTENT='${config.certContent}'
TRAVIX_KEY_CONTENT='${config.keyContent}'
TRAVIX_FIREBASE_API_KEY='${config.apiKey}'
TRAVIX_FIREBASE_DATABASE_URL='${config.databaseURL}'
TRAVIX_FIREBASE_STORAGE_BUCKET='${config.storageBucket}'
TRAVIX_FIREBASE_MESSAGING_SENDER_ID='${config.messagingSenderId}'
TRAVIX_FIREBASE_AUTH_DOMAIN='${config.authDomain}'
TRAVIX_FIREBASE_REFRESH_TOKEN_URL='https://securetoken.googleapis.com/v1/token?key='
TRAVIX_DEVELOPER_PROFILE_URL='https://developerprofile.${config.projectId.endsWith("development") ? "development." : config.projectId.endsWith("staging") ? "staging." : ""}travix.com/'
TRAVIX_LOGGER_URL='https://frogger.travix.com/'
TRAVIX_UPLOAD_URL='${config.uploadUrl}'`);
