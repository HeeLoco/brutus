from flask import Flask
from library.foo import test

app = Flask(__name__)


@app.route("/api/")
def hello():
    return test()


if __name__ == '__main__':
    app.run(debug=True)