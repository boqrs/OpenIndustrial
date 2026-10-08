import request from "./request";

export interface User {
  id: string;
  tenant_id: string;
  uuid?: string;
  email: string;
  name?: string;
  user_type?: string;
  status?: string;
}

export interface LoginRequest {
  tenant_code: string;
  email: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface RequestAccessRequest {
  tenant_code: string;
  email: string;
  name: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}

export interface AcceptInvitationRequest {
  token: string;
  password: string;
}

export async function login(data: LoginRequest): Promise<LoginResponse> {
  return request.post("/login", data);
}

export async function requestAccess(data: RequestAccessRequest): Promise<void> {
  await request.post("/identity/access-requests", data);
}

export async function refreshToken(
  data: RefreshTokenRequest,
): Promise<RefreshTokenResponse> {
  return request.post("/refresh", data);
}

export async function logout(): Promise<void> {
  await request.post("/logout", {});
}

export async function acceptInvitation(
  data: AcceptInvitationRequest,
): Promise<void> {
  await request.post("/identity/invitations/accept", data);
}
