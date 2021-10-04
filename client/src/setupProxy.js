const { createProxyMiddleware }= require('http-proxy-middleware')

module.exports = function(app) {
    app.use(
        '/api/auth',
        createProxyMiddleware({
            target: "http://auth-api:7777/v1",
            changeOrigin: true
        })
    )
}