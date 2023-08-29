from flask import redirect, session, Response

from . import blueprint, check_user_loggedin

@blueprint.route("/logout", methods=["GET", "POST"])
def logout():
    session.clear()

    return redirect("/login")
