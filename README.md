# LeetSignal
An AWS Lambda function that sends a notification whenever a user in the group solves a new Leetcode question. Use it for interview group-study motivation!

### Prerequisites
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html)
- [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)

### Configuration
The following configuration fields must be set in a ```cmd/lambda/config.json``` file. A template config is provided as ```config.example.json```.
- ```profiles```: list of leetcode IDs
- ```ntfy_topic```: [ntfy](https://docs.ntfy.sh/) notificaton topic. All ntfy topics are public, so choose something that is not easily guessable

### Deploy
1. **Build the project**
   ```
   $ make
   ```
2. **Deploy with SAM**
  - First time deployment
    ```
    $ sam deploy --guided
    ```
  - Subsequent deployments
    ```
    $ sam deploy
    ```

### Usage
- Install the [ntfy](https://ntfy.sh/) app (iOS/Android/Desktop).
- Subscribe to the ```ntfy_topic``` configured in config.json

## Acceptable Use Policy
LeetSignal is intended only for collaborative, group-study, and motivational purposes.
- You may only add LeetCode usernames to your configuration if you have the explicit consent of the individuals being tracked.
- Using this project to monitor users without their knowledge or permission is **STRICTLY** prohibited.
- By deploying or using LeetSignal, you agree to use it responsibly and in compliance with this policy.

## Other notes
API data is provided by [noworneverev/leetcode-api](https://github.com/noworneverev/leetcode-api)

LeetSignal uses the following AWS services and operates well under the free-tier limits for each service:
- [Lambda](https://aws.amazon.com/lambda/): runs hourly
- Eventbridge: schedules/triggers lambda function
- CloudFormation: infrastructure as code service
- Cloudwatch: minimal logging
- AWS Parameter Store: storing number of solved Leetcodes between Lambda invocations

