db.createUser(
    {
        user: "hasher",
        pwd: "password",
        roles: [
            { role: "readWrite", db: "hashes" }
        ]
    },
    {
        w: "majority",
        wtimeout: 5000
    }
);

