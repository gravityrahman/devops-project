# DevOps Technical Challenge

Below you'll find multiple user stories to complete.
Please complete as many user stories as feel comfortable with.


## Stories
### Make Deploys Faster
Our current CI works; however, we currently have to wait for the Go application to build regardless of if it was updated.

Given: The current Go application and pipeline  
When: The Go app gets no updates  
AND  
When: Changes have been made to Terraform  
Then: The app stages shall not run  
AND  
Then: The Terraform stages shall run

### Test the Code Before Deploy
Our code currently has tests that can be run locally. However, the pipeline doesn't run these test.
As such broken code can go to production. Make the pipeline test the code.

Given: The current Gitlab pipeline and Go application  
When: The tests in the Go app fail  
Then: The pipeline should not deploy the code

### Enable DynamoDB
Our current Go application can retrieve a value from DynamoDB and return it to the client.
Make the needed changes so that both local development and production work.

Given: The current Go application  
When: A requests is made to `/dynamo`  
Then: The request should return `This was fun!`

### Fix the Bug
Our code currently fails the test (see above). 
Since we now have a working pipeline we need to fix the bug.

Review the documentation for the project and fix the failing test appropriately.

## Help
If you have any questions or need any assistance please reach out.
