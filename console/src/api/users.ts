import request from "./request";

export type UserStatus = "init" | "invited" | "active" | "disabled";

export interface ManagedUser {
  id: number | string;
  uuid?: string;
  tenant_id?: number | string;
  role_id?: number | string;
  name: string;
  email: string;
  user_type: string;
  status: UserStatus | string;
  created_at?: string;
}

export interface UserListResult {
  users: ManagedUser[];
  total: number;
}

type UnknownRecord = Record<string, unknown>;

function isRecord(value: unknown): value is UnknownRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/**
 * 同时兼容 Go 默认字段名（ID、Email、UserType）
 * 和带 JSON tag 的字段名（id、email、user_type）。
 */
function getField(record: UnknownRecord, ...keys: string[]): unknown {
  for (const key of keys) {
    if (record[key] !== undefined && record[key] !== null) {
      return record[key];
    }
  }

  return undefined;
}

function readString(value: unknown): string {
  if (typeof value === "string" || typeof value === "number") {
    return String(value);
  }

  return "";
}

// function readOptionalId(value: unknown): number | string | undefined {
//   if (typeof value === "number" || typeof value === "string") {
//     return value;
//   }

//   return undefined;
// }

function extractUsers(payload: unknown): {
  users: unknown[];
  total?: number;
} {
  if (Array.isArray(payload)) {
    return {
      users: payload,
      total: undefined,
    };
  }

  if (!isRecord(payload)) {
    return { users: [] };
  }

  for (const key of ["users", "items", "list", "records"]) {
    const value = payload[key];

    if (Array.isArray(value)) {
      return {
        users: value,
        total: typeof payload.total === "number" ? payload.total : undefined,
      };
    }
  }

  return { users: [] };
}

function normalizeUser(value: unknown): ManagedUser | null {
  if (!isRecord(value)) {
    return null;
  }

  // 兼容 Go 默认 JSON 字段名与 snake_case 字段名。
  const get = (lower: string, upper: string): unknown =>
    value[lower] ?? value[upper];

  const email = readString(get("email", "Email"));

  if (!email) {
    return null;
  }

  const uuid = readString(get("uuid", "UUID"));

  return {
    id: readString(get("id", "ID")) || uuid || email,
    uuid: uuid || undefined,
    tenant_id: get("tenant_id", "TenantID") as number | string | undefined,
    role_id: get("role_id", "RoleID") as number | string | undefined,
    name: readString(get("name", "Name")) || "未设置姓名",
    email,
    user_type: readString(get("user_type", "UserType")) || "employee",
    status: readString(get("status", "Status")) || "unknown",
    created_at: readString(get("created_at", "CreatedAt")) || undefined,
  };
}

export async function listUsers(
  options: {
    status?: string;
    keyword?: string;
    limit?: number;
    offset?: number;
  } = {},
): Promise<UserListResult> {
  const params: Record<string, string | number> = {
    limit: options.limit ?? 100,
    offset: options.offset ?? 0,
  };

  if (options.status) {
    params.status = options.status;
  }

  if (options.keyword?.trim()) {
    params.keyword = options.keyword.trim();
  }

  const payload: unknown = await request.get("/users/lists", {
    params,
  });

  const result = extractUsers(payload);

  const users = result.users
    .map(normalizeUser)
    .filter((user): user is ManagedUser => user !== null);

  return {
    users,
    total: result.total ?? users.length,
  };
}

export async function inviteUser(data: {
  name: string;
  email: string;
  role_id: number;
}): Promise<void> {
  await request.post("/identity/invitations", data);
}

export interface IdentityRole {
  id: number;
  uuid?: string;
  name: string;
  description?: string;
}

export async function listRoles(): Promise<IdentityRole[]> {
  const payload: unknown = await request.get("/identity/roles");

  const values = Array.isArray(payload)
    ? payload
    : isRecord(payload) && Array.isArray(payload.roles)
      ? payload.roles
      : [];

  return values
    .filter(isRecord)
    .map((role) => ({
      id: Number(getField(role, "id", "ID")),
      uuid: readString(getField(role, "uuid", "UUID")) || undefined,
      name:
        readString(getField(role, "name", "Name")) ||
        readString(getField(role, "code", "Code")) ||
        "未命名角色",
      description:
        readString(getField(role, "description", "Description")) || undefined,
    }))
    .filter((role) => Number.isFinite(role.id) && role.id > 0);
}
