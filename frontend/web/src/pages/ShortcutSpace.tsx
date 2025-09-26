import { Button } from "@mui/joy";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";
import { useParams, useSearchParams } from "react-router-dom";
import CreateShortcutDrawer from "@/components/CreateShortcutDrawer";
import { isURL } from "@/helpers/utils";
import useNavigateTo from "@/hooks/useNavigateTo";
import { useShortcutStore, useUserStore, useWorkspaceStore } from "@/stores";
import { Shortcut } from "@/types/proto/api/v1/shortcut_service";

const ShortcutSpace = () => {
  const params = useParams();
  const [searchParams] = useSearchParams();
  const shortcutName = params["*"] || "";
  const navigateTo = useNavigateTo();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();
  const shortcutStore = useShortcutStore();
  const workspaceStore = useWorkspaceStore();
  const [shortcut, setShortcut] = useState<Shortcut>();
  const [loading, setLoading] = useState(true);
  const [showCreateShortcutDrawer, setShowCreateShortcutDrawer] = useState(false);

  // Check if the current route matches the workspace shortcut prefix
  const requestedPrefix = params["prefix"] || "s"; // "s" for legacy route, or the dynamic prefix
  const currentShortcutPrefix = workspaceStore.getShortcutPrefix();

  useEffect(() => {
    (async () => {
      // Only fetch shortcut if the requested prefix matches the current shortcut prefix
      if (requestedPrefix !== currentShortcutPrefix) {
        setLoading(false);
        return;
      }

      try {
        const shortcut = await shortcutStore.fetchShortcutByName(shortcutName);
        setShortcut(shortcut);
      } catch (error: any) {
        console.error(error);
        toast.error(error.details);
      }
      setLoading(false);
    })();
  }, [shortcutName, requestedPrefix, currentShortcutPrefix]);

  if (loading) {
    return null;
  }

  // If the requested prefix doesn't match the current shortcut prefix, show 404
  if (requestedPrefix !== currentShortcutPrefix) {
    navigateTo("/404");
    return null;
  }

  if (!shortcut) {
    if (!currentUser) {
      navigateTo("/404");
      return null;
    }

    // If shortcut is not found, prompt user to create it.
    return (
      <>
        <div className="w-full h-[100svh] flex flex-col justify-center items-center p-4">
          <p className="text-xl">
            Shortcut <span className="font-mono">{shortcutName}</span> Not Found.
          </p>
          <div className="mt-4">
            <Button variant="plain" size="sm" onClick={() => setShowCreateShortcutDrawer(true)}>
              👉 Click here to create it
            </Button>
          </div>
        </div>
        {showCreateShortcutDrawer && (
          <CreateShortcutDrawer
            initialShortcut={{ name: shortcutName }}
            onClose={() => setShowCreateShortcutDrawer(false)}
            onConfirm={() => navigateTo("/")}
          />
        )}
      </>
    );
  }

  // If shortcut is a URL, redirect to it directly.
  if (isURL(shortcut.link)) {
    window.document.title = "Redirecting...";
    const url = new URL(shortcut.link);
    searchParams.forEach((value, key) => {
      url.searchParams.append(key, value);
    });
    window.location.href = url.toString();
    return null;
  }

  // Otherwise, render the shortcut link as plain text.
  return <div>{shortcut.link}</div>;
};

export default ShortcutSpace;
