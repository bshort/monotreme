import { Card, Skeleton, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { workspaceServiceClient } from "@/grpcweb";
import { WorkspaceStats } from "@/types/proto/api/v1/workspace_service";
import Sparkline from "@/components/Sparkline";

const WorkspaceStatsSection = () => {
  const { t } = useTranslation();
  const [stats, setStats] = useState<WorkspaceStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        const workspaceStats = await workspaceServiceClient.getWorkspaceStats({});
        setStats(workspaceStats);
        setError(null);
      } catch (error) {
        console.error("Failed to fetch workspace stats:", error);
        setError("Failed to load workspace statistics");
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) {
    return (
      <div className="w-full flex flex-col justify-start items-start">
        <Typography level="title-lg" className="mb-4">
          {t("setting.workspace.stats")}
        </Typography>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
          {[...Array(4)].map((_, index) => (
            <Card key={index} className="p-4">
              <Skeleton variant="text" width="100%" height="1.5rem" />
              <Skeleton variant="text" width="60%" height="2rem" className="mt-2" />
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full flex flex-col justify-start items-start">
        <Typography level="title-lg" className="mb-4">
          {t("setting.workspace.stats")}
        </Typography>
        <Typography color="danger">{error}</Typography>
      </div>
    );
  }

  // Extract sparkline data from historical measurements
  const getSparklineData = (field: keyof typeof stats.historicalData[0]) => {
    if (!stats?.historicalData) return [];
    // Reverse to show chronological order (oldest to newest)
    return stats.historicalData.slice().reverse().map(measurement => measurement[field] || 0);
  };

  return (
    <div className="w-full flex flex-col justify-start items-start">
      <Typography level="title-lg" className="mb-4">
        {t("setting.workspace.stats")}
      </Typography>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 w-full">
        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("setting.workspace.stats.total-shortcuts")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalShortcuts || 0}
            </Typography>
            <Sparkline
              data={getSparklineData('shortcutsCount')}
              width={60}
              height={20}
              color="#3b82f6"
              className="ml-2"
            />
          </div>
        </Card>

        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("setting.workspace.stats.total-users")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalUsers || 0}
            </Typography>
            <Sparkline
              data={getSparklineData('usersCount')}
              width={60}
              height={20}
              color="#10b981"
              className="ml-2"
            />
          </div>
        </Card>

        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("setting.workspace.stats.total-collections")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalCollections || 0}
            </Typography>
            <Sparkline
              data={getSparklineData('collectionsCount')}
              width={60}
              height={20}
              color="#f59e0b"
              className="ml-2"
            />
          </div>
        </Card>

        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("setting.workspace.stats.total-hits")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalHits || 0}
            </Typography>
            <Sparkline
              data={getSparklineData('hitsCount')}
              width={60}
              height={20}
              color="#ef4444"
              className="ml-2"
            />
          </div>
        </Card>
      </div>
    </div>
  );
};

export default WorkspaceStatsSection;