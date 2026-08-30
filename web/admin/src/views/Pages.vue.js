import { onMounted, reactive, ref } from 'vue';
import { NButton, useMessage } from 'naive-ui';
import { api } from '../api';
const message = useMessage(), rows = ref([]), show = ref(false), form = reactive({ title: '', slug: '', content: '', status: 0, seo_description: '' });
async function load() { const x = await api.get('/pages'); rows.value = x.list; }
function edit(r) { Object.assign(form, r || { id: 0, title: '', slug: '', content: '', status: 0, seo_description: '' }); show.value = true; }
async function save() { try {
    form.id ? await api.put(`/pages/${form.id}`, form) : await api.post('/pages', form);
    show.value = false;
    load();
}
catch (e) {
    message.error(e.message);
} }
const columns = [{ title: '标题', key: 'title' }, { title: 'Slug', key: 'slug' }, { title: '状态', key: 'status', render: (r) => r.status === 1 ? '已发布' : '草稿' }, { title: '操作', key: 'x', render: (r) => window.Vue?.h?.(NButton, { onClick: () => edit(r) }, () => '编辑') || r.title }];
onMounted(load); // @ts-ignore
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
        return (__VLS_ctx.edit());
        // @ts-ignore
        [edit,];
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
for (const [r] of __VLS_vFor((__VLS_ctx.rows))) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        key: (r.id),
        ...{ class: "release" },
        ...{ style: {} },
    });
    /** @type {__VLS_StyleScopedClasses['release']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.b, __VLS_intrinsics.b)({});
    (r.title);
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ style: {} },
    });
    (r.slug);
    (r.status === 1 ? '已发布' : '草稿');
    let __VLS_8;
    /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
    nButton;
    // @ts-ignore
    const __VLS_9 = __VLS_asFunctionalComponent1(__VLS_8, new __VLS_8({
        ...{ 'onClick': {} },
    }));
    const __VLS_10 = __VLS_9({
        ...{ 'onClick': {} },
    }, ...__VLS_functionalComponentArgsRest(__VLS_9));
    let __VLS_13;
    const __VLS_14 = {
        /** @type {typeof __VLS_13.click} */
        onClick: (...[$event]) => {
            return (__VLS_ctx.edit(r));
            // @ts-ignore
            [edit, rows,];
        },
    };
    const { default: __VLS_15 } = __VLS_11.slots;
    // @ts-ignore
    [];
    var __VLS_11;
    var __VLS_12;
    // @ts-ignore
    [];
}
let __VLS_16;
/** @ts-ignore @type { | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal'] | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal']} */
nModal;
// @ts-ignore
const __VLS_17 = __VLS_asFunctionalComponent1(__VLS_16, new __VLS_16({
    show: (__VLS_ctx.show),
    preset: "card",
    title: "编辑单页",
    ...{ style: {} },
}));
const __VLS_18 = __VLS_17({
    show: (__VLS_ctx.show),
    preset: "card",
    title: "编辑单页",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
const { default: __VLS_21 } = __VLS_19.slots;
let __VLS_22;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_23 = __VLS_asFunctionalComponent1(__VLS_22, new __VLS_22({
    labelPlacement: "top",
}));
const __VLS_24 = __VLS_23({
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_23));
const { default: __VLS_27 } = __VLS_25.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-grid" },
});
/** @type {__VLS_StyleScopedClasses['form-grid']} */ ;
let __VLS_28;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_29 = __VLS_asFunctionalComponent1(__VLS_28, new __VLS_28({
    label: "标题",
}));
const __VLS_30 = __VLS_29({
    label: "标题",
}, ...__VLS_functionalComponentArgsRest(__VLS_29));
const { default: __VLS_33 } = __VLS_31.slots;
let __VLS_34;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_35 = __VLS_asFunctionalComponent1(__VLS_34, new __VLS_34({
    value: (__VLS_ctx.form.title),
}));
const __VLS_36 = __VLS_35({
    value: (__VLS_ctx.form.title),
}, ...__VLS_functionalComponentArgsRest(__VLS_35));
// @ts-ignore
[show, form,];
var __VLS_31;
let __VLS_39;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_40 = __VLS_asFunctionalComponent1(__VLS_39, new __VLS_39({
    label: "Slug",
}));
const __VLS_41 = __VLS_40({
    label: "Slug",
}, ...__VLS_functionalComponentArgsRest(__VLS_40));
const { default: __VLS_44 } = __VLS_42.slots;
let __VLS_45;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_46 = __VLS_asFunctionalComponent1(__VLS_45, new __VLS_45({
    value: (__VLS_ctx.form.slug),
}));
const __VLS_47 = __VLS_46({
    value: (__VLS_ctx.form.slug),
}, ...__VLS_functionalComponentArgsRest(__VLS_46));
// @ts-ignore
[form,];
var __VLS_42;
let __VLS_50;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_51 = __VLS_asFunctionalComponent1(__VLS_50, new __VLS_50({
    label: "状态",
}));
const __VLS_52 = __VLS_51({
    label: "状态",
}, ...__VLS_functionalComponentArgsRest(__VLS_51));
const { default: __VLS_55 } = __VLS_53.slots;
let __VLS_56;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent1(__VLS_56, new __VLS_56({
    value: (__VLS_ctx.form.status),
    options: ([{ label: '草稿', value: 0 }, { label: '发布', value: 1 }]),
}));
const __VLS_58 = __VLS_57({
    value: (__VLS_ctx.form.status),
    options: ([{ label: '草稿', value: 0 }, { label: '发布', value: 1 }]),
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
// @ts-ignore
[form,];
var __VLS_53;
let __VLS_61;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_62 = __VLS_asFunctionalComponent1(__VLS_61, new __VLS_61({
    ...{ class: "wide" },
    label: "内容（Markdown）",
}));
const __VLS_63 = __VLS_62({
    ...{ class: "wide" },
    label: "内容（Markdown）",
}, ...__VLS_functionalComponentArgsRest(__VLS_62));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_66 } = __VLS_64.slots;
let __VLS_67;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_68 = __VLS_asFunctionalComponent1(__VLS_67, new __VLS_67({
    value: (__VLS_ctx.form.content),
    type: "textarea",
    rows: (12),
}));
const __VLS_69 = __VLS_68({
    value: (__VLS_ctx.form.content),
    type: "textarea",
    rows: (12),
}, ...__VLS_functionalComponentArgsRest(__VLS_68));
// @ts-ignore
[form,];
var __VLS_64;
let __VLS_72;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent1(__VLS_72, new __VLS_72({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_74 = __VLS_73({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
let __VLS_77;
const __VLS_78 = {
    /** @type {typeof __VLS_77.click} */
    onClick: (__VLS_ctx.save),
};
const { default: __VLS_79 } = __VLS_75.slots;
// @ts-ignore
[save,];
var __VLS_75;
var __VLS_76;
// @ts-ignore
[];
var __VLS_25;
// @ts-ignore
[];
var __VLS_19;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
