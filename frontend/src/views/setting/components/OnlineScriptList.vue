<template>
  <el-dialog
    v-model="visible"
    :title="$t('settings.scripts.scriptMarket')"
    width="90%"
    :close-on-click-modal="false"
  >
    <div class="online-script-list">
      <!-- 搜索和筛选区域 -->
      <div class="filter-section">
        <el-input
          v-model="searchKeyword"
          :placeholder="$t('settings.scripts.searchPlaceholder')"
          clearable
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <el-select
          v-model="selectedCategory"
          :placeholder="$t('settings.scripts.categoryFilter')"
          filterable
          default-first-option
          :reserve-keyword="false"
          collapse-tags-tooltip
          clearable
          class="category-select"
        >
          <el-option :label="$t('settings.scripts.allCategories')" value="" />
          <el-option
            v-for="category in categories"
            :key="category"
            :label="category"
            :value="category"
          />
        </el-select>

        <el-select
          v-model="selectedTags"
          :placeholder="$t('settings.scripts.tagFilter')"
          multiple
          filterable
          default-first-option
          :reserve-keyword="false"
          collapse-tags-tooltip
          clearable
          collapse-tags
          class="tag-select"
        >
          <el-option
            v-for="tag in allTags"
            :key="tag"
            :label="tag"
            :value="tag"
          />
        </el-select>

        <el-button :icon="Refresh" @click="loadPlugins" :loading="loading">
          {{ $t("common.refresh") }}
        </el-button>
      </div>

      <!-- 统计信息 -->
      <div class="stats-section" v-if="filteredPlugins.length > 0">
        <span class="stats-text">
          {{
            $t("settings.scripts.totalPlugins", {
              total: filteredPlugins.length,
              all: plugins.length,
            })
          }}
        </span>
      </div>

      <!-- 插件列表 -->
      <div class="plugins-container" v-loading="loading">
        <div
          v-if="!loading && filteredPlugins.length === 0"
          class="empty-state"
        >
          <el-empty :description="$t('settings.scripts.noPluginsFound')" />
        </div>

        <div v-else class="plugins-grid">
          <el-card
            v-for="plugin in filteredPlugins"
            :key="plugin.id"
            class="plugin-card"
            shadow="hover"
          >
            <template #header>
              <div class="plugin-header">
                <div class="plugin-title-section">
                  <h3 class="plugin-name">{{ plugin.name }}</h3>
                  <div class="plugin-version-info">
                    <el-tag
                      size="small"
                      :type="hasUpdate(plugin) ? 'warning' : 'info'"
                      class="plugin-version"
                    >
                      v{{ plugin.version }}
                    </el-tag>
                    <span
                      v-if="
                        isInstalled(plugin) &&
                        installedMap.get(plugin.id)?.version
                      "
                      class="installed-version"
                    >
                      ({{ $t("settings.scripts.installedVersion") }}: v{{
                        installedMap.get(plugin.id)?.version
                      }})
                    </span>
                  </div>
                </div>
                <div class="plugin-actions">
                  <el-button
                    v-if="isInstalled(plugin) && hasUpdate(plugin)"
                    type="warning"
                    size="small"
                    :loading="updatingMap.get(plugin.id) || false"
                    @click="handleUpdate(plugin)"
                  >
                    <el-icon><Refresh /></el-icon>
                    {{ $t("settings.scripts.update") }}
                  </el-button>
                  <el-button
                    v-else-if="isInstalled(plugin)"
                    type="success"
                    size="small"
                    disabled
                  >
                    <el-icon><Check /></el-icon>
                    {{ $t("settings.scripts.installed") }}
                  </el-button>
                  <el-button
                    v-else
                    type="primary"
                    size="small"
                    :loading="installingMap.get(plugin.id) || false"
                    @click="handleInstall(plugin)"
                  >
                    <el-icon><Download /></el-icon>
                    {{ $t("settings.scripts.install") }}
                  </el-button>
                </div>
              </div>
            </template>

            <div class="plugin-content">
              <p class="plugin-description">{{ plugin.description }}</p>

              <div class="plugin-meta">
                <div class="meta-item">
                  <el-icon><User /></el-icon>
                  <span>{{ plugin.author }}</span>
                </div>
                <div class="meta-item">
                  <el-icon><FolderOpened /></el-icon>
                  <span>{{ plugin.category }}</span>
                </div>
              </div>

              <div
                class="plugin-tags"
                v-if="plugin.tags && plugin.tags.length > 0"
              >
                <el-tag
                  v-for="tag in plugin.tags"
                  :key="tag"
                  size="small"
                  class="tag-item"
                  @click="toggleTag(tag)"
                >
                  {{ tag }}
                </el-tag>
              </div>

              <div
                class="plugin-content-types"
                v-if="plugin.contentTypes && plugin.contentTypes.length > 0"
              >
                <span class="content-types-label"
                  >{{ $t("settings.scripts.supportedTypes") }}:</span
                >
                <el-tag
                  v-for="type in plugin.contentTypes"
                  :key="type"
                  size="small"
                  type="info"
                  class="content-type-tag"
                >
                  {{ type }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Search,
  Refresh,
  Download,
  Check,
  User,
  FolderOpened,
} from "@element-plus/icons-vue";
import { useI18n } from "vue-i18n";
import {
  SaveUserScript,
  GetAllUserScripts,
} from "../../../../wailsjs/go/main/App";

