from typing import Union
from flask import Flask, render_template, redirect, url_for, request, session
import pyotp
import json
import qrcode
import io
import base64
import model
import requests

from . import API_ADDRESS, blueprint, check_user_loggedin

REGISTER_ENDPOINT = "/register"

@blueprint.route("/register", methods=["GET", "POST"])
def register():
    if "user" in session:
        if not session["user"]["init"]:
            return redirect("/login")

        return redirect("/index")

    user_error = False
    password_error = False
    username_old = ""
    password_old = ""
    password_confirm_old = ""
    if request.method == "POST":
        username = request.form.get("username")
        password = request.form.get("password")
        password_confirm = request.form.get("password_confirm")

        if password == password_confirm:
            user = model.User(username, password).create()
            response = user.register(API_ADDRESS+REGISTER_ENDPOINT)
            if response.status_code != 200:
                user_error = True

            else:
                return render_template(
                    "register_success.html",
                    qr_code=user.qr_code(),
                    login_url="/login"
                )

        else:
            password_error=True

        username_old = username
        password_old = password
        password_confirm_old = password_confirm

    return render_template(
        "register.html",
        register_url="/register",
        user_error=user_error,
        password_error=password_error,
        username_old=username_old,
        password_old=password_old,
        password_confirm_old=password_confirm_old
    )