<script setup lang="ts">
import { onMounted, reactive, ref, watch } from "vue";
import { ElMessage } from "element-plus";

import { getAttributeValues, updateAttributeValues } from "../api";
import type { AttributeDefinition, AttributeValues } from "../types";

const props = defineProps<{
  productId: number;
  definitions: AttributeDefinition[];
}>();

const loading = ref(false);
const saving = ref(false);
const errorMessage = ref("");

const values = reactive<Record<string, unknown>>({});
const jsonTexts = reactive<Record<string, string>>({});

function isJson(value: string): boolean {
  try {
    JSON.parse(value);
    return true;
  } catch {
    return false;
  }
}

function displayValue(value: unknown): string {
  if (value === undefined || value === null) return "";
  if (typeof value === "object") return JSON.stringify(value, null, 2);
  return String(value);
}

function initializeValues(data: AttributeValues) {
  for (const key of Object.keys(values)) delete values[key];
  for (const key of Object.keys(jsonTexts)) delete jsonTexts[key];

  for (const definition of props.definitions) {
    if (Object.prototype.hasOwnProperty.call(data, definition.name)) {
      const value = data[definition.name];
      values[definition.name] = value;

      if (definition.data_type === "json") {
        jsonTexts[definition.name] = displayValue(value);
      }
    } else {
      values[definition.name] =
        definition.data_type === "boolean" ? false : undefined;

      if (definition.data_type === "json") {
        jsonTexts[definition.name] = "";
      }
    }
  }
}

async function loadValues() {
  if (!Number.isInteger(props.productId) || props.productId <= 0) return;

  loading.value = true;
  errorMessage.value = "";

  try {
    const data = await getAttributeValues(props.productId);
    initializeValues(data);
  } catch (error: unknown) {
    errorMessage.value =
      error instanceof Error ? error.message : "属性值加载失败";
  } finally {
    loading.value = false;
  }
}

async function saveValues() {
  if (saving.value) return;

  const payload: AttributeValues = {};

  for (const definition of props.definitions) {
    const name = definition.name;
    let value = values[name];

    if (definition.data_type === "json") {
      const raw = (jsonTexts[name] ?? "").trim();

      if (!raw) {
        value = undefined;
      } else {
        if (!isJson(raw)) {
          ElMessage.warning(`属性“${definition.label || name}”不是有效 JSON`);
          return;
        }

        value = JSON.parse(raw);
      }
    }

    if (value === undefined || value === "") {
      if (definition.required) {
        ElMessage.warning(`请填写必填属性：${definition.label || name}`);
        return;
      }

      continue;
    }

    if (
      definition.data_type === "integer" &&
      (typeof value !== "number" || !Number.isInteger(value))
    ) {
      ElMessage.warning(`属性“${definition.label || name}”必须是整数`);
      return;
    }

    if (
      definition.data_type === "float" &&
      (typeof value !== "number" || !Number.isFinite(value))
    ) {
      ElMessage.warning(`属性“${definition.label || name}”必须是有效数字`);
      return;
    }

    payload[name] = value;
  }

  saving.value = true;

  try {
    await updateAttributeValues(props.productId, { values: payload });
    ElMessage.success("属性值保存成功");
    await loadValues();
  } catch (error: unknown) {
    ElMessage.error(error instanceof Error ? error.message : "属性值保存失败");
  } finally {
    saving.value = false;
  }
}

watch(
  () => [props.productId, props.definitions] as const,
  () => void loadValues(),
  { deep: true },
);

onMounted(() => void loadValues());
</script>