const { t } = useI18n();

interface Plugin {
  id: string;
  name: string;
  description: string;
  author: string;
  version: string;
  category: string;
  tags: string[];
  scriptUrl: string;
  icon: string;
  trigger: string;
  contentTypes: string[];
  keywords: string[];
}

interface PluginsResponse {
  version: string;
  lastUpdated?: string;
  plugins: Plugin[];
}

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  "update:modelValue": [value: boolean];
  installed: [];
}>();

const visible = ref(false);
const loading = ref(false);
const plugins = ref<Plugin[]>([]);
const searchKeyword = ref("");
const selectedCategory = ref("");
const selectedTags = ref<string[]>([]);
const installingMap = ref<Map<string, boolean>>(new Map());
const updatingMap = ref<Map<string, boolean>>(new Map());
const installedMap = ref<Map<string, { version: string }>>(new Map());

// Cloudflare Pages URL
const PLUGINS_JSON_URL = "https://clip-save-plugins.pages.dev/plugins.json";

watch(
  () => props.modelValue,
  (val) => {
    visible.value = val;
    if (val) {
      loadPlugins();
      checkInstalledPlugins();
    }
  }
);

watch(visible, (val) => {
  emit("update:modelValue", val);
});

// 获取所有分类
const categories = computed(() => {
  const cats = new Set<string>();
  plugins.value.forEach((plugin) => {
    if (plugin.category) {
      cats.add(plugin.category);
    }
  });
  return Array.from(cats).sort();
});

// 获取所有标签
const allTags = computed(() => {
  const tags = new Set<string>();
  plugins.value.forEach((plugin) => {
    if (plugin.tags && plugin.tags.length > 0) {
      plugin.tags.forEach((tag) => tags.add(tag));
    }
  });
  return Array.from(tags).sort();
});

// 过滤插件
const filteredPlugins = computed(() => {
  let result = plugins.value;

  // 搜索关键词过滤
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter((plugin) => {
      return (
        plugin.name.toLowerCase().includes(keyword) ||
        plugin.description.toLowerCase().includes(keyword) ||
        plugin.author.toLowerCase().includes(keyword) ||
        (plugin.tags &&
          plugin.tags.some((tag) => tag.toLowerCase().includes(keyword)))
      );
    });
  }

  // 分类过滤
  if (selectedCategory.value) {
    result = result.filter(
      (plugin) => plugin.category === selectedCategory.value
    );
  }

  // 标签过滤
  if (selectedTags.value.length > 0) {
    result = result.filter((plugin) => {
      return (
        plugin.tags &&
        selectedTags.value.some((tag) => plugin.tags.includes(tag))
      );
    });
  }

  return result;
});

