import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { bootstrapAuth } from './api';
import { currentUser } from './state';
const route = useRoute(), router = useRouter(), ready = ref(false), collapsed = ref(false);
const menu = [{ label: '仪表盘', key: '/' }, { label: '应用管理', key: '/apps' }, { label: '存储管理', key: '/storages' }, { label: '主题外观', key: '/themes' }, { label: '单页管理', key: '/pages' }, { label: '站点设置', key: '/settings' }];
const active = computed(() => route.path.startsWith('/apps') ? '/apps' : route.path);
onMounted(async () => { currentUser.value = await bootstrapAuth(); ready.value = true; if (!currentUser.value && route.path != '/login')
    router.replace('/login'); });
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
if (!__VLS_ctx.ready) {
    let __VLS_0;
    /** @ts-ignore @type { | typeof __VLS_components.nSpin | typeof __VLS_components.NSpin | typeof __VLS_components['n-spin']} */
    nSpin;
    // @ts-ignore
    const __VLS_1 = __VLS_asFunctionalComponent1(__VLS_0, new __VLS_0({
        ...{ class: "boot" },
        size: "large",
    }));
    const __VLS_2 = __VLS_1({
        ...{ class: "boot" },
        size: "large",
    }, ...__VLS_functionalComponentArgsRest(__VLS_1));
    var __VLS_5;
    /** @type {__VLS_StyleScopedClasses['boot']} */ ;
    var __VLS_3;
}
else if (__VLS_ctx.route.path === '/login') {
    let __VLS_6;
    /** @ts-ignore @type { | typeof __VLS_components.routerView | typeof __VLS_components.RouterView | typeof __VLS_components['router-view']} */
    routerView;
    // @ts-ignore
    const __VLS_7 = __VLS_asFunctionalComponent1(__VLS_6, new __VLS_6({}));
    const __VLS_8 = __VLS_7({}, ...__VLS_functionalComponentArgsRest(__VLS_7));
    var __VLS_11;
    var __VLS_9;
}
else {
    let __VLS_12;
    /** @ts-ignore @type { | typeof __VLS_components.nLayout | typeof __VLS_components.NLayout | typeof __VLS_components['n-layout'] | typeof __VLS_components.nLayout | typeof __VLS_components.NLayout | typeof __VLS_components['n-layout']} */
    nLayout;
    // @ts-ignore
    const __VLS_13 = __VLS_asFunctionalComponent1(__VLS_12, new __VLS_12({
        hasSider: true,
        ...{ class: "shell" },
    }));
    const __VLS_14 = __VLS_13({
        hasSider: true,
        ...{ class: "shell" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_13));
    var __VLS_17;
    /** @type {__VLS_StyleScopedClasses['shell']} */ ;
    const { default: __VLS_18 } = __VLS_15.slots;
    let __VLS_19;
    /** @ts-ignore @type { | typeof __VLS_components.nLayoutSider | typeof __VLS_components.NLayoutSider | typeof __VLS_components['n-layout-sider'] | typeof __VLS_components.nLayoutSider | typeof __VLS_components.NLayoutSider | typeof __VLS_components['n-layout-sider']} */
    nLayoutSider;
    // @ts-ignore
    const __VLS_20 = __VLS_asFunctionalComponent1(__VLS_19, new __VLS_19({
        bordered: true,
        collapseMode: "width",
        collapsedWidth: (72),
        width: (244),
        collapsed: (__VLS_ctx.collapsed),
    }));
    const __VLS_21 = __VLS_20({
        bordered: true,
        collapseMode: "width",
        collapsedWidth: (72),
        width: (244),
        collapsed: (__VLS_ctx.collapsed),
    }, ...__VLS_functionalComponentArgsRest(__VLS_20));
    const { default: __VLS_24 } = __VLS_22.slots;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "logo" },
    });
    /** @type {__VLS_StyleScopedClasses['logo']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    if (!__VLS_ctx.collapsed) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.b, __VLS_intrinsics.b)({});
    }
    let __VLS_25;
    /** @ts-ignore @type { | typeof __VLS_components.nMenu | typeof __VLS_components.NMenu | typeof __VLS_components['n-menu']} */
    nMenu;
    // @ts-ignore
    const __VLS_26 = __VLS_asFunctionalComponent1(__VLS_25, new __VLS_25({
        ...{ 'onUpdate:value': {} },
        value: (__VLS_ctx.active),
        options: (__VLS_ctx.menu),
    }));
    const __VLS_27 = __VLS_26({
        ...{ 'onUpdate:value': {} },
        value: (__VLS_ctx.active),
        options: (__VLS_ctx.menu),
    }, ...__VLS_functionalComponentArgsRest(__VLS_26));
    let __VLS_30;
    const __VLS_31 = {
        /** @type {typeof __VLS_30.'update:value'} */
        'onUpdate:value': (__VLS_ctx.router.push),
    };
    var __VLS_28;
    var __VLS_29;
    __VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
        ...{ onClick: (...[$event]) => {
                if (!!(!__VLS_ctx.ready))
                    throw 0;
                if (!!(__VLS_ctx.route.path === '/login'))
                    throw 0;
                return (__VLS_ctx.collapsed = !__VLS_ctx.collapsed);
                // @ts-ignore
                [ready, route, collapsed, collapsed, collapsed, collapsed, active, menu, router,];
            } },
        ...{ class: "collapse" },
    });
    /** @type {__VLS_StyleScopedClasses['collapse']} */ ;
    (__VLS_ctx.collapsed ? '›' : '‹');
    // @ts-ignore
    [collapsed,];
    var __VLS_22;
    let __VLS_32;
    /** @ts-ignore @type { | typeof __VLS_components.nLayout | typeof __VLS_components.NLayout | typeof __VLS_components['n-layout'] | typeof __VLS_components.nLayout | typeof __VLS_components.NLayout | typeof __VLS_components['n-layout']} */
    nLayout;
    // @ts-ignore
    const __VLS_33 = __VLS_asFunctionalComponent1(__VLS_32, new __VLS_32({}));
    const __VLS_34 = __VLS_33({}, ...__VLS_functionalComponentArgsRest(__VLS_33));
    const { default: __VLS_37 } = __VLS_35.slots;
    let __VLS_38;
    /** @ts-ignore @type { | typeof __VLS_components.nLayoutHeader | typeof __VLS_components.NLayoutHeader | typeof __VLS_components['n-layout-header'] | typeof __VLS_components.nLayoutHeader | typeof __VLS_components.NLayoutHeader | typeof __VLS_components['n-layout-header']} */
    nLayoutHeader;
    // @ts-ignore
    const __VLS_39 = __VLS_asFunctionalComponent1(__VLS_38, new __VLS_38({
        bordered: true,
        ...{ class: "header" },
    }));
    const __VLS_40 = __VLS_39({
        bordered: true,
        ...{ class: "header" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_39));
    /** @type {__VLS_StyleScopedClasses['header']} */ ;
    const { default: __VLS_43 } = __VLS_41.slots;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.b, __VLS_intrinsics.b)({});
    (__VLS_ctx.menu.find(x => x.key === __VLS_ctx.active)?.label || '内容编辑');
    __VLS_asFunctionalElement1(__VLS_intrinsics.small, __VLS_intrinsics.small)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ class: "profile" },
    });
    /** @type {__VLS_StyleScopedClasses['profile']} */ ;
    (__VLS_ctx.currentUser?.nickname || __VLS_ctx.currentUser?.username);
    let __VLS_44;
    /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
    nButton;
    // @ts-ignore
    const __VLS_45 = __VLS_asFunctionalComponent1(__VLS_44, new __VLS_44({
        ...{ 'onClick': {} },
        quaternary: true,
    }));
    const __VLS_46 = __VLS_45({
        ...{ 'onClick': {} },
        quaternary: true,
    }, ...__VLS_functionalComponentArgsRest(__VLS_45));
    let __VLS_49;
    const __VLS_50 = {
        /** @type {typeof __VLS_49.click} */
        onClick: (...[$event]) => {
            if (!!(!__VLS_ctx.ready))
                throw 0;
            if (!!(__VLS_ctx.route.path === '/login'))
                throw 0;
            __VLS_ctx.currentUser = null;
            __VLS_ctx.router.push('/login');
            // @ts-ignore
            [active, menu, router, currentUser, currentUser, currentUser,];
        },
    };
    const { default: __VLS_51 } = __VLS_47.slots;
    // @ts-ignore
    [];
    var __VLS_47;
    var __VLS_48;
    // @ts-ignore
    [];
    var __VLS_41;
    let __VLS_52;
    /** @ts-ignore @type { | typeof __VLS_components.nLayoutContent | typeof __VLS_components.NLayoutContent | typeof __VLS_components['n-layout-content'] | typeof __VLS_components.nLayoutContent | typeof __VLS_components.NLayoutContent | typeof __VLS_components['n-layout-content']} */
    nLayoutContent;
    // @ts-ignore
    const __VLS_53 = __VLS_asFunctionalComponent1(__VLS_52, new __VLS_52({
        ...{ class: "content" },
    }));
    const __VLS_54 = __VLS_53({
        ...{ class: "content" },
    }, ...__VLS_functionalComponentArgsRest(__VLS_53));
    /** @type {__VLS_StyleScopedClasses['content']} */ ;
    const { default: __VLS_57 } = __VLS_55.slots;
    let __VLS_58;
    /** @ts-ignore @type { | typeof __VLS_components.routerView | typeof __VLS_components.RouterView | typeof __VLS_components['router-view']} */
    routerView;
    // @ts-ignore
    const __VLS_59 = __VLS_asFunctionalComponent1(__VLS_58, new __VLS_58({}));
    const __VLS_60 = __VLS_59({}, ...__VLS_functionalComponentArgsRest(__VLS_59));
    // @ts-ignore
    [];
    var __VLS_55;
    // @ts-ignore
    [];
    var __VLS_35;
    // @ts-ignore
    [];
    var __VLS_15;
}
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
