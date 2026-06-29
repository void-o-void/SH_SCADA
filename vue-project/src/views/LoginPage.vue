<script setup>
import { ref, computed } from 'vue'
import qrcode from '@/assets/images/erweima.png'
import selectedIcon from '@/assets/icons/selected.svg'
import unselectIcon from '@/assets/icons/unselect.svg'

// ====== 状态变量（类比 QML 的 property） ======
const tabIsLogin = ref(true)
const username = ref('')
const password = ref('')
const rememberPwd = ref(false)
const agreeTerms = ref(false)
const showSuggestions = ref(false)

// 注册表单
const regUsername = ref('')
const regPassword = ref('')
const regEmail = ref('')

// 登录历史（类比 QML 的 ListModel）
const userHistory = ref([])

// 过滤后的建议列表
const filteredSuggestions = computed(() => {
    if (!username.value) return userHistory.value
    return userHistory.value.filter(u => u.includes(username.value))
})

// ====== 操作函数 ======
function selectSuggestion(name) {
    username.value = name
    showSuggestions.value = false
}

function handleLogin() {
    if (!username.value || !password.value) {
        ElMessage.warning('请输入账号和密码')
        return
    }
    // 保存登录历史（去重）
    if (username.value && !userHistory.value.includes(username.value)) {
        userHistory.value.push(username.value)
    }
    ElMessage.success(`欢迎 ${username.value}！登录成功`)
    // 实际项目中这里跳转到主页面
}

function handleRegister() {
    if (!regUsername.value || !regPassword.value || !regEmail.value) {
        ElMessage.warning('请填写完整的注册信息')
        return
    }
    ElMessage.success('注册成功！请登录')
    tabIsLogin.value = true
}
</script>



<template>
    <div class="login-wrapper">
        <!--
      主卡片：类比 QML 的 ApplicationWindow
      1000×600，圆角 20px，背景图片
    -->
        <div class="login-card">
            <!-- 左侧：背景图 + 二维码 -->
            <div class="left-panel">
                <img :src="qrcode" class="qrcode-img" alt="二维码" />
            </div>

            <!-- 右侧：半透明白色面板 -->
            <div class="right-panel">
                <!-- ==================== 登录区域 ==================== -->
                <template v-if="tabIsLogin">
                    <!-- 右上角"去注册"链接 -->
                    <span class="toggle-link" @click="tabIsLogin = false">去注册账户？</span>

                    <!-- 标题：类比 QML login_title -->
                    <div class="title">SHiot平台</div>

                    <!-- 用户名提示 -->
                    <div class="input-tip">请输入账号</div>

                    <!-- 用户名输入框（与密码框同款 el-input）+ 历史建议下拉 -->
                    <div class="username-wrapper">
                        <el-input v-model="username" placeholder="请输入账号" class="login-input"
                            @focus="showSuggestions = true" @blur="showSuggestions = false" />
                        <!-- 历史建议下拉列表 -->
                        <ul v-if="showSuggestions && filteredSuggestions.length > 0" class="suggestions-dropdown">
                            <li v-for="(item, index) in filteredSuggestions" :key="index" class="suggestion-item"
                                @mousedown.prevent="selectSuggestion(item)">
                                {{ item }}
                            </li>
                        </ul>
                    </div>

                    <!-- 密码提示 -->
                    <div class="input-tip" style="margin-top: 20px">请输入密码</div>

                    <!-- 密码输入框：类比 QML 的 FramelessLineEdit + echoMode: Password -->
                    <el-input v-model="password" type="password" show-password placeholder="请输入密码"
                        class="login-input" />

                    <!-- 记住密码 -->
                    <div class="checkbox-row" style="margin-top: 35px">
                        <img :src="rememberPwd ? selectedIcon : unselectIcon" class="checkbox-icon"
                            @click="rememberPwd = !rememberPwd" />
                        <span class="checkbox-label">记住密码</span>
                    </div>

                    <!-- 同意协议 -->
                    <div class="checkbox-row" style="margin-top: 20px">
                        <img :src="agreeTerms ? selectedIcon : unselectIcon" class="checkbox-icon"
                            @click="agreeTerms = !agreeTerms" />
                        <span class="checkbox-label">已阅读并同意隐私保护协议以及互联网……</span>
                    </div>

                    <!-- 登录按钮：类比 QML 绿色按钮，圆角 10px -->
                    <button class="submit-btn" @click="handleLogin">登录</button>
                </template>

                <!-- ==================== 注册区域 ==================== -->
                <template v-else>
                    <!-- 右上角"去登录"链接 -->
                    <span class="toggle-link" @click="tabIsLogin = true">已有账户，去登录？</span>

                    <!-- 标题 -->
                    <div class="title">SHiot平台</div>

                    <!-- 注册副标题 -->
                    <div class="reg-subtitle">让物联更有趣，欢迎加入！</div>

                    <!-- 注册表单 -->
                    <div class="input-tip" style="margin-top: 40px">请输入账户</div>
                    <el-input v-model="regUsername" placeholder="请输入账户" class="login-input" />

                    <div class="input-tip" style="margin-top: 20px">请输入密码</div>
                    <el-input v-model="regPassword" type="password" show-password placeholder="请输入密码"
                        class="login-input" />

                    <div class="input-tip" style="margin-top: 20px">请输入邮箱</div>
                    <el-input v-model="regEmail" placeholder="请输入邮箱" class="login-input" />

                    <button class="submit-btn" style="margin-top: 60px" @click="handleRegister">提交</button>
                </template>
            </div>
        </div>
    </div>
