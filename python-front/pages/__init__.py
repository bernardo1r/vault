from flask import Blueprint

blueprint = Blueprint("pages", __name__)

from .utils import *
from .register import register
from .login import login
from .index import index
from .item import item, item_create, item_edit, item_delete
from .logout import logout