// 加载插件列表
async function loadPlugins() {
  loading.value = true;
  try {
    const response = await fetch(PLUGINS_JSON_URL);
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
    const data: PluginsResponse = await response.json();
    plugins.value = data.plugins || [];
    // ElMessage.success(
    //   t("settings.scripts.loadPluginsSuccess", { count: plugins.value.length })
    // );
  } catch (error: any) {
    console.error("加载插件列表失败:", error);
    ElMessage.error(
      t("settings.scripts.loadPluginsError", { error: error.message || error })
    );
  } finally {
    loading.value = false;
  }
}

// 检查已安装的插件
async function checkInstalledPlugins() {
  try {
    installedMap.value = new Map();
    const userScripts = await GetAllUserScripts();
    const installed = new Map<string, { version: string }>();

    // 直接使用 PluginID 和 PluginVersion 字段
    userScripts.forEach((script) => {
      if (script.PluginID) {
        installed.set(script.PluginID, {
          version: script.PluginVersion || "",
        });
      }
    });

    installedMap.value = installed;
  } catch (error: any) {
    console.error("检查已安装插件失败:", error);
  }
}

// 比较版本号（简单语义化版本比较）
function compareVersions(v1: string, v2: string): number {
  if (!v1 || !v2) return 0;

  const parts1 = v1.split(".").map(Number);
  const parts2 = v2.split(".").map(Number);
  const maxLen = Math.max(parts1.length, parts2.length);

  for (let i = 0; i < maxLen; i++) {
    const num1 = parts1[i] || 0;
    const num2 = parts2[i] || 0;
    if (num1 > num2) return 1;
    if (num1 < num2) return -1;
  }
  return 0;
}

// 检查是否有更新
function hasUpdate(plugin: Plugin): boolean {
  const installed = installedMap.value.get(plugin.id);
  if (!installed || !installed.version) return false;
  return compareVersions(plugin.version, installed.version) > 0;
}

// 检查是否已安装
function isInstalled(plugin: Plugin): boolean {
  return installedMap.value.has(plugin.id);
}

// 切换标签筛选
function toggleTag(tag: string) {
  const index = selectedTags.value.indexOf(tag);
  if (index > -1) {
    selectedTags.value.splice(index, 1);
  } else {
    selectedTags.value.push(tag);
  }
}

// 更新插件
async function handleUpdate(plugin: Plugin) {
  if (!isInstalled(plugin)) {
    ElMessage.warning(t("settings.scripts.notInstalled"));
    return;
  }

  const installed = installedMap.value.get(plugin.id);
  if (!hasUpdate(plugin)) {
    ElMessage.info(t("settings.scripts.alreadyLatestVersion"));
    return;
  }

  try {
    await ElMessageBox.confirm(
      t("settings.scripts.updateConfirm", {
        name: plugin.name,
        currentVersion: installed?.version || "",
        newVersion: plugin.version,
      }),
      t("settings.scripts.updateTitle"),
      {
        confirmButtonText: t("settings.scripts.update"),
        cancelButtonText: t("common.cancel"),
        type: "warning",
      }
    );

    updatingMap.value.set(plugin.id, true);

    // 下载新版本脚本内容
    const scriptResponse = await fetch(plugin.scriptUrl);
    if (!scriptResponse.ok) {
      throw new Error(`下载脚本失败: HTTP ${scriptResponse.status}`);
    }
    const scriptContent = await scriptResponse.text();

    // 获取已安装的脚本ID
    const userScripts = await GetAllUserScripts();
    const existingScript = userScripts.find((s) => s.PluginID === plugin.id);

    if (!existingScript) {
      throw new Error("未找到已安装的脚本");
    }

    // 更新脚本数据
    const scriptData = {
      ID: existingScript.ID,
      Name: plugin.name,
      Enabled: existingScript.Enabled, // 保持原有启用状态
      Trigger: plugin.trigger || "manual",
      ContentType: plugin.contentTypes || [],
      Keywords: plugin.keywords || [],
      Script: scriptContent,
      Description: plugin.description || "",
      PluginID: plugin.id,
      PluginVersion: plugin.version, // 更新版本号
      SortOrder: existingScript.SortOrder, // 保持原有排序
    };

    // 保存脚本
    await SaveUserScript(JSON.stringify(scriptData));

    installedMap.value.set(plugin.id, { version: plugin.version });
    ElMessage.success(
      t("settings.scripts.updateSuccess", {
        name: plugin.name,
        version: plugin.version,
      })
    );
    emit("installed");
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("更新插件失败:", error);
      ElMessage.error(
        t("settings.scripts.updateError", { error: error.message || error })
      );
    }
  } finally {
    updatingMap.value.set(plugin.id, false);
  }
}

