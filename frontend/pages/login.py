from flask import render_template, redirect, request, session
import model

import model
from . import blueprint, API_ADDRESS 

LOGIN_ENDPOINT = "/login"
AUTHENTICATOR_ENDPOINT = "/authenticator"

@blueprint.route("/login", methods=['GET', 'POST'])
def login():
    if "user" not in session:
        return login_password()

    elif session["user"]["init"]:
        return redirect("/index")

    else:
        return login_authenticator()

def login_password():
    username_old = ""
    if request.method == "GET":
        return render_template(
            "login.html",
            login_url="/login",
            register_url="/register",
            username_old=username_old,
            error=False
        )

    username = request.form.get("username")
    if username == None:
        username = ""

    password = request.form.get("password")
    if password == None:
        password = "" 

    user = model.User(username, password)
    if user.try_login(API_ADDRESS+LOGIN_ENDPOINT):
        session["user"] = user.to_json()
        return render_template(
            "authenticator.html",
            username=username,
            authenticator_url="/login",
            error=False
        )

    return render_template(
        "login.html",
        login_url="/login",
        register_url="/register",
        username_old=username,
        error=True
    )

def login_authenticator():
    user = model.User.from_json(session["user"])

    passcode = request.form.get("passcode")
    ok = user.login(API_ADDRESS+AUTHENTICATOR_ENDPOINT, passcode)
    if ok:
        session["user"] = user.to_json()
        return redirect("/index")


    return render_template(
        "authenticator.html",
        username=user.email,
        authenticator_url="/login",
        error=True
    )