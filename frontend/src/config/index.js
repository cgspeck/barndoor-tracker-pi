const config = {
  endpoint:
    process.env.NODE_ENV === 'development' ? 'http://localhost:5000' : '',
};

export default config;
