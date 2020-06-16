const config = {
  endpoint:
    process.env.NODE_ENV === 'development' ? 'http://localhost:8882' : ':5000',
};

export default config;
