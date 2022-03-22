config = {
  _id: "rs0",
  members: [
    { _id: 0, host: "primary:27017" },
    { _id: 1, host: "secondary1:27018" },
    { _id: 2, host: "secondary2:27019" },
  ],
};
rs.initiate(config);
