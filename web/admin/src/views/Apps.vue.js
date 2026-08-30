import { h, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { NButton, NPopconfirm, NTag, useMessage } from 'naive-ui';
import { api } from '../api';
import { openAppPreview } from '../appPreview';
const router = useRouter();
const message = useMessage();
const rows = ref([]);
const q = ref('');
const loading = ref(false);
const previewingId = ref(null);
async function load() {
    loading.value = true;
    try {
        const x = await api.get('/apps', { params: { q: q.value, page_size: 100 } });
        rows.value = x.list;
    }
    catch (e) {
        message.error(e.message);
    }
    finally {
        loading.value = false;
    }
}
async function preview(id) {
    previewingId.value = id;
    try {
        await openAppPreview(id);
    }
    catch (e) {
        message.error(e.message || '预览打开失败');
    }
    finally {
        previewingId.value = null;
    }
}
async function state(row, publish) {
    try {
        await api.post(`/apps/${row.id}/${publish ? 'publish' : 'unpublish'}`);
        message.success(publish ? '已发布' : '已下架');
        load();
    }
    catch (e) {
        message.error(e.message);
    }
}
async function remove(id) {
    await api.delete(`/apps/${id}`);
    load();
}
const columns = [
    { title: '应用', key: 'name', render: (r) => h('div', {}, [h('b', r.name), h('div', { style: 'color:#8b93a1' }, r.tagline || r.slug)]) },
    { title: '类型', key: 'type', render: (r) => h(NTag, { type: r.type === 1 ? 'success' : 'info' }, { default: () => r.type === 1 ? '自研' : '收录' }) },
    { title: '状态', key: 'status', render: (r) => h(NTag, { type: r.status === 1 ? 'success' : r.status === 2 ? 'warning' : 'default' }, { default: () => ['草稿', '已发布', '已下架'][r.status] }) },
    { title: '下载', key: 'download_count' },
    {
        title: '操作',
        key: 'actions',
        render: (r) => h('div', { style: 'display:flex;gap:8px' }, [
            h(NButton, { size: 'small', onClick: () => preview(r.id), loading: previewingId.value === r.id }, { default: () => '预览' }),
            h(NButton, { size: 'small', onClick: () => router.push(`/apps/${r.id}/edit`) }, { default: () => '编辑' }),
            h(NButton, { size: 'small', onClick: () => router.push(`/apps/${r.id}/releases`) }, { default: () => '版本' }),
            h(NButton, { size: 'small', type: r.status === 1 ? 'warning' : 'primary', onClick: () => state(r, r.status !== 1) }, { default: () => r.status === 1 ? '下架' : '发布' }),
            h(NPopconfirm, { onPositiveClick: () => remove(r.id) }, { trigger: () => h(NButton, { size: 'small', type: 'error', secondary: true }, { default: () => '删除' }), default: () => '确定删除该应用？' }),
        ]),
    },
];
onMounted(load);
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "page-head" },
});
/** @type {__VLS_StyleScopedClasses['page-head']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({});
let __VLS_0;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_2 = __VLS_1({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_5;
const __VLS_6 = {
    /** @type {typeof __VLS_5.click} */
    onClick: (...[$event]) => {
        return (__VLS_ctx.router.push('/apps/new'));
        // @ts-ignore
        [router,];
    },
};
const { default: __VLS_7 } = __VLS_3.slots;
// @ts-ignore
[];
var __VLS_3;
var __VLS_4;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "panel" },
});
/** @type {__VLS_StyleScopedClasses['panel']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ style: {} },
});
let __VLS_8;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent1(__VLS_8, new __VLS_8({
    ...{ 'onKeyup': {} },
    value: (__VLS_ctx.q),
    clearable: true,
    placeholder: "搜索名称或简介",
}));
const __VLS_10 = __VLS_9({
    ...{ 'onKeyup': {} },
    value: (__VLS_ctx.q),
    clearable: true,
    placeholder: "搜索名称或简介",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
let __VLS_13;
const __VLS_14 = {
    /** @type {typeof __VLS_13.keyup} */
    onKeyup: (__VLS_ctx.load),
};
var __VLS_11;
var __VLS_12;
let __VLS_15;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_16 = __VLS_asFunctionalComponent1(__VLS_15, new __VLS_15({
    ...{ 'onClick': {} },
}));
const __VLS_17 = __VLS_16({
    ...{ 'onClick': {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_16));
let __VLS_20;
const __VLS_21 = {
    /** @type {typeof __VLS_20.click} */
    onClick: (__VLS_ctx.load),
};
const { default: __VLS_22 } = __VLS_18.slots;
// @ts-ignore
[q, load, load,];
var __VLS_18;
var __VLS_19;
let __VLS_23;
/** @ts-ignore @type { | typeof __VLS_components.nDataTable | typeof __VLS_components.NDataTable | typeof __VLS_components['n-data-table']} */
nDataTable;
// @ts-ignore
const __VLS_24 = __VLS_asFunctionalComponent1(__VLS_23, new __VLS_23({
    columns: (__VLS_ctx.columns),
    data: (__VLS_ctx.rows),
    loading: (__VLS_ctx.loading),
    rowKey: ((r) => r.id),
}));
const __VLS_25 = __VLS_24({
    columns: (__VLS_ctx.columns),
    data: (__VLS_ctx.rows),
    loading: (__VLS_ctx.loading),
    rowKey: ((r) => r.id),
}, ...__VLS_functionalComponentArgsRest(__VLS_24));
// @ts-ignore
[columns, rows, loading,];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
