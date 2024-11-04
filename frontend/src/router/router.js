import {createRouter, createWebHashHistory} from "vue-router";
import Form3 from "../views/From3.vue";
import Friends from "../views/Friends.vue";
import Home from "../views/home.vue";
const routes = [
    {path:'/',name:'Form3',component:Form3},
    {path:'/Index',name:'Index',component:Home},
    {path:'/Friends',name:'Friend',component:Friends},
]
const router = createRouter({
    history:createWebHashHistory(),
    routes,
    linkActiveClass:'vue-school-active-link' // 这个是用于路由点击时出现的活动状态的类
})

router.beforeEach(async (to, from) => {
})
export default router