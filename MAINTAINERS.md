# Maintainers

## Travis-CI

We use Travis for generating our executables automatically. To successfully build it, we need to make sure these environment variables are set there:

```
export TRAVIX_FIREBASE_API_KEY=""
export TRAVIX_FIREBASE_AUTH_DOMAIN=""
export TRAVIX_FIREBASE_DATABASE_URL=""
export TRAVIX_FIREBASE_STORAGE_BUCKET=""
export TRAVIX_FIREBASE_MESSAGING_SENDER_ID=""

export TRAVIX_DEVELOPER_PROFILE_URL=""
```

You can set them here: [https://travis-ci.org/Travix-International/Travix.Core.Adk/settings](https://travis-ci.org/Travix-International/Travix.Core.Adk/settings)
