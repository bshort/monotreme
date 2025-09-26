import { useColorScheme } from "@mui/joy";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Outlet } from "react-router-dom";
import Footer from "@/components/Footer";
import Header from "@/components/Header";
import Navigator from "@/components/Navigator";
import useNavigateTo from "@/hooks/useNavigateTo";
import { useUserStore, useWorkspaceStore } from "@/stores";

const Root: React.FC = () => {
  const navigateTo = useNavigateTo();
  const { setMode } = useColorScheme();
  const { i18n } = useTranslation();
  const userStore = useUserStore();
  const workspaceStore = useWorkspaceStore();
  const currentUser = userStore.getCurrentUser();
  const currentUserSetting = userStore.getCurrentUserSetting();
  const isInitialized = Boolean(currentUser) && Boolean(currentUserSetting);

  useEffect(() => {
    // Wait for workspace profile to be loaded before making decisions
    if (!workspaceStore.profile.mode) {
      // Profile not loaded yet, don't redirect
      return;
    }

    // Always redirect non-logged-in users to landing page
    if (!currentUser) {
      navigateTo("/landing", {
        replace: true,
      });
      return;
    }

    // Prepare user setting.
    userStore.fetchUserSetting(currentUser.id);
  }, [workspaceStore.profile, currentUser]);

  useEffect(() => {
    if (!currentUserSetting) {
      return;
    }

    i18n.changeLanguage(currentUserSetting.general?.locale || "en");

    if (currentUserSetting.general?.colorTheme === "LIGHT") {
      setMode("light");
    } else if (currentUserSetting.general?.colorTheme === "DARK") {
      setMode("dark");
    } else {
      setMode("system");
    }
  }, [currentUserSetting]);

  return (
    isInitialized && (
      <div className="w-full min-h-screen flex flex-col justify-start items-start dark:bg-zinc-900">
        <Header />
        <Navigator />
        <div className="flex-1 w-full">
          <Outlet />
        </div>
        <Footer />
      </div>
    )
  );
};

export default Root;
