import { onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useMessage } from 'naive-ui';
import { api } from '../api';
const route = useRoute(), message = useMessage(), rows = ref([]), show = ref(false), assetFor = ref(null), sourceFor = ref(null), form = reactive({ version: '', channel: 1, title: '', changelog: '', min_required_version: '' }), asset = reactive({ name: '', file_name: '', os: 'windows', arch: 'amd64', kind: 1, size: 0, sha256: '' }), source = reactive({ name: '官方下载', source_type: 2, external_url: '', extract_code: '', priority: 0, is_enabled: true });
async function load() { const x = await api.get(`/apps/${route.params.id}/releases`); rows.value = x.list; }
async function addRelease() { await api.post(`/apps/${route.params.id}/releases`, form); show.value = false; Object.assign(form, { version: '', channel: 1, title: '', changelog: '', min_required_version: '' }); load(); }
async function publish(id) { try {
    await api.post(`/releases/${id}/publish`);
    message.success('版本已发布');
    load();
}
catch (e) {
    message.error(e.message);
} }
async function addAsset() { const x = await api.post(`/releases/${assetFor.value.id}/assets`, asset); assetFor.value = null; sourceFor.value = x; Object.assign(source, { name: '官方下载', source_type: 2, external_url: '', extract_code: '', priority: 0, is_enabled: true }); load(); }
async function addSource() { await api.post(`/assets/${sourceFor.value.id}/sources`, source); sourceFor.value = null; load(); }
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
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
    ...{ style: {} },
});
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
        return (__VLS_ctx.show = true);
        // @ts-ignore
        [show,];
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
    });
    /** @type {__VLS_StyleScopedClasses['release']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ style: {} },
    });
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    __VLS_asFunctionalElement1(__VLS_intrinsics.h3, __VLS_intrinsics.h3)({
        ...{ style: {} },
    });
    (r.version);
    let __VLS_8;
    /** @ts-ignore @type { | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag'] | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag']} */
    nTag;
    // @ts-ignore
    const __VLS_9 = __VLS_asFunctionalComponent1(__VLS_8, new __VLS_8({
        size: "small",
    }));
    const __VLS_10 = __VLS_9({
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_9));
    const { default: __VLS_13 } = __VLS_11.slots;
    (['', 'stable', 'beta', 'alpha'][r.channel]);
    // @ts-ignore
    [rows,];
    var __VLS_11;
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
        ...{ style: {} },
    });
    (r.title || '未填写标题');
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    let __VLS_14;
    /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
    nButton;
    // @ts-ignore
    const __VLS_15 = __VLS_asFunctionalComponent1(__VLS_14, new __VLS_14({
        ...{ 'onClick': {} },
        size: "small",
    }));
    const __VLS_16 = __VLS_15({
        ...{ 'onClick': {} },
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_15));
    let __VLS_19;
    const __VLS_20 = {
        /** @type {typeof __VLS_19.click} */
        onClick: (...[$event]) => {
            return (__VLS_ctx.assetFor = r);
            // @ts-ignore
            [assetFor,];
        },
    };
    const { default: __VLS_21 } = __VLS_17.slots;
    // @ts-ignore
    [];
    var __VLS_17;
    var __VLS_18;
    let __VLS_22;
    /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
    nButton;
    // @ts-ignore
    const __VLS_23 = __VLS_asFunctionalComponent1(__VLS_22, new __VLS_22({
        ...{ 'onClick': {} },
        size: "small",
        type: "primary",
        disabled: (r.status === 1),
    }));
    const __VLS_24 = __VLS_23({
        ...{ 'onClick': {} },
        size: "small",
        type: "primary",
        disabled: (r.status === 1),
    }, ...__VLS_functionalComponentArgsRest(__VLS_23));
    let __VLS_27;
    const __VLS_28 = {
        /** @type {typeof __VLS_27.click} */
        onClick: (...[$event]) => {
            return (__VLS_ctx.publish(r.id));
            // @ts-ignore
            [publish,];
        },
    };
    const { default: __VLS_29 } = __VLS_25.slots;
    (r.status === 1 ? '已发布' : '发布版本');
    // @ts-ignore
    [];
    var __VLS_25;
    var __VLS_26;
    __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({});
    (r.changelog);
    for (const [a] of __VLS_vFor((r.assets))) {
        __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
            key: (a.id),
            ...{ class: "asset-row" },
        });
        /** @type {__VLS_StyleScopedClasses['asset-row']} */ ;
        __VLS_asFunctionalElement1(__VLS_intrinsics.b, __VLS_intrinsics.b)({});
        (a.name);
        (a.os);
        (a.arch);
        (a.file_name);
        let __VLS_30;
        /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
        nButton;
        // @ts-ignore
        const __VLS_31 = __VLS_asFunctionalComponent1(__VLS_30, new __VLS_30({
            ...{ 'onClick': {} },
            text: true,
            type: "primary",
            ...{ style: {} },
        }));
        const __VLS_32 = __VLS_31({
            ...{ 'onClick': {} },
            text: true,
            type: "primary",
            ...{ style: {} },
        }, ...__VLS_functionalComponentArgsRest(__VLS_31));
        let __VLS_35;
        const __VLS_36 = {
            /** @type {typeof __VLS_35.click} */
            onClick: (...[$event]) => {
                return (__VLS_ctx.sourceFor = a);
                // @ts-ignore
                [sourceFor,];
            },
        };
        const { default: __VLS_37 } = __VLS_33.slots;
        // @ts-ignore
        [];
        var __VLS_33;
        var __VLS_34;
        if (a.sources?.length) {
            __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
                ...{ style: {} },
            });
            (a.sources.map((s) => s.name).join('、'));
        }
        // @ts-ignore
        [];
    }
    // @ts-ignore
    [];
}
if (!__VLS_ctx.rows.length) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ style: {} },
    });
}
let __VLS_38;
/** @ts-ignore @type { | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal'] | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal']} */
nModal;
// @ts-ignore
const __VLS_39 = __VLS_asFunctionalComponent1(__VLS_38, new __VLS_38({
    show: (__VLS_ctx.show),
    preset: "card",
    title: "新建版本",
    ...{ style: {} },
}));
const __VLS_40 = __VLS_39({
    show: (__VLS_ctx.show),
    preset: "card",
    title: "新建版本",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_39));
