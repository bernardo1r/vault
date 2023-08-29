from flask import redirect

ISSUER_OTP = "Cofre"

API_ADDRESS = "https://127.0.0.1:54321"

def check_user_loggedin(session):
    if "user" not in session or not session["user"]["init"]:
        return redirect("/login")

    return None