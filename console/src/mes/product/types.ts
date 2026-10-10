export type ProductStatus = "pending" | "active" | "inactive" | "archived";

export type AttributeDataType =
  | "string"
  | "text"
  | "integer"
  | "float"
  | "boolean"
  | "datetime"
  | "json";

export interface AttributeDefinitionInput {
  label?: string;
  description?: string;
  data_type: AttributeDataType;
  unit?: string;
  required?: boolean;
}

/** 列表接口返回的属性信息 */
export interface ProductAttribute {
  id: number;
  name: string;
  label?: string;
  description?: string;
  data_type: string;
  unit?: string;
  required?: boolean;
}

/** 产品型号列表项 */
export interface ProductModel {
  id: number;
  resource_id: number;
  name: string;
  code: string;
  version: string;
  category: string;
  description?: string;
  status: ProductStatus | string;
  attributes?: ProductAttribute[];
  created_at: string;
  updated_at: string;
}

/** 创建产品型号请求 */
export interface CreateProductModelRequest {
  name: string;
  code: string;
  version: string;
  category: string;
  description?: string;
  attributes?: Record<string, AttributeDefinitionInput>;
}

/** 更新产品型号请求 */
export interface UpdateProductModelRequest {
  name?: string;
  category?: string;
  description?: string;
}

/** 属性定义 */
export interface AttributeDefinition extends ProductAttribute {}

/** 产品型号详情 */
export interface ProductDetail extends ProductModel {
  attribute?: unknown[];
  attribute_definition?: AttributeDefinition[];
}

/** 列表查询条件 */
export interface ProductListQuery {
  category?: string;
  status?: string;
  code?: string;
  current_page?: number;
  page_size?: number;
}

/** 列表响应 */
export interface ProductListResponse {
  items: ProductModel[];
  total: number;
  next: boolean;
}
