<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-[#F9FAFB] py-8 px-2">
    <!-- Logo -->
    <img src="/src/logo.png" alt="CertMon Logo" class="w-20 h-20 mb-4 mx-auto" />
    <!-- 主標題 -->
    <h1 class="text-2xl sm:text-3xl md:text-4xl font-extrabold text-blue-400 mb-1 tracking-wide">監控列表</h1>
    <!-- 副標題 -->
    <div class="text-sm sm:text-base text-slate-500 mb-8">這些 domain 就交給我守護啦！</div>

    <!-- 卡片 -->
    <div class="bg-white rounded-2xl shadow-md w-full max-w-[800px] flex-col items-center py-8 px-2 sm:px-6 mb-6 overflow-x-auto">
      <!-- 空狀態說明 -->
      <div v-if="domains.length === 0" class="text-slate-400 text-base py-12 text-center">
        這裡會顯示你正在守護的網站，趕快加入第一個 domain 吧！
      </div>
      <!-- 表格 -->
      <div v-else class="overflow-x-auto">
        <table class="w-full min-w-[500px] text-center text-xs sm:text-sm">
          <thead>
            <tr class="text-slate-400 text-xs sm:text-sm h-10 sm:h-12">
              <th class="w-[160px] sm:w-[220px] font-semibold">Domain</th>
              <th class="w-[80px] sm:w-[100px] font-semibold">到期日</th>
              <th class="w-[80px] sm:w-[100px] font-semibold">剩餘天數</th>
              <th class="w-[80px] sm:w-[100px] font-semibold">狀態</th>
              <th class="w-[80px] sm:w-[100px] font-semibold">詳細</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in domains" :key="item.domain" class="h-10 sm:h-[52px] border-t border-slate-100">
              <td class="w-[160px] sm:w-[220px] text-slate-700 truncate">{{ item.domain }}</td>
              <td class="w-[80px] sm:w-[100px] text-slate-700">{{ item.expiry }}</td>
              <td class="w-[80px] sm:w-[100px] text-slate-700">{{ item.daysLeft }}</td>
              <td class="w-[80px] sm:w-[100px] flex items-center justify-center gap-2 mx-auto">
                <span :class="['inline-block w-3 h-3 rounded-full', statusDotColor(item.status)]"></span>
                <span class="text-slate-700">{{ statusText(item.status) }}</span>
              </td>
              <td class="w-[80px] sm:w-[100px]">
                <button class="text-blue-400 hover:underline font-semibold" @click="goToDetail(item)">詳細</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <!-- 新增監控按鈕 -->
    <div class="flex gap-4 justify-center w-full">
      <button class="bg-blue-400 hover:bg-blue-500 text-white text-base px-6 sm:px-8 py-3 rounded-[16px] font-semibold transition w-full max-w-xs sm:max-w-none sm:w-auto" @click="goToAdd">
        新增監控
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 範例資料，實際可由父層 props 或 API 傳入
const domains = ref([
  { domain: 'example.com', expiry: '2024/07/01', daysLeft: 23, status: 'warning' },
  { domain: 'www.test.com', expiry: '2024/12/20', daysLeft: 195, status: 'ok' },
  { domain: 'api.foo.com', expiry: '2025/03/15', daysLeft: 280, status: 'ok' },
])

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
function goToDetail(item) {
  router.push({ name: 'detail', params: { domain: item.domain } })
}
function goToAdd() {
  router.push({ name: 'home' })
}
</script>
