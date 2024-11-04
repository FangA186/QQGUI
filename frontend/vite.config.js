import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    // proxy: {
    //   // 将 '/api' 的请求代理到后端服务器
    //   '/api': {
    //     target: 'http://wails.localhost:8080/', // 后端 API 地址
    //     changeOrigin: true,
    //     rewrite: (path) => path.replace(/^\/api/, '')
    //   }
    // }
  }
});