// 安装插件
async function handleInstall(plugin: Plugin) {
  if (isInstalled(plugin)) {
    ElMessage.warning(t("settings.scripts.alreadyInstalled"));
    return;
  }

  try {
    await ElMessageBox.confirm(
      t("settings.scripts.installConfirm", { name: plugin.name }),
      t("settings.scripts.installTitle"),
      {
        confirmButtonText: t("settings.scripts.install"),
        cancelButtonText: t("common.cancel"),
        type: "info",
      }
    );

    installingMap.value.set(plugin.id, true);

    // 下载脚本内容
    const scriptResponse = await fetch(plugin.scriptUrl);
    if (!scriptResponse.ok) {
      throw new Error(`下载脚本失败: HTTP ${scriptResponse.status}`);
    }
    const scriptContent = await scriptResponse.text();

    // 创建脚本数据
    const scriptData = {
      ID: "",
      Name: plugin.name,
      Enabled: true,
      Trigger: plugin.trigger || "manual",
      ContentType: plugin.contentTypes || [],
      Keywords: plugin.keywords || [],
      Script: scriptContent,
      Description: plugin.description || "",
      PluginID: plugin.id, // 直接设置插件ID
      PluginVersion: plugin.version, // 保存版本号
      SortOrder: 0,
    };

    // 保存脚本
    await SaveUserScript(JSON.stringify(scriptData));

    installedMap.value.set(plugin.id, { version: plugin.version });
    ElMessage.success(
      t("settings.scripts.installSuccess", { name: plugin.name })
    );
    emit("installed");
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("安装插件失败:", error);
      ElMessage.error(
        t("settings.scripts.installError", { error: error.message || error })
      );
    }
  } finally {
    installingMap.value.set(plugin.id, false);
  }
}

onMounted(() => {
  if (visible.value) {
    loadPlugins();
    checkInstalledPlugins();
  }
});
</script>

<style scoped>
.online-script-list {
  min-height: 72vh;
}

.filter-section {
  display: flex;
  gap: 12px;
  margin-bottom: 2px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.category-select {
  width: 180px;
}

.tag-select {
  width: 200px;
}

.stats-section {
  margin-bottom: 10px;
  padding: 8px 0;
  border-bottom: 1px solid #ebeef5;
}

.stats-text {
  font-size: 14px;
  color: #909399;
}

.plugins-container {
  height: 62vh;
  overflow-y: auto;
}

.empty-state {
  padding: 40px 0;
}

.plugins-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
}

.plugin-card {
  transition: all 0.3s;
}

.plugin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.plugin-title-section {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  flex: 1;
}

.plugin-version-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.plugin-name {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.plugin-version-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.plugin-version {
  flex-shrink: 0;
}

.installed-version {
  font-size: 12px;
  color: #909399;
}

.plugin-description {
  margin: 0 0 6px 0;
  line-height: 1.6;
  font-size: 14px;
}

.plugin-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
  font-size: 13px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.plugin-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag-item {
  cursor: pointer;
  transition: all 0.2s;
}

.tag-item:hover {
  transform: scale(1.05);
}

.plugin-content-types {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}

.content-types-label {
  font-size: 12px;
  margin-right: 4px;
}

.content-type-tag {
  font-size: 11px;
}

@media (max-width: 768px) {
  .plugins-grid {
    grid-template-columns: 1fr;
  }

  .filter-section {
    flex-direction: column;
  }

  .search-input,
  .category-select,
  .tag-select {
    width: 100%;
  }
}
</style>
