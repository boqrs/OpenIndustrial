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

function readString(value: unknown): string {
  if (typeof value === "string" || typeof value === "number") {
    return String(value);
  }

  return "";
}

function extractUsers(payload: unknown): {
  users: unknown[];
  total?: number;
} {
  if (Array.isArray(payload)) {
    return { users: payload };
  }

  if (!isRecord(payload)) {
    return { users: [] };
  }

  for (const key of ["users", "items", "list", "records"]) {
    if (Array.isArray(payload[key])) {
      return {
        users: payload[key] as unknown[],
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

  const email = readString(value.email);

  if (!email) {
    return null;
  }

  const uuid = readString(value.uuid);

  return {
    id: readString(value.id) || uuid || email,
    uuid: uuid || undefined,
    tenant_id:
      typeof value.tenant_id === "number" || typeof value.tenant_id === "string"
        ? value.tenant_id
        : undefined,
    role_id:
      typeof value.role_id === "number" || typeof value.role_id === "string"
        ? value.role_id
        : undefined,
    name: readString(value.name) || "未设置姓名",
    email,
    user_type: readString(value.user_type) || "employee",
    status: readString(value.status) || "unknown",
    created_at: readString(value.created_at) || undefined,
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
      id: Number(role.id),
      uuid: readString(role.uuid) || undefined,
      name: readString(role.name) || readString(role.code) || "未命名角色",
      description: readString(role.description) || undefined,
    }))
    .filter((role) => Number.isFinite(role.id) && role.id > 0);
}
