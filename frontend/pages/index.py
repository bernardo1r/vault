from typing import Union
from flask import Flask, render_template, redirect, url_for, request, session
import pyotp
import json
import qrcode
import io
import base64
import model
import requests
import model
import comm

from . import blueprint, API_ADDRESS, check_user_loggedin


INDEX_ENDPOINT = "/index"

@blueprint.route("/index", methods=['GET','POST'])
def index():
    redir = check_user_loggedin(session)
    if redir != None:
        return redir

    user = model.User.from_json(session["user"])
    req = comm.RequestSigned(
        method="GET", base_url=API_ADDRESS, endpoint=INDEX_ENDPOINT,
        user=user.email, body=""
    )
    response = req.request(user.private_key_sign)
    if response.status_code != 200:
        return redirect("/logout")

    return render_template(
        "index.html",
        index=response.json()
    )
