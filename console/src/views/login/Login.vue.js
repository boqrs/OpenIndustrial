import { ref } from "vue";
const account = ref("");
const password = ref("");
const remember = ref(false);
const loading = ref(false);
const errorMessage = ref("");
const handleLogin = async () => {
    errorMessage.value = "";
    if (!account.value.trim()) {
        errorMessage.value = "请输入账户名";
        return;
    }
    if (!password.value) {
        errorMessage.value = "请输入登录密码";
        return;
    }
    loading.value = true;
    try {
        /*
         * TODO:
         * 后续在这里接入真实登录 API。
         *
         * 例如：
         *
         * const response = await login({
         *   email: account.value,
         *   password: password.value,
         * })
         *
         * 登录成功后：
         *
         * router.replace('/dashboard')
         */
        // 当前阶段仅用于验证页面交互。
        await new Promise((resolve) => {
            window.setTimeout(resolve, 500);
        });
    }
    finally {
        loading.value = false;
    }
};
const __VLS_ctx = {
    ...{},
    ...{},
};
let __VLS_components;
let __VLS_intrinsics;
let __VLS_directives;
/** @type {__VLS_StyleScopedClasses['brand-side']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-side']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-title']} */ ;
/** @type {__VLS_StyleScopedClasses['stage-dot']} */ ;
/** @type {__VLS_StyleScopedClasses['stage-label']} */ ;
/** @type {__VLS_StyleScopedClasses['field']} */ ;
/** @type {__VLS_StyleScopedClasses['field']} */ ;
/** @type {__VLS_StyleScopedClasses['field']} */ ;
/** @type {__VLS_StyleScopedClasses['field']} */ ;
/** @type {__VLS_StyleScopedClasses['field']} */ ;
/** @type {__VLS_StyleScopedClasses['checkbox-wrap']} */ ;
/** @type {__VLS_StyleScopedClasses['checkbox-wrap']} */ ;
/** @type {__VLS_StyleScopedClasses['link']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-login']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-login']} */ ;
/** @type {__VLS_StyleScopedClasses['btn-login']} */ ;
/** @type {__VLS_StyleScopedClasses['login-page']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-side']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-title']} */ ;
/** @type {__VLS_StyleScopedClasses['lifecycle']} */ ;
/** @type {__VLS_StyleScopedClasses['stage-line']} */ ;
/** @type {__VLS_StyleScopedClasses['form-side']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-side']} */ ;
/** @type {__VLS_StyleScopedClasses['logo-lockup']} */ ;
/** @type {__VLS_StyleScopedClasses['logo-mark']} */ ;
/** @type {__VLS_StyleScopedClasses['logo-text-main']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-title']} */ ;
/** @type {__VLS_StyleScopedClasses['brand-desc']} */ ;
/** @type {__VLS_StyleScopedClasses['lifecycle']} */ ;
/** @type {__VLS_StyleScopedClasses['stage-line']} */ ;
/** @type {__VLS_StyleScopedClasses['code-sample']} */ ;
/** @type {__VLS_StyleScopedClasses['code-sample-value']} */ ;
/** @type {__VLS_StyleScopedClasses['form-side']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "login-page" },
});
/** @type {__VLS_StyleScopedClasses['login-page']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.section, __VLS_intrinsics.section)({
    ...{ class: "brand-side" },
});
/** @type {__VLS_StyleScopedClasses['brand-side']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "brand-content" },
});
/** @type {__VLS_StyleScopedClasses['brand-content']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "logo-lockup" },
});
/** @type {__VLS_StyleScopedClasses['logo-lockup']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.svg, __VLS_intrinsics.svg)({
    ...{ class: "logo-mark" },
    viewBox: "0 0 60 60",
    fill: "none",
    xmlns: "http://www.w3.org/2000/svg",
    'aria-hidden': "true",
});
/** @type {__VLS_StyleScopedClasses['logo-mark']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.path)({
    d: "M30 4 A26 26 0 1 1 4 30",
    stroke: "#F0A030",
    'stroke-width': "2.4",
    fill: "none",
    'stroke-linecap': "round",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.path)({
    d: "M4 30 A26 26 0 0 1 12 12",
    stroke: "#F0A030",
    'stroke-width': "2.4",
    fill: "none",
    'stroke-linecap': "round",
    opacity: "0.35",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.rect)({
    x: "20",
    y: "20",
    width: "20",
    height: "16",
    rx: "3",
    stroke: "#E8EDF5",
    'stroke-width': "1.8",
    fill: "none",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.circle)({
    cx: "30",
    cy: "28",
    r: "2.4",
    fill: "#F0A030",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.line)({
    x1: "24",
    y1: "36",
    x2: "24",
    y2: "40",
    stroke: "#E8EDF5",
    'stroke-width': "1.6",
    'stroke-linecap': "round",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.line)({
    x1: "36",
    y1: "36",
    x2: "36",
    y2: "40",
    stroke: "#E8EDF5",
    'stroke-width': "1.6",
    'stroke-linecap': "round",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.line)({
    x1: "22",
    y1: "44",
    x2: "38",
    y2: "44",
    stroke: "#F0A030",
    'stroke-width': "1.6",
    'stroke-linecap': "round",
    opacity: "0.7",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.line)({
    x1: "26",
    y1: "47",
    x2: "34",
    y2: "47",
    stroke: "#F0A030",
    'stroke-width': "1.6",
    'stroke-linecap': "round",
    opacity: "0.4",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.circle)({
    cx: "30",
    cy: "4",
    r: "3",
    fill: "#F0A030",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.circle)({
    cx: "56",
    cy: "30",
    r: "3",
    fill: "#F0A030",
    opacity: "0.55",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.circle)({
    cx: "4",
    cy: "30",
    r: "3",
    fill: "#F0A030",
    opacity: "0.55",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "logo-text" },
});
/** @type {__VLS_StyleScopedClasses['logo-text']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "logo-text-main" },
});
/** @type {__VLS_StyleScopedClasses['logo-text-main']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "logo-text-sub" },
});
/** @type {__VLS_StyleScopedClasses['logo-text-sub']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h1, __VLS_intrinsics.h1)({
    ...{ class: "brand-title" },
});
/** @type {__VLS_StyleScopedClasses['brand-title']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.br)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.em, __VLS_intrinsics.em)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "brand-desc" },
});
/** @type {__VLS_StyleScopedClasses['brand-desc']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.br)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "lifecycle" },
});
/** @type {__VLS_StyleScopedClasses['lifecycle']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage" },
});
/** @type {__VLS_StyleScopedClasses['stage']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-dot" },
});
/** @type {__VLS_StyleScopedClasses['stage-dot']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-label active" },
});
/** @type {__VLS_StyleScopedClasses['stage-label']} */ ;
/** @type {__VLS_StyleScopedClasses['active']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-line" },
});
/** @type {__VLS_StyleScopedClasses['stage-line']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage" },
});
/** @type {__VLS_StyleScopedClasses['stage']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-dot dim" },
});
/** @type {__VLS_StyleScopedClasses['stage-dot']} */ ;
/** @type {__VLS_StyleScopedClasses['dim']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-label" },
});
/** @type {__VLS_StyleScopedClasses['stage-label']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-line" },
});
/** @type {__VLS_StyleScopedClasses['stage-line']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage" },
});
/** @type {__VLS_StyleScopedClasses['stage']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-dot dim" },
});
/** @type {__VLS_StyleScopedClasses['stage-dot']} */ ;
/** @type {__VLS_StyleScopedClasses['dim']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "stage-label" },
});
/** @type {__VLS_StyleScopedClasses['stage-label']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "code-sample" },
});
/** @type {__VLS_StyleScopedClasses['code-sample']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
    ...{ class: "code-sample-label" },
});
/** @type {__VLS_StyleScopedClasses['code-sample-label']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({
    ...{ class: "code-sample-value" },
});
/** @type {__VLS_StyleScopedClasses['code-sample-value']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.section, __VLS_intrinsics.section)({
    ...{ class: "form-side" },
});
/** @type {__VLS_StyleScopedClasses['form-side']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-box" },
});
/** @type {__VLS_StyleScopedClasses['form-box']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.h2, __VLS_intrinsics.h2)({
    ...{ class: "form-title" },
});
/** @type {__VLS_StyleScopedClasses['form-title']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "form-subtitle" },
});
/** @type {__VLS_StyleScopedClasses['form-subtitle']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.form, __VLS_intrinsics.form)({
    ...{ onSubmit: (__VLS_ctx.handleLogin) },
});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "field" },
});
/** @type {__VLS_StyleScopedClasses['field']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.label, __VLS_intrinsics.label)({
    for: "account",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.input)({
    id: "account",
    value: (__VLS_ctx.account),
    type: "text",
    autocomplete: "username",
    placeholder: "请输入企业账户或邮箱",
    disabled: (__VLS_ctx.loading),
});
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "field" },
});
/** @type {__VLS_StyleScopedClasses['field']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.label, __VLS_intrinsics.label)({
    for: "password",
});
__VLS_asFunctionalElement1(__VLS_intrinsics.input)({
    ...{ onKeyup: (__VLS_ctx.handleLogin) },
    id: "password",
    type: "password",
    autocomplete: "current-password",
    placeholder: "请输入登录密码",
    disabled: (__VLS_ctx.loading),
});
(__VLS_ctx.password);
__VLS_asFunctionalElement1(__VLS_intrinsics.div, __VLS_intrinsics.div)({
    ...{ class: "form-extra" },
});
/** @type {__VLS_StyleScopedClasses['form-extra']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.label, __VLS_intrinsics.label)({
    ...{ class: "checkbox-wrap" },
});
/** @type {__VLS_StyleScopedClasses['checkbox-wrap']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.input)({
    type: "checkbox",
    disabled: (__VLS_ctx.loading),
});
(__VLS_ctx.remember);
__VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ onClick: () => { } },
    ...{ class: "link" },
    href: "#",
});
/** @type {__VLS_StyleScopedClasses['link']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.button, __VLS_intrinsics.button)({
    ...{ class: "btn-login" },
    type: "submit",
    disabled: (__VLS_ctx.loading),
});
/** @type {__VLS_StyleScopedClasses['btn-login']} */ ;
if (__VLS_ctx.loading) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
}
else {
    __VLS_asFunctionalElement1(__VLS_intrinsics.span, __VLS_intrinsics.span)({});
}
if (__VLS_ctx.errorMessage) {
    __VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
        ...{ class: "error-message" },
    });
    /** @type {__VLS_StyleScopedClasses['error-message']} */ ;
    (__VLS_ctx.errorMessage);
}
__VLS_asFunctionalElement1(__VLS_intrinsics.p, __VLS_intrinsics.p)({
    ...{ class: "form-footer" },
});
/** @type {__VLS_StyleScopedClasses['form-footer']} */ ;
__VLS_asFunctionalElement1(__VLS_intrinsics.a, __VLS_intrinsics.a)({
    ...{ onClick: () => { } },
    ...{ class: "link" },
    href: "#",
});
/** @type {__VLS_StyleScopedClasses['link']} */ ;
// @ts-ignore
[handleLogin, handleLogin, account, loading, loading, loading, loading, loading, password, remember, errorMessage, errorMessage,];
const __VLS_export = (await import('vue')).defineComponent({});
export default {};
