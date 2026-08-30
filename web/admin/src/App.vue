<script setup lang="ts">
import {computed,onMounted,ref} from 'vue'
import {useRoute,useRouter} from 'vue-router'
import {NButton,NLayout,NLayoutContent,NLayoutHeader,NLayoutSider,NMenu,NSpin} from 'naive-ui'
import {bootstrapAuth} from './api';import {currentUser} from './state'
const route=useRoute(),router=useRouter(),ready=ref(false),collapsed=ref(false)
const menu=[{label:'仪表盘',key:'/'},{label:'应用管理',key:'/apps'},{label:'存储管理',key:'/storages'},{label:'主题外观',key:'/themes'},{label:'单页管理',key:'/pages'},{label:'站点设置',key:'/settings'}]
const active=computed(()=>route.path.startsWith('/apps')?'/apps':route.path)
onMounted(async()=>{currentUser.value=await bootstrapAuth();ready.value=true;if(!currentUser.value&&route.path!='/login')router.replace('/login')})
</script>
<template><n-spin v-if="!ready" class="boot" size="large"/><router-view v-else-if="route.path==='/login'"/><n-layout v-else has-sider class="shell"><n-layout-sider bordered collapse-mode="width" :collapsed-width="72" :width="244" :collapsed="collapsed"><div class="logo"><span>N</span><b v-if="!collapsed">NetUpDown</b></div><n-menu :value="active" :options="menu" @update:value="router.push"/><button class="collapse" @click="collapsed=!collapsed">{{collapsed?'›':'‹'}}</button></n-layout-sider><n-layout><n-layout-header bordered class="header"><div><b>{{menu.find(x=>x.key===active)?.label||'内容编辑'}}</b><small>管理你的发布内容与下载来源</small></div><div class="profile">{{currentUser?.nickname||currentUser?.username}} <n-button quaternary @click="currentUser=null;router.push('/login')">退出</n-button></div></n-layout-header><n-layout-content class="content"><router-view/></n-layout-content></n-layout></n-layout></template>
