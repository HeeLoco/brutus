
def create(name, subscription):
    return "adls id"



# retrun all settings that should be possible to set for this storage
def get_input():
    return [
        {
            "id": "name",
            "display": "Name",
            "required": True,
            "type": "string",            
            "validation": "\\w\\w\\w+", # at least three characters, nothing special
        },
        {
            "id": "subscription",
            "display": "Subscription",
            "required": True,
            "type": "string",            
            "validation": "" # guid?
        }
    ]

if(__name__ == "__main__"):
    print("this is a library, but you can execute it as if it were a script, e.g. for testing?")