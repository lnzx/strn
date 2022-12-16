import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import VueRouter from 'unplugin-vue-router/vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    VueRouter(),
    vue()
  ],
  resolve: {
    alias: {
      '@src': resolve(__dirname, 'src'),
      '@views': resolve(__dirname, 'src/views'),
      '@comps': resolve(__dirname, 'src/components')
    }
  }
})
