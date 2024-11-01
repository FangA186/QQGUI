// store/messages.js
import { defineStore } from 'pinia';

export const useMessageStore = defineStore('message', {
    state: () => ({
        messages: [],
    }),
    actions: {
        addMessage(message) {
            this.messages.push(message);
        },
    },
});
