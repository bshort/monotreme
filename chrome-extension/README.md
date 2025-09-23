# Monotreme URL Sender Chrome Extension

A simple Chrome extension that allows you to send the current tab's URL to your Monotreme backend service to create shortcuts.

## Features

- **One-click URL sending**: Send the current tab's URL to your Monotreme server with a single click
- **Configurable server URL**: Set your Monotreme instance URL in the extension settings
- **Automatic shortcut naming**: Generates meaningful shortcut names based on the website hostname
- **Visual feedback**: Shows loading states and success/error messages
- **Private by default**: Creates shortcuts with private visibility

## Installation

### Development Installation

1. Open Chrome and navigate to `chrome://extensions/`
2. Enable "Developer mode" (toggle in the top right)
3. Click "Load unpacked" and select the `chrome-extension` directory
4. The extension will appear in your extensions list

### Production Installation

Package the extension using Chrome's built-in packaging or submit to the Chrome Web Store.

## Setup

1. **Configure your Monotreme server URL**:
   - Click the extension icon in your toolbar
   - Enter your Monotreme instance URL (e.g., `http://localhost:5231` or `https://your-monotreme-server.com`)
   - Click "Save Configuration"

2. **Ensure you're authenticated**:
   - Open your Monotreme instance in a browser tab
   - Sign in to your account
   - The extension uses cookies for authentication, so you need to be logged in

## Usage

1. Navigate to any webpage you want to create a shortcut for
2. Click the Monotreme extension icon in your toolbar
3. Review the current URL displayed in the popup
4. Click "Send Current URL"
5. Wait for the success confirmation

The extension will:
- Generate a unique shortcut name based on the hostname
- Use the page title as the shortcut title
- Add tags for easy identification (`chrome-extension`, `auto-generated`)
- Set the visibility to `PRIVATE`

## API Integration

The extension communicates with your Monotreme backend using the following API:

- **Endpoint**: `POST /api/v1/shortcuts`
- **Authentication**: Cookies (requires active session)
- **Payload**:
  ```json
  {
    "shortcut": {
      "name": "generated-name-1234",
      "link": "https://example.com/page",
      "title": "Page Title",
      "description": "Sent from Chrome extension",
      "tags": ["chrome-extension", "auto-generated"],
      "visibility": "PRIVATE"
    }
  }
  ```

## Files Structure

```
chrome-extension/
├── manifest.json       # Extension manifest (v3)
├── background.js       # Background service worker
├── popup.html         # Extension popup UI
├── popup.js           # Popup logic and event handling
├── icon16.png         # Extension icon (16x16)
├── icon48.png         # Extension icon (48x48)
├── icon128.png        # Extension icon (128x128)
└── README.md          # This file
```

## Permissions

The extension requires the following permissions:

- `activeTab`: To access the current tab's URL and title
- `storage`: To store configuration (server URL)

## Error Handling

The extension handles various error scenarios:

- **Network errors**: Shows network connection issues
- **Authentication errors**: Indicates when you need to sign in
- **Server errors**: Displays server response errors
- **Configuration errors**: Prompts to configure server URL

## Development

To modify the extension:

1. Make changes to the source files
2. Go to `chrome://extensions/`
3. Click the refresh button for the extension
4. Test your changes

## Security Notes

- The extension only sends URLs to your configured Monotreme server
- Authentication is handled through browser cookies
- No sensitive data is stored locally except for the server URL
- All shortcuts are created as private by default

## Troubleshooting

**Extension not working?**
- Ensure you're signed in to your Monotreme instance
- Check that the server URL is configured correctly
- Verify your Monotreme server is running and accessible

**Network errors?**
- Check if your Monotreme server is running
- Verify the server URL is correct and accessible
- Check for any firewall or CORS issues

**Authentication errors?**
- Sign in to your Monotreme instance in a browser tab
- The extension relies on browser cookies for authentication