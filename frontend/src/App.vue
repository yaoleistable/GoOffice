<script setup>
import { ref, computed } from 'vue'
import { SelectFiles, ExtractPages } from '../wailsjs/go/main/App'

const selectedFiles = ref([])
const pageNumbers = ref('')
const results = ref([])

const canExtract = computed(() => {
  return selectedFiles.value.length > 0 && 
         pageNumbers.value.trim() !== '' &&
         /^[0-9,\s-]+$/.test(pageNumbers.value)
})

async function selectFiles() {
  try {
    const files = await SelectFiles()
    selectedFiles.value = files
    results.value = []
  } catch (err) {
    console.error('选择文件失败:', err)
  }
}

function removeFile(file) {
  selectedFiles.value = selectedFiles.value.filter(f => f.path !== file.path)
  results.value = []
}

async function extractPages() {
  if (!canExtract.value) return

  try {
    const pageRange = pageNumbers.value.trim()
    const result = await ExtractPages(selectedFiles.value.map(f => f.path), pageRange)
    results.value = result
  } catch (err) {
    console.error('提取页面失败:', err)
    results.value = selectedFiles.value.map(f => ({
      file: f.name,
      success: false,
      message: '处理失败：' + err.message
    }))
  }
}
</script>

<template>
  <div class="app-container">
    <header class="app-header">
      <h1>PDF页面提取工具</h1>
      <p class="subtitle">选择PDF文件并提取指定页面</p>
    </header>

    <main class="app-main">
      <!-- 文件选择区域 -->
      <section class="file-section">
        <div class="section-header">
          <h2>选择PDF文件</h2>
          <button class="btn-select" @click="selectFiles">
            <span class="icon">📄</span>
            选择文件
          </button>
        </div>

        <!-- 文件列表 -->
        <div v-if="selectedFiles.length > 0" class="file-list">
          <div v-for="file in selectedFiles" :key="file.path" class="file-item">
            <div class="file-info">
              <span class="file-name">{{ file.name }}</span>
              <span class="file-pages">共 {{ file.pages }} 页</span>
            </div>
            <button class="btn-remove" @click="removeFile(file)">×</button>
          </div>
        </div>
      </section>

      <!-- 页面选择区域 -->
      <section class="pages-section">
        <div class="section-header">
          <h2>选择要提取的页面</h2>
        </div>
        <div class="input-group">
          <input 
            type="text" 
            v-model="pageNumbers" 
            placeholder="输入页码，用逗号分隔（如：1,3,5或1-3）"
            :disabled="selectedFiles.length === 0"
            class="page-input"
          />
          <div class="input-hint">支持输入多个页码，用逗号分隔</div>
        </div>
      </section>

      <!-- 操作按钮 -->
      <section class="action-section">
        <button 
          class="btn-extract" 
          @click="extractPages"
          :disabled="!canExtract"
          :class="{ 'btn-disabled': !canExtract }"
        >
          <span class="icon">✂️</span>
          提取页面
        </button>
      </section>

      <!-- 结果展示区域 -->
      <section v-if="results.length > 0" class="results-section">
        <div class="section-header">
          <h2>处理结果</h2>
        </div>
        <div class="results-list">
          <div v-for="result in results" :key="result.file" class="result-item">
            <div class="result-status" :class="result.success ? 'success' : 'error'">
              {{ result.success ? '✓' : '✗' }}
            </div>
            <div class="result-info">
              <div class="result-file">{{ result.file }}</div>
              <div class="result-message">{{ result.message }}</div>
            </div>
          </div>
        </div>
      </section>
    </main>

    <footer class="app-footer">
      <p>提取后的文件将保存在原文件所在目录，"output"文件夹</p>
    </footer>
  </div>
</template>

<style>
.app-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.app-header {
  text-align: center;
  margin-bottom: 2rem;
}

.app-header h1 {
  color: #2c3e50;
  margin: 0;
  font-size: 2rem;
}

.subtitle {
  color: #666;
  margin-top: 0.5rem;
}

.app-main {
  flex: 1;
}

section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.section-header h2 {
  margin: 0;
  color: #2c3e50;
  font-size: 1.25rem;
}

.btn-select {
  background: #4CAF50;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: background-color 0.2s;
}

.btn-select:hover {
  background: #45a049;
}

.file-list {
  margin-top: 1rem;
}

.file-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 0.5rem;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.file-name {
  font-weight: 500;
  color: #2c3e50;
}

.file-pages {
  color: #666;
  font-size: 0.9rem;
}

.btn-remove {
  background: none;
  border: none;
  color: #dc3545;
  font-size: 1.25rem;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.btn-remove:hover {
  background: rgba(220, 53, 69, 0.1);
}

.input-group {
  margin-top: 1rem;
}

.page-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.page-input:focus {
  outline: none;
  border-color: #4CAF50;
}

.input-hint {
  margin-top: 0.5rem;
  color: #666;
  font-size: 0.9rem;
}

.btn-extract {
  width: 100%;
  background: #2196F3;
  color: white;
  border: none;
  padding: 1rem;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: background-color 0.2s;
}

.btn-extract:hover:not(.btn-disabled) {
  background: #1976D2;
}

.btn-disabled {
  background: #ccc;
  cursor: not-allowed;
}

.results-list {
  margin-top: 1rem;
}

.result-item {
  display: flex;
  align-items: flex-start;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 0.5rem;
}

.result-status {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 1rem;
  font-weight: bold;
}

.result-status.success {
  background: #d4edda;
  color: #155724;
}

.result-status.error {
  background: #f8d7da;
  color: #721c24;
}

.result-info {
  flex: 1;
}

.result-file {
  font-weight: 500;
  color: #2c3e50;
  margin-bottom: 0.25rem;
}

.result-message {
  color: #666;
  font-size: 0.9rem;
}

.app-footer {
  text-align: center;
  color: #666;
  font-size: 0.9rem;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid #dee2e6;
}
</style>

