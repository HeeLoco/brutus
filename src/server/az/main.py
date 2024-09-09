import os
import uvicorn

from azure.identity import DefaultAzureCredential
from azure.mgmt.resource import ResourceManagementClient
from enum import Enum
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

#local
from library.resource_groups import resource_group_create_or_update, resource_group_show

class testEndpoints(str, Enum):
    me         = "me"
    connection = "connection"
    clientinfo = "clientdata"
    # all        = "all"

class resourceClientInfo(BaseModel):
    credential: str
    subscription_id: str 
    api_version: str = "2022-09-01"
    base_url: str = "https://management.azure.com"

app = FastAPI()

# Acquire a credential object using DevaultAzureCredential.
# credential = DefaultAzureCredential()

# Retrieve subscription ID from environment variable.
# subscription_id = "xxx" #os.environ["AZURE_SUBSCRIPTION_ID"]

# Obtain the management object for resources.
#  resource_client = ResourceManagementClient(credential, subscription_id, "2022-09-01", )

def gen_resource_client(resource_client_info: resourceClientInfo):
    #custom validation(s)
    return ResourceManagementClient(resource_client_info.credential, resource_client_info.subscription_id, resource_client_info.api_version, resource_client_info.base_url)

@app.post("/test/{test_endpoint}", tags=["test"])
def tests(test_endpoint: testEndpoints, resource_client_info: resourceClientInfo | None = None):
    """
    Endpoints for testing purposes are:

    - **me**: own sample response
    - **connection**: connection / reachablility test
    - **clientinfo**: needed client information to get access to desired connection
    """
    match test_endpoint:
        case testEndpoints.me:
            return {"result":"nice"}
        
        case testEndpoints.connection:
            #test connection(s)
            return {"connected":True}
        
        case testEndpoints.clientinfo:
            if not resource_client_info:
                # return {"result":"no body"}
                raise HTTPException(status_code=420, detail="body is missing")
            #send back plain 
            return resource_client_info.model_dump()
        
    # if test_endpoint is testEndpoints.me:
    #     return {"result":"nice"}
    
    # if test_endpoint is testEndpoints.connection:
    #     return {"connected":True}





@app.post("/resource_group/create", 
          tags=["resource group"],
          response_description="The created group",)
def create_resource_group(name, location, tags={}):
    return resource_group_create_or_update(name, location, resource_client, tags)

@app.post("/resource_group/update", tags=["resource group"])
def update_resource_group(name, location, tags={}):
    #additional checks before updating?
    return resource_group_create_or_update(name, location, resource_client, tags)

@app.get("/resource_group/list", tags=["resource group"])
def list_resource_groups(resource_client):
    return resource_client.resource_groups.list()

@app.get("/resource_group/show/{name}", tags=["resource group"])
def resource_group_show(name: str, resource_client):
    return resource_group_show(name, resource_client)

@app.get("/resource_group/get_resources_expandBy", tags=["resource group"])
def resource_group_show(name, resource_client, expandBy):
    #check given expandBy argument
    return resource_group_show(name, resource_client, expandBy)


if __name__ == "__main__":
    uvicorn.run(app, port=8000, host="0.0.0.0")