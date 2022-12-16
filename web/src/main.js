import { createApp } from 'vue'
import App from './App.vue'
import VueApexCharts from "vue3-apexcharts"

import './assets/main.scss'

import { createRouter, createWebHistory } from 'vue-router/auto'
import { createApi } from '@src/utils/useApi'

const router = createRouter({
    history: createWebHistory(),
    // You don't need to pass the routes anymore, the plugin writes it for you 🤖
})

const app = createApp(App)
app.provide('api', createApi())

app.use(router).use(VueApexCharts).mount('#app')


