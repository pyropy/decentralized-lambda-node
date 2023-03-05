# Decentralized Lambda Node

Goal of this project is to provide an interface for invoking WASM serverless functions on decentralized networks.

Currently, only Bacalhau network is supported but the aim is to support other networks in future as well.

## Running

To start the node run:

```bash
go run cmd/main.go node start
```

## Example usage

In this example we will run a lambda function that returns a greeting message. This function is written
using [decentralized-lambda-runtime](https://github.com/pyropy/decentralized-lambda-runtime). Here's the source code:

```rust
use decentralized_lambda_runtime::{Error, LambdaEvent};
use serde::{Deserialize, Serialize};


#[derive(Deserialize)]
struct Request {
    name: String,
}


#[derive(Serialize)]
struct Response {
    msg: String,
}

fn main() -> Result<(), Error> {
    decentralised_lambda_runtime::run(my_handler)?;
    Ok(())
}

pub(crate) fn my_handler(event: LambdaEvent<Request>) -> Result<Response, Error> {
    let command = event.payload.name;

    let resp = Response {
        msg: format!("Hello, {}", command),
    };

    Ok(resp)
}
```

This function is compiled to WASM and then uploaded to the IPFS networka as part of the job spec. To upload

```bash
go run cmd/main.go wasm deploy --file <path to your wasm file>
```

This should produce CID result logged in terminal, which is in fact job invocation spec CID which contains CID to the binary. CID contents should look something like this:

```json
{
  "binary": { "/": "QmdoUoHi31JJb72Y2PUerzk1aYvtmcrrx9cVRrgMK9DqCd" },
  "executionLayer": "bacalhau",
  "persistanceLayer": "ipfs"
}
```

The CID of the job invocation spec is then used to invoke the function.

```bash
# Invoke a function
curl --request POST \
  --url http://localhost:8080/invoke/QmQYEqnbCCMshBCG68ou56jW7LBscBWYd6fo13EDDCAFU4 \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Lambda"
}'
```

And the response is:

```json
{
  "result": {
    "msg": "Hello, Lambda"
  }
}
```
