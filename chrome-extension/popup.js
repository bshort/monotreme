// Popup script for Monotreme URL Sender

document.addEventListener('DOMContentLoaded', async () => {
    const currentUrlElement = document.getElementById('currentUrl');
    const sendButton = document.getElementById('sendButton');
    const sendButtonText = document.getElementById('sendButtonText');
    const spinner = document.getElementById('spinner');
    const statusElement = document.getElementById('status');
    const serverUrlInput = document.getElementById('serverUrl');
    const saveConfigButton = document.getElementById('saveConfig');

    // Load configuration
    async function loadConfig() {
        return new Promise((resolve) => {
            chrome.runtime.sendMessage({ action: 'getConfig' }, (response) => {
                resolve(response);
            });
        });
    }

    // Save configuration
    async function saveConfig() {
        const serverUrl = serverUrlInput.value.trim();
        if (!serverUrl) {
            showStatus('Please enter a server URL', 'error');
            return;
        }

        return new Promise((resolve) => {
            chrome.runtime.sendMessage({
                action: 'setConfig',
                config: { serverUrl }
            }, (response) => {
                if (response.success) {
                    showStatus('Configuration saved successfully', 'success');
                } else {
                    showStatus('Failed to save configuration', 'error');
                }
                resolve(response);
            });
        });
    }

    // Get current tab info
    async function getCurrentTab() {
        return new Promise((resolve) => {
            chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
                resolve(tabs[0]);
            });
        });
    }

    // Show status message
    function showStatus(message, type) {
        statusElement.textContent = message;
        statusElement.className = `status ${type}`;
        statusElement.classList.remove('hidden');

        // Auto-hide success messages after 3 seconds
        if (type === 'success') {
            setTimeout(() => {
                statusElement.classList.add('hidden');
            }, 3000);
        }
    }

    // Set loading state
    function setLoading(loading) {
        sendButton.disabled = loading;
        if (loading) {
            sendButtonText.textContent = 'Sending...';
            spinner.classList.remove('hidden');
        } else {
            sendButtonText.textContent = 'Send Current URL';
            spinner.classList.add('hidden');
        }
    }

    // Send current URL
    async function sendCurrentUrl() {
        setLoading(true);
        statusElement.classList.add('hidden');

        try {
            const response = await new Promise((resolve) => {
                chrome.runtime.sendMessage({ action: 'sendCurrentUrl' }, resolve);
            });

            if (response.success) {
                showStatus(`Shortcut created: ${response.data.name || 'Success'}`, 'success');
            } else {
                showStatus(`Error: ${response.error}`, 'error');
            }
        } catch (error) {
            showStatus(`Error: ${error.message}`, 'error');
        } finally {
            setLoading(false);
        }
    }

    // Initialize
    try {
        // Load and display current tab URL
        const currentTab = await getCurrentTab();
        if (currentTab && currentTab.url) {
            currentUrlElement.textContent = currentTab.url;
        } else {
            currentUrlElement.textContent = 'Unable to get current URL';
        }

        // Load current configuration
        const config = await loadConfig();
        serverUrlInput.value = config.serverUrl || 'http://localhost:5231';

        // Event listeners
        sendButton.addEventListener('click', sendCurrentUrl);
        saveConfigButton.addEventListener('click', saveConfig);

        // Allow Enter key to save config
        serverUrlInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                saveConfig();
            }
        });

    } catch (error) {
        console.error('Failed to initialize popup:', error);
        showStatus('Failed to initialize extension', 'error');
    }
});