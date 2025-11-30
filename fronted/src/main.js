import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'nprogress/nprogress.css' // 引入 nprogress 样式
import router from './router'
import './style.css'
import App from './App.vue'
import { initSiteConfig } from './store/site'

const app = createApp(App)

app.use(ElementPlus)
app.use(router)

// 初始化网站配置
initSiteConfig().then(() => {
  app.mount('#app')
})