</template>

<style>
/* 全局重置 */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

/* 引入 QML 同款标题字体（YouSheBiaoTiHei，有神标题黑） */
@font-face {
    font-family: 'YouSheBiaoTiHei';
    src: url('@/assets/fonts/YouSheBiaoTiHei-2.ttf') format('truetype');
    font-weight: normal;
    font-style: normal;
}

/* 引入 QML 注册副标题字体（SmileySans，得意黑斜体） */
@font-face {
    font-family: 'SmileySans';
    src: url('@/assets/fonts/SmileySans-Oblique.otf') format('opentype');
    font-weight: normal;
    font-style: oblique;
}

body {
    width: 100vw;
    height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
    font-family: 'Microsoft YaHei', sans-serif;
}
</style>

<style scoped>
/* ====== 外层容器：居中显示 ====== */
.login-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100vw;
    height: 100vh;
}

/* ====== 主卡片：1000×600，圆角，类比 QML Rectangle with radius: 20 ====== */
.login-card {
    width: 1000px;
    height: 600px;
    border-radius: 20px;
    overflow: hidden;
    position: relative;
    background: url('@/assets/images/login_bg4.png') center/cover no-repeat;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

/* ====== 左侧面板：二维码区域 ====== */
.left-panel {
    position: absolute;
    left: 0;
    top: 0;
    width: 50%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}

.qrcode-img {
    width: 40%;
    height: auto;
}

/* ====== 右侧面板：半透明白色 ======
  类比 QML color: "#F2FFFFFF" (约 95% 不透明度) ====== */
.right-panel {
    position: absolute;
    right: 0;
    top: 0;
    width: 50%;
    height: 100%;
    background: rgba(255, 255, 255, 0.95);
    padding: 0 11%;
    display: flex;
    flex-direction: column;
    justify-content: center;
}

/* ====== 切换登录/注册链接（右上角） ====== */
.toggle-link {
    position: absolute;
    top: 20px;
    right: 24px;
    color: #0066CC;
    font-size: 16px;
    text-decoration: underline;
    cursor: pointer;
}

/* ====== 标题：绿色 #21C491，QML 同款字体 + 斜体 ====== */
.title {
    color: #21C491;
    font-family: 'YouSheBiaoTiHei', 'Microsoft YaHei', sans-serif;
    font-size: 48px;
    font-style: italic;
    text-align: center;
    margin-bottom: 40px;
}

/* 注册副标题：QML 同款得意黑斜体，pixelSize 30 */
.reg-subtitle {
    font-family: 'YouSheBiaoTiHei', 'Microsoft YaHei', sans-serif;
    font-size: 30px;
    color: #000;
    text-align: left;
    white-space: nowrap;
    margin-top: 15px;
    margin-bottom: 10px;
}

/* ====== 输入提示文字 ====== */
.input-tip {
    color: rgba(26, 26, 26, 0.3);
    font-size: 14px;
    margin-bottom: 4px;
}

/* ====== 输入框：去边框风格，类比 QML FramelessLineEdit ====== */
.login-input {
    width: 100%;
}

.login-input :deep(.el-input__wrapper) {
    border: none;
    border-bottom: 2px solid #e0e0e0;
    border-radius: 0;
    box-shadow: none;
    padding: 4px 0;
    background: transparent;
    font-size: 18px;
    transition: border-color 0.3s;
}

.login-input :deep(.el-input__wrapper:hover) {
    border-bottom-color: #21C491;
}

.login-input :deep(.el-input__wrapper.is-focus) {
    border-bottom-color: #21C491;
    box-shadow: none;
}

.login-input :deep(.el-input__inner) {
    font-size: 18px;
    font-family: 'Microsoft YaHei', sans-serif;
}

/* ====== 用户名字段 + 建议下拉列表 ====== */
.username-wrapper {
    position: relative;
    width: 100%;
}

.suggestions-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 1000;
    max-height: 180px;
    overflow-y: auto;
    margin: 4px 0 0 0;
    padding: 0;
    list-style: none;
    background: #fff;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
}

.suggestion-item {
    padding: 10px 16px;
    font-size: 16px;
    color: #333;
    cursor: pointer;
    transition: background 0.15s;
}

.suggestion-item:hover {
    background: #ecf5ff;
    color: #21C491;
}

/* ====== 复选框行 ====== */
.checkbox-row {
    display: flex;
    align-items: center;
    gap: 8px;
}

.checkbox-icon {
    width: 16px;
    height: 16px;
    cursor: pointer;
    user-select: none;
}

.checkbox-label {
    font-size: 14px;
    color: #333;
}

/* ====== 提交按钮：绿色圆角，类比 QML loginButton ====== */
.submit-btn {
    width: 100%;
    height: 45px;
    border: none;
    border-radius: 10px;
    background: #5CBF60;
    color: white;
    font-size: 20px;
    font-weight: bold;
    font-family: 'Microsoft YaHei', sans-serif;
    cursor: pointer;
    transition: background 0.2s;
    margin-top: 30px;
}

.submit-btn:hover {
    background: #4CAF50;
}

.submit-btn:active {
    background: #3d8b40;
}
</style>
