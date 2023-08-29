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


ITEM_ENDPOINT = "/item"

@blueprint.route("/item", methods=['GET','POST'])
def item():
    redir = check_user_loggedin(session)
    if redir != None:
        return redir

    id = request.args.get("id")
    if id == None:
        return redirect("/index")

    user = model.User.from_json(session["user"])
    req = comm.RequestSigned(
        method="GET", base_url=API_ADDRESS, endpoint=ITEM_ENDPOINT,
        user=user.email, body="", query={"id": id},
    )
    response = req.request(user.private_key_sign)
    item = comm.Item.from_json_string(response.text, user.private_key_encrypt, user.public_key_sign)
    return render_template(
        "item.html",
        index_url="/index",
        item=item
    )

@blueprint.route("/item/create", methods=["GET", "POST"])
def item_create():
    redir = check_user_loggedin(session)
    if redir != None:
        return redir

    if request.method == "GET":
        return render_template(
            "item_edit.html",
            title="Item Creation",
            save_item_url="/item/create",
            label_old="",
            username_old="",
            password_old=""
        )

    user = model.User.from_json(session["user"])
    label = request.form.get("label")
    username = request.form.get("username")
    password = request.form.get("password")
    id = request.args.get("id")
    if id == None:
        id = ""
    item = comm.Item(label, username, password, id=id)
    body = item.to_json_string(user.public_key_encrypt, user.private_key_sign)
    req = comm.RequestSigned(
        method="POST", base_url=API_ADDRESS, endpoint=ITEM_ENDPOINT,
        user=user.email, body=body
    )
    _ = req.request(user.private_key_sign)

    return redirect("/index")

@blueprint.route("/item/edit", methods=["GET"])
def item_edit():
    redir = check_user_loggedin(session)
    if redir != None:
        return redir

    id = request.args.get("id")
    if id == None:
        return redirect("/index")
    user = model.User.from_json(session["user"])
    req = comm.RequestSigned(
        method="GET", base_url=API_ADDRESS, endpoint=ITEM_ENDPOINT,
        user=user.email, body="", query={"id": id}
    )
    response = req.request(user.private_key_sign)
    item = comm.Item.from_json_string(response.text, user.private_key_encrypt, user.public_key_sign)
    if request.method == "GET":
        return render_template(
            "item_edit.html",
            title="Item Edit",
            save_item_url=f"/item/create?id={id}",
            label_old=item.label,
            username_old=item.username,
            password_old=item.password
        )

@blueprint.route("/item/delete", methods=["GET"])
def item_delete():
    redir = check_user_loggedin(session)
    if redir != None:
        return redir
    
    id = request.args.get("id")
    if id == None:
        return redirect("/index")
    user = model.User.from_json(session["user"])
    req = comm.RequestSigned(
        method="DELETE", base_url=API_ADDRESS, endpoint=ITEM_ENDPOINT,
        user=user.email, body="", query={"id": id}
    )
    req.request(user.private_key_sign)

    return redirect("/index")
