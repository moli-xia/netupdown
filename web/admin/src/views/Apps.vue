<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NDataTable, NInput, NPopconfirm, NTag, useMessage } from 'naive-ui'
import { api } from '../api'
import { openAppPreview } from '../appPreview'

const router = useRouter()
const message = useMessage()
const rows = ref<any[]>([])
const q = ref('')
const loading = ref(false)
const previewingId = ref<number | null>(null)

async function load() {
  loading.value = true
  try {
    const x: any = await api.get('/apps', { params: { q: q.value, page_size: 100 } })
    rows.value = x.list
  } catch (e: any) {
    message.error(e.message)
  } finally {
    loading.value = false
  }
}

async function preview(id: number) {
  previewingId.value = id
  try {
    await openAppPreview(id)
  } catch (e: any) {
    message.error(e.message || '预览打开失败')
  } finally {
    previewingId.value = null
  }
}

async function state(row: any, publish: boolean) {
  try {
    await api.post(`/apps/${row.id}/${publish ? 'publish' : 'unpublish'}`)
    message.success(publish ? '已发布' : '已下架')
    load()
  } catch (e: any) {
    message.error(e.message)
  }
}

async function remove(id: number) {
  await api.delete(`/apps/${id}`)
  load()
}

const columns = [
  { title: '应用', key: 'name', render: (r: any) => h('div', {}, [h('b', r.name), h('div', { style: 'color:#8b93a1' }, r.tagline || r.slug)]) },
  { title: '类型', key: 'type', render: (r: any) => h(NTag, { type: r.type === 1 ? 'success' : 'info' }, { default: () => r.type === 1 ? '自研' : '收录' }) },
  { title: '状态', key: 'status', render: (r: any) => h(NTag, { type: r.status === 1 ? 'success' : r.status === 2 ? 'warning' : 'default' }, { default: () => ['草稿', '已发布', '已下架'][r.status] }) },
  { title: '下载', key: 'download_count' },
  {
    title: '操作',
    key: 'actions',
    render: (r: any) => h('div', { style: 'display:flex;gap:8px' }, [
      h(NButton, { size: 'small', onClick: () => preview(r.id), loading: previewingId.value === r.id }, { default: () => '预览' }),
      h(NButton, { size: 'small', onClick: () => router.push(`/apps/${r.id}/edit`) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', onClick: () => router.push(`/apps/${r.id}/releases`) }, { default: () => '版本' }),
      h(NButton, { size: 'small', type: r.status === 1 ? 'warning' : 'primary', onClick: () => state(r, r.status !== 1) }, { default: () => r.status === 1 ? '下架' : '发布' }),
      h(NPopconfirm, { onPositiveClick: () => remove(r.id) }, { trigger: () => h(NButton, { size: 'small', type: 'error', secondary: true }, { default: () => '删除' }), default: () => '确定删除该应用？' }),
    ]),
  },
]

onMounted(load)
</script>

<template>
  <div class="page-head"><h1>应用管理</h1><n-button type="primary" @click="router.push('/apps/new')">新建应用</n-button></div>
  <div class="panel">
    <div style="display:flex;gap:10px;margin-bottom:16px"><n-input v-model:value="q" clearable placeholder="搜索名称或简介" @keyup.enter="load" /><n-button @click="load">搜索</n-button></div>
    <n-data-table :columns="columns" :data="rows" :loading="loading" :row-key="(r: any) => r.id" />
  </div>
</template>
