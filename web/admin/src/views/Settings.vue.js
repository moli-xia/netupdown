import { onMounted, reactive } from 'vue';
import { useMessage } from 'naive-ui';
import { api } from '../api';
const message = useMessage(), form = reactive({});
onMounted(async () => Object.assign(form, await api.get('/settings')));
async function save() { try {
    await api.put('/settings', form);
    message.success('设置已保存，前台立即生效');
}
catch (e) {
    message.error(e.message);
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
    onClick: (__VLS_ctx.save),
};
const { default: __VLS_7 } = __VLS_3.slots;
// @ts-ignore
[save,];
var __VLS_3;
var __VLS_4;
let __VLS_8;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent1(__VLS_8, new __VLS_8({
    ...{ class: "panel" },
    labelPlacement: "top",
}));
const __VLS_10 = __VLS_9({
    ...{ class: "panel" },
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
/** @type {__VLS_StyleScopedClasses['panel']} */ ;
const { default: __VLS_13 } = __VLS_11.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-grid" },
});
/** @type {__VLS_StyleScopedClasses['form-grid']} */ ;
let __VLS_14;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_15 = __VLS_asFunctionalComponent1(__VLS_14, new __VLS_14({
    label: "站点名称",
}));
const __VLS_16 = __VLS_15({
    label: "站点名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_15));
const { default: __VLS_19 } = __VLS_17.slots;
let __VLS_20;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_21 = __VLS_asFunctionalComponent1(__VLS_20, new __VLS_20({
    value: (__VLS_ctx.form['site.title']),
}));
const __VLS_22 = __VLS_21({
    value: (__VLS_ctx.form['site.title']),
}, ...__VLS_functionalComponentArgsRest(__VLS_21));
// @ts-ignore
[form,];
var __VLS_17;
let __VLS_25;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_26 = __VLS_asFunctionalComponent1(__VLS_25, new __VLS_25({
    label: "副标题",
}));
const __VLS_27 = __VLS_26({
    label: "副标题",
}, ...__VLS_functionalComponentArgsRest(__VLS_26));
const { default: __VLS_30 } = __VLS_28.slots;
let __VLS_31;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_32 = __VLS_asFunctionalComponent1(__VLS_31, new __VLS_31({
    value: (__VLS_ctx.form['site.subtitle']),
}));
const __VLS_33 = __VLS_32({
    value: (__VLS_ctx.form['site.subtitle']),
}, ...__VLS_functionalComponentArgsRest(__VLS_32));
// @ts-ignore
[form,];
var __VLS_28;
let __VLS_36;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_37 = __VLS_asFunctionalComponent1(__VLS_36, new __VLS_36({
    ...{ class: "wide" },
    label: "站点描述",
}));
const __VLS_38 = __VLS_37({
    ...{ class: "wide" },
    label: "站点描述",
}, ...__VLS_functionalComponentArgsRest(__VLS_37));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_41 } = __VLS_39.slots;
let __VLS_42;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_43 = __VLS_asFunctionalComponent1(__VLS_42, new __VLS_42({
    value: (__VLS_ctx.form['site.description']),
    type: "textarea",
}));
const __VLS_44 = __VLS_43({
    value: (__VLS_ctx.form['site.description']),
    type: "textarea",
}, ...__VLS_functionalComponentArgsRest(__VLS_43));
// @ts-ignore
[form,];
var __VLS_39;
let __VLS_47;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_48 = __VLS_asFunctionalComponent1(__VLS_47, new __VLS_47({
    label: "备案号",
}));
const __VLS_49 = __VLS_48({
    label: "备案号",
}, ...__VLS_functionalComponentArgsRest(__VLS_48));
const { default: __VLS_52 } = __VLS_50.slots;
let __VLS_53;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_54 = __VLS_asFunctionalComponent1(__VLS_53, new __VLS_53({
    value: (__VLS_ctx.form['site.icp']),
}));
const __VLS_55 = __VLS_54({
    value: (__VLS_ctx.form['site.icp']),
}, ...__VLS_functionalComponentArgsRest(__VLS_54));
// @ts-ignore
[form,];
var __VLS_50;
let __VLS_58;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_59 = __VLS_asFunctionalComponent1(__VLS_58, new __VLS_58({
    label: "列表每页数量",
}));
const __VLS_60 = __VLS_59({
    label: "列表每页数量",
}, ...__VLS_functionalComponentArgsRest(__VLS_59));
const { default: __VLS_63 } = __VLS_61.slots;
let __VLS_64;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_65 = __VLS_asFunctionalComponent1(__VLS_64, new __VLS_64({
    value: (__VLS_ctx.form['content.apps_page_size']),
}));
const __VLS_66 = __VLS_65({
    value: (__VLS_ctx.form['content.apps_page_size']),
}, ...__VLS_functionalComponentArgsRest(__VLS_65));
// @ts-ignore
[form,];
var __VLS_61;
let __VLS_69;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_70 = __VLS_asFunctionalComponent1(__VLS_69, new __VLS_69({
    ...{ class: "wide" },
    label: "页脚内容",
}));
const __VLS_71 = __VLS_70({
    ...{ class: "wide" },
    label: "页脚内容",
}, ...__VLS_functionalComponentArgsRest(__VLS_70));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_74 } = __VLS_72.slots;
let __VLS_75;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_76 = __VLS_asFunctionalComponent1(__VLS_75, new __VLS_75({
    value: (__VLS_ctx.form['site.footer']),
    type: "textarea",
}));
const __VLS_77 = __VLS_76({
    value: (__VLS_ctx.form['site.footer']),
    type: "textarea",
}, ...__VLS_functionalComponentArgsRest(__VLS_76));
// @ts-ignore
[form,];
var __VLS_72;
let __VLS_80;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_81 = __VLS_asFunctionalComponent1(__VLS_80, new __VLS_80({
    ...{ class: "wide" },
    label: "Head 注入",
}));
const __VLS_82 = __VLS_81({
    ...{ class: "wide" },
    label: "Head 注入",
}, ...__VLS_functionalComponentArgsRest(__VLS_81));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_85 } = __VLS_83.slots;
let __VLS_86;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_87 = __VLS_asFunctionalComponent1(__VLS_86, new __VLS_86({
    value: (__VLS_ctx.form['custom.head']),
    type: "textarea",
}));
const __VLS_88 = __VLS_87({
    value: (__VLS_ctx.form['custom.head']),
    type: "textarea",
}, ...__VLS_functionalComponentArgsRest(__VLS_87));
// @ts-ignore
[form,];
var __VLS_83;
// @ts-ignore
[];
var __VLS_11;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
