// store/userStore.js
import { defineStore } from 'pinia';

export const useUserStore = defineStore('user', {
    state: () => ({
        userId: localStorage.getItem('userID') || null, // 从本地缓存获取用户ID
    }),
    actions: {
        setUserId(id) {
            this.userId = id;
            localStorage.setItem('userID', id); // 登录时存储到本地
        },
        logout() {
            this.userId = null;
            localStorage.removeItem('userID'); // 登出时清除缓存
        },
    },
});
