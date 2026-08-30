import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage } from 'naive-ui';
import { api, setToken } from '../api';
import { currentUser } from '../state';
const router = useRouter(), message = useMessage(), loading = ref(false), form = reactive({ username: '', password: '' });
async function submit() { loading.value = true; try {
    const data = await api.post('/auth/login', form);
    setToken(data.access_token);
    currentUser.value = data.user;
    router.push('/');
}
catch (e) {
    message.error(e.message);
}
finally {
    loading.value = false;
} } // @ts-ignore
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "login" },
});
/** @type {__VLS_StyleScopedClasses['login']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.section, __VLS_intrinsics.section)({
    ...{ class: "login-art" },
});
/** @type {__VLS_StyleScopedClasses['login-art']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.section, __VLS_intrinsics.section)({
    ...{ class: "login-box" },
});
/** @type {__VLS_StyleScopedClasses['login-box']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "login-card" },
});
/** @type {__VLS_StyleScopedClasses['login-card']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
    ...{ style: {} },
});
__VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ style: {} },
});
let __VLS_0;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
    ...{ 'onSubmit': {} },
    size: "large",
}));
const __VLS_2 = __VLS_1({
    ...{ 'onSubmit': {} },
    size: "large",
}, ...__VLS_functionalComponentArgsRest(__VLS_1));
let __VLS_5;
const __VLS_6 = {
    /** @type {typeof __VLS_5.submit} */
    onSubmit: (__VLS_ctx.submit),
};
const { default: __VLS_7 } = __VLS_3.slots;
let __VLS_8;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_9 = __VLS_asFunctionalComponent1(__VLS_8, new __VLS_8({
    label: "用户名",
}));
const __VLS_10 = __VLS_9({
    label: "用户名",
}, ...__VLS_functionalComponentArgsRest(__VLS_9));
const { default: __VLS_13 } = __VLS_11.slots;
let __VLS_14;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_15 = __VLS_asFunctionalComponent1(__VLS_14, new __VLS_14({
    value: (__VLS_ctx.form.username),
    autofocus: true,
}));
const __VLS_16 = __VLS_15({
    value: (__VLS_ctx.form.username),
    autofocus: true,
}, ...__VLS_functionalComponentArgsRest(__VLS_15));
// @ts-ignore
[submit, form,];
var __VLS_11;
let __VLS_19;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_20 = __VLS_asFunctionalComponent1(__VLS_19, new __VLS_19({
    label: "密码",
}));
const __VLS_21 = __VLS_20({
    label: "密码",
}, ...__VLS_functionalComponentArgsRest(__VLS_20));
const { default: __VLS_24 } = __VLS_22.slots;
let __VLS_25;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_26 = __VLS_asFunctionalComponent1(__VLS_25, new __VLS_25({
    ...{ 'onKeyup': {} },
    value: (__VLS_ctx.form.password),
    type: "password",
    showPasswordOn: "click",
}));
const __VLS_27 = __VLS_26({
    ...{ 'onKeyup': {} },
    value: (__VLS_ctx.form.password),
    type: "password",
    showPasswordOn: "click",
}, ...__VLS_functionalComponentArgsRest(__VLS_26));
let __VLS_30;
const __VLS_31 = {
    /** @type {typeof __VLS_30.keyup} */
    onKeyup: (__VLS_ctx.submit),
};
var __VLS_28;
var __VLS_29;
// @ts-ignore
[submit, form,];
var __VLS_22;
let __VLS_32;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_33 = __VLS_asFunctionalComponent1(__VLS_32, new __VLS_32({
    ...{ 'onClick': {} },
    block: true,
    type: "primary",
    loading: (__VLS_ctx.loading),
}));
const __VLS_34 = __VLS_33({
    ...{ 'onClick': {} },
    block: true,
    type: "primary",
    loading: (__VLS_ctx.loading),
}, ...__VLS_functionalComponentArgsRest(__VLS_33));
let __VLS_37;
const __VLS_38 = {
    /** @type {typeof __VLS_37.click} */
    onClick: (__VLS_ctx.submit),
};
const { default: __VLS_39 } = __VLS_35.slots;
// @ts-ignore
[submit, loading,];
var __VLS_35;
var __VLS_36;
// @ts-ignore
[];
var __VLS_3;
var __VLS_4;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
