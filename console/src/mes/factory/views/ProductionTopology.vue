<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Plus,
  Refresh,
  OfficeBuilding,
  Operation,
  Monitor,
  Connection,
} from "@element-plus/icons-vue";

import {
  createTopologyNode,
  deleteTopologyNode,
  getFactoryTopology,
  listFactories,
  updateTopologyNode,
} from "../api";

import type {
  CreateTopologyNodeRequest,
  Factory,
  FactoryTopology,
  TopologyNode,
  TopologyNodeType,
} from "../types";

const route = useRoute();
const router = useRouter();

const loadingFactories = ref(false);
const loadingTopology = ref(false);
const saving = ref(false);

const factories = ref<Factory[]>([]);
const selectedFactoryId = ref<number | undefined>();
const topology = ref<FactoryTopology | null>(null);

const dialogVisible = ref(false);
const formRef = ref();

const editingNode = ref<TopologyNode | null>(null);

interface NodeForm {
  type: TopologyNodeType;
  name: string;
  code: string;
  parent_resource_id?: number;
}

const emptyForm = (): NodeForm => ({
  type: "production_line",
  name: "",
  code: "",
});

const form = ref<NodeForm>(emptyForm());

const rules = {
  name: [
    {
      required: true,
      message: "请输入节点名称",
      trigger: "blur",
    },
  ],
};

/**
 * 统一前后端节点类型。
 *
 * 后端 Resource 内部类型通常为：
 *   PRODUCTION_LINE
 *   WORK_STATION
 *
 * 前端使用：
 *   production_line
 *   work_station
 *
 * 所有拓扑展示和操作统一使用归一化后的类型，
 * 避免统计错误、树节点错误以及按钮不显示。
 */
function normalizeNodeType(type: string): string {
  switch (type.trim().toUpperCase()) {
    case "PRODUCTION_LINE":
      return "production_line";
    case "WORK_STATION":
      return "work_station";
    default:
      return type;
  }
}

const selectedFactory = computed(
  () =>
    factories.value.find((item) => item.id === selectedFactoryId.value) ?? null,
);

/**
 * 统一处理后端返回的节点类型。
 * 后续统计、树结构、按钮判断均使用该数据。
 */
const flatNodes = computed<TopologyNode[]>(() =>
  (topology.value?.nodes ?? []).map((node) => ({
    ...node,
    type: normalizeNodeType(node.type),
  })),
);

const productionLines = computed(() =>
  flatNodes.value.filter((node) => node.type === "production_line"),
);

const workStations = computed(() =>
  flatNodes.value.filter((node) => node.type === "work_station"),
);

/**
 * 根据后端返回的扁平节点构建树形结构。
 *
 * parent_resource_id 指向父节点的 resource_id。
 * 没有找到父节点的节点作为根节点展示。
 */
const treeNodes = computed<TopologyNode[]>(() => {
  const nodes = flatNodes.value.map((node) => ({
    ...node,
    children: [] as TopologyNode[],
  }));

  const byId = new Map<number, TopologyNode>();

  nodes.forEach((node) => {
    byId.set(node.resource_id, node);
  });

  const roots: TopologyNode[] = [];

  for (const node of nodes) {
    const parentId = node.parent_resource_id;
    const parent = parentId ? byId.get(parentId) : undefined;

    if (parent) {
      parent.children = parent.children ?? [];
      parent.children.push(node);
    } else {
      roots.push(node);
    }
  }

  return roots;
});

/**
 * 创建工作站时的生产线选项。
 */
const parentLineOptions = computed(() =>
  productionLines.value.map((node) => ({
    label: node.name,
    value: node.resource_id,
  })),
);

/**
 * 加载当前租户的工厂列表。
 *
 * 优先选择 URL 中指定的工厂，
 * 其次保留当前选中的工厂，
 * 最后默认选择列表中的第一家工厂。
 */
