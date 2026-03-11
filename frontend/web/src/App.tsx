import { useColorScheme } from "@mui/joy";
import { useEffect } from "react";
import { Outlet, useLocation } from "react-router-dom";
import Footer from "@/components/Footer";
import { useWorkspaceStore } from "@/stores";
import { FeatureType } from "./stores/workspace";

function App() {
  const { mode: colorScheme } = useColorScheme();
  const workspaceStore = useWorkspaceStore();
  const location = useLocation();

  // Check if we're on a page that uses the Root layout (which already has a footer)
  const isRootLayoutPage = location.pathname.startsWith('/') &&
    !location.pathname.startsWith('/landing') &&
    !location.pathname.startsWith('/auth') &&
    !location.pathname.startsWith('/quick-save') &&
    !location.pathname.match(/^\/[^\/]+\/[^\/]+$/); // Not shortcut/collection routes


  useEffect(() => {
    const styleEl = document.createElement("style");
    styleEl.innerHTML = workspaceStore.setting.customStyle;
    styleEl.setAttribute("type", "text/css");
    document.body.insertAdjacentElement("beforeend", styleEl);
  }, [workspaceStore.setting.customStyle]);

  useEffect(() => {
    const hasCustomBranding = workspaceStore.checkFeatureAvailable(FeatureType.CustomeBranding);
    if (!hasCustomBranding || !workspaceStore.setting.branding) {
      return;
    }

    const favicon = document.querySelector("link[rel='icon']") as HTMLLinkElement;
    favicon.href = new TextDecoder().decode(workspaceStore.setting.branding);
  }, [workspaceStore.setting.branding]);

  useEffect(() => {
    const root = document.documentElement;
    if (colorScheme === "light") {
      root.classList.remove("dark");
    } else if (colorScheme === "dark") {
      root.classList.add("dark");
    } else {
      const darkMediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
      if (darkMediaQuery.matches) {
        root.classList.add("dark");
      } else {
        root.classList.remove("dark");
      }

      const handleColorSchemeChange = (e: MediaQueryListEvent) => {
        if (e.matches) {
          root.classList.add("dark");
        } else {
          root.classList.remove("dark");
        }
      };
      try {
        darkMediaQuery.addEventListener("change", handleColorSchemeChange);
      } catch (error) {
        console.error("failed to initial color scheme listener", error);
      }

      return () => {
        darkMediaQuery.removeEventListener("change", handleColorSchemeChange);
      };
    }
  }, [colorScheme]);

  return (
    <div className="min-h-screen flex flex-col">
      <div className="flex-1">
        <Outlet />
      </div>
    </div>
  );
}

export default App;
