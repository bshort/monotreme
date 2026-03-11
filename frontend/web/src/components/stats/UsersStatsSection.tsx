import { Card, Skeleton, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { workspaceServiceClient } from "@/grpcweb";
import { WorkspaceStats } from "@/types/proto/api/v1/workspace_service";
import Sparkline from "@/components/Sparkline";

const UsersStatsSection = () => {
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
        setError("Failed to load statistics");
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  const getUsersSparklineData = () => {
    if (!stats?.historicalData) return [];
    return stats.historicalData.slice().reverse().map(measurement => measurement.usersCount || 0);
  };

  if (loading) {
    return (
      <div className="w-full">
        <Typography level="title-lg" className="mb-4">
          {t("stats.users.title")}
        </Typography>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(3)].map((_, index) => (
            <Card key={index} className="p-6">
              <Skeleton variant="text" width="60%" height="1.5rem" />
              <Skeleton variant="text" width="40%" height="2.5rem" className="mt-2" />
              <Skeleton variant="rectangular" width="100%" height="60px" className="mt-4" />
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full">
        <Typography level="title-lg" className="mb-4">
          {t("stats.users.title")}
        </Typography>
        <Typography color="danger">{error}</Typography>
      </div>
    );
  }

  return (
    <div className="w-full">
      <Typography level="title-lg" className="mb-4">
        {t("stats.users.title")}
      </Typography>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card className="p-6">
          <Typography level="body-sm" color="neutral" className="mb-2">
            {t("stats.users.total")}
          </Typography>
          <Typography level="h1" className="mb-4">
            {stats?.totalUsers || 0}
          </Typography>
          <div className="h-16">
            <Sparkline
              data={getUsersSparklineData()}
              width={240}
              height={64}
              color="#10b981"
            />
          </div>
        </Card>

        <Card className="p-6">
          <Typography level="body-sm" color="neutral" className="mb-2">
            {t("stats.users.growth")}
          </Typography>
          <Typography level="h2" className="mb-4" color="success">
            {getUsersSparklineData().length >= 2
              ? `+${Math.max(0, (getUsersSparklineData()[getUsersSparklineData().length - 1] || 0) - (getUsersSparklineData()[0] || 0))}`
              : "0"
            }
          </Typography>
          <Typography level="body-sm" color="neutral">
            {t("stats.users.growth-period")}
          </Typography>
        </Card>

        <Card className="p-6">
          <Typography level="body-sm" color="neutral" className="mb-2">
            {t("stats.users.activity")}
          </Typography>
          <Typography level="h2" className="mb-4">
            {getUsersSparklineData().length > 0 ? getUsersSparklineData()[getUsersSparklineData().length - 1] : 0}
          </Typography>
          <Typography level="body-sm" color="neutral">
            {t("stats.users.current-count")}
          </Typography>
        </Card>
      </div>
    </div>
  );
};

export default UsersStatsSection;