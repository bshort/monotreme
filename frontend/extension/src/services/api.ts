import { Storage } from "@plasmohq/storage";

const storage = new Storage();

export interface CreateShortcutRequest {
  name: string;
  link: string;
  title?: string;
  description?: string;
  tags?: string[];
  visibility?: 'PUBLIC' | 'WORKSPACE' | 'PRIVATE';
}

export interface CreateShortcutResponse {
  id: number;
  name: string;
  link: string;
  title?: string;
  description?: string;
  tags?: string[];
  visibility?: string;
}

export class ApiService {
  private async getInstanceUrl(): Promise<string> {
    const instanceUrl = await storage.getItem<string>("instance_url");
    if (!instanceUrl) {
      throw new Error("Instance URL not configured. Please set it in the extension options.");
    }
    return instanceUrl;
  }

  private async getAuthHeaders(): Promise<HeadersInit> {
    // For now, we'll assume the user is authenticated via cookies
    // In a more complete implementation, you might store auth tokens
    return {
      'Content-Type': 'application/json',
    };
  }

  async createShortcut(request: CreateShortcutRequest): Promise<CreateShortcutResponse> {
    const instanceUrl = await this.getInstanceUrl();
    const headers = await this.getAuthHeaders();

    const response = await fetch(`${instanceUrl}/api/v1/shortcuts`, {
      method: 'POST',
      headers,
      credentials: 'include', // Include cookies for authentication
      body: JSON.stringify({
        shortcut: request
      })
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to create shortcut: ${response.status} ${errorText}`);
    }

    return response.json();
  }

  async sendCurrentUrl(url: string, title?: string): Promise<CreateShortcutResponse> {
    // Generate a simple name from the URL
    const urlObj = new URL(url);
    const hostname = urlObj.hostname.replace('www.', '');
    const name = this.generateShortcutName(hostname);

    const request: CreateShortcutRequest = {
      name,
      link: url,
      title: title || document.title || hostname,
      description: `Sent from browser extension`,
      tags: ['extension', 'auto-generated'],
      visibility: 'PRIVATE'
    };

    return this.createShortcut(request);
  }

  private generateShortcutName(hostname: string): string {
    // Generate a simple name based on hostname
    // Remove common TLDs and make it more readable
    const cleanHostname = hostname
      .replace(/\.(com|org|net|edu|gov|io|co)$/, '')
      .replace(/[^a-zA-Z0-9]/g, '-')
      .toLowerCase();

    // Add a timestamp to make it unique
    const timestamp = Date.now().toString().slice(-4);
    return `${cleanHostname}-${timestamp}`;
  }
}

export const apiService = new ApiService();