async function loadFactories() {
  loadingFactories.value = true;

  try {
    factories.value = await listFactories();

    const queryId = Number(route.query.factory_id);

    const requestedFactory = factories.value.find(
      (factory) => factory.id === queryId,
    );

    const currentFactory = factories.value.find(
      (factory) => factory.id === selectedFactoryId.value,
    );

    selectedFactoryId.value =
      requestedFactory?.id ?? currentFactory?.id ?? factories.value[0]?.id;

    if (selectedFactoryId.value) {
      await loadTopology();
    } else {
      topology.value = null;
    }
  } catch (error) {
    ElMessage.error(
      error instanceof Error ? error.message : "获取工厂列表失败",
    );
  } finally {
    loadingFactories.value = false;
  }
}

/**
 * 加载当前工厂的生产拓扑。
 */
async function loadTopology() {
  if (!selectedFactoryId.value) {
    topology.value = null;
    return;
  }

  loadingTopology.value = true;

  try {
    topology.value = await getFactoryTopology(selectedFactoryId.value);
  } catch (error) {
    topology.value = null;

    ElMessage.error(
      error instanceof Error ? error.message : "获取生产拓扑失败",
    );
  } finally {
    loadingTopology.value = false;
  }
}

/**
 * 切换工厂。
 */
function onFactoryChange(value: number | undefined) {
  selectedFactoryId.value = value;

  if (value) {
    router.replace({
      name: "ProductionTopology",
      query: {
        factory_id: String(value),
      },
    });
  }

  void loadTopology();
}

/**
 * 打开新建生产线对话框。
 *
 * 生产线直接归属于当前工厂，
 * 不需要在请求中传入 parent_resource_id。
 */
function openCreateLine() {
  editingNode.value = null;

  form.value = {
    ...emptyForm(),
    type: "production_line",
  };

  dialogVisible.value = true;
}

/**
 * 打开新建工作站对话框。
 *
 * 工作站必须归属于一条生产线。
 */
function openCreateWorkStation(parent: TopologyNode) {
  if (normalizeNodeType(parent.type) !== "production_line") {
    ElMessage.warning("只能在生产线下创建工作站");
    return;
  }

  editingNode.value = null;

  form.value = {
    type: "work_station",
    name: "",
    code: "",
    parent_resource_id: parent.resource_id,
  };

  dialogVisible.value = true;
}

/**
 * 打开编辑节点对话框。
 *
 * 编辑时只允许修改节点名称；
 * 节点类型和父节点关系不在此处修改。
 */
function openEditNode(node: TopologyNode) {
  editingNode.value = node;

  form.value = {
    type: normalizeNodeType(node.type) as TopologyNodeType,
    name: node.name,
    code: "",
    parent_resource_id: node.parent_resource_id,
  };

  dialogVisible.value = true;
}

/**
 * 创建或更新拓扑节点。
 */
async function submitNode() {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return;

    if (!selectedFactoryId.value) {
      ElMessage.warning("请先选择工厂");
      return;
    }

    saving.value = true;

    try {
      if (editingNode.value) {
        await updateTopologyNode(editingNode.value.resource_id, {
          name: form.value.name.trim(),
        });

        ElMessage.success("节点已更新");
      } else {
        const payload: CreateTopologyNodeRequest = {
          factory_id: selectedFactoryId.value,
          type: form.value.type,
          name: form.value.name.trim(),
          code: form.value.code.trim() || undefined,
        };

        if (form.value.type === "work_station") {
          if (!form.value.parent_resource_id) {
            ElMessage.warning("请先选择所属生产线");
            return;
          }

          payload.parent_resource_id = form.value.parent_resource_id;
        }

        await createTopologyNode(payload);

        ElMessage.success(
          form.value.type === "production_line"
            ? "生产线创建成功"
            : "工作站创建成功",
        );
      }

      dialogVisible.value = false;
      await loadTopology();
    } catch (error) {
      ElMessage.error(error instanceof Error ? error.message : "保存节点失败");
    } finally {
      saving.value = false;
    }
  });
}

