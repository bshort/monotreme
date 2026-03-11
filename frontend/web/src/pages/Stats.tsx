import { Divider, Typography } from "@mui/joy";
import { useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useUserStore } from "@/stores";
import { Role } from "@/types/proto/api/v1/user_service";
import UsersStatsSection from "@/components/stats/UsersStatsSection";
import ShortcutsStatsSection from "@/components/stats/ShortcutsStatsSection";
import CollectionsStatsSection from "@/components/stats/CollectionsStatsSection";
import SiteMetricsStatsSection from "@/components/stats/SiteMetricsStatsSection";

const Stats = () => {
  const { t } = useTranslation();
  const currentUser = useUserStore().getCurrentUser();
  const isAdmin = currentUser.role === Role.ADMIN;

  useEffect(() => {
    if (!isAdmin) {
      window.location.href = "/";
    }
  }, [isAdmin]);

  if (!isAdmin) {
    return null;
  }

  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 py-6 flex flex-col justify-start items-start gap-y-8">
      <div className="w-full">
        <Typography level="h2" className="mb-2">
          {t("stats.title")}
        </Typography>
        <Typography level="body-lg" color="neutral">
          {t("stats.description")}
        </Typography>
      </div>

      <UsersStatsSection />
      <Divider />

      <ShortcutsStatsSection />
      <Divider />

      <CollectionsStatsSection />
      <Divider />

      <SiteMetricsStatsSection />
    </div>
  );
};

export default Stats;