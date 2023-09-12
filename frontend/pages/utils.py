from flask import redirect
import os

ISSUER_OTP = "Cofre"

API_ADDRESS = "https://"+os.environ.get("API_ADDR", "")

def check_user_loggedin(session):
    if "user" not in session or not session["user"]["init"]:
        return redirect("/login")

    return None
