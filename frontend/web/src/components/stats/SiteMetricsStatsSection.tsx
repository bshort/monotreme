import { Card, Skeleton, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { workspaceServiceClient } from "@/grpcweb";
import { WorkspaceStats } from "@/types/proto/api/v1/workspace_service";
import Sparkline from "@/components/Sparkline";

const SiteMetricsStatsSection = () => {
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

  const getHitsSparklineData = () => {
    if (!stats?.historicalData) return [];
    return stats.historicalData.slice().reverse().map(measurement => measurement.hitsCount || 0);
  };

  const getTotalItemsSparklineData = () => {
    if (!stats?.historicalData) return [];
    return stats.historicalData.slice().reverse().map(measurement =>
      (measurement.shortcutsCount || 0) + (measurement.collectionsCount || 0)
    );
  };

  if (loading) {
    return (
      <div className="w-full">
        <Typography level="title-lg" className="mb-4">
          {t("stats.site-metrics.title")}
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
          {t("stats.site-metrics.title")}
        </Typography>
        <Typography color="danger">{error}</Typography>
      </div>
    );
  }

  return (
    <div className="w-full">
      <Typography level="title-lg" className="mb-4">
        {t("stats.site-metrics.title")}
      </Typography>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card className="p-6">
          <Typography level="body-sm" color="neutral" className="mb-2">
            {t("stats.site-metrics.total-hits")}
          </Typography>
          <Typography level="h1" className="mb-4">
            {stats?.totalHits || 0}
          </Typography>
          <div className="h-16">
            <Sparkline
              data={getHitsSparklineData()}
              width={240}
              height={64}
              color="#ef4444"
            />
          </div>
        </Card>

        <Card className="p-6">
          <Typography level="body-sm" color="neutral" className="mb-2">
            {t("stats.site-metrics.total-content")}
          </Typography>
          <Typography level="h2" className="mb-4">
            {((stats?.totalShortcuts || 0) + (stats?.totalCollections || 0))}
          </Typography>
          <div className="h-16">
            <Sparkline
              data={getTotalItemsSparklineData()}
              width={240}
              height={64}
              color="#8b5cf6"
            />
          </div>
        </Card>

        <Card className="p-6">
          <Typography level="body-sm" color="neutral" className="mb-2">
            {t("stats.site-metrics.engagement")}
          </Typography>
          <Typography level="h2" className="mb-4">
            {stats?.totalShortcuts && stats.totalShortcuts > 0
              ? Math.round((stats?.totalHits || 0) / stats.totalShortcuts * 100) / 100
              : 0
            }
          </Typography>
          <Typography level="body-sm" color="neutral">
            {t("stats.site-metrics.hits-per-shortcut")}
          </Typography>
        </Card>
      </div>

      <div className="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("stats.site-metrics.total-shortcuts")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalShortcuts || 0}
            </Typography>
            <Sparkline
              data={stats?.historicalData?.slice().reverse().map(m => m.shortcutsCount || 0) || []}
              width={60}
              height={20}
              color="#3b82f6"
              className="ml-2"
            />
          </div>
        </Card>

        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("stats.site-metrics.total-users")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalUsers || 0}
            </Typography>
            <Sparkline
              data={stats?.historicalData?.slice().reverse().map(m => m.usersCount || 0) || []}
              width={60}
              height={20}
              color="#10b981"
              className="ml-2"
            />
          </div>
        </Card>

        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("stats.site-metrics.total-collections")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalCollections || 0}
            </Typography>
            <Sparkline
              data={stats?.historicalData?.slice().reverse().map(m => m.collectionsCount || 0) || []}
              width={60}
              height={20}
              color="#f59e0b"
              className="ml-2"
            />
          </div>
        </Card>

        <Card className="p-4">
          <Typography level="body-sm" color="neutral">
            {t("stats.site-metrics.total-hits")}
          </Typography>
          <div className="flex items-center justify-between mt-2">
            <Typography level="title-lg">
              {stats?.totalHits || 0}
            </Typography>
            <Sparkline
              data={stats?.historicalData?.slice().reverse().map(m => m.hitsCount || 0) || []}
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

export default SiteMetricsStatsSection;