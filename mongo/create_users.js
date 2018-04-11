db.createUser({ user: 'jsmith', pwd: 'password', roles: [ { role: "userAdminAnyDatabase", db: "admin" } ] });

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

