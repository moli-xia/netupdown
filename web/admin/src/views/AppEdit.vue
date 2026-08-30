<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NCheckbox, NForm, NFormItem, NInput, NSelect, useMessage } from 'naive-ui'
import { api } from '../api'
import { openAppPreview } from '../appPreview'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const saving = ref(false)
const previewing = ref(false)
const cats = ref<any[]>([])
const form = reactive<any>({
  name: '',
  slug: '',
  type: 1,
  category_id: null,
  tagline: '',
  description: '',
  icon: '',
  screenshots: [],
  official_url: '',
  repo_url: '',
  developer: '',
  license: '',
  platforms: [],
  seo_title: '',
  seo_description: '',
  seo_keywords: '',
})
const platforms = ['windows', 'macos', 'linux', 'android', 'ios', 'web']

onMounted(async () => {
  const x: any = await api.get('/categories')
  cats.value = x.list
  if (route.params.id) Object.assign(form, await api.get(`/apps/${route.params.id}`))
})

function slug() {
  if (!route.params.id && !form.slug) {
    form.slug = form.name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
  }
}

async function save() {
  saving.value = true
  try {
    const data = route.params.id ? await api.put(`/apps/${route.params.id}`, form) : await api.post('/apps', form)
    message.success('已保存')
    if (!route.params.id) router.replace(`/apps/${(data as any).id}/edit`)
  } catch (e: any) {
    message.error(e.message)
  } finally {
    saving.value = false
  }
}

async function preview() {
  if (!route.params.id) return
  previewing.value = true
  try {
    await openAppPreview(Number(route.params.id))
  } catch (e: any) {
    message.error(e.message || '预览打开失败')
  } finally {
    previewing.value = false
  }
}
</script>

<template>
  <div class="page-head">
    <div>
      <h1>{{ route.params.id ? '编辑应用' : '新建应用' }}</h1>
      <span v-if="route.params.id" class="preview-hint">预览使用已保存内容，草稿也可以预览</span>
    </div>
    <div>
      <n-button :disabled="!route.params.id" :loading="previewing" @click="preview">预览</n-button>
      <n-button @click="router.push('/apps')">返回</n-button>
      <n-button type="primary" :loading="saving" @click="save">保存</n-button>
    </div>
  </div>
  <n-form class="panel" label-placement="top">
    <div class="form-grid">
      <n-form-item label="应用名称"><n-input v-model:value="form.name" @blur="slug" /></n-form-item>
      <n-form-item label="URL Slug"><n-input v-model:value="form.slug" placeholder="my-app" /></n-form-item>
      <n-form-item label="类型"><n-select v-model:value="form.type" :options="[{ label: '自研', value: 1 }, { label: '收录', value: 2 }]" /></n-form-item>
      <n-form-item label="分类"><n-select v-model:value="form.category_id" clearable :options="cats.map(x => ({ label: x.name, value: x.id }))" /></n-form-item>
      <n-form-item class="wide" label="一句话简介"><n-input v-model:value="form.tagline" maxlength="200" show-count /></n-form-item>
      <n-form-item class="wide" label="详细介绍（Markdown）"><n-input v-model:value="form.description" type="textarea" :rows="10" /></n-form-item>
      <n-form-item label="图标 URL"><n-input v-model:value="form.icon" placeholder="/uploads/..." /></n-form-item>
      <n-form-item label="开发者"><n-input v-model:value="form.developer" /></n-form-item>
      <n-form-item label="官方网站"><n-input v-model:value="form.official_url" /></n-form-item>
      <n-form-item label="代码仓库"><n-input v-model:value="form.repo_url" /></n-form-item>
      <n-form-item label="许可"><n-input v-model:value="form.license" /></n-form-item>
      <n-form-item class="wide" label="支持平台">
        <n-checkbox v-for="p in platforms" :key="p" :checked="form.platforms.includes(p)" @update:checked="(v: boolean) => form.platforms = v ? [...form.platforms, p] : form.platforms.filter((x: string) => x !== p)">{{ p }}</n-checkbox>
      </n-form-item>
    </div>
  </n-form>
</template>
