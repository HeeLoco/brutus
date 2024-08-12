from flask import Flask
from library.storage.azure.adls import create, get_input

app = Flask(__name__)


@app.route("/api/")
def hello():
    return test()


@app.route("/api/get_endpoints")
def get_endpoints():
    return [
        {
            "Id": "Azure",
            "Name": "Azure",
            "Entities": [
                {
                    "Id": "adls",
                    "Name": "Data Lake Storage",
                    "Create": "/api/storage/create",
                    "CreateMethod": "POST",
                    "GetConfig": "/api/storage/config",
                    "GetConfigMethod": "GET",
                }
            ]
        }
    ]


@app.route("/api/storage/create", methods=["POST"])
def storage_config_create(name, subscription):
    return create()


@app.route("/api/storage/config", methods=["GET"])
def storage_config_get():
    return get_input()


if __name__ == "__main__":
    app.run(debug=True)
