db.createUser({user: "adder", pwd: "password", roles: [{ role: "readWrite", db: "images" }]}, {w: "majority", wtimeout: 5000});

