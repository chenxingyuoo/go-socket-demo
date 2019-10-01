
module.exports = {
  // 开发服务器配置
  devServer: {
    port: `8082`,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8000',
        changeOrigin: true,
        secure: false,
        pathRewrite: {
          '/api': ''
        }
      }
    }
  }
}
