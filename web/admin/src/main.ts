import {createApp,h,ref} from 'vue'
import {createPinia} from 'pinia'
import {createRouter,createWebHistory} from 'vue-router'
import {NConfigProvider,NDialogProvider,NMessageProvider,zhCN,dateZhCN} from 'naive-ui'
import App from './App.vue'
import Login from './views/Login.vue'
import Dashboard from './views/Dashboard.vue'
import Apps from './views/Apps.vue'
import AppEdit from './views/AppEdit.vue'
import Releases from './views/Releases.vue'
import Storages from './views/Storages.vue'
import Settings from './views/Settings.vue'
import Pages from './views/Pages.vue'
import Themes from './views/Themes.vue'
import './styles.css'
import {currentUser} from './state'
const router=createRouter({history:createWebHistory('/admin/'),routes:[{path:'/login',component:Login},{path:'/',component:Dashboard},{path:'/apps',component:Apps},{path:'/apps/new',component:AppEdit},{path:'/apps/:id/edit',component:AppEdit},{path:'/apps/:id/releases',component:Releases},{path:'/storages',component:Storages},{path:'/themes',component:Themes},{path:'/settings',component:Settings},{path:'/pages',component:Pages}]})
router.beforeEach(to=>{if(to.path!='/login'&&!currentUser.value)return '/login';if(to.path=='/login'&&currentUser.value)return '/'})
createApp({render:()=>h(NConfigProvider,{locale:zhCN,dateLocale:dateZhCN},{default:()=>h(NDialogProvider,null,{default:()=>h(NMessageProvider,null,{default:()=>h(App)})})})}).use(createPinia()).use(router).mount('#app')
