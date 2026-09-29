# Deploying to AWS App Runner

This ships a Dockerfile-based deploy via [AWS App Runner](https://aws.amazon.com/apprunner/), pulling the image
from ECR. Fill in the placeholders in `apprunner-service.json` (`<ACCOUNT_ID>`, `<REGION>`, `<ECR_ACCESS_ROLE_ARN>`,
`<SECRETS_MANAGER_ARN:...>`) before running these.

1. Create the ECR repository and push the image:

   ```bash
   aws ecr create-repository --repository-name {{ cookiecutter.project_slug }}

   aws ecr get-login-password --region <REGION> | \
     docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com

   docker build -t {{ cookiecutter.project_slug }} .
   docker tag {{ cookiecutter.project_slug }}:latest <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/{{ cookiecutter.project_slug }}:latest
   docker push <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/{{ cookiecutter.project_slug }}:latest
   ```

2. Store the secrets referenced in `apprunner-service.json` (`DB_SOURCE`, `TOKEN_SYMMETRIC_KEY`, `SESSION_SECRET`) in
   AWS Secrets Manager, and create an IAM role App Runner can use to pull from ECR (`AccessRoleArn`).

3. Create the service:

   ```bash
   aws apprunner create-service --cli-input-json file://deploy/aws/apprunner-service.json
   ```

Database: provision Postgres separately (RDS is the standard choice) and put its connection string in
`DB_SOURCE` via Secrets Manager — App Runner has no built-in database like Railway/DigitalOcean App Platform do.
