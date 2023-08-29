import argon2
import hashlib

ARGON2_TIME: int = 1
ARGON2_MEMORY: int = 2**21
ARGON2_THREADS: int = 4
ARGON2_KEYLEN: int = 32
ARGON2_TYPE: argon2.Type = argon2.Type.ID

def generate_key(password: str, salt: str) -> bytes:
    salt = hashlib.sha3_256(salt.encode()).digest()
    return argon2.low_level.hash_secret_raw(
        secret=password.encode(),
        salt=salt,
        time_cost=ARGON2_TIME,
        memory_cost=ARGON2_MEMORY,
        parallelism=ARGON2_THREADS,
        hash_len=ARGON2_KEYLEN,
        type=ARGON2_TYPE
    )