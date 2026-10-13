export interface Factory {
  id: number;
  resource_id: number;
  name: string;
  code: string;
  address?: string;
  timezone: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export type CreateFactoryRequest = {
  name: string;
  code: string;
  address?: string;
  timezone?: string;
};

export type UpdateFactoryRequest = {
  name?: string;
  code?: string;
  address?: string;
  timezone?: string;
};

export type TopologyNodeType = "production_line" | "work_station";

export interface TopologyNode {
  resource_id: number;
  type: TopologyNodeType | string;
  name: string;
  status: string;
  parent_resource_id?: number;
  metadata?: Record<string, unknown>;
  children?: TopologyNode[];
}

export interface FactoryTopology {
  factory: Factory;
  nodes: TopologyNode[];
}

export interface CreateTopologyNodeRequest {
  factory_id: number;
  parent_resource_id?: number;
  type: TopologyNodeType;
  name: string;
  code?: string;
  metadata?: Record<string, unknown>;
}

export interface UpdateTopologyNodeRequest {
  name?: string;
  code?: string;
  metadata?: Record<string, unknown>;
}
