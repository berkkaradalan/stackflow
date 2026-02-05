import { api } from "../api-client";

export interface ModelConfig {
  id: string;
  name: string;
  description: string;
  max_tokens: number;
  supports_streaming: boolean;
  supports_vision: boolean;
  input_price_per_m_token: number;
  output_price_per_m_token: number;
}

export interface ProviderConfig {
  name: string;
  display_name: string;
  base_url: string;
  health_check_path: string;
  chat_completion_path: string;
  requires_api_key: boolean;
  models: ModelConfig[];
}

export interface ProviderListResponse {
  providers: ProviderConfig[];
}

export interface ModelListResponse {
  provider: string;
  models: ModelConfig[];
}

export const providersApi = {
  /**
   * Get all providers
   */
  getAll: async (): Promise<ProviderListResponse> => {
    return api.get<ProviderListResponse>("/api/providers");
  },

  /**
   * Get a specific provider by name
   */
  getByName: async (name: string): Promise<ProviderConfig> => {
    return api.get<ProviderConfig>(`/api/providers/${name}`);
  },

  /**
   * Get models for a specific provider
   */
  getModels: async (providerName: string): Promise<ModelListResponse> => {
    return api.get<ModelListResponse>(`/api/providers/${providerName}/models`);
  },
};
