# NeatoBotLib

This is an unofficial API client which can help you to interact with the Neato cloudservices which are used to control you Neato Connected vacuum robot.

Thanks to [Lars Brillert @kangguru](https://github.com/kangguru) who reverse engineered the Neato API from which this library is ported from. Port is based on https://github.com/kangguru/botvac

## Usage
Check the examples to get a hint on how to use the library, most is self explanatory.

Currently the following methods are available in the NeatoBotLib class

* Auth
* GetDashboard
* GetRobotState

The method names should give you an idea what the specific action will cause. Still this is not all, but that's what is available for the moment.

## OAuth2
Make a token request using the following curl command:

```
curl -X "POST" "https://mykobold.eu.auth0.com/passwordless/start" \
     -H 'Content-Type: application/json' \
     -d $'{
  "send": "code",
  "email": "ENTER_YOUR_EMAIL_HERE",
  "client_id": "KY4YbVAvtgB7lp8vIbWQ7zLk3hssZlhR",
  "connection": "email"
}'
```

After receiving the OTP:
```
curl -X "POST" "https://mykobold.eu.auth0.com/oauth/token" \
     -H 'Content-Type: application/json' \
     -d $'{
  "prompt": "login",
  "grant_type": "http://auth0.com/oauth/grant-type/passwordless/otp",
  "scope": "openid email profile read:current_user",
  "locale": "en",
  "otp": "ENTER_OTP_HERE",
  "source": "vorwerk_auth0",
  "platform": "ios",
  "audience": "https://mykobold.eu.auth0.com/userinfo",
  "username": "ENTER_YOUR_EMAIL_HERE",
  "client_id": "KY4YbVAvtgB7lp8vIbWQ7zLk3hssZlhR",
  "realm": "email",
  "country_code": "DE"
}'
```

From the output, the id_token value is the token you need.
(OAuth OTP original from https://github.com/nicoh88/node-kobold/)
