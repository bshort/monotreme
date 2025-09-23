// Background script for Monotreme URL Sender

// Default configuration
const DEFAULT_CONFIG = {
  serverUrl: 'http://localhost:5231',
  apiEndpoint: '/api/v1/shortcuts'
};

// Get configuration from storage
async function getConfig() {
  const result = await chrome.storage.sync.get(['serverUrl']);
  return {
    serverUrl: result.serverUrl || DEFAULT_CONFIG.serverUrl
  };
}

// Send URL to Monotreme backend
async function sendUrlToServer(url, title) {
  const config = await getConfig();
  const endpoint = `${config.serverUrl}${DEFAULT_CONFIG.apiEndpoint}`;

  // Generate a simple name from the URL
  const urlObj = new URL(url);
  const hostname = urlObj.hostname.replace('www.', '');
  const timestamp = Date.now().toString().slice(-4);
  const name = `${hostname.replace(/[^a-zA-Z0-9]/g, '-').toLowerCase()}-${timestamp}`;

  const shortcutData = {
    shortcut: {
      name: name,
      link: url,
      title: title || hostname,
      description: 'Sent from Chrome extension',
      tags: ['chrome-extension', 'auto-generated'],
      visibility: 'PRIVATE'
    }
  };

  try {
    const response = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // Include cookies for authentication
      body: JSON.stringify(shortcutData)
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`HTTP ${response.status}: ${errorText}`);
    }

    const result = await response.json();
    return { success: true, data: result };
  } catch (error) {
    console.error('Failed to send URL to server:', error);
    return { success: false, error: error.message };
  }
}

// Handle messages from popup
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === 'sendCurrentUrl') {
    // Get current active tab and send its URL
    chrome.tabs.query({ active: true, currentWindow: true }, async (tabs) => {
      if (tabs[0]) {
        const result = await sendUrlToServer(tabs[0].url, tabs[0].title);
        sendResponse(result);
      } else {
        sendResponse({ success: false, error: 'No active tab found' });
      }
    });

    // Return true to indicate we'll send response asynchronously
    return true;
  }

  if (request.action === 'getConfig') {
    getConfig().then(config => {
      sendResponse(config);
    });
    return true;
  }

  if (request.action === 'setConfig') {
    chrome.storage.sync.set(request.config, () => {
      sendResponse({ success: true });
    });
    return true;
  }
});

// Install event
chrome.runtime.onInstalled.addListener(() => {
  console.log('Monotreme URL Sender extension installed');
});