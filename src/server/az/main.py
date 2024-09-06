import uvicorn
from fastapi import FastAPI

import os

from azure.identity import DefaultAzureCredential
from azure.mgmt.resource import ResourceManagementClient

from library.resource_groups import resource_group_create_or_update, resource_group_show

app = FastAPI()

if __name__ == "__main__":
    uvicorn.run(app, port=8000, host="0.0.0.0")

# Acquire a credential object using DevaultAzureCredential.
credential = DefaultAzureCredential()

# Retrieve subscription ID from environment variable.
subscription_id = os.environ["AZURE_SUBSCRIPTION_ID"]

# Obtain the management object for resources.
resource_client = ResourceManagementClient(credential, subscription_id)

@app.get("/testme")
def testme_function():
    return {"result":"nice"}

@app.post("/resource_group/create")
def resource_group_create(name, location, tags={}):
    return resource_group_create_or_update(name, location, resource_client, tags)

@app.post("/resource_group/update")
def resource_group_create(name, location, tags={}):
    #additional checks before updating?
    return resource_group_create_or_update(name, location, resource_client, tags)

@app.get("/resource_group/list")
def resource_group_list(resource_client):
    return resource_client.resource_groups.list()

@app.get("/resource_group/get_resources")
def resource_group_show(name, resource_client):
    return resource_group_show(name, resource_client)

@app.get("/resource_group/get_resources_expandBy")
def resource_group_show(name, resource_client, expandBy):
    #check given expandBy argument
    return resource_group_show(name, resource_client, expandBy)
