import { Button, CssVarsProvider } from "@mui/joy";
import { Toaster, toast } from "react-hot-toast";
import { useState } from "react";
import Icon from "@/components/Icon";
import Logo from "@/components/Logo";
import { StorageContextProvider, useStorageContext } from "./context";
import { apiService } from "./services/api";
import "./style.css";

const IndexPopup = () => {
  const context = useStorageContext();
  const isInitialized = context.instanceUrl;
  const [isLoading, setIsLoading] = useState(false);

  const handleSettingButtonClick = () => {
    chrome.runtime.openOptionsPage();
  };

  const handleRefreshButtonClick = () => {
    chrome.runtime.reload();
    chrome.browserAction.setPopup({ popup: "" });
  };

  const handleSendUrlButtonClick = async () => {
    if (!isInitialized) {
      toast.error("Please set your instance URL first");
      return;
    }

    setIsLoading(true);
    try {
      // Get the current active tab
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

      if (!tab.url) {
        throw new Error("Could not get current tab URL");
      }

      // Send the URL to the backend
      const result = await apiService.sendCurrentUrl(tab.url, tab.title);

      toast.success(`Shortcut created: ${result.name}`);
    } catch (error) {
      console.error("Failed to send URL:", error);
      toast.error(error instanceof Error ? error.message : "Failed to send URL");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full min-w-[512px] px-4 pt-4">
      <div className="w-full flex flex-row justify-between items-center">
        <div className="flex flex-row justify-start items-center dark:text-gray-400">
          <Logo className="w-6 h-auto mr-1" />
          <span className="">Monotreme</span>
        </div>
      </div>

      <div className="w-full mt-4">
        {isInitialized ? (
          <>
            <p className="w-full mb-2">
              <span>Your instance URL is </span>
              <a
                className="inline-flex flex-row justify-start items-center underline text-blue-600 hover:opacity-80"
                href={context.instanceUrl}
                target="_blank"
              >
                <span className="mr-1">{context.instanceUrl}</span>
                <Icon.ExternalLink className="w-4 h-auto" />
              </a>
            </p>
            <div className="w-full mb-3">
              <Button
                size="md"
                color="primary"
                disabled={isLoading}
                onClick={handleSendUrlButtonClick}
                sx={{ width: '100%' }}
              >
                {isLoading ? (
                  <Icon.Loader2 className="w-4 h-auto mr-2 animate-spin" />
                ) : (
                  <Icon.Share className="w-4 h-auto mr-2" />
                )}
                {isLoading ? 'Sending...' : 'Send Current URL'}
              </Button>
            </div>
            <div className="w-full flex flex-row justify-between items-center mb-2">
              <div className="flex flex-row justify-start items-center gap-2">
                <Button size="sm" variant="outlined" color="neutral" onClick={handleSettingButtonClick}>
                  <Icon.Settings className="w-5 h-auto text-gray-500 dark:text-gray-400 mr-1" />
                  Setting
                </Button>
                <Button
                  size="sm"
                  variant="outlined"
                  color="neutral"
                  component="a"
                  href="https://github.com/bshort/monotreme"
                  target="_blank"
                >
                  <Icon.Github className="w-5 h-auto text-gray-500 dark:text-gray-400 mr-1" />
                  GitHub
                </Button>
              </div>
            </div>
          </>
        ) : (
          <div className="w-full flex flex-col justify-start items-center">
            <Icon.Cookie strokeWidth={1} className="w-20 h-auto mb-4 text-gray-400" />
            <p className="dark:text-gray-400">Please set your instance URL first.</p>
            <div className="w-full flex flex-row justify-center items-center py-4">
              <Button size="sm" color="primary" onClick={handleSettingButtonClick}>
                <Icon.Settings className="w-5 h-auto mr-1" /> Go to Setting
              </Button>
              <span className="mx-2 dark:text-gray-400">Or</span>
              <Button size="sm" variant="outlined" color="neutral" onClick={handleRefreshButtonClick}>
                <Icon.RefreshCcw className="w-5 h-auto mr-1" /> Refresh
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

const Popup = () => {
  return (
    <StorageContextProvider>
      <CssVarsProvider>
        <IndexPopup />
        <Toaster position="top-center" />
      </CssVarsProvider>
    </StorageContextProvider>
  );
};

export default Popup;
