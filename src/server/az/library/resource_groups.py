def resource_group_create_or_update(name, location, resource_client, tags={}):
    rg_result = resource_client.resource_groups.create_or_update(
        {name}, 
        {
            "location": f'"{location}"',
            "tags": {tags},
        },
    )
    return rg_result

def resource_group_show(name, resource_client, expandBy=""):
    rg_result = resource_client.resources.list_by_resource_group(
        name,
        expand = {expandBy}
    )
    return rg_result