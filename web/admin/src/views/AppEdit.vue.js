import { onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useMessage } from 'naive-ui';
import { api } from '../api';
import { openAppPreview } from '../appPreview';
const route = useRoute();
const router = useRouter();
const message = useMessage();
const saving = ref(false);
const previewing = ref(false);
const cats = ref([]);
const form = reactive({
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
});
const platforms = ['windows', 'macos', 'linux', 'android', 'ios', 'web'];
onMounted(async () => {
    const x = await api.get('/categories');
    cats.value = x.list;
    if (route.params.id)
        Object.assign(form, await api.get(`/apps/${route.params.id}`));
});
function slug() {
    if (!route.params.id && !form.slug) {
        form.slug = form.name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
    }
}
async function save() {
    saving.value = true;
    try {
        const data = route.params.id ? await api.put(`/apps/${route.params.id}`, form) : await api.post('/apps', form);
        message.success('已保存');
        if (!route.params.id)
            router.replace(`/apps/${data.id}/edit`);
    }
    catch (e) {
        message.error(e.message);
    }
    finally {
        saving.value = false;
    }
}
async function preview() {
    if (!route.params.id)
        return;
    previewing.value = true;
    try {
        await openAppPreview(Number(route.params.id));
    }
    catch (e) {
        message.error(e.message || '预览打开失败');
    }
    finally {
        previewing.value = false;
    }
}
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
(__VLS_ctx.route.params.id ? '编辑应用' : '新建应用');
if (__VLS_ctx.route.params.id) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ class: "preview-hint" },
    });
    /** @type {__VLS_StyleScopedClasses['preview-hint']} */ ;
}
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
let __VLS_0;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
    ...{ 'onClick': {} },
    disabled: (!__VLS_ctx.route.params.id),
    loading: (__VLS_ctx.previewing),
}));
const __VLS_2 = __VLS_1({
    ...{ 'onClick': {} },
    disabled: (!__VLS_ctx.route.params.id),
    loading: (__VLS_ctx.previewing),
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_5;
const __VLS_6 = {
    /** @type {typeof __VLS_5.click} */
    onClick: (__VLS_ctx.preview),
};
const { default: __VLS_7 } = __VLS_3.slots;
// @ts-ignore
[route, route, route, previewing, preview,];
var __VLS_3;
var __VLS_4;
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
        return (__VLS_ctx.router.push('/apps'));
        // @ts-ignore
        [router,];
    },
};
const { default: __VLS_15 } = __VLS_11.slots;
// @ts-ignore
[];
var __VLS_11;
var __VLS_12;
let __VLS_16;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_17 = __VLS_asFunctionalComponent1(__VLS_16, new __VLS_16({
    ...{ 'onClick': {} },
    type: "primary",
    loading: (__VLS_ctx.saving),
}));
const __VLS_18 = __VLS_17({
    ...{ 'onClick': {} },
    type: "primary",
    loading: (__VLS_ctx.saving),
}, ...__VLS_functionalComponentArgsRest(__VLS_17));
let __VLS_21;
const __VLS_22 = {
    /** @type {typeof __VLS_21.click} */
    onClick: (__VLS_ctx.save),
};
const { default: __VLS_23 } = __VLS_19.slots;
// @ts-ignore
[saving, save,];
var __VLS_19;
var __VLS_20;
let __VLS_24;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_25 = __VLS_asFunctionalComponent1(__VLS_24, new __VLS_24({
    ...{ class: "panel" },
    labelPlacement: "top",
}));
const __VLS_26 = __VLS_25({
    ...{ class: "panel" },
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_25));
/** @type {__VLS_StyleScopedClasses['panel']} */ ;
const { default: __VLS_29 } = __VLS_27.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-grid" },
});
/** @type {__VLS_StyleScopedClasses['form-grid']} */ ;
let __VLS_30;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_31 = __VLS_asFunctionalComponent1(__VLS_30, new __VLS_30({
    label: "应用名称",
}));
const __VLS_32 = __VLS_31({
    label: "应用名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_31));
const { default: __VLS_35 } = __VLS_33.slots;
let __VLS_36;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_37 = __VLS_asFunctionalComponent1(__VLS_36, new __VLS_36({
    ...{ 'onBlur': {} },
    value: (__VLS_ctx.form.name),
}));
const __VLS_38 = __VLS_37({
    ...{ 'onBlur': {} },
    value: (__VLS_ctx.form.name),
}, ...__VLS_functionalComponentArgsRest(__VLS_37));
let __VLS_41;
const __VLS_42 = {
    /** @type {typeof __VLS_41.blur} */
    onBlur: (__VLS_ctx.slug),
};
var __VLS_39;
var __VLS_40;
// @ts-ignore
[form, slug,];
var __VLS_33;
let __VLS_43;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_44 = __VLS_asFunctionalComponent1(__VLS_43, new __VLS_43({
    label: "URL Slug",
}));
const __VLS_45 = __VLS_44({
    label: "URL Slug",
}, ...__VLS_functionalComponentArgsRest(__VLS_44));
const { default: __VLS_48 } = __VLS_46.slots;
let __VLS_49;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_50 = __VLS_asFunctionalComponent1(__VLS_49, new __VLS_49({
    value: (__VLS_ctx.form.slug),
    placeholder: "my-app",
}));
const __VLS_51 = __VLS_50({
    value: (__VLS_ctx.form.slug),
    placeholder: "my-app",
}, ...__VLS_functionalComponentArgsRest(__VLS_50));
// @ts-ignore
[form,];
var __VLS_46;
let __VLS_54;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_55 = __VLS_asFunctionalComponent1(__VLS_54, new __VLS_54({
    label: "类型",
}));
const __VLS_56 = __VLS_55({
    label: "类型",
}, ...__VLS_functionalComponentArgsRest(__VLS_55));
const { default: __VLS_59 } = __VLS_57.slots;
let __VLS_60;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_61 = __VLS_asFunctionalComponent1(__VLS_60, new __VLS_60({
    value: (__VLS_ctx.form.type),
    options: ([{ label: '自研', value: 1 }, { label: '收录', value: 2 }]),
}));
const __VLS_62 = __VLS_61({
    value: (__VLS_ctx.form.type),
    options: ([{ label: '自研', value: 1 }, { label: '收录', value: 2 }]),
}, ...__VLS_functionalComponentArgsRest(__VLS_61));
// @ts-ignore
[form,];
var __VLS_57;
let __VLS_65;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_66 = __VLS_asFunctionalComponent1(__VLS_65, new __VLS_65({
    label: "分类",
}));
const __VLS_67 = __VLS_66({
    label: "分类",
}, ...__VLS_functionalComponentArgsRest(__VLS_66));
const { default: __VLS_70 } = __VLS_68.slots;
let __VLS_71;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_72 = __VLS_asFunctionalComponent1(__VLS_71, new __VLS_71({
    value: (__VLS_ctx.form.category_id),
    clearable: true,
    options: (__VLS_ctx.cats.map(x => ({ label: x.name, value: x.id }))),
}));
const __VLS_73 = __VLS_72({
    value: (__VLS_ctx.form.category_id),
    clearable: true,
    options: (__VLS_ctx.cats.map(x => ({ label: x.name, value: x.id }))),
}, ...__VLS_functionalComponentArgsRest(__VLS_72));
// @ts-ignore
[form, cats,];
var __VLS_68;
let __VLS_76;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_77 = __VLS_asFunctionalComponent1(__VLS_76, new __VLS_76({
    ...{ class: "wide" },
    label: "一句话简介",
}));
const __VLS_78 = __VLS_77({
    ...{ class: "wide" },
    label: "一句话简介",
}, ...__VLS_functionalComponentArgsRest(__VLS_77));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_81 } = __VLS_79.slots;
let __VLS_82;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_83 = __VLS_asFunctionalComponent1(__VLS_82, new __VLS_82({
    value: (__VLS_ctx.form.tagline),
    maxlength: "200",
    showCount: true,
}));
const __VLS_84 = __VLS_83({
    value: (__VLS_ctx.form.tagline),
    maxlength: "200",
    showCount: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_83));
// @ts-ignore
[form,];
var __VLS_79;
let __VLS_87;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_88 = __VLS_asFunctionalComponent1(__VLS_87, new __VLS_87({
    ...{ class: "wide" },
    label: "详细介绍（Markdown）",
}));
const __VLS_89 = __VLS_88({
    ...{ class: "wide" },
    label: "详细介绍（Markdown）",
}, ...__VLS_functionalComponentArgsRest(__VLS_88));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_92 } = __VLS_90.slots;
let __VLS_93;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_94 = __VLS_asFunctionalComponent1(__VLS_93, new __VLS_93({
    value: (__VLS_ctx.form.description),
    type: "textarea",
    rows: (10),
}));
const __VLS_95 = __VLS_94({
    value: (__VLS_ctx.form.description),
    type: "textarea",
    rows: (10),
}, ...__VLS_functionalComponentArgsRest(__VLS_94));
// @ts-ignore
[form,];
var __VLS_90;
let __VLS_98;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_99 = __VLS_asFunctionalComponent1(__VLS_98, new __VLS_98({
    label: "图标 URL",
}));
const __VLS_100 = __VLS_99({
    label: "图标 URL",
}, ...__VLS_functionalComponentArgsRest(__VLS_99));
const { default: __VLS_103 } = __VLS_101.slots;
let __VLS_104;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_105 = __VLS_asFunctionalComponent1(__VLS_104, new __VLS_104({
    value: (__VLS_ctx.form.icon),
    placeholder: "/uploads/...",
}));
const __VLS_106 = __VLS_105({
    value: (__VLS_ctx.form.icon),
    placeholder: "/uploads/...",
}, ...__VLS_functionalComponentArgsRest(__VLS_105));
// @ts-ignore
[form,];
var __VLS_101;
let __VLS_109;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_110 = __VLS_asFunctionalComponent1(__VLS_109, new __VLS_109({
    label: "开发者",
}));
const __VLS_111 = __VLS_110({
    label: "开发者",
}, ...__VLS_functionalComponentArgsRest(__VLS_110));
const { default: __VLS_114 } = __VLS_112.slots;
let __VLS_115;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_116 = __VLS_asFunctionalComponent1(__VLS_115, new __VLS_115({
    value: (__VLS_ctx.form.developer),
}));
const __VLS_117 = __VLS_116({
    value: (__VLS_ctx.form.developer),
}, ...__VLS_functionalComponentArgsRest(__VLS_116));
// @ts-ignore
[form,];
var __VLS_112;
let __VLS_120;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_121 = __VLS_asFunctionalComponent1(__VLS_120, new __VLS_120({
    label: "官方网站",
}));
const __VLS_122 = __VLS_121({
    label: "官方网站",
}, ...__VLS_functionalComponentArgsRest(__VLS_121));
const { default: __VLS_125 } = __VLS_123.slots;
let __VLS_126;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_127 = __VLS_asFunctionalComponent1(__VLS_126, new __VLS_126({
    value: (__VLS_ctx.form.official_url),
}));
const __VLS_128 = __VLS_127({
    value: (__VLS_ctx.form.official_url),
}, ...__VLS_functionalComponentArgsRest(__VLS_127));
// @ts-ignore
[form,];
var __VLS_123;
let __VLS_131;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_132 = __VLS_asFunctionalComponent1(__VLS_131, new __VLS_131({
    label: "代码仓库",
}));
const __VLS_133 = __VLS_132({
    label: "代码仓库",
}, ...__VLS_functionalComponentArgsRest(__VLS_132));
const { default: __VLS_136 } = __VLS_134.slots;
let __VLS_137;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_138 = __VLS_asFunctionalComponent1(__VLS_137, new __VLS_137({
    value: (__VLS_ctx.form.repo_url),
}));
const __VLS_139 = __VLS_138({
    value: (__VLS_ctx.form.repo_url),
}, ...__VLS_functionalComponentArgsRest(__VLS_138));
// @ts-ignore
[form,];
var __VLS_134;
let __VLS_142;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_143 = __VLS_asFunctionalComponent1(__VLS_142, new __VLS_142({
    label: "许可",
}));
const __VLS_144 = __VLS_143({
    label: "许可",
}, ...__VLS_functionalComponentArgsRest(__VLS_143));
const { default: __VLS_147 } = __VLS_145.slots;
let __VLS_148;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_149 = __VLS_asFunctionalComponent1(__VLS_148, new __VLS_148({
    value: (__VLS_ctx.form.license),
}));
const __VLS_150 = __VLS_149({
    value: (__VLS_ctx.form.license),
}, ...__VLS_functionalComponentArgsRest(__VLS_149));
// @ts-ignore
[form,];
var __VLS_145;
let __VLS_153;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_154 = __VLS_asFunctionalComponent1(__VLS_153, new __VLS_153({
    ...{ class: "wide" },
    label: "支持平台",
}));
const __VLS_155 = __VLS_154({
    ...{ class: "wide" },
    label: "支持平台",
}, ...__VLS_functionalComponentArgsRest(__VLS_154));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_158 } = __VLS_156.slots;
for (const [p] of __VLS_vFor((__VLS_ctx.platforms))) {
    let __VLS_159;
    /** @ts-ignore @type { | typeof __VLS_components.nCheckbox | typeof __VLS_components.NCheckbox | typeof __VLS_components['n-checkbox'] | typeof __VLS_components.nCheckbox | typeof __VLS_components.NCheckbox | typeof __VLS_components['n-checkbox']} */
    nCheckbox;
    // @ts-ignore
    const __VLS_160 = __VLS_asFunctionalComponent1(__VLS_159, new __VLS_159({
        ...{ 'onUpdate:checked': {} },
        key: (p),
        checked: (__VLS_ctx.form.platforms.includes(p)),
    }));
    const __VLS_161 = __VLS_160({
        ...{ 'onUpdate:checked': {} },
        key: (p),
        checked: (__VLS_ctx.form.platforms.includes(p)),
    }, ...__VLS_functionalComponentArgsRest(__VLS_160));
    let __VLS_164;
    const __VLS_165 = {
        /** @type {typeof __VLS_164.'update:checked'} */
        'onUpdate:checked': ((v) => __VLS_ctx.form.platforms = v ? [...__VLS_ctx.form.platforms, p] : __VLS_ctx.form.platforms.filter((x) => x !== p)),
    };
    const { default: __VLS_166 } = __VLS_162.slots;
    (p);
    // @ts-ignore
    [form, form, form, form, platforms,];
    var __VLS_162;
    var __VLS_163;
    // @ts-ignore
    [];
}
// @ts-ignore
[];
var __VLS_156;
// @ts-ignore
[];
var __VLS_27;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
