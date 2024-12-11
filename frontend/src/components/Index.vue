<script setup>
import LeftNav from "./LeftNav.vue";
import Search from "./Search.vue";
import {CrateRoomReceive, CreateRoom, FriendList, QueryRoom} from "../../wailsjs/go/api/MyApp.js"
import {onMounted, ref} from "vue";

const showmessages = ref(true)
const dialogVisible = ref(false)
const avatar = ref('')
const friends = ref([])
const choiceFriendList = ref([])
const showGroup = ref()
const createGroup = async () => {
  dialogVisible.value = true
  const res = await FriendList(localStorage.getItem("userID"))
  for (let i = 0; i < res.InfoList.length; i++) {
    for (let j = i + 1; j < res.InfoList.length; j++) {
      if (res.InfoList[i].id === res.InfoList[j].id) {         //第一个等同于第二个，splice方法删除第二个
        res.InfoList.splice(j, 1);
        j--;
      }
    }
  }

  friends.value.forEach(item => {
    item.selected = false; // 默认所有的选择框为未选中
  });
  friends.value = res.InfoList
}
const choice = (item, event) => {
  if (event.target.checked) {
    item.selected = true
    choiceFriendList.value.push(item)
  } else {
    item.selected = false
    choiceFriendList.value = choiceFriendList.value.filter(newItem => newItem.id !== item.id)
  }
  console.log(choiceFriendList.value)
}
// 切换 checkbox 状态的函数
const toggleCheckbox = (item) => {
  // 切换选中状态
  if (item.selected) {
    choiceFriendList.value = choiceFriendList.value.filter(newItem => newItem.id !== item.id)
  } else {
    choiceFriendList.value.push(item)
  }
  item.selected = !item.selected;
  console.log(choiceFriendList.value)

};

// 处理接收到的事件
function handleClear(data) {
  if (data.message === false) {
    showmessages.value = data.message
    showGroup.value = false
  } else if (data.message === true) {
    showmessages.value = data.message
    showGroup.value = false
  }else if(data.message1 === false){
    showmessages.value = true
    showGroup.value = data.message1
  }else if(data.message1 === true){
    showmessages.value = false
    showGroup.value = data.message1
  }

}

onMounted(() => {
  avatar.value = localStorage.getItem("avatar")
})
const closeD = () => {
  dialogVisible.value = false
  choiceFriendList.value = []
  friends.value.forEach(item => {
    item.selected = false;
  });
}

const CreateRoomSure = async () => {
  console.log(choiceFriendList.value);

  // 构建参数
  const paramsMap = {}; // 使用 map 而非数组
  for (let i = 0; i < choiceFriendList.value.length; i++) {
    paramsMap[`user${i}`] = {
      userID: choiceFriendList.value[i].id.toString(), // 确保是字符串
      userName: choiceFriendList.value[i].username,
    };
  }

  // 添加当前用户信息
  paramsMap["suser"] = {
    userID: localStorage.getItem("userID"),
    userName: "",
  };

  console.log("发送的数据:", paramsMap);

  try {
    // 调用后端接口
    const res = await CreateRoom(paramsMap);
    let list = []
    console.log(res)
    list.push(res.groupInfo)
    console.log("后端返回的数据:", list);
    // 转换为后端期望的格式
    const formattedList = [];
    for (const key in list[0]) {
      if (list[0].hasOwnProperty(key)) {
        formattedList.push(list[0][key]); // 提取每个对象
      }
    }
    console.log("转换后的数据:", formattedList);
    const res1 = await CrateRoomReceive(res.superUserID, formattedList, res.roomName, res.superName, res.roomID)
    const res2 = await QueryRoom(parseInt(localStorage.getItem("userID")))
  } catch (error) {
    console.error("创建房间出错:", error);
  }
}
</script>

<template>
  <div id="index">
    <div class="navigation_bar">
      <div class="img_div">
        <img :src="avatar" alt="">
      </div>
      <LeftNav></LeftNav>
    </div>
    <div class="chat_search">
      <Search @clear="handleClear">
        <template #action-icon>
          <svg class="icon iconadd" aria-hidden="true" font-size="25" @click="createGroup">
            <use xlink:href="#icon-jiahao"></use>
          </svg>
        </template>
      </Search>
      <div class="messages" v-if="showmessages">
        <slot name="two"></slot>
      </div>
      <div class="messages" v-else-if="showGroup">
        <slot name="group"></slot>
      </div>
    </div>
    <el-dialog v-model="dialogVisible" width="400" @close="closeD">
      <template #footer>
        <div>
          <ul class="searchList">
            <li style="cursor: pointer;" v-for="(item,index) in friends" @click="toggleCheckbox(item)">
              <input type="checkbox" v-model="item.selected" :checked="item.selected" @change="choice(item, $event)"
                     @click.stop>
              <img :src="item.avatar" alt="">
              <span style="text-indent: 0.5vw;">{{ item.username }}</span>
            </li>

          </ul>
          <button style="margin-right: 0.5vh" @click="CreateRoomSure">确定</button>
          <button style="margin-left: 0.5vh" @click="closeD">取消</button>
        </div>
      </template>
    </el-dialog>
    <div class="chat_frame">
      <slot name="three">

      </slot>
    </div>
  </div>
</template>

<style scoped lang="less">
#index {
  width: 100%;
  height: 100%;
  background-color: #cfd9df;
  display: flex;

  .navigation_bar {
    width: 3%;
    height: 100%;
    background-color: #2e2e2e;

    .img_div {
      width: 100%;

      img {
        width: 80%;
        height: 2.3vw;
        margin-top: 1vh;
        image-rendering: -moz-crisp-edges; /* Firefox */
        image-rendering: -o-crisp-edges; /* Opera */
        image-rendering: -webkit-optimize-contrast;
        image-rendering: crisp-edges;
        -ms-interpolation-mode: nearest-neighbor; /* IE (non-standard property) */
      }
    }
  }

  .chat_search {
    display: flex;
    flex-direction: column;
    width: 18%;
    background-color: #f7f7f7;

    .messages {
      height: 95vh;
    }
  }

  button {
    border: none;
    cursor: pointer;
    width: 20%;
    padding: 1vh;
    border-radius: 1vh;
    transition: all 0.5s;
  }

  button:hover {
    background-color: black;
    color: white;
    box-shadow: 0 0 10px black;
    transition: all 0.5s;
  }

  .searchList {
    width: 100%;
    padding: 0;
    margin: 0;

    li {
      text-align: left;
      list-style: none;
      padding: 0.5vh;
      position: relative;
      height: 5vh;
      line-height: 5vh;
      color: black;

      img {
        width: 13%;
      //border-radius: 1vh; position: absolute; left: 6%;
        top: 50%;
        image-rendering: -moz-crisp-edges; /* Firefox */
        image-rendering: -o-crisp-edges; /* Opera */
        image-rendering: -webkit-optimize-contrast;
        image-rendering: crisp-edges;
        -ms-interpolation-mode: nearest-neighbor; /* IE (non-standard property) */
        height: 5vh;
        transform: translateY(-50%);
      }

      span {
        position: absolute;
        color: black;
        top: 35%;
        left: 19%;
        transform: translateY(-50%);
      }
    }

  }

  .chat_frame {
    display: flex;
    flex-direction: column;
    width: 79%;
  }
}
</style>