/**
 * 删除生产线或工作站。
 *
 * 如果生产线还有工作站，后端应拒绝删除。
 */
async function removeNode(node: TopologyNode) {
  const normalizedType = normalizeNodeType(node.type);

  const label = normalizedType === "production_line" ? "生产线" : "工作站";

  try {
    await ElMessageBox.confirm(
      `确定删除${label}「${node.name}」吗？存在子节点时，后端会拒绝删除。`,
      "删除节点",
      {
        type: "warning",
        confirmButtonText: "确认删除",
        cancelButtonText: "取消",
      },
    );

    await deleteTopologyNode(node.resource_id);

    ElMessage.success("节点已删除");

    await loadTopology();
  } catch (error) {
    if (error === "cancel" || error === "close") {
      return;
    }

    ElMessage.error(error instanceof Error ? error.message : "删除节点失败");
  }
}

/**
 * 节点类型显示名称。
 */
function nodeTypeLabel(type: string) {
  const normalizedType = normalizeNodeType(type);

  if (normalizedType === "production_line") {
    return "生产线";
  }

  if (normalizedType === "work_station") {
    return "工作站";
  }

  return type;
}

/**
 * 节点状态显示名称。
 */
function nodeStatusLabel(status: string) {
  if (status === "active") return "启用";

  if (status === "inactive" || status === "disabled") {
    return "停用";
  }

  return status || "未知";
}

/**
 * 判断节点是否有子节点。
 */
function hasChildren(node: TopologyNode) {
  return Boolean(node.children?.length);
}

/**
 * 监听 URL 中的工厂 ID。
 *
 * 支持从其他页面跳转到指定工厂的生产拓扑。
 */
watch(
  () => route.query.factory_id,
  (value) => {
    const id = Number(value);

    if (
      id > 0 &&
      factories.value.some((factory) => factory.id === id) &&
      selectedFactoryId.value !== id
    ) {
      selectedFactoryId.value = id;
      void loadTopology();
    }
  },
);

onMounted(loadFactories);
</script>

