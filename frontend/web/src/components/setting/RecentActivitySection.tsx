import { Card, Typography, Chip, Divider, Avatar } from "@mui/joy";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useUserStore } from "@/stores";
import { activityServiceClient } from "@/grpcweb";
import { GetRecentActivityResponse, RecentUser, RecentShortcut, RecentCollection, RecentClick, MostClickedShortcut } from "@/types/proto/api/v1/activity_service";
import Icon from "@/components/Icon";

const RecentActivitySection = () => {
  const userStore = useUserStore();
  const [loading, setLoading] = useState(true);
  const [recentActivity, setRecentActivity] = useState<GetRecentActivityResponse | null>(null);

  useEffect(() => {
    const loadRecentActivity = async () => {
      try {
        setLoading(true);
        const response = await activityServiceClient.getRecentActivity({ limit: 5 });
        setRecentActivity(response);
      } catch (error) {
        console.error("Failed to load recent activity:", error);
      } finally {
        setLoading(false);
      }
    };

    loadRecentActivity();
  }, []);

  const formatTimeAgo = (date: Date) => {
    const now = new Date();
    const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60));

    if (diffInHours < 1) return "Just now";
    if (diffInHours < 24) return `${diffInHours}h ago`;
    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 7) return `${diffInDays}d ago`;
    return date.toLocaleDateString();
  };

  const getInitials = (name: string) => {
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
  };

  if (loading) {
    return (
      <Card className="w-full p-6">
        <Typography level="title-md" className="mb-3">
          Recent Activity
        </Typography>
        <div className="flex items-center justify-center py-8">
          <Icon.Loader className="w-5 h-5 animate-spin mr-2" />
          <span className="text-sm text-gray-500">Loading recent activity...</span>
        </div>
      </Card>
    );
  }

  return (
    <Card className="w-full p-6">
      <Typography level="title-md" className="mb-4">
        Recent Activity
      </Typography>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {/* Recent Users */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.Users className="w-4 h-4 text-blue-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Users
            </Typography>
          </div>
          <div className="space-y-2">
            {recentActivity?.recentUsers && recentActivity.recentUsers.length > 0 ? (
              recentActivity.recentUsers.map((user) => (
                <div key={user.id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                  <Avatar size="sm" className="bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400">
                    {getInitials(user.nickname || user.email)}
                  </Avatar>
                  <div className="flex-1 min-w-0">
                    <Typography level="body-sm" className="truncate font-medium">
                      {user.nickname || user.email}
                    </Typography>
                    <Typography level="body-xs" className="text-gray-500 truncate">
                      {formatTimeAgo(new Date(user.createdTime))}
                    </Typography>
                  </div>
                </div>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent users
              </Typography>
            )}
          </div>
        </div>

        {/* Recent Shortcuts */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.Link className="w-4 h-4 text-green-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Shortcuts
            </Typography>
          </div>
          <div className="space-y-2">
            {recentActivity?.recentShortcuts && recentActivity.recentShortcuts.length > 0 ? (
              recentActivity.recentShortcuts.map((shortcut) => (
                <Link key={shortcut.id} to={`/shortcut/${shortcut.id}`} className="block">
                  <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <div className="w-8 h-8 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                      <Icon.Link className="w-4 h-4 text-green-600 dark:text-green-400" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <Typography level="body-sm" className="truncate font-medium">
                        {shortcut.title || shortcut.name}
                      </Typography>
                      <Typography level="body-xs" className="text-gray-500 truncate">
                        {formatTimeAgo(new Date(shortcut.createdTime))}
                      </Typography>
                    </div>
                  </div>
                </Link>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent shortcuts
              </Typography>
            )}
          </div>
        </div>

        {/* Recent Collections */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.FolderOpen className="w-4 h-4 text-purple-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Collections
            </Typography>
          </div>
          <div className="space-y-2">
            {recentActivity?.recentCollections && recentActivity.recentCollections.length > 0 ? (
              recentActivity.recentCollections.map((collection) => (
                <div key={collection.id} className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                  <div className="w-8 h-8 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                    <Icon.FolderOpen className="w-4 h-4 text-purple-600 dark:text-purple-400" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <Typography level="body-sm" className="truncate font-medium">
                      {collection.title || collection.name}
                    </Typography>
                    <Typography level="body-xs" className="text-gray-500 truncate">
                      {formatTimeAgo(new Date(collection.createdTime))}
                    </Typography>
                  </div>
                </div>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent collections
              </Typography>
            )}
          </div>
        </div>

        {/* Recent Clicks */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.Clock className="w-4 h-4 text-indigo-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Recent Clicks
            </Typography>
          </div>
          <div className="space-y-2">
            {recentActivity?.recentClicks && recentActivity.recentClicks.length > 0 ? (
              recentActivity.recentClicks.map((click) => (
                <Link key={click.shortcutId} to={`/shortcut/${click.shortcutId}`} className="block">
                  <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <div className="w-8 h-8 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                      <Icon.Clock className="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <Typography level="body-sm" className="truncate font-medium">
                        {click.shortcutTitle || click.shortcutName}
                      </Typography>
                      <Typography level="body-xs" className="text-gray-500 truncate">
                        {formatTimeAgo(new Date(click.clickedTime))}
                      </Typography>
                    </div>
                  </div>
                </Link>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No recent clicks
              </Typography>
            )}
          </div>
        </div>

        {/* Most Clicked */}
        <div>
          <div className="flex items-center gap-2 mb-3">
            <Icon.MousePointer className="w-4 h-4 text-orange-600" />
            <Typography level="title-sm" className="text-gray-900 dark:text-gray-100">
              Most Clicked
            </Typography>
          </div>
          <div className="space-y-2">
            {recentActivity?.mostClickedShortcuts && recentActivity.mostClickedShortcuts.length > 0 ? (
              recentActivity.mostClickedShortcuts.map((shortcut) => (
                <Link key={shortcut.id} to={`/shortcut/${shortcut.id}`} className="block">
                  <div className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <div className="w-8 h-8 bg-orange-100 dark:bg-orange-900/30 rounded-lg flex items-center justify-center flex-shrink-0">
                      <Icon.MousePointer className="w-4 h-4 text-orange-600 dark:text-orange-400" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <Typography level="body-sm" className="truncate font-medium">
                        {shortcut.title || shortcut.name}
                      </Typography>
                      <Typography level="body-xs" className="text-gray-500 truncate">
                        {shortcut.viewCount} clicks
                      </Typography>
                    </div>
                  </div>
                </Link>
              ))
            ) : (
              <Typography level="body-xs" className="text-gray-400 italic">
                No clicks yet
              </Typography>
            )}
          </div>
        </div>
      </div>
    </Card>
  );
};

export default RecentActivitySection;