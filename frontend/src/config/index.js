const config = {
  endpoint: process.env.NODE_ENV == "development" ? "http://localhost:8882" : ""
};

export default config;
