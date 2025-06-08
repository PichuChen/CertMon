<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-[#F9FAFB] px-2">
    <!-- Logo -->
    <img src="/src/logo.png" alt="CertMon Logo" class="w-32 h-32 mb-4 mx-auto" />

    <!-- 主標題 -->
    <h1 class="text-3xl sm:text-4xl font-extrabold text-blue-400 mb-2 tracking-wide">CertMon</h1>
    <!-- 副標題 -->
    <div class="text-base sm:text-lg text-slate-500 mb-8">SSL 憑證到期自動監控</div>

    <!-- 輸入區塊 -->
    <div class="flex flex-col sm:flex-row items-center w-full max-w-md mb-2 gap-2 sm:gap-0">
      <input
        v-model="inputDomain"
        type="text"
        placeholder="輸入你的 domain，讓我來守護！"
        class="flex-1 rounded-[16px] px-4 py-3 border border-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-200 bg-white text-base shadow-sm placeholder:text-slate-400"
        @keyup.enter="addDomain"
        :disabled="loading"
      />
      <button
        class="sm:ml-4 bg-blue-400 hover:bg-blue-500 text-white text-base px-6 py-3 rounded-[16px] font-semibold transition w-full sm:w-auto"
        @click="addDomain"
        :disabled="loading"
      >
        <span v-if="loading">新增中...</span>
        <span v-else>新增監控</span>
      </button>
    </div>
    <div v-if="errorMsg" class="text-red-500 text-xs mb-2">{{ errorMsg }}</div>
    <!-- 範例說明 -->
    <div class="text-xs text-slate-400 mb-8">
      範例：<span class="text-slate-500">example.com</span>
    </div>

    <!-- 導覽列 -->
    <div class="flex flex-col sm:flex-row sm:space-x-6 mb-2 w-full max-w-xs">
      <button class="text-blue-500 font-semibold border-b-2 border-blue-500 pb-1 cursor-pointer bg-transparent w-full sm:w-auto" @click="goToList">
        監控列表
      </button>
      <button class="text-slate-400 font-semibold pb-1 cursor-not-allowed bg-transparent w-full sm:w-auto" disabled>
        通知設定
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()

const inputDomain = ref('')
const loading = ref(false)
const errorMsg = ref('')

async function addDomain() {
  errorMsg.value = ''
  if (!inputDomain.value.trim()) {
    errorMsg.value = '請輸入 domain'
    return
  }
  loading.value = true
  try {
    const res = await fetch('/api/domains', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ domain: inputDomain.value.trim() })
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      errorMsg.value = err.error || '新增失敗'
      loading.value = false
      return
    }
    inputDomain.value = ''
    router.push({ name: 'list' })
  } catch (e) {
    errorMsg.value = '網路錯誤'
  } finally {
    loading.value = false
  }
}

function goToList() {
  router.push({ name: 'list' })
}
</script>