const { default: __VLS_43 } = __VLS_41.slots;
let __VLS_44;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_45 = __VLS_asFunctionalComponent1(__VLS_44, new __VLS_44({
    labelPlacement: "top",
}));
const __VLS_46 = __VLS_45({
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_45));
const { default: __VLS_49 } = __VLS_47.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-grid" },
});
/** @type {__VLS_StyleScopedClasses['form-grid']} */ ;
let __VLS_50;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_51 = __VLS_asFunctionalComponent1(__VLS_50, new __VLS_50({
    label: "版本号",
}));
const __VLS_52 = __VLS_51({
    label: "版本号",
}, ...__VLS_functionalComponentArgsRest(__VLS_51));
const { default: __VLS_55 } = __VLS_53.slots;
let __VLS_56;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_57 = __VLS_asFunctionalComponent1(__VLS_56, new __VLS_56({
    value: (__VLS_ctx.form.version),
    placeholder: "1.0.0",
}));
const __VLS_58 = __VLS_57({
    value: (__VLS_ctx.form.version),
    placeholder: "1.0.0",
}, ...__VLS_functionalComponentArgsRest(__VLS_57));
// @ts-ignore
[show, rows, form,];
var __VLS_53;
let __VLS_61;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_62 = __VLS_asFunctionalComponent1(__VLS_61, new __VLS_61({
    label: "渠道",
}));
const __VLS_63 = __VLS_62({
    label: "渠道",
}, ...__VLS_functionalComponentArgsRest(__VLS_62));
const { default: __VLS_66 } = __VLS_64.slots;
let __VLS_67;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_68 = __VLS_asFunctionalComponent1(__VLS_67, new __VLS_67({
    value: (__VLS_ctx.form.channel),
    options: ([{ label: 'stable', value: 1 }, { label: 'beta', value: 2 }, { label: 'alpha', value: 3 }]),
}));
const __VLS_69 = __VLS_68({
    value: (__VLS_ctx.form.channel),
    options: ([{ label: 'stable', value: 1 }, { label: 'beta', value: 2 }, { label: 'alpha', value: 3 }]),
}, ...__VLS_functionalComponentArgsRest(__VLS_68));
// @ts-ignore
[form,];
var __VLS_64;
let __VLS_72;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_73 = __VLS_asFunctionalComponent1(__VLS_72, new __VLS_72({
    ...{ class: "wide" },
    label: "标题",
}));
const __VLS_74 = __VLS_73({
    ...{ class: "wide" },
    label: "标题",
}, ...__VLS_functionalComponentArgsRest(__VLS_73));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_77 } = __VLS_75.slots;
let __VLS_78;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_79 = __VLS_asFunctionalComponent1(__VLS_78, new __VLS_78({
    value: (__VLS_ctx.form.title),
}));
const __VLS_80 = __VLS_79({
    value: (__VLS_ctx.form.title),
}, ...__VLS_functionalComponentArgsRest(__VLS_79));
// @ts-ignore
[form,];
var __VLS_75;
let __VLS_83;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_84 = __VLS_asFunctionalComponent1(__VLS_83, new __VLS_83({
    ...{ class: "wide" },
    label: "更新日志（Markdown）",
}));
const __VLS_85 = __VLS_84({
    ...{ class: "wide" },
    label: "更新日志（Markdown）",
}, ...__VLS_functionalComponentArgsRest(__VLS_84));
/** @type {__VLS_StyleScopedClasses['wide']} */ ;
const { default: __VLS_88 } = __VLS_86.slots;
let __VLS_89;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_90 = __VLS_asFunctionalComponent1(__VLS_89, new __VLS_89({
    value: (__VLS_ctx.form.changelog),
    type: "textarea",
    rows: (7),
}));
const __VLS_91 = __VLS_90({
    value: (__VLS_ctx.form.changelog),
    type: "textarea",
    rows: (7),
}, ...__VLS_functionalComponentArgsRest(__VLS_90));
// @ts-ignore
[form,];
var __VLS_86;
let __VLS_94;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_95 = __VLS_asFunctionalComponent1(__VLS_94, new __VLS_94({
    label: "强制更新阈值",
}));
const __VLS_96 = __VLS_95({
    label: "强制更新阈值",
}, ...__VLS_functionalComponentArgsRest(__VLS_95));
const { default: __VLS_99 } = __VLS_97.slots;
let __VLS_100;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_101 = __VLS_asFunctionalComponent1(__VLS_100, new __VLS_100({
    value: (__VLS_ctx.form.min_required_version),
}));
const __VLS_102 = __VLS_101({
    value: (__VLS_ctx.form.min_required_version),
}, ...__VLS_functionalComponentArgsRest(__VLS_101));
// @ts-ignore
[form,];
var __VLS_97;
let __VLS_105;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_106 = __VLS_asFunctionalComponent1(__VLS_105, new __VLS_105({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_107 = __VLS_106({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_106));
let __VLS_110;
const __VLS_111 = {
    /** @type {typeof __VLS_110.click} */
    onClick: (__VLS_ctx.addRelease),
};
const { default: __VLS_112 } = __VLS_108.slots;
// @ts-ignore
[addRelease,];
var __VLS_108;
var __VLS_109;
// @ts-ignore
[];
var __VLS_47;
// @ts-ignore
[];
var __VLS_41;
let __VLS_113;
/** @ts-ignore @type { | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal'] | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal']} */
nModal;
// @ts-ignore
const __VLS_114 = __VLS_asFunctionalComponent1(__VLS_113, new __VLS_113({
    ...{ 'onUpdate:show': {} },
    show: (!!__VLS_ctx.assetFor),
    preset: "card",
    title: "添加版本文件",
    ...{ style: {} },
}));
const __VLS_115 = __VLS_114({
    ...{ 'onUpdate:show': {} },
    show: (!!__VLS_ctx.assetFor),
    preset: "card",
    title: "添加版本文件",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_114));
let __VLS_118;
const __VLS_119 = {
    /** @type {typeof __VLS_118.'update:show'} */
    'onUpdate:show': (v => { if (!v)
        __VLS_ctx.assetFor = null; }),
};
const { default: __VLS_120 } = __VLS_116.slots;
let __VLS_121;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_122 = __VLS_asFunctionalComponent1(__VLS_121, new __VLS_121({
    labelPlacement: "top",
}));
const __VLS_123 = __VLS_122({
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_122));
const { default: __VLS_126 } = __VLS_124.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-grid" },
});
/** @type {__VLS_StyleScopedClasses['form-grid']} */ ;
let __VLS_127;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_128 = __VLS_asFunctionalComponent1(__VLS_127, new __VLS_127({
    label: "显示名称",
}));
const __VLS_129 = __VLS_128({
    label: "显示名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_128));
const { default: __VLS_132 } = __VLS_130.slots;
let __VLS_133;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_134 = __VLS_asFunctionalComponent1(__VLS_133, new __VLS_133({
    value: (__VLS_ctx.asset.name),
}));
const __VLS_135 = __VLS_134({
    value: (__VLS_ctx.asset.name),
}, ...__VLS_functionalComponentArgsRest(__VLS_134));
// @ts-ignore
[assetFor, assetFor, asset,];
var __VLS_130;
let __VLS_138;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_139 = __VLS_asFunctionalComponent1(__VLS_138, new __VLS_138({
    label: "文件名",
}));
const __VLS_140 = __VLS_139({
    label: "文件名",
}, ...__VLS_functionalComponentArgsRest(__VLS_139));
const { default: __VLS_143 } = __VLS_141.slots;
let __VLS_144;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_145 = __VLS_asFunctionalComponent1(__VLS_144, new __VLS_144({
    value: (__VLS_ctx.asset.file_name),
}));
const __VLS_146 = __VLS_145({
    value: (__VLS_ctx.asset.file_name),
}, ...__VLS_functionalComponentArgsRest(__VLS_145));
// @ts-ignore
[asset,];
var __VLS_141;
let __VLS_149;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_150 = __VLS_asFunctionalComponent1(__VLS_149, new __VLS_149({
    label: "系统",
}));
const __VLS_151 = __VLS_150({
    label: "系统",
}, ...__VLS_functionalComponentArgsRest(__VLS_150));
const { default: __VLS_154 } = __VLS_152.slots;
let __VLS_155;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_156 = __VLS_asFunctionalComponent1(__VLS_155, new __VLS_155({
    value: (__VLS_ctx.asset.os),
    options: (['windows', 'macos', 'linux', 'android', 'ios', 'any'].map(x => ({ label: x, value: x }))),
}));
const __VLS_157 = __VLS_156({
    value: (__VLS_ctx.asset.os),
    options: (['windows', 'macos', 'linux', 'android', 'ios', 'any'].map(x => ({ label: x, value: x }))),
}, ...__VLS_functionalComponentArgsRest(__VLS_156));
// @ts-ignore
[asset,];
var __VLS_152;
let __VLS_160;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_161 = __VLS_asFunctionalComponent1(__VLS_160, new __VLS_160({
    label: "架构",
}));
const __VLS_162 = __VLS_161({
    label: "架构",
}, ...__VLS_functionalComponentArgsRest(__VLS_161));
const { default: __VLS_165 } = __VLS_163.slots;
let __VLS_166;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_167 = __VLS_asFunctionalComponent1(__VLS_166, new __VLS_166({
    value: (__VLS_ctx.asset.arch),
    options: (['amd64', 'arm64', 'universal', 'any'].map(x => ({ label: x, value: x }))),
}));
const __VLS_168 = __VLS_167({
    value: (__VLS_ctx.asset.arch),
    options: (['amd64', 'arm64', 'universal', 'any'].map(x => ({ label: x, value: x }))),
}, ...__VLS_functionalComponentArgsRest(__VLS_167));
// @ts-ignore
[asset,];
var __VLS_163;
let __VLS_171;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_172 = __VLS_asFunctionalComponent1(__VLS_171, new __VLS_171({
    label: "大小（字节）",
}));
const __VLS_173 = __VLS_172({
    label: "大小（字节）",
}, ...__VLS_functionalComponentArgsRest(__VLS_172));
const { default: __VLS_176 } = __VLS_174.slots;
let __VLS_177;
/** @ts-ignore @type { | typeof __VLS_components.nInputNumber | typeof __VLS_components.NInputNumber | typeof __VLS_components['n-input-number']} */
nInputNumber;
// @ts-ignore
const __VLS_178 = __VLS_asFunctionalComponent1(__VLS_177, new __VLS_177({
    value: (__VLS_ctx.asset.size),
    ...{ style: {} },
}));
const __VLS_179 = __VLS_178({
    value: (__VLS_ctx.asset.size),
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_178));
// @ts-ignore
[asset,];
var __VLS_174;
let __VLS_182;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_183 = __VLS_asFunctionalComponent1(__VLS_182, new __VLS_182({
    label: "SHA256",
}));
const __VLS_184 = __VLS_183({
    label: "SHA256",
}, ...__VLS_functionalComponentArgsRest(__VLS_183));
const { default: __VLS_187 } = __VLS_185.slots;
let __VLS_188;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_189 = __VLS_asFunctionalComponent1(__VLS_188, new __VLS_188({
    value: (__VLS_ctx.asset.sha256),
}));
const __VLS_190 = __VLS_189({
    value: (__VLS_ctx.asset.sha256),
}, ...__VLS_functionalComponentArgsRest(__VLS_189));
// @ts-ignore
[asset,];
var __VLS_185;
let __VLS_193;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_194 = __VLS_asFunctionalComponent1(__VLS_193, new __VLS_193({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_195 = __VLS_194({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_194));
let __VLS_198;
const __VLS_199 = {
    /** @type {typeof __VLS_198.click} */
    onClick: (__VLS_ctx.addAsset),
};
const { default: __VLS_200 } = __VLS_196.slots;
// @ts-ignore
[addAsset,];
var __VLS_196;
var __VLS_197;
// @ts-ignore
[];
var __VLS_124;
// @ts-ignore
[];
var __VLS_116;
var __VLS_117;
let __VLS_201;
/** @ts-ignore @type { | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal'] | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal']} */
nModal;
// @ts-ignore
const __VLS_202 = __VLS_asFunctionalComponent1(__VLS_201, new __VLS_201({
    ...{ 'onUpdate:show': {} },
    show: (!!__VLS_ctx.sourceFor),
    preset: "card",
    title: "添加下载来源",
    ...{ style: {} },
}));
const __VLS_203 = __VLS_202({
    ...{ 'onUpdate:show': {} },
    show: (!!__VLS_ctx.sourceFor),
    preset: "card",
    title: "添加下载来源",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_202));
let __VLS_206;
const __VLS_207 = {
    /** @type {typeof __VLS_206.'update:show'} */
    'onUpdate:show': (v => { if (!v)
        __VLS_ctx.sourceFor = null; }),
};
const { default: __VLS_208 } = __VLS_204.slots;
let __VLS_209;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_210 = __VLS_asFunctionalComponent1(__VLS_209, new __VLS_209({
    labelPlacement: "top",
}));
const __VLS_211 = __VLS_210({
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_210));
const { default: __VLS_214 } = __VLS_212.slots;
let __VLS_215;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_216 = __VLS_asFunctionalComponent1(__VLS_215, new __VLS_215({
    label: "来源类型",
}));
const __VLS_217 = __VLS_216({
    label: "来源类型",
}, ...__VLS_functionalComponentArgsRest(__VLS_216));
const { default: __VLS_220 } = __VLS_218.slots;
let __VLS_221;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_222 = __VLS_asFunctionalComponent1(__VLS_221, new __VLS_221({
    value: (__VLS_ctx.source.source_type),
    options: ([{ label: '外链 / 网盘', value: 2 }, { label: '托管对象', value: 1 }]),
}));
const __VLS_223 = __VLS_222({
    value: (__VLS_ctx.source.source_type),
    options: ([{ label: '外链 / 网盘', value: 2 }, { label: '托管对象', value: 1 }]),
}, ...__VLS_functionalComponentArgsRest(__VLS_222));
// @ts-ignore
[sourceFor, sourceFor, source,];
var __VLS_218;
let __VLS_226;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_227 = __VLS_asFunctionalComponent1(__VLS_226, new __VLS_226({
    label: "显示名称",
}));
const __VLS_228 = __VLS_227({
    label: "显示名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_227));
const { default: __VLS_231 } = __VLS_229.slots;
let __VLS_232;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_233 = __VLS_asFunctionalComponent1(__VLS_232, new __VLS_232({
    value: (__VLS_ctx.source.name),
}));
const __VLS_234 = __VLS_233({
    value: (__VLS_ctx.source.name),
}, ...__VLS_functionalComponentArgsRest(__VLS_233));
// @ts-ignore
[source,];
var __VLS_229;
if (__VLS_ctx.source.source_type === 2) {
    let __VLS_237;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_238 = __VLS_asFunctionalComponent1(__VLS_237, new __VLS_237({
        label: "外链 URL",
    }));
    const __VLS_239 = __VLS_238({
        label: "外链 URL",
    }, ...__VLS_functionalComponentArgsRest(__VLS_238));
    const { default: __VLS_242 } = __VLS_240.slots;
    let __VLS_243;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_244 = __VLS_asFunctionalComponent1(__VLS_243, new __VLS_243({
        value: (__VLS_ctx.source.external_url),
    }));
    const __VLS_245 = __VLS_244({
        value: (__VLS_ctx.source.external_url),
    }, ...__VLS_functionalComponentArgsRest(__VLS_244));
    // @ts-ignore
    [source, source,];
    var __VLS_240;
    let __VLS_248;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_249 = __VLS_asFunctionalComponent1(__VLS_248, new __VLS_248({
        label: "提取码",
    }));
    const __VLS_250 = __VLS_249({
        label: "提取码",
    }, ...__VLS_functionalComponentArgsRest(__VLS_249));
    const { default: __VLS_253 } = __VLS_251.slots;
    let __VLS_254;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_255 = __VLS_asFunctionalComponent1(__VLS_254, new __VLS_254({
        value: (__VLS_ctx.source.extract_code),
    }));
    const __VLS_256 = __VLS_255({
        value: (__VLS_ctx.source.extract_code),
    }, ...__VLS_functionalComponentArgsRest(__VLS_255));
    // @ts-ignore
    [source,];
    var __VLS_251;
}
else {
    let __VLS_259;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_260 = __VLS_asFunctionalComponent1(__VLS_259, new __VLS_259({
        label: "存储 ID",
    }));
    const __VLS_261 = __VLS_260({
        label: "存储 ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_260));
    const { default: __VLS_264 } = __VLS_262.slots;
    let __VLS_265;
    /** @ts-ignore @type { | typeof __VLS_components.nInputNumber | typeof __VLS_components.NInputNumber | typeof __VLS_components['n-input-number']} */
    nInputNumber;
    // @ts-ignore
    const __VLS_266 = __VLS_asFunctionalComponent1(__VLS_265, new __VLS_265({
        value: (__VLS_ctx.source.storage_id),
    }));
    const __VLS_267 = __VLS_266({
        value: (__VLS_ctx.source.storage_id),
    }, ...__VLS_functionalComponentArgsRest(__VLS_266));
    // @ts-ignore
    [source,];
    var __VLS_262;
    let __VLS_270;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_271 = __VLS_asFunctionalComponent1(__VLS_270, new __VLS_270({
        label: "对象键",
    }));
    const __VLS_272 = __VLS_271({
        label: "对象键",
    }, ...__VLS_functionalComponentArgsRest(__VLS_271));
    const { default: __VLS_275 } = __VLS_273.slots;
    let __VLS_276;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_277 = __VLS_asFunctionalComponent1(__VLS_276, new __VLS_276({
        value: (__VLS_ctx.source.object_key),
        placeholder: "app/version/file.exe",
    }));
    const __VLS_278 = __VLS_277({
        value: (__VLS_ctx.source.object_key),
        placeholder: "app/version/file.exe",
    }, ...__VLS_functionalComponentArgsRest(__VLS_277));
    // @ts-ignore
    [source,];
    var __VLS_273;
}
let __VLS_281;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_282 = __VLS_asFunctionalComponent1(__VLS_281, new __VLS_281({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_283 = __VLS_282({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_282));
let __VLS_286;
const __VLS_287 = {
    /** @type {typeof __VLS_286.click} */
    onClick: (__VLS_ctx.addSource),
};
const { default: __VLS_288 } = __VLS_284.slots;
// @ts-ignore
[addSource,];
var __VLS_284;
var __VLS_285;
// @ts-ignore
[];
var __VLS_212;
// @ts-ignore
[];
var __VLS_204;
var __VLS_205;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
