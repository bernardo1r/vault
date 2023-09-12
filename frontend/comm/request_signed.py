from typing import Union
import json
from collections import OrderedDict
import datetime
from Cryptodome.PublicKey import RSA
from Cryptodome.Signature import pss
from Cryptodome.Hash import SHA3_256
import base64
import requests
import ssl

from crypto import SERVER_CERT, CLIENT_CERT


def verify_cert(url: str):
    url = url.split("/")[-2]
    url, port = url.split(":")
    try:
        server_cert = ssl.get_server_certificate((url, port), ca_certs=SERVER_CERT)
    except ssl.SSLCertVerificationError:
        print("invalid server certificate", file=sys.stderr) 
        os.exit(1)

class RequestSigned:
    method: str
    base_url: str
    endpoint: str
    user: str
    date: str

    def __init__(
        self,
        method: str = "GET",
        base_url: str = "",
        endpoint: str = "",
        user: str = "",
        body: str = "",
        query: dict[str,str] = dict()
    ):
        self.method = method
        self.base_url = base_url
        self.endpoint = endpoint
        self.user = user
        self.date = datetime.datetime.now().astimezone().isoformat(timespec="seconds")
        self.body = body
        self.query=query

    def dump(self) -> str:
        data = OrderedDict()
        data["method"] = self.method
        data["endpoint"] = self.endpoint
        data["user"] = self.user
        data["date"] = self.date
        data["body"] = self.body
        
        return json.dumps(data, ensure_ascii=False, separators=(",", ":"))

    def sign(self, private_key_sign: RSA.RsaKey, encoded=True) -> Union[bytes, str]:
        data = self.dump()
        hash = SHA3_256.new(data.encode())
        sig = pss.new(private_key_sign).sign(hash)

        if encoded:
            sig = base64.b64encode(sig).decode()

        return sig

    def request(self, private_key_sign: RSA.RsaKey) -> requests.Response:
        headers = dict()
        headers["User"] = self.user
        headers["Signature"] = self.sign(private_key_sign)
        headers["Date"] = self.date
        url = self.base_url+self.endpoint
        if len(self.query) > 0:
            url += "?"
        for key, value in self.query.items():
            url += f"{key}={value}"

        verify_cert(url)
        return requests.request(
            method=self.method,
            url=url,
            headers=headers,
            data=self.body,
            cert=CLIENT_CERT,
            verify=False
        )

