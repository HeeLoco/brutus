from flask import Flask
from library.storage.azure.adls import create, get_input

app = Flask(__name__)


@app.route("/api/")
def hello():
    return test()
  
  
@app.route("/api/storage/create", methods=["POST"])
def storage_config_create(name, subscription):
    return create()
  
@app.route("/api/storage/config")
def storage_config_get():
    return get_input()


if __name__ == "__main__":
    app.run(debug=True)
