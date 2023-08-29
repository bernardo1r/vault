CREATE TABLE public."user" (
    email VARCHAR PRIMARY KEY,
    public_key_enc VARCHAR NOT NULL,
    private_key_enc VARCHAR NOT NULL,
    public_key_sign VARCHAR NOT NULL,
    private_key_sign VARCHAR NOT NULL,
    totp VARCHAR NOT NULL
);

CREATE TABLE public.item (
    email VARCHAR NOT NULL,
    id CHARACTER(24) NOT NULL,
    "label" VARCHAR NOT NULL,
    "key" VARCHAR NOT NULL,
    credential VARCHAR NOT NULL,
    PRIMARY KEY(email, id),
    FOREIGN KEY(email)
        REFERENCES public.user(email)
);

CREATE FUNCTION public."index"("user" VARCHAR)
RETURNS TABLE ("id" CHARACTER(24), "label" VARCHAR)
LANGUAGE plpgsql AS
$body$
BEGIN
    RETURN QUERY
    SELECT item."id", item."label"
    FROM public.item
    WHERE email = "user";
END
$body$;