# Policy Inference Decider (PID)

Microservice responsible for policy inference.

## Load tests

Artillery + AWS credentials. From project root:

```bash
npm install -g artillery
./loadtest/run-lambda.sh "https://YOUR_LAMBDA_FUNCTION_URL" [test-file.yaml]
```

Example (50 RPS): `./loadtest/run-lambda.sh "https://...lambda-url.us-east-1.on.aws" test-50rps.yaml`  
Example (100 RPS mixed): `./loadtest/run-lambda.sh "https://...lambda-url.us-east-1.on.aws" test-100rps-mixed.yaml`