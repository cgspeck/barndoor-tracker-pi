const config = {
  port: process.env.NODE_ENV == 'development' ? '8882' : '80',
  endpoint: process.env.NODE_ENV == 'development' ? 'http://localhost:8882' : ''
}

export default config;
