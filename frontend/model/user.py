from Cryptodome.PublicKey import RSA
from Cryptodome.Cipher import AES
import crypto
import pyotp
import qrcode
import io
import base64
import hashlib
import requests

from crypto import CLIENT_CERT
from comm import verify_cert

class User:
    KEYPAIR_BITS: int = 2048
    PROGRAM_NAME: str = "cofre"
    AESGCMNONCE_LEN: int = 12

    email: str
    password: str
    _init: bool
    key: bytes
    public_key_encrypt: RSA.RsaKey
    private_key_encrypt: RSA.RsaKey
    public_key_sign: RSA.RsaKey
    private_key_sign: RSA.RsaKey
    totp: pyotp.TOTP

    def __init__(self, email: str, password: str):
        self.email = email
        self.password = password
        self.key = crypto.generate_key(self.password, self.email)
        self._init = False

    def create(self):
        self._init = True
        self.private_key_encrypt = RSA.generate(User.KEYPAIR_BITS)
        self.public_key_encrypt = self.private_key_encrypt.public_key()
        self.private_key_sign = RSA.generate(User.KEYPAIR_BITS)
        self.public_key_sign = self.private_key_sign.public_key()
        secret = pyotp.random_base32()
        self.totp = pyotp.TOTP(secret)

        return self

    def initialized(self) -> bool:
        return self._init

    def qr_code(self) -> str:
        url = self.totp.provisioning_uri(name=self.email, issuer_name=User.PROGRAM_NAME)
        qr_code = qrcode.make(url).resize((200, 200))
        buffer = io.BytesIO()
        qr_code.save(buffer, format="PNG")
        qr_code_b64 = base64.b64encode(buffer.getvalue()).decode()

        return qr_code_b64

    def to_json(self):
        user = dict()
        user["email"] = self.email
        user["password"] = self.password
        user["key"] = base64.b64encode(self.key).decode()
        user["init"] = self._init
        
        if self._init:
            user["public_key_encrypt"] = base64.b64encode(self.public_key_encrypt.export_key("DER")).decode()
            user["private_key_encrypt"] = base64.b64encode(self.private_key_encrypt.export_key("DER")).decode()
            user["public_key_sign"] = base64.b64encode(self.public_key_sign.export_key("DER")).decode()
            user["private_key_sign"] = base64.b64encode(self.private_key_sign.export_key("DER")).decode()

        return user

    @staticmethod
    def from_json(data: dict):
        user = User.__new__(User)
        user.email = data["email"]
        user.password = data["password"]
        user.key = base64.b64decode(data["key"].encode())
        user._init = data["init"]

        if not user._init:
            return user

        user.public_key_encrypt = RSA.import_key(base64.b64decode(data["public_key_encrypt"].encode()))
        user.private_key_encrypt = RSA.import_key(base64.b64decode(data["private_key_encrypt"].encode()))
        user.public_key_sign = RSA.import_key(base64.b64decode(data["public_key_sign"].encode()))
        user.private_key_sign = RSA.import_key(base64.b64decode(data["private_key_sign"].encode()))

        return user

    def _encrypt_key(self, encryption_key: bool, encoded: bool) -> tuple[str, bytes]:
        if encryption_key:
            payload = self.private_key_encrypt
            nonce = bytes([0]*User.AESGCMNONCE_LEN)

        else:
            payload = self.private_key_sign
            nonce = bytes([255]*User.AESGCMNONCE_LEN)

        cipher = AES.new(key=self.key, mode=AES.MODE_GCM, nonce=nonce)
        ciphertext, tag = cipher.encrypt_and_digest(payload.export_key("DER"))
        ciphertext = ciphertext + tag
        
        if not encoded:
            return ciphertext
        
        return base64.b64encode(ciphertext).decode()

    def _decrypt_key(self, key, encryption_key: bool, encoded: bool) -> RSA.RsaKey:
        if encryption_key:
            nonce = bytes([0]*User.AESGCMNONCE_LEN)

        else:
            nonce = bytes([255]*User.AESGCMNONCE_LEN)

        if encoded:
            key = base64.b64decode(key.encode())

        ciphertext, tag = key[:-16], key[-16:]
        cipher = AES.new(key=self.key, mode=AES.MODE_GCM, nonce=nonce)
        try:
            data = cipher.decrypt_and_verify(ciphertext, tag)

        except ValueError:
            return None

        return RSA.import_key(data)

    def _make_register_payload(self) -> dict:
        payload = dict()
        payload["publicKeyEnc"] = base64.b64encode(self.public_key_encrypt.export_key("DER")).decode()
        payload["privateKeyEnc"] = self._encrypt_key(True, True)
        payload["publicKeySign"] = base64.b64encode(self.public_key_sign.export_key("DER")).decode()
        payload["privateKeySign"] = self._encrypt_key(False, True)
        payload["totp"] = self.totp.secret

        return payload

    def _api_password(self) -> str:
        data = self.key + self.password.encode()
        password = base64.b64encode(hashlib.sha3_256(data).digest()).decode()
        return password

    def register(self, endpoint: str) -> requests.Response:
        password = self._api_password()
        payload = self._make_register_payload()
        response = requests.request(
            method="POST", url=endpoint, 
            auth=(self.email, password), json=payload,
            cert=CLIENT_CERT,
            verify=False
        )
        verify_cert(endpoint)

        return response

    def try_login(self, endpoint: str) -> bool:
        if self._init:
            return True

        password = self._api_password()
        response = requests.request(
            method="GET", url=endpoint,
            auth=((self.email, password)),
            cert=CLIENT_CERT,
            verify=False
        )
        verify_cert(endpoint)
        if response.status_code == 200:
            return True

        return False

    def login(self, endpoint: str, passcode: str) -> bool:
        if self._init:
            return True

        password = self._api_password()
        response = requests.request(
            method="GET", url=endpoint,
            headers={"Passcode": passcode},
            auth=((self.email, password)),
            cert=CLIENT_CERT,
            verify=False
        )
        verify_cert(endpoint)
        if response.status_code != 200:
            return False

        self._init = True
        data = response.json()
        self.public_key_encrypt = RSA.import_key(base64.b64decode(data["publicKeyEnc"].encode()))
        self.private_key_encrypt = self._decrypt_key(data["privateKeyEnc"], True, True)
        self.public_key_sign = RSA.import_key(base64.b64decode(data["publicKeySign"].encode()))
        self.private_key_sign = self._decrypt_key(data["privateKeySign"], False, True)

        return True
