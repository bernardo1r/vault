# import the Flask class from the flask module
from flask import Flask, redirect

import pages


app = Flask(__name__)
app.config["SESSION_TYPE"] = "filesystem"
app.secret_key = "test"
app.register_blueprint(pages.blueprint)

@app.route('/')
def home():
    return redirect("/login")

if __name__ == '__main__':
    app.run(debug=False)
