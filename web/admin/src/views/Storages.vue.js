import { onMounted, reactive, ref } from 'vue';
import { useMessage } from 'naive-ui';
import { api } from '../api';
const rows = ref([]), show = ref(false), message = useMessage(), driver = ref('local'), base = reactive({ name: '本地存储', is_default: false, is_enabled: true, remark: '' }), local = reactive({ root: 'files' }), s3 = reactive({ endpoint: '', region: 'auto', bucket: '', access_key_id: '', secret_access_key: '', base_path: 'files/', force_path_style: false, public_base_url: '', presign_expire_minutes: 30 });
async function load() { rows.value = await api.get('/storages'); }
function open() { driver.value = 'local'; Object.assign(base, { name: '本地存储', is_default: false, is_enabled: true, remark: '' }); show.value = true; }
async function save() { try {
    await api.post('/storages', { ...base, driver: driver.value, config: JSON.stringify(driver.value === 'local' ? local : s3) });
    show.value = false;
    message.success('存储已保存');
    load();
}
catch (e) {
    message.error(e.message);
} }
async function test(id) { try {
    const x = await api.post(`/storages/${id}/test`);
    message.success(`连接正常，${x.latency_ms} ms`);
}
catch (e) {
    message.error(e.message);
} }
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
    onClick: (__VLS_ctx.open),
};
const { default: __VLS_7 } = __VLS_3.slots;
// @ts-ignore
[open,];
var __VLS_3;
var __VLS_4;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "metric-grid" },
});
/** @type {__VLS_StyleScopedClasses['metric-grid']} */ ;
for (const [s] of __VLS_vFor((__VLS_ctx.rows))) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        key: (s.id),
        ...{ class: "metric" },
    });
    /** @type {__VLS_StyleScopedClasses['metric']} */ ;
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({});
    let __VLS_8;
    /** @ts-ignore @type { | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag'] | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag']} */
    nTag;
    // @ts-ignore
    const __VLS_9 = __VLS_asFunctionalComponent1(__VLS_8, new __VLS_8({}));
    const __VLS_10 = __VLS_9({}, ...__VLS_functionalComponentArgsRest(__VLS_9));
    const { default: __VLS_13 } = __VLS_11.slots;
    (s.driver);
    // @ts-ignore
    [rows,];
    var __VLS_11;
    if (s.is_default) {
        let __VLS_14;
        /** @ts-ignore @type { | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag'] | typeof __VLS_components.nTag | typeof __VLS_components.NTag | typeof __VLS_components['n-tag']} */
        nTag;
        // @ts-ignore
        const __VLS_15 = __VLS_asFunctionalComponent1(__VLS_14, new __VLS_14({
            type: "success",
        }));
        const __VLS_16 = __VLS_15({
            type: "success",
        }, ...__VLS_functionalComponentArgsRest(__VLS_15));
        const { default: __VLS_19 } = __VLS_17.slots;
        // @ts-ignore
        [];
        var __VLS_17;
    }
    __VLS_asFunctionalElement1(__VLS_intrinsics.strong, __VLS_intrinsics.strong)({
        ...{ style: {} },
    });
    (s.name);
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
    (s.remark || s.config);
    __VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
        ...{ style: {} },
    });
    let __VLS_20;
    /** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
    nButton;
    // @ts-ignore
    const __VLS_21 = __VLS_asFunctionalComponent1(__VLS_20, new __VLS_20({
        ...{ 'onClick': {} },
        size: "small",
    }));
    const __VLS_22 = __VLS_21({
        ...{ 'onClick': {} },
        size: "small",
    }, ...__VLS_functionalComponentArgsRest(__VLS_21));
    let __VLS_25;
    const __VLS_26 = {
        /** @type {typeof __VLS_25.click} */
        onClick: (...[$event]) => {
            return (__VLS_ctx.test(s.id));
            // @ts-ignore
            [test,];
        },
    };
    const { default: __VLS_27 } = __VLS_23.slots;
    // @ts-ignore
    [];
    var __VLS_23;
    var __VLS_24;
    // @ts-ignore
    [];
}
let __VLS_28;
/** @ts-ignore @type { | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal'] | typeof __VLS_components.nModal | typeof __VLS_components.NModal | typeof __VLS_components['n-modal']} */
nModal;
// @ts-ignore
const __VLS_29 = __VLS_asFunctionalComponent1(__VLS_28, new __VLS_28({
    show: (__VLS_ctx.show),
    preset: "card",
    title: "添加存储",
    ...{ style: {} },
}));
const __VLS_30 = __VLS_29({
    show: (__VLS_ctx.show),
    preset: "card",
    title: "添加存储",
    ...{ style: {} },
}, ...__VLS_functionalComponentArgsRest(__VLS_29));
const { default: __VLS_33 } = __VLS_31.slots;
let __VLS_34;
/** @ts-ignore @type { | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form'] | typeof __VLS_components.nForm | typeof __VLS_components.NForm | typeof __VLS_components['n-form']} */
nForm;
// @ts-ignore
const __VLS_35 = __VLS_asFunctionalComponent1(__VLS_34, new __VLS_34({
    labelPlacement: "top",
}));
const __VLS_36 = __VLS_35({
    labelPlacement: "top",
}, ...__VLS_functionalComponentArgsRest(__VLS_35));
const { default: __VLS_39 } = __VLS_37.slots;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-grid" },
});
/** @type {__VLS_StyleScopedClasses['form-grid']} */ ;
let __VLS_40;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_41 = __VLS_asFunctionalComponent1(__VLS_40, new __VLS_40({
    label: "驱动",
}));
const __VLS_42 = __VLS_41({
    label: "驱动",
}, ...__VLS_functionalComponentArgsRest(__VLS_41));
const { default: __VLS_45 } = __VLS_43.slots;
let __VLS_46;
/** @ts-ignore @type { | typeof __VLS_components.nSelect | typeof __VLS_components.NSelect | typeof __VLS_components['n-select']} */
nSelect;
// @ts-ignore
const __VLS_47 = __VLS_asFunctionalComponent1(__VLS_46, new __VLS_46({
    value: (__VLS_ctx.driver),
    options: ([{ label: '本地磁盘', value: 'local' }, { label: 'S3 兼容对象存储', value: 's3' }]),
}));
const __VLS_48 = __VLS_47({
    value: (__VLS_ctx.driver),
    options: ([{ label: '本地磁盘', value: 'local' }, { label: 'S3 兼容对象存储', value: 's3' }]),
}, ...__VLS_functionalComponentArgsRest(__VLS_47));
// @ts-ignore
[show, driver,];
var __VLS_43;
let __VLS_51;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_52 = __VLS_asFunctionalComponent1(__VLS_51, new __VLS_51({
    label: "名称",
}));
const __VLS_53 = __VLS_52({
    label: "名称",
}, ...__VLS_functionalComponentArgsRest(__VLS_52));
const { default: __VLS_56 } = __VLS_54.slots;
let __VLS_57;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_58 = __VLS_asFunctionalComponent1(__VLS_57, new __VLS_57({
    value: (__VLS_ctx.base.name),
}));
const __VLS_59 = __VLS_58({
    value: (__VLS_ctx.base.name),
}, ...__VLS_functionalComponentArgsRest(__VLS_58));
// @ts-ignore
[base,];
var __VLS_54;
if (__VLS_ctx.driver === 'local') {
    let __VLS_62;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_63 = __VLS_asFunctionalComponent1(__VLS_62, new __VLS_62({
        ...{ class: "wide" },
        label: "数据目录下的相对路径或绝对路径",
    }));
    const __VLS_64 = __VLS_63({
        ...{ class: "wide" },
        label: "数据目录下的相对路径或绝对路径",
    }, ...__VLS_functionalComponentArgsRest(__VLS_63));
    /** @type {__VLS_StyleScopedClasses['wide']} */ ;
    const { default: __VLS_67 } = __VLS_65.slots;
    let __VLS_68;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_69 = __VLS_asFunctionalComponent1(__VLS_68, new __VLS_68({
        value: (__VLS_ctx.local.root),
    }));
    const __VLS_70 = __VLS_69({
        value: (__VLS_ctx.local.root),
    }, ...__VLS_functionalComponentArgsRest(__VLS_69));
    // @ts-ignore
    [driver, local,];
    var __VLS_65;
}
else {
    let __VLS_73;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_74 = __VLS_asFunctionalComponent1(__VLS_73, new __VLS_73({
        ...{ class: "wide" },
        label: "Endpoint",
    }));
    const __VLS_75 = __VLS_74({
        ...{ class: "wide" },
        label: "Endpoint",
    }, ...__VLS_functionalComponentArgsRest(__VLS_74));
    /** @type {__VLS_StyleScopedClasses['wide']} */ ;
    const { default: __VLS_78 } = __VLS_76.slots;
    let __VLS_79;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_80 = __VLS_asFunctionalComponent1(__VLS_79, new __VLS_79({
        value: (__VLS_ctx.s3.endpoint),
        placeholder: "https://<account>.r2.cloudflarestorage.com",
    }));
    const __VLS_81 = __VLS_80({
        value: (__VLS_ctx.s3.endpoint),
        placeholder: "https://<account>.r2.cloudflarestorage.com",
    }, ...__VLS_functionalComponentArgsRest(__VLS_80));
    // @ts-ignore
    [s3,];
    var __VLS_76;
    let __VLS_84;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_85 = __VLS_asFunctionalComponent1(__VLS_84, new __VLS_84({
        label: "Region",
    }));
    const __VLS_86 = __VLS_85({
        label: "Region",
    }, ...__VLS_functionalComponentArgsRest(__VLS_85));
    const { default: __VLS_89 } = __VLS_87.slots;
    let __VLS_90;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_91 = __VLS_asFunctionalComponent1(__VLS_90, new __VLS_90({
        value: (__VLS_ctx.s3.region),
    }));
    const __VLS_92 = __VLS_91({
        value: (__VLS_ctx.s3.region),
    }, ...__VLS_functionalComponentArgsRest(__VLS_91));
    // @ts-ignore
    [s3,];
    var __VLS_87;
    let __VLS_95;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_96 = __VLS_asFunctionalComponent1(__VLS_95, new __VLS_95({
        label: "Bucket",
    }));
    const __VLS_97 = __VLS_96({
        label: "Bucket",
    }, ...__VLS_functionalComponentArgsRest(__VLS_96));
    const { default: __VLS_100 } = __VLS_98.slots;
    let __VLS_101;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_102 = __VLS_asFunctionalComponent1(__VLS_101, new __VLS_101({
        value: (__VLS_ctx.s3.bucket),
    }));
    const __VLS_103 = __VLS_102({
        value: (__VLS_ctx.s3.bucket),
    }, ...__VLS_functionalComponentArgsRest(__VLS_102));
    // @ts-ignore
    [s3,];
    var __VLS_98;
    let __VLS_106;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_107 = __VLS_asFunctionalComponent1(__VLS_106, new __VLS_106({
        label: "Access Key ID",
    }));
    const __VLS_108 = __VLS_107({
        label: "Access Key ID",
    }, ...__VLS_functionalComponentArgsRest(__VLS_107));
    const { default: __VLS_111 } = __VLS_109.slots;
    let __VLS_112;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_113 = __VLS_asFunctionalComponent1(__VLS_112, new __VLS_112({
        value: (__VLS_ctx.s3.access_key_id),
    }));
    const __VLS_114 = __VLS_113({
        value: (__VLS_ctx.s3.access_key_id),
    }, ...__VLS_functionalComponentArgsRest(__VLS_113));
    // @ts-ignore
    [s3,];
    var __VLS_109;
    let __VLS_117;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_118 = __VLS_asFunctionalComponent1(__VLS_117, new __VLS_117({
        label: "Secret Access Key",
    }));
    const __VLS_119 = __VLS_118({
        label: "Secret Access Key",
    }, ...__VLS_functionalComponentArgsRest(__VLS_118));
    const { default: __VLS_122 } = __VLS_120.slots;
    let __VLS_123;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_124 = __VLS_asFunctionalComponent1(__VLS_123, new __VLS_123({
        value: (__VLS_ctx.s3.secret_access_key),
        type: "password",
    }));
    const __VLS_125 = __VLS_124({
        value: (__VLS_ctx.s3.secret_access_key),
        type: "password",
    }, ...__VLS_functionalComponentArgsRest(__VLS_124));
    // @ts-ignore
    [s3,];
    var __VLS_120;
    let __VLS_128;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_129 = __VLS_asFunctionalComponent1(__VLS_128, new __VLS_128({
        label: "桶内前缀",
    }));
    const __VLS_130 = __VLS_129({
        label: "桶内前缀",
    }, ...__VLS_functionalComponentArgsRest(__VLS_129));
    const { default: __VLS_133 } = __VLS_131.slots;
    let __VLS_134;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_135 = __VLS_asFunctionalComponent1(__VLS_134, new __VLS_134({
        value: (__VLS_ctx.s3.base_path),
    }));
    const __VLS_136 = __VLS_135({
        value: (__VLS_ctx.s3.base_path),
    }, ...__VLS_functionalComponentArgsRest(__VLS_135));
    // @ts-ignore
    [s3,];
    var __VLS_131;
    let __VLS_139;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_140 = __VLS_asFunctionalComponent1(__VLS_139, new __VLS_139({
        label: "公共下载域名",
    }));
    const __VLS_141 = __VLS_140({
        label: "公共下载域名",
    }, ...__VLS_functionalComponentArgsRest(__VLS_140));
    const { default: __VLS_144 } = __VLS_142.slots;
    let __VLS_145;
    /** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
    nInput;
    // @ts-ignore
    const __VLS_146 = __VLS_asFunctionalComponent1(__VLS_145, new __VLS_145({
        value: (__VLS_ctx.s3.public_base_url),
        placeholder: "https://dl.example.com",
    }));
    const __VLS_147 = __VLS_146({
        value: (__VLS_ctx.s3.public_base_url),
        placeholder: "https://dl.example.com",
    }, ...__VLS_functionalComponentArgsRest(__VLS_146));
    // @ts-ignore
    [s3,];
    var __VLS_142;
    let __VLS_150;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_151 = __VLS_asFunctionalComponent1(__VLS_150, new __VLS_150({
        label: "预签名有效分钟",
    }));
    const __VLS_152 = __VLS_151({
        label: "预签名有效分钟",
    }, ...__VLS_functionalComponentArgsRest(__VLS_151));
    const { default: __VLS_155 } = __VLS_153.slots;
    let __VLS_156;
    /** @ts-ignore @type { | typeof __VLS_components.nInputNumber | typeof __VLS_components.NInputNumber | typeof __VLS_components['n-input-number']} */
    nInputNumber;
    // @ts-ignore
    const __VLS_157 = __VLS_asFunctionalComponent1(__VLS_156, new __VLS_156({
        value: (__VLS_ctx.s3.presign_expire_minutes),
    }));
    const __VLS_158 = __VLS_157({
        value: (__VLS_ctx.s3.presign_expire_minutes),
    }, ...__VLS_functionalComponentArgsRest(__VLS_157));
    // @ts-ignore
    [s3,];
    var __VLS_153;
    let __VLS_161;
    /** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
    nFormItem;
    // @ts-ignore
    const __VLS_162 = __VLS_asFunctionalComponent1(__VLS_161, new __VLS_161({
        label: "Path-style",
    }));
    const __VLS_163 = __VLS_162({
        label: "Path-style",
    }, ...__VLS_functionalComponentArgsRest(__VLS_162));
    const { default: __VLS_166 } = __VLS_164.slots;
    let __VLS_167;
    /** @ts-ignore @type { | typeof __VLS_components.nSwitch | typeof __VLS_components.NSwitch | typeof __VLS_components['n-switch']} */
    nSwitch;
    // @ts-ignore
    const __VLS_168 = __VLS_asFunctionalComponent1(__VLS_167, new __VLS_167({
        value: (__VLS_ctx.s3.force_path_style),
    }));
    const __VLS_169 = __VLS_168({
        value: (__VLS_ctx.s3.force_path_style),
    }, ...__VLS_functionalComponentArgsRest(__VLS_168));
    // @ts-ignore
    [s3,];
    var __VLS_164;
}
let __VLS_172;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_173 = __VLS_asFunctionalComponent1(__VLS_172, new __VLS_172({
    label: "备注",
}));
const __VLS_174 = __VLS_173({
    label: "备注",
}, ...__VLS_functionalComponentArgsRest(__VLS_173));
const { default: __VLS_177 } = __VLS_175.slots;
let __VLS_178;
/** @ts-ignore @type { | typeof __VLS_components.nInput | typeof __VLS_components.NInput | typeof __VLS_components['n-input']} */
nInput;
// @ts-ignore
const __VLS_179 = __VLS_asFunctionalComponent1(__VLS_178, new __VLS_178({
    value: (__VLS_ctx.base.remark),
}));
const __VLS_180 = __VLS_179({
    value: (__VLS_ctx.base.remark),
}, ...__VLS_functionalComponentArgsRest(__VLS_179));
// @ts-ignore
[base,];
var __VLS_175;
let __VLS_183;
/** @ts-ignore @type { | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item'] | typeof __VLS_components.nFormItem | typeof __VLS_components.NFormItem | typeof __VLS_components['n-form-item']} */
nFormItem;
// @ts-ignore
const __VLS_184 = __VLS_asFunctionalComponent1(__VLS_183, new __VLS_183({
    label: "设为默认",
}));
const __VLS_185 = __VLS_184({
    label: "设为默认",
}, ...__VLS_functionalComponentArgsRest(__VLS_184));
const { default: __VLS_188 } = __VLS_186.slots;
let __VLS_189;
/** @ts-ignore @type { | typeof __VLS_components.nSwitch | typeof __VLS_components.NSwitch | typeof __VLS_components['n-switch']} */
nSwitch;
// @ts-ignore
const __VLS_190 = __VLS_asFunctionalComponent1(__VLS_189, new __VLS_189({
    value: (__VLS_ctx.base.is_default),
}));
const __VLS_191 = __VLS_190({
    value: (__VLS_ctx.base.is_default),
}, ...__VLS_functionalComponentArgsRest(__VLS_190));
// @ts-ignore
[base,];
var __VLS_186;
let __VLS_194;
/** @ts-ignore @type { | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button'] | typeof __VLS_components.nButton | typeof __VLS_components.NButton | typeof __VLS_components['n-button']} */
nButton;
// @ts-ignore
const __VLS_195 = __VLS_asFunctionalComponent1(__VLS_194, new __VLS_194({
    ...{ 'onClick': {} },
    type: "primary",
}));
const __VLS_196 = __VLS_195({
    ...{ 'onClick': {} },
    type: "primary",
}, ...__VLS_functionalComponentArgsRest(__VLS_195));
let __VLS_199;
const __VLS_200 = {
    /** @type {typeof __VLS_199.click} */
    onClick: (__VLS_ctx.save),
};
const { default: __VLS_201 } = __VLS_197.slots;
// @ts-ignore
[save,];
var __VLS_197;
var __VLS_198;
// @ts-ignore
[];
var __VLS_37;
// @ts-ignore
[];
var __VLS_31;
// @ts-ignore
[];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
