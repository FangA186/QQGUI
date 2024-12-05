import {createApp, watch} from 'vue';
import App from './App.vue';
import './style.css';
import { EventsOn } from '../wailsjs/runtime/runtime.js';
import router from "./router/router.js";
import ElementPlus, {ElNotification} from 'element-plus';
import {createPinia} from 'pinia';
import {useUserStore} from './store/userStore.js'; // 确保你有 userStore
// import { useMessageStore } from './store/messages.js';
import 'element-plus/dist/index.css';
const app = createApp(App);

// 初始化 Pinia 并注册
const pinia = createPinia();
app.use(pinia).use(router).use(ElementPlus);


// 必须在应用挂载之后调用 useUserStore
const userStore = useUserStore();
// const messageStore = useMessageStore();
// 监听用户 ID 变化，建立 SSE 连接
watch(
    () => userStore.userId, // 监听用户ID变化
    (userId) => {
        if (userId) {
            // 监听后端事件
            EventsOn('event_name', (data) => {
                console.log('Received event data:', data);
                // 处理接收到的数据
            });
            // 用户ID存在时，创建 SSE 连接
            const eventSource = new EventSource(`http://wails.localhost:8080/sse?user_id=${userId}`);
            // 连接成功时的提示
            console.log(eventSource)
            eventSource.onopen = function () {
                console.log('SSE 连接成功');
                console.log(eventSource)
            };
            eventSource.onmessage = function (event) {
                // messageStore.addMessage(JSON.parse(event.data)); // 将消息存储到 store 中
                console.log(JSON.parse(event.data))
                const info = JSON.parse(event.data)
                if (info.i===0){
                    open1(JSON.parse(event.data))
                }else if (info.i===1){
                    open2(JSON.parse(event.data))
                }else {
                    open3(JSON.parse(event.data))
                }
            };

            eventSource.onerror = function (event) {
                console.error('SSE 连接发生错误:', event);
                eventSource.close(); // 关闭连接
            };

            // 当用户注销时关闭 SSE
            watch(
                () => userStore.userId,
                (newUserId) => {
                    if (!newUserId) {
                        eventSource.close(); // 关闭 SSE 连接
                    }
                }
            );
        }
    },
    {immediate: true} // 立即执行，确保应用启动时也会检测
);
const open1 = (obj) => {
    const svgIcon = `
        <svg class="icon" aria-hidden="true" style="width: 20px; height: 20px;">
            <use xlink:href="#${obj.icon}"></use>
        </svg>
    `;
    ElNotification({
        title: '申请通知',
        dangerouslyUseHTMLString: true,
        message: ` <div style="display: flex;justify-content: center;align-items: center;">
 <img alt="" src="${obj.friend.Avatar}" style="border-radius: 0.2vw;width: 2.5vw;height: 5vh;margin-right:0.5vh">
 <span>${obj.friend.Username}${obj.message}</span>
 </div>
`,
        duration: 0,
    })
}
const open2 = (obj) => {
    const svgIcon = `
        <svg class="icon" aria-hidden="true" style="width: 20px; height: 20px;">
            <use xlink:href="#${obj.icon}"></use>
        </svg>
    `;
    ElNotification({
        title: '消息通知',
        dangerouslyUseHTMLString: true,
        message: ` <div style="display: flex;justify-content: center;align-items: center;">
 <img alt="" src="${obj.avatar}" style="border-radius: 0.2vw;width: 2.5vw;height: 5vh;margin-right:0.5vh">
 <span>${obj.username}向您发来一条信息</span>
 </div>
`,
        duration: 0,
    })
}
const open3 = (obj) => {
    const svgIcon = `
        <svg class="icon" aria-hidden="true" style="width: 20px; height: 20px;">
            <use xlink:href="#${obj.icon}"></use>
        </svg>
    `;
    ElNotification({
        title: '群聊邀请通知',
        dangerouslyUseHTMLString: true,
        message: ` <div style="display: flex;justify-content: center;align-items: center;">
 <span>${obj.createName}已邀请你进入${obj.RoomName}群聊</span>
 </div>
`,
        duration: 0,
    })
}
// 应用挂载之后再执行逻辑
app.mount('#app');