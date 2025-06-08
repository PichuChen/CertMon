<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-[#F9FAFB] py-8 px-2">
    <!-- Logo -->
    <img src="/src/logo.png" alt="CertMon Logo" class="w-20 h-20 mb-4 mx-auto" />
    <!-- 主標題 -->
    <h1 class="text-2xl sm:text-3xl md:text-4xl font-extrabold text-blue-400 mb-1 tracking-wide">憑證詳細</h1>
    <!-- 副標題 -->
    <div class="text-sm sm:text-base text-slate-500 mb-8">
      <template v-if="cert && cert.domain">
        這是 {{ cert.domain }} 的 SSL 憑證詳細資訊喔！
      </template>
      <template v-else-if="!error">
        載入中...
      </template>
    </div>

    <!-- 卡片 -->
    <div class="bg-white rounded-[16px] shadow-md w-full max-w-[800px] px-2 sm:px-8 py-8 mb-6">
      <template v-if="cert !== null && !error">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 sm:gap-x-8 gap-y-4">
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">網域 Domain</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal break-all">{{ cert.domain }}</span>
          </div>
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">簽發單位</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal break-all">{{ cert.issuer }}</span>
          </div>
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">簽發日</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal">{{ cert.startDate }}</span>
          </div>
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">到期日</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal">{{ cert.expiry }}</span>
          </div>
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">剩餘天數</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal">{{ cert.daysLeft }}</span>
          </div>
          <div class="flex flex-col sm:flex-row items-start sm:items-center">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">狀態</span>
            <span :class="['inline-block w-3 h-3 rounded-full mr-2', statusDotColor(cert.status)]"></span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal">{{ statusText(cert.status) }}</span>
          </div>
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">憑證序號</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal break-all">{{ cert.serial }}</span>
          </div>
          <div class="flex flex-col sm:flex-row">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">憑證類型</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal">{{ cert.type }}</span>
          </div>
          <div class="flex flex-col sm:flex-row col-span-1 sm:col-span-2">
            <span class="w-32 text-slate-500 text-[15px] sm:text-[16px] font-medium">SAN</span>
            <span class="text-slate-800 text-[15px] sm:text-[16px] font-normal break-all">{{ cert.san }}</span>
          </div>
        </div>
      </template>
      <template v-else-if="error">
        <div class="text-slate-500 text-base text-center py-12 font-medium">
          找不到憑證資料，請重新整理或回到列表頁
        </div>
      </template>
      <template v-else>
        <div class="text-slate-400 text-base text-center py-12 font-medium">
          載入中...
        </div>
      </template>
    </div>
    <!-- 按鈕區 -->
    <div class="flex flex-col sm:flex-row gap-4 w-full max-w-xs sm:max-w-none sm:w-auto justify-center">
      <button class="bg-blue-400 hover:bg-blue-500 text-white text-base px-6 sm:px-8 py-3 rounded-[16px] font-semibold transition w-full sm:w-auto" @click="goBack">
        返回監控列表
      </button>
      <button class="bg-slate-200 text-slate-500 text-base px-6 sm:px-8 py-3 rounded-[16px] font-semibold transition w-full sm:w-auto" disabled>
        複製憑證
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const cert = ref(null)
const error = ref(false)

onMounted(async () => {
  try {
    const res = await fetch(`/api/domains/${route.params.domain}`)
    if (!res.ok) throw new Error('API error')
    const data = await res.json()
    cert.value = {
      domain: data.domain,
      issuer: data.issuer || 'Let\'s Encrypt',
      startDate: data.valid_from ? data.valid_from.slice(0, 10).replace(/-/g, '/') : '',
      expiry: data.valid_to ? data.valid_to.slice(0, 10).replace(/-/g, '/') : '',
      daysLeft: data.days_left,
      status: convertStatus(data.status),
      serial: data.serial || '',
      type: data.type || '',
      san: data.san || data.domain
    }
  } catch (e) {
    error.value = true
  }
})

function convertStatus(status) {
  if (status === 'valid') return 'ok'
  if (status === 'expiring') return 'warning'
  return 'error'
}

function statusDotColor(status) {
  if (status === 'ok') return 'bg-green-400'
  if (status === 'warning') return 'bg-yellow-400'
  return 'bg-red-400'
}
function statusText(status) {
  if (status === 'ok') return '正常'
  if (status === 'warning') return '快到期'
  return '已過期/斷線'
}
function goBack() {
  router.push({ name: 'list' })
}
</script>
