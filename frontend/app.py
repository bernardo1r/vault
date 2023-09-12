# import the Flask class from the flask module
from flask import Flask, redirect
import os
import sys

import pages

VARIABLES = ("API_ADDR", "FLASK_RUN_HOST", "FLASK_RUN_PORT")

app = Flask(__name__)
app.config["SESSION_TYPE"] = "filesystem"
app.secret_key = "test"
app.register_blueprint(pages.blueprint)

@app.route('/')
def home():
    return redirect("/login")

def check_variables():
    for variable in VARIABLES:
        value = os.environ.get(variable)
        if variable == None:
            print(f"variable {variable} not found", file=sys.stderr)
            sys.exit(1)

if __name__ == '__main__':
    check_variables()
    app.run(debug=False, host=os.environ.get("FLASK_RUN_HOST"), port=os.environ.get("FLASK_RUN_PORT"))
