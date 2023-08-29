import secrets
import base64
from Cryptodome.PublicKey import RSA
from Cryptodome.Cipher import PKCS1_OAEP, AES
from Cryptodome.Signature import pss
from Cryptodome.Hash import SHA3_256
from typing import Union
import json
from collections import OrderedDict


class Item:
    IDBYTES_LEN: int = 16
    SYMMETRICKEY_LEN: int = 32
    GCMNONCE_LEN: int = 12
    GCMTAG_LEN: int = 16
    id: str
    label: str
    key: bytes
    username: str
    password: str

    def __init__(self, label: str, username: str, password: str, id: str=""):
        if id == "":
            id = self._gen_id()

        self.id = id
        self.label = label
        self.username = username
        self.password = password
        self._make_key()

    def _gen_id(self) -> str:
        return base64.urlsafe_b64encode(secrets.token_bytes(self.IDBYTES_LEN)).decode()

    def _make_key(self) -> None:
        self.key = secrets.token_bytes(self.SYMMETRICKEY_LEN)

    def _decrypt_key(self, private_key_encrypt: RSA.RsaKey) -> bool:
        cipher = PKCS1_OAEP.new(private_key_encrypt)
        try:
            self.key = cipher.decrypt(self.key)
        except ValueError:
            return False

        return True

    def _decrypt_signature(self, sig: bytes) -> Union[bytes, None]:
        sig, tag = sig[:-self.GCMTAG_LEN], sig[-self.GCMTAG_LEN:]
        nonce = bytes([1]*self.GCMNONCE_LEN)
        cipher = AES.new(self.key, mode=AES.MODE_GCM, nonce=nonce)
        try:
            return cipher.decrypt_and_verify(sig, tag)
        except ValueError:
            return None

    def _verify_sig_key(self, sig: bytes, public_key_sign: RSA.RsaKey) -> bool:
        hash = SHA3_256.new(self.key)
        verifier = pss.new(public_key_sign)
        try:
            verifier.verify(hash, sig)
        except ValueError:
            return False

        return True

    def _decrypt_credential(self, data: bytes) -> bool:
        ciphertext, tag = data[:-self.GCMTAG_LEN], data[-self.GCMTAG_LEN:]
        nonce = bytes([0]*self.GCMNONCE_LEN)
        cipher = AES.new(self.key, mode=AES.MODE_GCM, nonce=nonce)
        try:
            plaintext = cipher.decrypt_and_verify(ciphertext, tag)
        except ValueError:
            return False

        data = json.loads(plaintext)
        self.username = data["username"]
        self.password = data["password"]

        return True

    @staticmethod
    def from_json_string(
        payload: str,
        private_key_encrypt: RSA.RsaKey,
        public_key_sign: RSA.RsaKey
    ):
        data = json.loads(payload)
        item = Item.__new__(Item)
        item.id = data["id"]
        item.label = data["label"]
        item.key = base64.b64decode(data["key"].encode())
        item.key, sig = item.key[:-public_key_sign.size_in_bytes()-item.GCMTAG_LEN], item.key[-public_key_sign.size_in_bytes()-item.GCMTAG_LEN:]
        if not item._decrypt_key(private_key_encrypt):
            return None

        sig = item._decrypt_signature(sig)
        if sig == None:
            return None

        if not item._verify_sig_key(sig, public_key_sign):
            return None

        if item._decrypt_credential(base64.b64decode(data["credential"].encode())):
            return item

        return None

    def sign_key(self, private_key_sign: RSA.RsaKey) -> bytes:
        hash = SHA3_256.new(self.key)
        sig = pss.new(private_key_sign).sign(hash)

        return sig

    def encrypt_key(self, public_key_encrypt: RSA.RsaKey) -> bytes:
        cipher = PKCS1_OAEP.new(public_key_encrypt)
        ciphertext = cipher.encrypt(self.key)

        return ciphertext

    def encrypt_signature(self, sig: bytes) -> bytes:
        nonce = bytes([1]*self.GCMNONCE_LEN)
        cipher = AES.new(self.key, mode=AES.MODE_GCM, nonce=nonce)
        ciphertext, tag = cipher.encrypt_and_digest(sig)

        return ciphertext + tag

    def _encrypt_credential(self) -> bytes:
        data = OrderedDict()
        data["username"] = self.username
        data["password"] = self.password
        payload = json.dumps(data, ensure_ascii=False, separators=(",",":"))
        payload = payload.encode()
        nonce = bytes([0]*self.GCMNONCE_LEN)
        cipher = AES.new(self.key, mode=AES.MODE_GCM, nonce=nonce)
        ciphertext, tag = cipher.encrypt_and_digest(payload)
        return ciphertext + tag

    def to_json_string(
        self,
        public_key_encrypt: RSA.RsaKey,
        private_key_sign: RSA.RsaKey
    ) -> str:
        data = OrderedDict()
        credential = base64.b64encode(self._encrypt_credential()).decode()
        sig = self.sign_key(private_key_sign)
        sig = self.encrypt_signature(sig)
        key = self.encrypt_key(public_key_encrypt)
        key = key + sig
        key = base64.b64encode(key).decode()
        data["id"] = self.id
        data["label"] = self.label
        data["key"] = key
        data["credential"] = credential

        return json.dumps(data, ensure_ascii=False, separators=(",", ":"))