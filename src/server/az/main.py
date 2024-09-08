import os
import uvicorn

from azure.identity import DefaultAzureCredential
from azure.mgmt.resource import ResourceManagementClient
from enum import Enum
from fastapi import FastAPI

#local
from library.resource_groups import resource_group_create_or_update, resource_group_show

class testEndpoints(str, Enum):
    me         = "me"
    connection = "connection"
    # all        = "all"

app = FastAPI()

# Acquire a credential object using DevaultAzureCredential.
# credential = DefaultAzureCredential()

# Retrieve subscription ID from environment variable.
# subscription_id = "xxx" #os.environ["AZURE_SUBSCRIPTION_ID"]

# Obtain the management object for resources.
# resource_client = ResourceManagementClient(credential, subscription_id, "2022-09-01")

@app.get("/test/{test_endpoint}")
def tests(test_endpoint: testEndpoints):
    match test_endpoint:
        case testEndpoints.me:
            return {"result":"nice"}
        
        case testEndpoints.connection:
            return {"connected":True}
    # if test_endpoint is testEndpoints.me:
    #     return {"result":"nice"}
    
    # if test_endpoint is testEndpoints.connection:
    #     return {"connected":True}

@app.post("/resource_group/create")
def create_resource_group(name, location, tags={}):
    return resource_group_create_or_update(name, location, resource_client, tags)

@app.post("/resource_group/update")
def update_resource_group(name, location, tags={}):
    #additional checks before updating?
    return resource_group_create_or_update(name, location, resource_client, tags)

@app.get("/resource_group/list")
def list_resource_groups(resource_client):
    return resource_client.resource_groups.list()

@app.get("/resource_group/show/{name}")
def resource_group_show(name: str, resource_client):
    return resource_group_show(name, resource_client)

@app.get("/resource_group/get_resources_expandBy")
def resource_group_show(name, resource_client, expandBy):
    #check given expandBy argument
    return resource_group_show(name, resource_client, expandBy)


if __name__ == "__main__":
    uvicorn.run(app, port=8000, host="0.0.0.0")