<template>
  <article class="content-panel values-panel" v-loading="loading">
    <div class="panel-heading">
      <div>
        <h2>属性值</h2>
        <p>维护当前产品型号的属性数据，保存后写入 Cloud 资源属性。</p>
      </div>

      <el-button
        type="primary"
        :loading="saving"
        :disabled="loading || definitions.length === 0"
        @click="saveValues"
      >
        保存属性值
      </el-button>
    </div>

    <el-alert
      v-if="errorMessage"
      :title="errorMessage"
      type="error"
      show-icon
      :closable="false"
      class="values-error"
    />

    <div v-if="definitions.length === 0" class="values-empty">
      当前产品型号没有属性定义，请先完成属性结构配置。
    </div>

    <el-form v-else label-position="top" class="values-form">
      <div
        v-for="definition in definitions"
        :key="definition.id || definition.name"
        class="value-field"
      >
        <el-form-item>
          <template #label>
            <span>{{ definition.label || definition.name }}</span>
            <el-tag
              v-if="definition.required"
              size="small"
              type="warning"
              class="required-tag"
            >
              必填
            </el-tag>
            <span v-if="definition.unit" class="unit-label">
              （{{ definition.unit }}）
            </span>
          </template>

          <el-input
            v-if="definition.data_type === 'string'"
            v-model="values[definition.name] as string"
            :placeholder="`请输入${definition.label || definition.name}`"
            clearable
          />

          <el-input
            v-else-if="definition.data_type === 'text'"
            v-model="values[definition.name] as string"
            type="textarea"
            :rows="3"
          />

          <el-input-number
            v-else-if="definition.data_type === 'integer'"
            v-model="values[definition.name] as number"
            :precision="0"
            controls-position="right"
            class="full-width"
          />

          <el-input-number
            v-else-if="definition.data_type === 'float'"
            v-model="values[definition.name] as number"
            :precision="6"
            controls-position="right"
            class="full-width"
          />

          <el-switch
            v-else-if="definition.data_type === 'boolean'"
            v-model="values[definition.name] as boolean"
            active-text="是"
            inactive-text="否"
          />

          <el-date-picker
            v-else-if="definition.data_type === 'datetime'"
            v-model="values[definition.name] as string"
            type="datetime"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            placeholder="选择日期时间"
            class="full-width"
          />

          <el-input
            v-else-if="definition.data_type === 'json'"
            v-model="jsonTexts[definition.name]"
            type="textarea"
            :rows="5"
            placeholder='请输入 JSON，例如 {"key":"value"}'
          />

          <el-input
            v-else
            :model-value="displayValue(values[definition.name])"
            disabled
          />

          <div class="field-meta">
            <code>{{ definition.name }}</code>
            <span>{{ definition.data_type }}</span>
          </div>
        </el-form-item>
      </div>
    </el-form>
  </article>
</template>

<style scoped>
.values-panel {
  min-width: 0;
  margin-top: 22px;
  overflow: hidden;
  border: 1px solid #eaecf0;
  border-radius: 10px;
  background: #fff;
}

.panel-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 21px 22px;
  border-bottom: 1px solid #f0f1f3;
}

.panel-heading h2 {
  margin: 0;
  color: #0e1f33;
  font-size: 15px;
  font-weight: 750;
}

.panel-heading p {
  margin: 7px 0 0;
  color: #667085;
  font-size: 11px;
  line-height: 1.7;
}

.values-form {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 24px;
  padding: 22px;
}

.value-field {
  min-width: 0;
}

.required-tag {
  margin-left: 8px;
}

.unit-label {
  margin-left: 4px;
  color: #667085;
  font-size: 12px;
}

.field-meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-top: 7px;
  color: #98a2b3;
  font-size: 11px;
}

.field-meta code {
  color: #8a5b16;
}

.full-width {
  width: 100%;
}

.values-error {
  margin: 16px 22px 0;
  width: auto;
}

.values-empty {
  margin: 22px;
  padding: 28px 12px;
  border: 1px dashed #d0d5dd;
  border-radius: 8px;
  color: #98a2b3;
  font-size: 12px;
  text-align: center;
}

@media (max-width: 680px) {
  .panel-heading {
    align-items: flex-start;
    flex-direction: column;
  }

  .values-form {
    grid-template-columns: 1fr;
    padding: 18px;
  }
}
</style>
