import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://192.168.142.132:8080', // 后端服务地址
        changeOrigin: true,
        secure: false,
      }
    },
    host: true,
    port: 3000,
  }
})