<template>
  <div class="topology-page">
    <section class="page-hero">
      <div class="hero-icon">
        <el-icon :size="26">
          <Connection />
        </el-icon>
      </div>

      <div class="hero-copy">
        <div class="eyebrow">PRODUCTION TOPOLOGY</div>
        <h1>生产拓扑</h1>
        <p>维护工厂、生产线和工作站之间的组织关系。</p>
      </div>

      <div class="hero-actions">
        <el-button :icon="ArrowLeft" @click="router.push('/mes/factories')">
          返回工厂列表
        </el-button>

        <el-button
          :icon="Refresh"
          :loading="loadingTopology"
          @click="loadTopology"
        >
          刷新拓扑
        </el-button>
      </div>
    </section>

    <section class="selector-card">
      <div class="selector-copy">
        <span class="selector-label">当前工厂</span>

        <strong>
          {{ selectedFactory?.name ?? "请选择工厂" }}
        </strong>

        <span class="selector-code">
          {{ selectedFactory?.code ?? "尚未选择" }}
        </span>
      </div>

      <el-select
        :model-value="selectedFactoryId"
        :loading="loadingFactories"
        placeholder="选择工厂"
        class="factory-select"
        @update:model-value="onFactoryChange"
      >
        <el-option
          v-for="factory in factories"
          :key="factory.id"
          :label="`${factory.name}（${factory.code}）`"
          :value="factory.id"
        />
      </el-select>
    </section>

    <section class="stat-grid">
      <div class="stat-card">
        <span class="stat-icon stat-icon--factory">
          <el-icon>
            <OfficeBuilding />
          </el-icon>
        </span>

        <div>
          <span class="stat-label">生产工厂</span>
          <strong>{{ selectedFactory ? 1 : 0 }}</strong>
        </div>
      </div>

      <div class="stat-card">
        <span class="stat-icon stat-icon--line">
          <el-icon>
            <Operation />
          </el-icon>
        </span>

        <div>
          <span class="stat-label">生产线</span>
          <strong>{{ productionLines.length }}</strong>
        </div>
      </div>

      <div class="stat-card">
        <span class="stat-icon stat-icon--station">
          <el-icon>
            <Monitor />
          </el-icon>
        </span>

        <div>
          <span class="stat-label">工作站</span>
          <strong>{{ workStations.length }}</strong>
        </div>
      </div>
    </section>

    <section class="topology-card">
      <div class="card-heading">
        <div>
          <h2>组织结构</h2>
          <p>工厂 → 生产线 → 工作站</p>
        </div>

        <el-button
          type="primary"
          :icon="Plus"
          :disabled="!selectedFactory"
          @click="openCreateLine"
        >
          新建生产线
        </el-button>
      </div>

      <el-empty
        v-if="!selectedFactory && !loadingTopology"
        description="请先选择工厂"
      />

      <el-table
        v-else
        v-loading="loadingTopology"
        :data="treeNodes"
        row-key="resource_id"
        default-expand-all
        :tree-props="{ children: 'children' }"
        empty-text="该工厂暂无生产线，请新建生产线"
      >
        <el-table-column label="节点名称" min-width="260">
          <template #default="{ row }">
            <div class="node-name">
              <span
                class="node-icon"
                :class="
                  normalizeNodeType(row.type) === 'production_line'
                    ? 'node-icon--line'
                    : 'node-icon--station'
                "
              >
                <el-icon>
                  <Operation
                    v-if="normalizeNodeType(row.type) === 'production_line'"
                  />
                  <Monitor v-else />
                </el-icon>
              </span>

              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="节点类型" width="130">
          <template #default="{ row }">
            <el-tag
              :type="
                normalizeNodeType(row.type) === 'production_line'
                  ? 'primary'
                  : 'success'
              "
              effect="light"
            >
              {{ nodeTypeLabel(row.type) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="资源 ID" min-width="120">
          <template #default="{ row }">
            <span class="resource-id">
              {{ row.resource_id }}
            </span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'active' ? 'success' : 'info'"
              effect="light"
            >
              {{ nodeStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" fixed="right" width="250">
          <template #default="{ row }">
            <el-button
              v-if="normalizeNodeType(row.type) === 'production_line'"
              link
              type="primary"
              :icon="Plus"
              @click="openCreateWorkStation(row)"
            >
              添加工作站
            </el-button>

            <el-button link type="primary" @click="openEditNode(row)">
              编辑
            </el-button>

            <el-button
              link
              type="danger"
              :disabled="hasChildren(row)"
              @click="removeNode(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-note">
        <span class="note-dot" />
        节点关系由后端校验；有子节点的生产线需要先处理其工作站，再删除生产线。
      </div>
    </section>

    <el-dialog
      v-model="dialogVisible"
      :title="
        editingNode
          ? '编辑' + nodeTypeLabel(editingNode.type)
          : form.type === 'production_line'
            ? '新建生产线'
            : '新建工作站'
      "
      width="520px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="节点类型">
          <el-tag v-if="form.type === 'production_line'" type="primary">
            生产线
          </el-tag>

          <el-tag v-else type="success"> 工作站 </el-tag>
        </el-form-item>

        <el-form-item
          v-if="form.type === 'work_station' && !editingNode"
          label="所属生产线"
          prop="parent_resource_id"
          :rules="[
            {
              required: true,
              message: '请选择所属生产线',
              trigger: 'change',
            },
          ]"
        >
          <el-select
            v-model="form.parent_resource_id"
            placeholder="选择生产线"
            class="full-width"
          >
            <el-option
              v-for="line in parentLineOptions"
              :key="line.value"
              :label="line.label"
              :value="line.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="节点名称" prop="name">
          <el-input
            v-model="form.name"
            maxlength="100"
            :placeholder="
              form.type === 'production_line'
                ? '例如：总装一线'
                : '例如：装配工作站 01'
            "
          />
        </el-form-item>

        <el-form-item v-if="!editingNode" label="节点编码">
          <el-input
            v-model="form.code"
            maxlength="64"
            placeholder="可选，例如：LINE-01"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false"> 取消 </el-button>

        <el-button type="primary" :loading="saving" @click="submitNode">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.topology-page {
  min-height: 100%;
  padding: 26px 30px;
  background: #f5f7fa;
  color: #172b42;
}

.page-hero {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 26px 30px;
  border-radius: 14px;
  color: #fff;
  background: linear-gradient(115deg, #0e1f33, #203a55);
}

.hero-icon {
  display: grid;
  width: 56px;
  height: 56px;
  flex: 0 0 auto;
  place-items: center;
  border: 1px solid rgb(240 160 48 / 35%);
  border-radius: 14px;
  color: #f0a030;
  background: rgb(255 255 255 / 5%);
}

.hero-copy {
  flex: 1;
  min-width: 0;
}

.eyebrow {
  color: #a9bdcf;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 1.8px;
}

.hero-copy h1 {
  margin: 9px 0;
  font-size: 25px;
}

.hero-copy p {
  margin: 0;
  color: #c4d0dd;
  font-size: 12px;
  line-height: 1.7;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.selector-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-top: 22px;
  padding: 20px 24px;
  border: 1px solid #e6ebf1;
  border-radius: 12px;
  background: #fff;
}

.selector-copy {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.selector-label {
  color: #8290a0;
  font-size: 11px;
}

.selector-copy strong {
  color: #20344b;
  font-size: 17px;
}

.selector-code {
  color: #71869c;
  font-size: 11px;
}

.factory-select {
  width: min(360px, 50%);
}

.stat-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
  margin-top: 18px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  border: 1px solid #e6ebf1;
  border-radius: 12px;
  background: #fff;
}

.stat-icon {
  display: grid;
  width: 44px;
  height: 44px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 11px;
  font-size: 20px;
}

.stat-icon--factory {
  color: #547ba1;
  background: #edf3f9;
}

.stat-icon--line {
  color: #7861bd;
  background: #f0edfb;
}

.stat-icon--station {
  color: #329572;
  background: #e8f7f0;
}

.stat-card > div {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.stat-label {
  color: #8290a0;
  font-size: 11px;
}

.stat-card strong {
  color: #20344b;
  font-size: 25px;
  line-height: 1;
}

.topology-card {
  margin-top: 20px;
  padding: 22px;
  border: 1px solid #e6ebf1;
  border-radius: 12px;
  background: #fff;
}

.card-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.card-heading h2 {
  margin: 0;
  font-size: 17px;
}

.card-heading p {
  margin: 7px 0 0;
  color: #8290a0;
  font-size: 12px;
}

.node-name {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
}

.node-icon {
  display: grid;
  width: 30px;
  height: 30px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 8px;
}

.node-icon--line {
  color: #7861bd;
  background: #f0edfb;
}

.node-icon--station {
  color: #329572;
  background: #e8f7f0;
}

.resource-id {
  color: #71869c;
  font-family: monospace;
  font-size: 12px;
}

.table-note {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-top: 17px;
  color: #8290a0;
  font-size: 11px;
  line-height: 1.7;
}

.note-dot {
  width: 6px;
  height: 6px;
  flex: 0 0 auto;
  margin-top: 5px;
  border-radius: 50%;
  background: #f0a030;
}

.full-width {
  width: 100%;
}

@media (max-width: 900px) {
  .topology-page {
    padding: 18px;
  }

  .page-hero {
    align-items: flex-start;
    flex-wrap: wrap;
    padding: 22px;
  }

  .hero-actions {
    width: 100%;
  }

  .stat-grid {
    grid-template-columns: 1fr;
  }

  .selector-card {
    align-items: flex-start;
    flex-direction: column;
  }

  .factory-select {
    width: 100%;
  }
}
</style>
