# Decentralised Lambda Node

Goal of this project is to provide an interface for invoking WASM serverless functions on decentralised networks.

Currently, only Bacalhau network is supported but the aim is to support other networks in future as well.

## Running

```bash
go run cmd/main.go
```

## Example usage

In this example we will run a lambda function that returns a greeting message. This function is written using [decentralised-lambda-runtime](https://github.com/pyropy/decentralised-lambda-runtime). Here's the source code:
```rust
use decentralised_lambda_runtime::{Error, LambdaEvent};
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

This function is compiled to WASM and then uploaded to the IPFS network. The CID of the WASM file is then used to invoke the function.

```bash
# Invoke a function
curl --request POST \
  --url http://localhost:8080/invoke/QmY4t7ih7fqGtwdxRfrYPFWSwZrxoQaH38EUMLFEDe2Y6S \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Lambda"
}'
```