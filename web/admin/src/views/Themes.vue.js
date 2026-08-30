import { ref } from 'vue';
import { useMessage } from 'naive-ui';
import { api } from '../api';
const rows = ref([]), message = useMessage(), uploading = ref(false);
async function load() { rows.value = await api.get('/themes'); }
async function activate(id) { try {
    await api.post(`/themes/${id}/activate`);
    message.success('主题已启用');
    load();
}
catch (e) {
    message.error(e.message);
} }
async function upload(ev) { const file = ev.target.files?.[0]; if (!file)
    return; uploading.value = true; try {
    const data = new FormData();
    data.append('file', file);
    await api.post('/themes/upload', data);
    message.success('主题安装成功');
    load();
}
catch (e) {
    message.error(e.message);
}
finally {
    uploading.value = false;
    ev.target.value = '';
} } // @ts-ignore
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
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
    ...{ style: {} },
});
__VLS_asFunctionalElement1(__VLS_intrinsics.label, __VLS_intrinsics.label)({
    ...{ class: "upload-btn" },
});
/** @type {__VLS_StyleScopedClasses['upload-btn']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.input)({
    ...{ onChange: (__VLS_ctx.upload) },
    type: "file",
    accept: ".zip",
    hidden: true,
});
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
(__VLS_ctx.uploading ? '安装中…' : '上传主题 ZIP');
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "metric-grid" },
});
/** @type {__VLS_StyleScopedClasses['metric-grid']} */ ;
for (const [t] of __VLS_vFor((__VLS_ctx.rows))) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.article, __VLS_intrinsics.article)({
        key: (t.id),
        ...{ class: "metric" },
    });
    /** @type {__VLS_StyleScopedClasses['metric']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    if (t.builtin) {
        let __VLS_0;
        /** @ts-ignore @type { | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag'] | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag']} */
        nTag;
        // @ts-ignore
        const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({}));
        const __VLS_2 = __VLS_1({}, ...__VLS_functionalComponentArgsRest(__VLS_1));
        const { default: __VLS_5 } = __VLS_3.slots;
        // @ts-ignore
        [upload, uploading, rows,];
        var __VLS_3;
    }
    if (t.active) {
        let __VLS_6;
        /** @ts-ignore @type { | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag'] | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag']} */
        nTag;
        // @ts-ignore
        const __VLS_7 = __VLS_asFunctionalComponent1(__VLS_6, new __VLS_6({
            type: "success",
        }));
        const __VLS_8 = __VLS_7({
            type: "success",
        }, ...__VLS_functionalComponentArgsRest(__VLS_7));
        const { default: __VLS_11 } = __VLS_9.slots;
        // @ts-ignore
        [];
        var __VLS_9;
    }
    __VLS_asFunctionalElement1(__VLS_intrinsics.strong, __VLS_intrinsics.strong)({
        ...{ style: {} },
    });
    (t.name);
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    (t.version);
    (t.author);
    __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
    (t.description);
    if (!t.active) {
        let __VLS_12;
        /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
        nButton;
        // @ts-ignore
        const __VLS_13 = __VLS_asFunctionalComponent1(__VLS_12, new __VLS_12({
            ...{ 'onClick': {} },
            type: "primary",
            size: "small",
        }));
        const __VLS_14 = __VLS_13({
            ...{ 'onClick': {} },
            type: "primary",
            size: "small",
        }, ...__VLS_functionalComponentArgsRest(__VLS_13));
        let __VLS_17;
        const __VLS_18 = {
            /** @type {typeof __VLS_17.click} */
            onClick: (...[$event]) => {
                if (!(!t.active))
                    throw 0;
                return (__VLS_ctx.activate(t.id));
                // @ts-ignore
                [activate,];
            },
        };
        const { default: __VLS_19 } = __VLS_15.slots;
        // @ts-ignore
        [];
        var __VLS_15;
        var __VLS_16;
    }
    // @ts-ignore
    [];
}
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
