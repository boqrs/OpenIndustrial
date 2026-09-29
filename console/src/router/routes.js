export const routes = [
    {
        path: "/",
        redirect: "/login",
    },
    {
        path: "/login",
        name: "Login",
        component: () => import("../views/login/Login.vue"),
    },
    {
        path: "/dashboard",
        name: "Dashboard",
        component: () => import("../views/dashboard/Dashboard.vue"),
    },
];
