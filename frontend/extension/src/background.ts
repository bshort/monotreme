import { Storage } from "@plasmohq/storage";
import { apiService } from "./services/api";
import { getShortcutUrl } from "./utils/shortcut";

const storage = new Storage();
const urlRegex = /https?:\/\/s\/(.+)/;

chrome.webRequest.onBeforeRequest.addListener(
  (param) => {
    (async () => {
      if (!param.url) {
        return;
      }

      const shortcutName = getShortcutNameFromUrl(param.url);
      if (shortcutName) {
        const instanceUrl = (await storage.getItem<string>("instance_url")) || "";
        const shortcutPath = await getShortcutUrl(shortcutName);
        const url = new URL(shortcutPath, instanceUrl);
        return chrome.tabs.update({ url: url.toString() });
      }
    })();
  },
  { urls: ["*://s/*", "*://*/search*", "*://*/s*", "*://duckduckgo.com/*"] },
);

const getShortcutNameFromUrl = (urlString: string) => {
  const matchResult = urlRegex.exec(urlString);
  if (matchResult === null) {
    return getShortcutNameFromSearchUrl(urlString);
  }
  return matchResult[1];
};

const getShortcutNameFromSearchUrl = (urlString: string) => {
  const url = new URL(urlString);
  if ((url.hostname.endsWith("google.com") || url.hostname.endsWith("bing.com")) && url.pathname === "/search") {
    const params = new URLSearchParams(url.search);
    const shortcutName = params.get("q");
    if (typeof shortcutName === "string" && shortcutName.startsWith("s/")) {
      return shortcutName.slice(2);
    }
  } else if (url.hostname.endsWith("baidu.com") && url.pathname === "/s") {
    const params = new URLSearchParams(url.search);
    const shortcutName = params.get("wd");
    if (typeof shortcutName === "string" && shortcutName.startsWith("s/")) {
      return shortcutName.slice(2);
    }
  } else if (url.hostname.endsWith("duckduckgo.com") && url.pathname === "/") {
    const params = new URLSearchParams(url.search);
    const shortcutName = params.get("q");
    if (typeof shortcutName === "string" && shortcutName.startsWith("s/")) {
      return shortcutName.slice(2);
    }
  }
  return "";
};

// Handle messages from popup or other parts of the extension
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "sendCurrentUrl") {
    (async () => {
      try {
        const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

        if (!tab.url) {
          throw new Error("Could not get current tab URL");
        }

        const result = await apiService.sendCurrentUrl(tab.url, tab.title);
        sendResponse({ success: true, data: result });
      } catch (error) {
        sendResponse({
          success: false,
          error: error instanceof Error ? error.message : "Unknown error"
        });
      }
    })();

    // Return true to indicate we'll respond asynchronously
    return true;
  }
});
